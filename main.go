// Language .ky - Created and Developed by Farhan Kertadiwangsa
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"ky/evaluator"
	"ky/lexer"
	"ky/object"
	"ky/parser"
)

const versi = "0.2.0"

const banner = `
 ██╗  ██╗██╗   ██╗
 ██║ ██╔╝╚██╗ ██╔╝   Bahasa Pemrograman .ky  v` + versi + `
 █████╔╝  ╚████╔╝    Kreator: Farhan Kertadiwangsa
 ██╔═██╗   ╚██╔╝     Ketik 'keluar' untuk keluar
 ██║  ██╗   ██║      Ketik 'bantuan' untuk daftar fungsi
 ╚═╝  ╚═╝   ╚═╝

`

const pesanBantuan = `
═══════════════════════════════════════════════════════════════
  KATA KUNCI ky-lang
═══════════════════════════════════════════════════════════════
  buat      — deklarasi variabel         buat x = 10
  fungsi    — definisi fungsi            fungsi(a, b) { ... }
  balik     — kembalikan nilai           balik hasil
  jika      — percabangan                jika x > 0 { ... }
  lain      — else                       } lain { ... }
  ulang     — perulangan for             ulang buat i=0; i<5; i=i+1 { ... }
  selama    — perulangan while           selama x > 0 { ... }
  henti     — keluar dari loop (break)
  lanjut    — lanjut iterasi (continue)
  benar     — true     salah — false     kosong — null

═══════════════════════════════════════════════════════════════
  FUNGSI BAWAAN 
═══════════════════════════════════════════════════════════════
  I/O:      tampilkan, cetak, tulis, format, baca, baca_berkas, tulis_berkas
  Tipe:     jenis, adalah_bilangan, adalah_teks, adalah_larik, adalah_kamus
  Konversi: bilangan, desimal, teks, logika
  String:   panjang, huruf_besar, huruf_kecil, potong, ganti, ada, indeks,
            awalan, akhiran, hitung, ulang_teks, balik_teks, pisah, gabung_teks,
            irisan_teks, judul
  Larik:    masukkan, hapus, sisipkan, pertama, terakhir, ekor, gabung,
            irisan, balik_larik, urutkan, unik, rentang
  HOF:      petakan, saring, lipat, setiap, salah_satu, zip
  Matematika: akar, mutlak, lantai, langit, bulat, sin, cos, tan,
              asin, acos, atan, atan2, log, log2, log10, pangkat,
              eksponen, maks, min, jumlah, rata, ganjil, genap
  Acak:     acak, acak_rentang, acak_pilih, acak_campur
  Waktu:    waktu_sekarang, tidur
  Kamus:    kunci, nilai_kamus, hapus_kamus, ada_kunci
  Utilitas: galat, tegas, kosongkan, PI, E, MAKS_BILANGAN, cetak_galat
═══════════════════════════════════════════════════════════════
`

func main() {
	args := os.Args[1:]

	switch {
	case len(args) == 0:
		runREPL()
	case len(args) == 1:
		runFile(args[0])
	default:
		fmt.Fprintln(os.Stderr, "Penggunaan: ky [file.ky]")
		os.Exit(1)
	}
}


func runREPL() {
	fmt.Print(banner)
	env := object.NewEnvironment()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(">>> ")
		if !scanner.Scan() {
			fmt.Println("\nSampai!")
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		switch line {
		case "keluar", "exit", "quit":
			fmt.Println("Sampai jumpa, Kreator Farhan!")
			return
		case "bantuan", "help":
			fmt.Print(pesanBantuan)
			continue
		}

		result := jalankan(line, env)
		if result != nil && result.Type() != object.KOSONG_OBJ {
			fmt.Println(result.Inspect())
		}
	}
}

// run file

func runFile(path string) {
	if !strings.HasSuffix(path, ".ky") {
		fmt.Fprintf(os.Stderr, "PERINGATAN: '%s' tidak berekstensi .ky\n", path)
	}
	src, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Gagal membaca file: %v\n", err)
		os.Exit(1)
	}
	env := object.NewEnvironment()
	result := jalankan(string(src), env)
	if result != nil && result.Type() == object.ERROR_OBJ {
		fmt.Fprintln(os.Stderr, evaluator.FormatError(result.(*object.Error), string(src)))
		os.Exit(1)
	}
}

//Core

func jalankan(src string, env *object.Environment) object.Object {
	l := lexer.New(src)
	p := parser.New(l)
	program := p.ParseProgram()

	if errs := p.Errors(); len(errs) > 0 {
		fmt.Fprintln(os.Stderr, "\n[.ky PARSE ERROR]")
		for _, e := range errs {
			fmt.Fprintln(os.Stderr, "  -", e)
		}
		fmt.Fprintln(os.Stderr)
		return nil
	}
	return evaluator.Eval(program, env)
}
