// Created and Developed by Farhan Kertadiwangsa
package evaluator

import (
        "bufio"
        "fmt"
        "math"
        "math/rand"
        "os"
        "sort"
        "strconv"
        "strings"
        "time"
        "unicode/utf8"

        "ky/object"
)

var Builtins map[string]*object.Builtin

func init() {
        Builtins = map[string]*object.Builtin{

        "tampilkan": {Name: "tampilkan", Fn: func(args ...object.Object) object.Object {
                parts := make([]string, len(args))
                for i, a := range args {
                        parts[i] = a.Inspect()
                }
                fmt.Println(strings.Join(parts, " "))
                return KOSONG
        }},

        // cetak(nilai...) — alias tampilkan
        "cetak": {Name: "cetak", Fn: func(args ...object.Object) object.Object {
                parts := make([]string, len(args))
                for i, a := range args {
                        parts[i] = a.Inspect()
                }
                fmt.Println(strings.Join(parts, " "))
                return KOSONG
        }},

        // tulis(nilai...) — cetak tanpa baris baru
        "tulis": {Name: "tulis", Fn: func(args ...object.Object) object.Object {
                parts := make([]string, len(args))
                for i, a := range args {
                        parts[i] = a.Inspect()
                }
                fmt.Print(strings.Join(parts, " "))
                return KOSONG
        }},

        // format(template, nilai...) — seperti fmt.Sprintf
        "format": {Name: "format", Fn: func(args ...object.Object) object.Object {
                if len(args) < 1 {
                        return object.NewError("format() butuh minimal 1 argumen")
                }
                tmpl, ok := args[0].(*object.String)
                if !ok {
                        return object.NewError("format() argumen pertama harus teks")
                }
                fmtArgs := make([]interface{}, len(args)-1)
                for i, a := range args[1:] {
                        switch v := a.(type) {
                        case *object.Integer:
                                fmtArgs[i] = v.Value
                        case *object.Float:
                                fmtArgs[i] = v.Value
                        case *object.String:
                                fmtArgs[i] = v.Value
                        case *object.Boolean:
                                fmtArgs[i] = v.Value
                        default:
                                fmtArgs[i] = a.Inspect()
                        }
                }
                return &object.String{Value: fmt.Sprintf(tmpl.Value, fmtArgs...)}
        }},

        // baca() — baca satu baris input dari pengguna
        "baca": {Name: "baca", Fn: func(args ...object.Object) object.Object {
                if len(args) > 0 {
                        if prompt, ok := args[0].(*object.String); ok {
                                fmt.Print(prompt.Value)
                        }
                }
                scanner := bufio.NewScanner(os.Stdin)
                if scanner.Scan() {
                        return &object.String{Value: scanner.Text()}
                }
                return &object.String{Value: ""}
        }},

        // baca_berkas(path) — baca isi file
        "baca_berkas": {Name: "baca_berkas", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("baca_berkas() butuh 1 argumen")
                }
                path, ok := args[0].(*object.String)
                if !ok {
                        return object.NewError("baca_berkas() butuh path teks")
                }
                data, err := os.ReadFile(path.Value)
                if err != nil {
                        return object.NewError("gagal baca file '%s': %v", path.Value, err)
                }
                return &object.String{Value: string(data)}
        }},

        // tulis_berkas(path, isi) — tulis string ke file
        "tulis_berkas": {Name: "tulis_berkas", Fn: func(args ...object.Object) object.Object {
                if len(args) != 2 {
                        return object.NewError("tulis_berkas() butuh 2 argumen")
                }
                path, ok1 := args[0].(*object.String)
                isi, ok2 := args[1].(*object.String)
                if !ok1 || !ok2 {
                        return object.NewError("tulis_berkas() butuh path dan isi berupa teks")
                }
                err := os.WriteFile(path.Value, []byte(isi.Value), 0644)
                if err != nil {
                        return object.NewError("gagal tulis file: %v", err)
                }
                return KOSONG
        }},


        // jenis(nilai) — kembalikan nama tipe sebagai teks
        "jenis": {Name: "jenis", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("jenis() butuh 1 argumen")
                }
                tipeMap := map[object.ObjectType]string{
                        object.INTEGER_OBJ:  "bilangan",
                        object.FLOAT_OBJ:    "desimal",
                        object.STRING_OBJ:   "teks",
                        object.BOOLEAN_OBJ:  "logika",
                        object.KOSONG_OBJ:   "kosong",
                        object.ARRAY_OBJ:    "larik",
                        object.MAP_OBJ:      "kamus",
                        object.FUNCTION_OBJ: "fungsi",
                        object.BUILTIN_OBJ:  "fungsi_bawaan",
                }
                if nama, ok := tipeMap[args[0].Type()]; ok {
                        return &object.String{Value: nama}
                }
                return &object.String{Value: string(args[0].Type())}
        }},

        // adalah_bilangan(nilai) — cek apakah integer
        "adalah_bilangan": {Name: "adalah_bilangan", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("adalah_bilangan() butuh 1 argumen")
                }
                return nativeBool(args[0].Type() == object.INTEGER_OBJ)
        }},

        // adalah_desimal(nilai) — cek apakah float
        "adalah_desimal": {Name: "adalah_desimal", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("adalah_desimal() butuh 1 argumen")
                }
                return nativeBool(args[0].Type() == object.FLOAT_OBJ)
        }},

        // adalah_teks(nilai) — cek apakah string
        "adalah_teks": {Name: "adalah_teks", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("adalah_teks() butuh 1 argumen")
                }
                return nativeBool(args[0].Type() == object.STRING_OBJ)
        }},

        // adalah_logika(nilai) — cek apakah boolean
        "adalah_logika": {Name: "adalah_logika", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("adalah_logika() butuh 1 argumen")
                }
                return nativeBool(args[0].Type() == object.BOOLEAN_OBJ)
        }},

        // adalah_larik(nilai) — cek apakah array
        "adalah_larik": {Name: "adalah_larik", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("adalah_larik() butuh 1 argumen")
                }
                return nativeBool(args[0].Type() == object.ARRAY_OBJ)
        }},

        // adalah_kamus(nilai) — cek apakah map
        "adalah_kamus": {Name: "adalah_kamus", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("adalah_kamus() butuh 1 argumen")
                }
                return nativeBool(args[0].Type() == object.MAP_OBJ)
        }},

        // adalah_fungsi(nilai) — cek apakah fungsi
        "adalah_fungsi": {Name: "adalah_fungsi", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("adalah_fungsi() butuh 1 argumen")
                }
                return nativeBool(args[0].Type() == object.FUNCTION_OBJ || args[0].Type() == object.BUILTIN_OBJ)
        }},

        // adalah_kosong(nilai) — cek apakah kosong/null
        "adalah_kosong": {Name: "adalah_kosong", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("adalah_kosong() butuh 1 argumen")
                }
                return nativeBool(args[0].Type() == object.KOSONG_OBJ)
        }},

        // bilangan(nilai) — konversi ke integer
        "bilangan": {Name: "bilangan", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("bilangan() butuh 1 argumen")
                }
                switch a := args[0].(type) {
                case *object.Integer:
                        return a
                case *object.Float:
                        return &object.Integer{Value: int64(a.Value)}
                case *object.String:
                        v, err := strconv.ParseInt(strings.TrimSpace(a.Value), 10, 64)
                        if err != nil {
                                // coba float dulu
                                f, err2 := strconv.ParseFloat(strings.TrimSpace(a.Value), 64)
                                if err2 != nil {
                                        return object.NewError("tidak bisa konversi '%s' ke bilangan", a.Value)
                                }
                                return &object.Integer{Value: int64(f)}
                        }
                        return &object.Integer{Value: v}
                case *object.Boolean:
                        if a.Value {
                                return &object.Integer{Value: 1}
                        }
                        return &object.Integer{Value: 0}
                default:
                        return object.NewError("bilangan() tidak mendukung tipe %s", args[0].Type())
                }
        }},

        // desimal(nilai) — konversi ke float
        "desimal": {Name: "desimal", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("desimal() butuh 1 argumen")
                }
                switch a := args[0].(type) {
                case *object.Float:
                        return a
                case *object.Integer:
                        return &object.Float{Value: float64(a.Value)}
                case *object.String:
                        v, err := strconv.ParseFloat(strings.TrimSpace(a.Value), 64)
                        if err != nil {
                                return object.NewError("tidak bisa konversi '%s' ke desimal", a.Value)
                        }
                        return &object.Float{Value: v}
                case *object.Boolean:
                        if a.Value {
                                return &object.Float{Value: 1.0}
                        }
                        return &object.Float{Value: 0.0}
                default:
                        return object.NewError("desimal() tidak mendukung tipe %s", args[0].Type())
                }
        }},

        // teks(nilai) — konversi ke string
        "teks": {Name: "teks", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("teks() butuh 1 argumen")
                }
                return &object.String{Value: args[0].Inspect()}
        }},

        // logika(nilai) — konversi ke boolean
        "logika": {Name: "logika", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("logika() butuh 1 argumen")
                }
                return nativeBool(isTruthy(args[0]))
        }},

      
        // panjang(nilai) — panjang string atau array
        "panjang": {Name: "panjang", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("panjang() butuh 1 argumen")
                }
                switch arg := args[0].(type) {
                case *object.String:
                        return &object.Integer{Value: int64(utf8.RuneCountInString(arg.Value))}
                case *object.Array:
                        return &object.Integer{Value: int64(len(arg.Elements))}
                case *object.Map:
                        return &object.Integer{Value: int64(len(arg.Pairs))}
                default:
                        return object.NewError("panjang() tidak mendukung tipe %s", args[0].Type())
                }
        }},

        // huruf_besar(teks) — UPPERCASE
        "huruf_besar": {Name: "huruf_besar", Fn: wrapString1("huruf_besar", strings.ToUpper)},

        // huruf_kecil(teks) — lowercase
        "huruf_kecil": {Name: "huruf_kecil", Fn: wrapString1("huruf_kecil", strings.ToLower)},

        // judul(teks) — Title Case
        "judul": {Name: "judul", Fn: wrapString1("judul", strings.Title)},

        // potong(teks) — trim spasi di kedua sisi
        "potong": {Name: "potong", Fn: wrapString1("potong", strings.TrimSpace)},

        // potong_kiri(teks, karakter) — trim kiri
        "potong_kiri": {Name: "potong_kiri", Fn: func(args ...object.Object) object.Object {
                if len(args) != 2 {
                        return object.NewError("potong_kiri() butuh 2 argumen")
                }
                s, ok1 := args[0].(*object.String)
                cut, ok2 := args[1].(*object.String)
                if !ok1 || !ok2 {
                        return object.NewError("potong_kiri() butuh dua teks")
                }
                return &object.String{Value: strings.TrimLeft(s.Value, cut.Value)}
        }},

        // potong_kanan(teks, karakter) — trim kanan
        "potong_kanan": {Name: "potong_kanan", Fn: func(args ...object.Object) object.Object {
                if len(args) != 2 {
                        return object.NewError("potong_kanan() butuh 2 argumen")
                }
                s, ok1 := args[0].(*object.String)
                cut, ok2 := args[1].(*object.String)
                if !ok1 || !ok2 {
                        return object.NewError("potong_kanan() butuh dua teks")
                }
                return &object.String{Value: strings.TrimRight(s.Value, cut.Value)}
        }},

        // ganti(teks, lama, baru) — replace substring
        "ganti": {Name: "ganti", Fn: func(args ...object.Object) object.Object {
                if len(args) < 3 || len(args) > 4 {
                        return object.NewError("ganti() butuh 3 argumen (atau 4 dengan batas)")
                }
                s, ok1 := args[0].(*object.String)
                old, ok2 := args[1].(*object.String)
                baru, ok3 := args[2].(*object.String)
                if !ok1 || !ok2 || !ok3 {
                        return object.NewError("ganti() butuh teks")
                }
                n := -1
                if len(args) == 4 {
                        if batas, ok := args[3].(*object.Integer); ok {
                                n = int(batas.Value)
                        }
                }
                return &object.String{Value: strings.Replace(s.Value, old.Value, baru.Value, n)}
        }},

        // ada(teks, sub) — cek apakah sub ada dalam teks
        "ada": {Name: "ada", Fn: func(args ...object.Object) object.Object {
                if len(args) != 2 {
                        return object.NewError("ada() butuh 2 argumen")
                }
                switch s := args[0].(type) {
                case *object.String:
                        if sub, ok := args[1].(*object.String); ok {
                                return nativeBool(strings.Contains(s.Value, sub.Value))
                        }
                case *object.Array:
                        for _, el := range s.Elements {
                                if el.Inspect() == args[1].Inspect() {
                                        return BENAR
                                }
                        }
                        return SALAH
                }
                return object.NewError("ada() butuh teks atau larik sebagai argumen pertama")
        }},

        // indeks(teks, sub) — posisi sub dalam teks (-1 jika tidak ada)
        "indeks": {Name: "indeks", Fn: func(args ...object.Object) object.Object {
                if len(args) != 2 {
                        return object.NewError("indeks() butuh 2 argumen")
                }
                switch s := args[0].(type) {
                case *object.String:
                        if sub, ok := args[1].(*object.String); ok {
                                return &object.Integer{Value: int64(strings.Index(s.Value, sub.Value))}
                        }
                case *object.Array:
                        for i, el := range s.Elements {
                                if el.Inspect() == args[1].Inspect() {
                                        return &object.Integer{Value: int64(i)}
                                }
                        }
                        return &object.Integer{Value: -1}
                }
                return object.NewError("indeks() butuh teks atau larik")
        }},

        // awalan(teks, prefix) — cek apakah teks diawali prefix
        "awalan": {Name: "awalan", Fn: func(args ...object.Object) object.Object {
                if len(args) != 2 {
                        return object.NewError("awalan() butuh 2 argumen")
                }
                s, ok1 := args[0].(*object.String)
                p, ok2 := args[1].(*object.String)
                if !ok1 || !ok2 {
                        return object.NewError("awalan() butuh dua teks")
                }
                return nativeBool(strings.HasPrefix(s.Value, p.Value))
        }},

        // akhiran(teks, suffix) — cek apakah teks diakhiri suffix
        "akhiran": {Name: "akhiran", Fn: func(args ...object.Object) object.Object {
                if len(args) != 2 {
                        return object.NewError("akhiran() butuh 2 argumen")
                }
                s, ok1 := args[0].(*object.String)
                sf, ok2 := args[1].(*object.String)
                if !ok1 || !ok2 {
                        return object.NewError("akhiran() butuh dua teks")
                }
                return nativeBool(strings.HasSuffix(s.Value, sf.Value))
        }},

        // hitung(teks, sub) — hitung kemunculan sub dalam teks
        "hitung": {Name: "hitung", Fn: func(args ...object.Object) object.Object {
                if len(args) != 2 {
                        return object.NewError("hitung() butuh 2 argumen")
                }
                s, ok1 := args[0].(*object.String)
                sub, ok2 := args[1].(*object.String)
                if !ok1 || !ok2 {
                        return object.NewError("hitung() butuh dua teks")
                }
                return &object.Integer{Value: int64(strings.Count(s.Value, sub.Value))}
        }},

        // ulang_teks(teks, n) — ulangi teks sebanyak n kali
        "ulang_teks": {Name: "ulang_teks", Fn: func(args ...object.Object) object.Object {
                if len(args) != 2 {
                        return object.NewError("ulang_teks() butuh 2 argumen")
                }
                s, ok1 := args[0].(*object.String)
                n, ok2 := args[1].(*object.Integer)
                if !ok1 || !ok2 {
                        return object.NewError("ulang_teks() butuh teks dan bilangan")
                }
                return &object.String{Value: strings.Repeat(s.Value, int(n.Value))}
        }},

        // balik_teks(teks) — balik urutan karakter
        "balik_teks": {Name: "balik_teks", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("balik_teks() butuh 1 argumen")
                }
                s, ok := args[0].(*object.String)
                if !ok {
                        return object.NewError("balik_teks() butuh teks")
                }
                r := []rune(s.Value)
                for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
                        r[i], r[j] = r[j], r[i]
                }
                return &object.String{Value: string(r)}
        }},

        // pisah(teks, pemisah) — split string
        "pisah": {Name: "pisah", Fn: func(args ...object.Object) object.Object {
                if len(args) != 2 {
                        return object.NewError("pisah() butuh 2 argumen")
                }
                s, ok1 := args[0].(*object.String)
                sep, ok2 := args[1].(*object.String)
                if !ok1 || !ok2 {
                        return object.NewError("pisah() butuh dua teks")
                }
                parts := strings.Split(s.Value, sep.Value)
                elems := make([]object.Object, len(parts))
                for i, p := range parts {
                        elems[i] = &object.String{Value: p}
                }
                return &object.Array{Elements: elems}
        }},

        // gabung_teks(pemisah, larik) — join array dengan pemisah
        "gabung_teks": {Name: "gabung_teks", Fn: func(args ...object.Object) object.Object {
                if len(args) != 2 {
                        return object.NewError("gabung_teks() butuh 2 argumen")
                }
                sep, ok1 := args[0].(*object.String)
                arr, ok2 := args[1].(*object.Array)
                if !ok1 || !ok2 {
                        return object.NewError("gabung_teks(pemisah, larik) butuh teks dan larik")
                }
                parts := make([]string, len(arr.Elements))
                for i, el := range arr.Elements {
                        s, ok := el.(*object.String)
                        if !ok {
                                return object.NewError("semua elemen larik harus teks")
                        }
                        parts[i] = s.Value
                }
                return &object.String{Value: strings.Join(parts, sep.Value)}
        }},

        // irisan_teks(teks, mulai, akhir) — substring
        "irisan_teks": {Name: "irisan_teks", Fn: func(args ...object.Object) object.Object {
                if len(args) < 2 || len(args) > 3 {
                        return object.NewError("irisan_teks() butuh 2 atau 3 argumen")
                }
                s, ok := args[0].(*object.String)
                if !ok {
                        return object.NewError("irisan_teks() argumen pertama harus teks")
                }
                r := []rune(s.Value)
                mulai, ok2 := args[1].(*object.Integer)
                if !ok2 {
                        return object.NewError("irisan_teks() argumen kedua harus bilangan")
                }
                start := int(mulai.Value)
                end := len(r)
                if len(args) == 3 {
                        akhir, ok3 := args[2].(*object.Integer)
                        if !ok3 {
                                return object.NewError("irisan_teks() argumen ketiga harus bilangan")
                        }
                        end = int(akhir.Value)
                }
                if start < 0 {
                        start = 0
                }
                if end > len(r) {
                        end = len(r)
                }
                if start > end {
                        return &object.String{Value: ""}
                }
                return &object.String{Value: string(r[start:end])}
        }},


        // masukkan(larik, nilai) — append elemen ke array baru
        "masukkan": {Name: "masukkan", Fn: func(args ...object.Object) object.Object {
                if len(args) != 2 {
                        return object.NewError("masukkan() butuh 2 argumen")
                }
                arr, ok := args[0].(*object.Array)
                if !ok {
                        return object.NewError("masukkan() argumen pertama harus larik")
                }
                elems := make([]object.Object, len(arr.Elements)+1)
                copy(elems, arr.Elements)
                elems[len(arr.Elements)] = args[1]
                return &object.Array{Elements: elems}
        }},

        // hapus(larik, indeks) — hapus elemen pada indeks
        "hapus": {Name: "hapus", Fn: func(args ...object.Object) object.Object {
                if len(args) != 2 {
                        return object.NewError("hapus() butuh 2 argumen")
                }
                arr, ok := args[0].(*object.Array)
                if !ok {
                        return object.NewError("hapus() argumen pertama harus larik")
                }
                idx, ok2 := args[1].(*object.Integer)
                if !ok2 {
                        return object.NewError("hapus() argumen kedua harus bilangan")
                }
                i := int(idx.Value)
                if i < 0 {
                        i = len(arr.Elements) + i
                }
                if i < 0 || i >= len(arr.Elements) {
                        return object.NewError("hapus() indeks %d di luar batas", idx.Value)
                }
                elems := make([]object.Object, 0, len(arr.Elements)-1)
                elems = append(elems, arr.Elements[:i]...)
                elems = append(elems, arr.Elements[i+1:]...)
                return &object.Array{Elements: elems}
        }},

        // sisipkan(larik, indeks, nilai) — sisipkan elemen pada indeks
        "sisipkan": {Name: "sisipkan", Fn: func(args ...object.Object) object.Object {
                if len(args) != 3 {
                        return object.NewError("sisipkan() butuh 3 argumen")
                }
                arr, ok1 := args[0].(*object.Array)
                idx, ok2 := args[1].(*object.Integer)
                if !ok1 || !ok2 {
                        return object.NewError("sisipkan(larik, indeks, nilai)")
                }
                i := int(idx.Value)
                if i < 0 || i > len(arr.Elements) {
                        return object.NewError("sisipkan() indeks %d di luar batas", i)
                }
                elems := make([]object.Object, len(arr.Elements)+1)
                copy(elems, arr.Elements[:i])
                elems[i] = args[2]
                copy(elems[i+1:], arr.Elements[i:])
                return &object.Array{Elements: elems}
        }},

        // pertama(larik) — elemen pertama
        "pertama": {Name: "pertama", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("pertama() butuh 1 argumen")
                }
                arr, ok := args[0].(*object.Array)
                if !ok {
                        return object.NewError("pertama() butuh larik")
                }
                if len(arr.Elements) == 0 {
                        return KOSONG
                }
                return arr.Elements[0]
        }},

        // terakhir(larik) — elemen terakhir
        "terakhir": {Name: "terakhir", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("terakhir() butuh 1 argumen")
                }
                arr, ok := args[0].(*object.Array)
                if !ok {
                        return object.NewError("terakhir() butuh larik")
                }
                if len(arr.Elements) == 0 {
                        return KOSONG
                }
                return arr.Elements[len(arr.Elements)-1]
        }},

        // ekor(larik) — semua kecuali elemen pertama
        "ekor": {Name: "ekor", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("ekor() butuh 1 argumen")
                }
                arr, ok := args[0].(*object.Array)
                if !ok {
                        return object.NewError("ekor() butuh larik")
                }
                if len(arr.Elements) == 0 {
                        return KOSONG
                }
                elems := make([]object.Object, len(arr.Elements)-1)
                copy(elems, arr.Elements[1:])
                return &object.Array{Elements: elems}
        }},

        // gabung(larik1, larik2) — gabungkan dua array
        "gabung": {Name: "gabung", Fn: func(args ...object.Object) object.Object {
                if len(args) != 2 {
                        return object.NewError("gabung() butuh 2 argumen")
                }
                a1, ok1 := args[0].(*object.Array)
                a2, ok2 := args[1].(*object.Array)
                if !ok1 || !ok2 {
                        return object.NewError("gabung() butuh dua larik")
                }
                elems := make([]object.Object, len(a1.Elements)+len(a2.Elements))
                copy(elems, a1.Elements)
                copy(elems[len(a1.Elements):], a2.Elements)
                return &object.Array{Elements: elems}
        }},

        // irisan(larik, mulai, akhir) — slice array
        "irisan": {Name: "irisan", Fn: func(args ...object.Object) object.Object {
                if len(args) < 2 || len(args) > 3 {
                        return object.NewError("irisan() butuh 2 atau 3 argumen")
                }
                arr, ok := args[0].(*object.Array)
                if !ok {
                        return object.NewError("irisan() argumen pertama harus larik")
                }
                mulai, ok2 := args[1].(*object.Integer)
                if !ok2 {
                        return object.NewError("irisan() argumen kedua harus bilangan")
                }
                start := int(mulai.Value)
                end := len(arr.Elements)
                if len(args) == 3 {
                        akhir, ok3 := args[2].(*object.Integer)
                        if !ok3 {
                                return object.NewError("irisan() argumen ketiga harus bilangan")
                        }
                        end = int(akhir.Value)
                }
                if start < 0 {
                        start = 0
                }
                if end > len(arr.Elements) {
                        end = len(arr.Elements)
                }
                elems := make([]object.Object, end-start)
                copy(elems, arr.Elements[start:end])
                return &object.Array{Elements: elems}
        }},

        // balik_larik(larik) — balik urutan elemen
        "balik_larik": {Name: "balik_larik", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("balik_larik() butuh 1 argumen")
                }
                arr, ok := args[0].(*object.Array)
                if !ok {
                        return object.NewError("balik_larik() butuh larik")
                }
                elems := make([]object.Object, len(arr.Elements))
                for i, e := range arr.Elements {
                        elems[len(arr.Elements)-1-i] = e
                }
                return &object.Array{Elements: elems}
        }},

        // urutkan(larik) — urutkan array bilangan atau teks
        "urutkan": {Name: "urutkan", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("urutkan() butuh 1 argumen")
                }
                arr, ok := args[0].(*object.Array)
                if !ok {
                        return object.NewError("urutkan() butuh larik")
                }
                elems := make([]object.Object, len(arr.Elements))
                copy(elems, arr.Elements)
                sort.Slice(elems, func(i, j int) bool {
                        switch a := elems[i].(type) {
                        case *object.Integer:
                                switch b := elems[j].(type) {
                                case *object.Integer:
                                        return a.Value < b.Value
                                case *object.Float:
                                        return float64(a.Value) < b.Value
                                }
                        case *object.Float:
                                switch b := elems[j].(type) {
                                case *object.Float:
                                        return a.Value < b.Value
                                case *object.Integer:
                                        return a.Value < float64(b.Value)
                                }
                        case *object.String:
                                if b, ok := elems[j].(*object.String); ok {
                                        return a.Value < b.Value
                                }
                        }
                        return false
                })
                return &object.Array{Elements: elems}
        }},

        // unik(larik) — hapus duplikat
        "unik": {Name: "unik", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("unik() butuh 1 argumen")
                }
                arr, ok := args[0].(*object.Array)
                if !ok {
                        return object.NewError("unik() butuh larik")
                }
                seen := map[string]bool{}
                var elems []object.Object
                for _, e := range arr.Elements {
                        key := e.Inspect()
                        if !seen[key] {
                                seen[key] = true
                                elems = append(elems, e)
                        }
                }
                return &object.Array{Elements: elems}
        }},

        // rentang(mulai, akhir, langkah?) — buat array integer dari mulai ke akhir
        "rentang": {Name: "rentang", Fn: func(args ...object.Object) object.Object {
                if len(args) < 2 || len(args) > 3 {
                        return object.NewError("rentang() butuh 2 atau 3 argumen: rentang(mulai, akhir[, langkah])")
                }
                mulai, ok1 := args[0].(*object.Integer)
                akhir, ok2 := args[1].(*object.Integer)
                if !ok1 || !ok2 {
                        return object.NewError("rentang() butuh bilangan bulat")
                }
                langkah := int64(1)
                if len(args) == 3 {
                        if l, ok := args[2].(*object.Integer); ok {
                                langkah = l.Value
                        }
                }
                if langkah == 0 {
                        return object.NewError("rentang() langkah tidak boleh 0")
                }
                var elems []object.Object
                for i := mulai.Value; (langkah > 0 && i < akhir.Value) || (langkah < 0 && i > akhir.Value); i += langkah {
                        elems = append(elems, &object.Integer{Value: i})
                }
                return &object.Array{Elements: elems}
        }},

        // petakan(larik, fungsi) — map: terapkan fungsi ke setiap elemen
        "petakan": {Name: "petakan", Fn: func(args ...object.Object) object.Object {
                if len(args) != 2 {
                        return object.NewError("petakan() butuh 2 argumen: petakan(larik, fungsi)")
                }
                arr, ok1 := args[0].(*object.Array)
                fn := args[1]
                if !ok1 {
                        return object.NewError("petakan() argumen pertama harus larik")
                }
                elems := make([]object.Object, len(arr.Elements))
                for i, el := range arr.Elements {
                        result := applyFunction(fn, []object.Object{el})
                        if object.IsError(result) {
                                return result
                        }
                        elems[i] = result
                }
                return &object.Array{Elements: elems}
        }},

        // saring(larik, fungsi) — filter: ambil elemen yang memenuhi kondisi
        "saring": {Name: "saring", Fn: func(args ...object.Object) object.Object {
                if len(args) != 2 {
                        return object.NewError("saring() butuh 2 argumen: saring(larik, fungsi)")
                }
                arr, ok1 := args[0].(*object.Array)
                fn := args[1]
                if !ok1 {
                        return object.NewError("saring() argumen pertama harus larik")
                }
                var elems []object.Object
                for _, el := range arr.Elements {
                        result := applyFunction(fn, []object.Object{el})
                        if object.IsError(result) {
                                return result
                        }
                        if isTruthy(result) {
                                elems = append(elems, el)
                        }
                }
                return &object.Array{Elements: elems}
        }},

        // lipat(larik, awal, fungsi) — reduce/fold
        "lipat": {Name: "lipat", Fn: func(args ...object.Object) object.Object {
                if len(args) != 3 {
                        return object.NewError("lipat() butuh 3 argumen: lipat(larik, awal, fungsi)")
                }
                arr, ok := args[0].(*object.Array)
                if !ok {
                        return object.NewError("lipat() argumen pertama harus larik")
                }
                akum := args[1]
                fn := args[2]
                for _, el := range arr.Elements {
                        result := applyFunction(fn, []object.Object{akum, el})
                        if object.IsError(result) {
                                return result
                        }
                        akum = result
                }
                return akum
        }},

        // setiap(larik, fungsi) — cek apakah SEMUA elemen memenuhi kondisi
        "setiap": {Name: "setiap", Fn: func(args ...object.Object) object.Object {
                if len(args) != 2 {
                        return object.NewError("setiap() butuh 2 argumen")
                }
                arr, ok := args[0].(*object.Array)
                if !ok {
                        return object.NewError("setiap() argumen pertama harus larik")
                }
                fn := args[1]
                for _, el := range arr.Elements {
                        result := applyFunction(fn, []object.Object{el})
                        if object.IsError(result) {
                                return result
                        }
                        if !isTruthy(result) {
                                return SALAH
                        }
                }
                return BENAR
        }},

        // salah_satu(larik, fungsi) — cek apakah ADA elemen yang memenuhi kondisi
        "salah_satu": {Name: "salah_satu", Fn: func(args ...object.Object) object.Object {
                if len(args) != 2 {
                        return object.NewError("salah_satu() butuh 2 argumen")
                }
                arr, ok := args[0].(*object.Array)
                if !ok {
                        return object.NewError("salah_satu() argumen pertama harus larik")
                }
                fn := args[1]
                for _, el := range arr.Elements {
                        result := applyFunction(fn, []object.Object{el})
                        if object.IsError(result) {
                                return result
                        }
                        if isTruthy(result) {
                                return BENAR
                        }
                }
                return SALAH
        }},

        // zip(larik1, larik2) — gabungkan dua larik menjadi larik pasangan
        "zip": {Name: "zip", Fn: func(args ...object.Object) object.Object {
                if len(args) != 2 {
                        return object.NewError("zip() butuh 2 argumen")
                }
                a1, ok1 := args[0].(*object.Array)
                a2, ok2 := args[1].(*object.Array)
                if !ok1 || !ok2 {
                        return object.NewError("zip() butuh dua larik")
                }
                n := len(a1.Elements)
                if len(a2.Elements) < n {
                        n = len(a2.Elements)
                }
                elems := make([]object.Object, n)
                for i := 0; i < n; i++ {
                        elems[i] = &object.Array{Elements: []object.Object{a1.Elements[i], a2.Elements[i]}}
                }
                return &object.Array{Elements: elems}
        }},

        // akar(x) — square root (math.Sqrt)
        "akar": {Name: "akar", Fn: wrapMath1("akar", math.Sqrt)},

        // mutlak(x) — absolute value (math.Abs)
        "mutlak": {Name: "mutlak", Fn: wrapMath1("mutlak", math.Abs)},

        // lantai(x) — floor (math.Floor)
        "lantai": {Name: "lantai", Fn: wrapMath1("lantai", math.Floor)},

        // langit(x) — ceiling (math.Ceil)
        "langit": {Name: "langit", Fn: wrapMath1("langit", math.Ceil)},

        // bulat(x) — round (math.Round)
        "bulat": {Name: "bulat", Fn: wrapMath1("bulat", math.Round)},

        // sin(x) — sine dalam radian (math.Sin)
        "sin": {Name: "sin", Fn: wrapMath1("sin", math.Sin)},

        // cos(x) — cosine dalam radian (math.Cos)
        "cos": {Name: "cos", Fn: wrapMath1("cos", math.Cos)},

        // tan(x) — tangent dalam radian (math.Tan)
        "tan": {Name: "tan", Fn: wrapMath1("tan", math.Tan)},

        // asin(x) — arc sine (math.Asin)
        "asin": {Name: "asin", Fn: wrapMath1("asin", math.Asin)},

        // acos(x) — arc cosine (math.Acos)
        "acos": {Name: "acos", Fn: wrapMath1("acos", math.Acos)},

        // atan(x) — arc tangent (math.Atan)
        "atan": {Name: "atan", Fn: wrapMath1("atan", math.Atan)},

        // atan2(y, x) — arc tangent of y/x (math.Atan2)
        "atan2": {Name: "atan2", Fn: wrapMath2("atan2", math.Atan2)},

        // log(x) — natural log (math.Log)
        "log": {Name: "log", Fn: wrapMath1("log", math.Log)},

        // log2(x) — log base 2 (math.Log2)
        "log2": {Name: "log2", Fn: wrapMath1("log2", math.Log2)},

        // log10(x) — log base 10 (math.Log10)
        "log10": {Name: "log10", Fn: wrapMath1("log10", math.Log10)},

        // pangkat(basis, eksponen) — pow (math.Pow)
        "pangkat": {Name: "pangkat", Fn: wrapMath2("pangkat", math.Pow)},

        // eksponen(x) — e^x (math.Exp)
        "eksponen": {Name: "eksponen", Fn: wrapMath1("eksponen", math.Exp)},

        // maks(a, b, ...) — nilai terbesar
        "maks": {Name: "maks", Fn: func(args ...object.Object) object.Object {
                if len(args) == 0 {
                        return object.NewError("maks() butuh minimal 1 argumen")
                }
                // Cek apakah argumen pertama adalah larik
                if arr, ok := args[0].(*object.Array); ok && len(args) == 1 {
                        if len(arr.Elements) == 0 {
                                return KOSONG
                        }
                        max := toFloat64(arr.Elements[0])
                        maxObj := arr.Elements[0]
                        for _, el := range arr.Elements[1:] {
                                if v := toFloat64(el); v > max {
                                        max = v
                                        maxObj = el
                                }
                        }
                        return maxObj
                }
                max := toFloat64(args[0])
                maxObj := args[0]
                for _, a := range args[1:] {
                        if v := toFloat64(a); v > max {
                                max = v
                                maxObj = a
                        }
                }
                return maxObj
        }},

        // min(a, b, ...) — nilai terkecil
        "min": {Name: "min", Fn: func(args ...object.Object) object.Object {
                if len(args) == 0 {
                        return object.NewError("min() butuh minimal 1 argumen")
                }
                if arr, ok := args[0].(*object.Array); ok && len(args) == 1 {
                        if len(arr.Elements) == 0 {
                                return KOSONG
                        }
                        mn := toFloat64(arr.Elements[0])
                        mnObj := arr.Elements[0]
                        for _, el := range arr.Elements[1:] {
                                if v := toFloat64(el); v < mn {
                                        mn = v
                                        mnObj = el
                                }
                        }
                        return mnObj
                }
                mn := toFloat64(args[0])
                mnObj := args[0]
                for _, a := range args[1:] {
                        if v := toFloat64(a); v < mn {
                                mn = v
                                mnObj = a
                        }
                }
                return mnObj
        }},

        // jumlah(larik) — total semua elemen
        "jumlah": {Name: "jumlah", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("jumlah() butuh 1 argumen (larik)")
                }
                arr, ok := args[0].(*object.Array)
                if !ok {
                        return object.NewError("jumlah() butuh larik")
                }
                var total float64
                isInt := true
                for _, el := range arr.Elements {
                        switch v := el.(type) {
                        case *object.Integer:
                                total += float64(v.Value)
                        case *object.Float:
                                total += v.Value
                                isInt = false
                        default:
                                return object.NewError("jumlah() semua elemen harus bilangan")
                        }
                }
                if isInt {
                        return &object.Integer{Value: int64(total)}
                }
                return &object.Float{Value: total}
        }},

        // rata(larik) — rata-rata semua elemen
        "rata": {Name: "rata", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("rata() butuh 1 argumen (larik)")
                }
                arr, ok := args[0].(*object.Array)
                if !ok || len(arr.Elements) == 0 {
                        return object.NewError("rata() butuh larik dengan minimal 1 elemen")
                }
                var total float64
                for _, el := range arr.Elements {
                        total += toFloat64(el)
                }
                return &object.Float{Value: total / float64(len(arr.Elements))}
        }},

        // ganjil(n) — cek apakah bilangan ganjil
        "ganjil": {Name: "ganjil", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("ganjil() butuh 1 argumen")
                }
                n, ok := args[0].(*object.Integer)
                if !ok {
                        return object.NewError("ganjil() butuh bilangan bulat")
                }
                return nativeBool(n.Value%2 != 0)
        }},

        // genap(n) — cek apakah bilangan genap
        "genap": {Name: "genap", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("genap() butuh 1 argumen")
                }
                n, ok := args[0].(*object.Integer)
                if !ok {
                        return object.NewError("genap() butuh bilangan bulat")
                }
                return nativeBool(n.Value%2 == 0)
        }},


        // acak(n) — angka acak antara 0 dan n-1
        "acak": {Name: "acak", Fn: func(args ...object.Object) object.Object {
                if len(args) == 0 {
                        return &object.Float{Value: rand.Float64()}
                }
                n, ok := args[0].(*object.Integer)
                if !ok {
                        return object.NewError("acak() butuh bilangan bulat")
                }
                if n.Value <= 0 {
                        return object.NewError("acak() n harus > 0")
                }
                return &object.Integer{Value: rand.Int63n(n.Value)}
        }},

        // acak_rentang(min, maks) — angka acak dalam rentang [min, maks]
        "acak_rentang": {Name: "acak_rentang", Fn: func(args ...object.Object) object.Object {
                if len(args) != 2 {
                        return object.NewError("acak_rentang() butuh 2 argumen")
                }
                mn, ok1 := args[0].(*object.Integer)
                mx, ok2 := args[1].(*object.Integer)
                if !ok1 || !ok2 {
                        return object.NewError("acak_rentang() butuh dua bilangan bulat")
                }
                if mx.Value <= mn.Value {
                        return object.NewError("acak_rentang() maks harus > min")
                }
                return &object.Integer{Value: mn.Value + rand.Int63n(mx.Value-mn.Value+1)}
        }},

        // acak_pilih(larik) — pilih elemen acak dari larik
        "acak_pilih": {Name: "acak_pilih", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("acak_pilih() butuh 1 argumen")
                }
                arr, ok := args[0].(*object.Array)
                if !ok || len(arr.Elements) == 0 {
                        return object.NewError("acak_pilih() butuh larik dengan minimal 1 elemen")
                }
                return arr.Elements[rand.Intn(len(arr.Elements))]
        }},

        // acak_campur(larik) — acak urutan elemen larik (Fisher-Yates shuffle)
        "acak_campur": {Name: "acak_campur", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("acak_campur() butuh 1 argumen")
                }
                arr, ok := args[0].(*object.Array)
                if !ok {
                        return object.NewError("acak_campur() butuh larik")
                }
                elems := make([]object.Object, len(arr.Elements))
                copy(elems, arr.Elements)
                rand.Shuffle(len(elems), func(i, j int) { elems[i], elems[j] = elems[j], elems[i] })
                return &object.Array{Elements: elems}
        }},


        // waktu_sekarang() — unix timestamp milidetik
        "waktu_sekarang": {Name: "waktu_sekarang", Fn: func(args ...object.Object) object.Object {
                return &object.Integer{Value: time.Now().UnixMilli()}
        }},

        // tidur(ms) — jeda eksekusi selama ms milidetik
        "tidur": {Name: "tidur", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("tidur() butuh 1 argumen (milidetik)")
                }
                ms, ok := args[0].(*object.Integer)
                if !ok {
                        return object.NewError("tidur() butuh bilangan bulat (milidetik)")
                }
                time.Sleep(time.Duration(ms.Value) * time.Millisecond)
                return KOSONG
        }},


        // galat(pesan) — buat error
        "galat": {Name: "galat", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("galat() butuh 1 argumen")
                }
                return object.NewError("%s", args[0].Inspect())
        }},

        // tegas(kondisi, pesan?) — assert, error jika kondisi salah
        "tegas": {Name: "tegas", Fn: func(args ...object.Object) object.Object {
                if len(args) < 1 {
                        return object.NewError("tegas() butuh minimal 1 argumen")
                }
                if !isTruthy(args[0]) {
                        msg := "assertion gagal"
                        if len(args) > 1 {
                                msg = args[1].Inspect()
                        }
                        return object.NewError(msg)
                }
                return KOSONG
        }},

        // kunci(kamus) — ambil semua kunci dari kamus sebagai larik
        "kunci": {Name: "kunci", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("kunci() butuh 1 argumen")
                }
                m, ok := args[0].(*object.Map)
                if !ok {
                        return object.NewError("kunci() butuh kamus")
                }
                var elems []object.Object
                for _, pair := range m.Pairs {
                        elems = append(elems, pair.Key)
                }
                return &object.Array{Elements: elems}
        }},

        // nilai_kamus(kamus) — ambil semua nilai dari kamus sebagai larik
        "nilai_kamus": {Name: "nilai_kamus", Fn: func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("nilai_kamus() butuh 1 argumen")
                }
                m, ok := args[0].(*object.Map)
                if !ok {
                        return object.NewError("nilai_kamus() butuh kamus")
                }
                var elems []object.Object
                for _, pair := range m.Pairs {
                        elems = append(elems, pair.Value)
                }
                return &object.Array{Elements: elems}
        }},

        // hapus_kamus(kamus, kunci) — hapus entri dari kamus
        "hapus_kamus": {Name: "hapus_kamus", Fn: func(args ...object.Object) object.Object {
                if len(args) != 2 {
                        return object.NewError("hapus_kamus() butuh 2 argumen")
                }
                m, ok := args[0].(*object.Map)
                if !ok {
                        return object.NewError("hapus_kamus() argumen pertama harus kamus")
                }
                key, ok2 := args[1].(object.Hashable)
                if !ok2 {
                        return object.NewError("hapus_kamus() kunci tidak valid")
                }
                newPairs := make(map[object.HashKey]object.MapPair)
                for k, v := range m.Pairs {
                        if k != key.HashKey() {
                                newPairs[k] = v
                        }
                }
                return &object.Map{Pairs: newPairs}
        }},

        // ada_kunci(kamus, kunci) — cek apakah kunci ada di kamus
        "ada_kunci": {Name: "ada_kunci", Fn: func(args ...object.Object) object.Object {
                if len(args) != 2 {
                        return object.NewError("ada_kunci() butuh 2 argumen")
                }
                m, ok := args[0].(*object.Map)
                if !ok {
                        return object.NewError("ada_kunci() argumen pertama harus kamus")
                }
                key, ok2 := args[1].(object.Hashable)
                if !ok2 {
                        return SALAH
                }
                _, exists := m.Pairs[key.HashKey()]
                return nativeBool(exists)
        }},

        // kosongkan() — buat larik kosong baru
        "kosongkan": {Name: "kosongkan", Fn: func(args ...object.Object) object.Object {
                return &object.Array{Elements: []object.Object{}}
        }},

        // PI — konstanta matematika Pi
        "PI": {Name: "PI", Fn: func(args ...object.Object) object.Object {
                return &object.Float{Value: math.Pi}
        }},

        // E — konstanta matematika Euler
        "E": {Name: "E", Fn: func(args ...object.Object) object.Object {
                return &object.Float{Value: math.E}
        }},

        // MAKS_BILANGAN — nilai integer maksimum
        "MAKS_BILANGAN": {Name: "MAKS_BILANGAN", Fn: func(args ...object.Object) object.Object {
                return &object.Integer{Value: math.MaxInt64}
        }},

        // cetak_galat(nilai...) — cetak ke stderr
        "cetak_galat": {Name: "cetak_galat", Fn: func(args ...object.Object) object.Object {
                parts := make([]string, len(args))
                for i, a := range args {
                        parts[i] = a.Inspect()
                }
                fmt.Fprintln(os.Stderr, strings.Join(parts, " "))
                return KOSONG
        }},
        }
}

//Helper

func wrapString1(name string, fn func(string) string) object.BuiltinFunction {
        return func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("%s() butuh 1 argumen", name)
                }
                s, ok := args[0].(*object.String)
                if !ok {
                        return object.NewError("%s() butuh teks", name)
                }
                return &object.String{Value: fn(s.Value)}
        }
}

func wrapMath1(name string, fn func(float64) float64) object.BuiltinFunction {
        return func(args ...object.Object) object.Object {
                if len(args) != 1 {
                        return object.NewError("%s() butuh 1 argumen", name)
                }
                v := toFloat64(args[0])
                if args[0].Type() != object.INTEGER_OBJ && args[0].Type() != object.FLOAT_OBJ {
                        return object.NewError("%s() butuh bilangan", name)
                }
                return &object.Float{Value: fn(v)}
        }
}

func wrapMath2(name string, fn func(float64, float64) float64) object.BuiltinFunction {
        return func(args ...object.Object) object.Object {
                if len(args) != 2 {
                        return object.NewError("%s() butuh 2 argumen", name)
                }
                a := toFloat64(args[0])
                b := toFloat64(args[1])
                return &object.Float{Value: fn(a, b)}
        }
}

func toFloat64(obj object.Object) float64 {
        switch v := obj.(type) {
        case *object.Integer:
                return float64(v.Value)
        case *object.Float:
                return v.Value
        default:
                return 0
        }
}
