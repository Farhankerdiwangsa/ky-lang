# Dokumentasi Bahasa Pemrograman .ky

> **Language ky-lang** — Bahasa pemrograman edukatif berbahasa Indonesia  
> Diciptakan dan dikembangkan oleh **Farhan Kertadiwangsa**

---

## Daftar Isi

1. [Tentang .ky](#tentang-ky)
2. [Instalasi & Menjalankan](#instalasi--menjalankan)
3. [Sintaks Dasar](#sintaks-dasar)
4. [Tipe Data](#tipe-data)
5. [Operator](#operator)
6. [Kata Kunci & Alias](#kata-kunci--alias)
7. [Struktur Kontrol](#struktur-kontrol)
8. [Fungsi](#fungsi)
9. [Larik (Array)](#larik-array)
10. [Kamus (Map)](#kamus-map)
11. [Fungsi Bawaan Lengkap](#fungsi-bawaan-lengkap)
12. [Benchmark Kecepatan](#benchmark-kecepatan)
13. [Filosofi & Desain](#filosofi--desain)

---

## Tentang .ky

**kylang** adalah bahasa pemrograman interpretatif yang dirancang khusus untuk pelajar Indonesia — mulai dari tingkat SD, SMP, hingga SMA. Bahasa ini menggunakan kata kunci berbahasa Indonesia sehingga konsep pemrograman lebih mudah diserap tanpa hambatan bahasa asing.

**Karakteristik utama:**

- Kata kunci dan fungsi bawaan dalam Bahasa Indonesia
- Sintaks bersih dan konsisten — cocok untuk pemula
- Mendukung paradigma fungsional (higher-order functions, closure)
- 60+ fungsi bawaan siap pakai tanpa `import`
- Dibangun dengan Go — stabil, cepat dikompilasi, tanpa dependensi eksternal
- REPL interaktif untuk eksplorasi langsung

---

## Instalasi & Menjalankan

### Build Interpreter

```bash
cd ky-lang
bash --login -c "go build -o ky ."
```

### Menjalankan File `.ky`

```bash
./ky-lang/ky namafile.ky
```

### REPL Interaktif

```bash
./ky-lang/ky
```

Ketik ekspresi atau pernyataan langsung di terminal. Keluar dengan `Ctrl+C` atau `Ctrl+D`.

### Contoh Bawaan

```bash
./ky-lang/ky ky-contoh/halo.ky
./ky-lang/ky ky-contoh/fibonacci.ky
./ky-lang/ky ky-contoh/logika.ky
./ky-lang/ky ky-contoh/fungsi_tinggi.ky
./ky-lang/ky ky-contoh/matematika.ky
```

---

## Sintaks Dasar

### Komentar

```ky
// Ini adalah komentar satu baris
```

Komentar dimulai dengan `//` dan berlaku hingga akhir baris. Belum mendukung komentar multi-baris.

### Deklarasi Variabel

Variabel dideklarasikan dengan kata kunci `buat`:

```ky
buat nama   = "Farhan"
buat umur   = 17
buat tinggi = 1.75
buat lulus  = benar
buat data   = kosong
```

> **Catatan:** Variabel harus dideklarasikan dengan `buat` sebelum digunakan. Assignment ulang (`x = nilai`) hanya berlaku untuk variabel yang sudah ada.

### Assignment Ulang

```ky
buat skor = 100
skor = skor + 50    // skor sekarang 150
```

### Tanda Kurung & Blok

Blok kode dibungkus dengan `{ }`. Tanda kurung `( )` digunakan untuk argumen fungsi dan pengelompokan ekspresi.

---

## Tipe Data

| Tipe        | Nama dalam .ky | Contoh Nilai                         |
|-------------|----------------|--------------------------------------|
| Integer     | `bilangan`     | `42`, `-7`, `0`                      |
| Float       | `desimal`      | `3.14`, `-0.5`, `2.718`              |
| String      | `teks`         | `"halo"`, `"Bahasa .ky"`            |
| Boolean     | `logika`       | `benar`, `salah`                     |
| Null        | `kosong`       | `kosong`, `nil`                      |
| Array       | `larik`        | `[1, 2, 3]`, `["a", "b"]`           |
| Map         | `kamus`        | `{"kunci": "nilai", "angka": 42}`    |
| Fungsi      | `fungsi`       | `fungsi(x) { balik x * 2 }`         |

### Nilai Truthy dan Falsy

| Nilai                          | Evaluasi |
|--------------------------------|----------|
| `benar`                        | truthy   |
| `salah`                        | falsy    |
| `kosong`                       | falsy    |
| `0`                            | falsy    |
| `0.0`                          | falsy    |
| `""` (string kosong)           | falsy    |
| `[]` (array kosong)            | falsy    |
| Semua nilai lainnya            | truthy   |

---

## Operator

### Aritmatika

| Operator | Fungsi       | Contoh        | Hasil |
|----------|--------------|---------------|-------|
| `+`      | Penjumlahan  | `10 + 3`      | `13`  |
| `-`      | Pengurangan  | `10 - 3`      | `7`   |
| `*`      | Perkalian    | `10 * 3`      | `30`  |
| `/`      | Pembagian    | `10 / 3`      | `3`   |
| `%`      | Sisa bagi    | `10 % 3`      | `1`   |

> Pembagian dua `bilangan` menghasilkan `bilangan` (integer division). Gunakan `desimal()` jika butuh hasil pecahan.

Operator `+` juga berlaku untuk penggabungan string:

```ky
buat salam = "Halo, " + "Dunia!"   // "Halo, Dunia!"
```

### Perbandingan

| Operator | Fungsi              | Contoh       | Hasil  |
|----------|---------------------|--------------|--------|
| `==`     | Sama dengan         | `5 == 5`     | `benar`|
| `!=`     | Tidak sama dengan   | `5 != 3`     | `benar`|
| `<`      | Kurang dari         | `3 < 5`      | `benar`|
| `>`      | Lebih dari          | `5 > 3`      | `benar`|
| `<=`     | Kurang dari / sama  | `3 <= 3`     | `benar`|
| `>=`     | Lebih dari / sama   | `5 >= 5`     | `benar`|

### Logika

| Operator       | Fungsi      | Contoh               | Hasil   |
|----------------|-------------|----------------------|---------|
| `dan`          | AND logika  | `benar dan salah`    | `salah` |
| `atau`         | OR logika   | `benar atau salah`   | `benar` |
| `tidak` / `bukan` | NOT logika | `tidak benar`     | `salah` |

### Prefix

| Operator | Fungsi             | Contoh  | Hasil   |
|----------|--------------------|---------|---------|
| `-`      | Negasi numerik     | `-5`    | `-5`    |
| `!`      | Negasi logika      | `!benar`| `salah` |

### Pengaksesan Elemen

```ky
buat arr  = [10, 20, 30]
buat map  = {"a": 1}

arr[0]      // 10
arr[-1]     // 30  (indeks negatif didukung)
map["a"]    // 1
```

---

## Kata Kunci & Alias

Bahasa ky-
lang mendukung kata kunci utama berbahasa Indonesia beserta alias berbahasa Inggris untuk kemudahan transisi.

| Kata Kunci Utama | Alias                       | Fungsi                      |
|------------------|-----------------------------|-----------------------------|
| `buat`           | —                           | Deklarasi variabel          |
| `fungsi`         | `def`, `fn`                 | Definisi fungsi             |
| `balik`          | `kembali`, `return`         | Kembalikan nilai (return)   |
| `jika`           | `if`                        | Percabangan kondisi         |
| `lain`           | `else`                      | Cabang alternatif           |
| `ulang`          | `for`                       | Perulangan gaya for         |
| `selama`         | `while`                     | Perulangan gaya while       |
| `henti`          | `stop`, `break`             | Keluar dari loop            |
| `lanjut`         | `skip`, `continue`          | Lanjut ke iterasi berikutnya|
| `benar`          | `true`                      | Nilai logika benar          |
| `salah`          | `false`                     | Nilai logika salah          |
| `kosong`         | `nil`                       | Nilai null / tidak ada nilai|
| `dan`            | —                           | Operator logika AND         |
| `atau`           | —                           | Operator logika OR          |
| `tidak`          | `bukan`                     | Operator logika NOT         |

---

## Struktur Kontrol

### Kondisi — `jika` / `lain`

```ky
jika kondisi {
    // blok jika kondisi benar
}
```

```ky
jika kondisi {
    // blok jika benar
} lain {
    // blok jika salah
}
```

**Kondisi bertingkat:**

```ky
jika nilai >= 90 {
    tampilkan("A")
} lain {
    jika nilai >= 80 {
        tampilkan("B")
    } lain {
        jika nilai >= 70 {
            tampilkan("C")
        } lain {
            tampilkan("D")
        }
    }
}
```

### Perulangan For — `ulang`

```ky
ulang buat i = 0; i < 10; i = i + 1 {
    tampilkan(i)
}
```

Format: `ulang <inisialisasi>; <kondisi>; <pembaruan> { ... }`

### Perulangan While — `selama`

```ky
buat n = 5
selama n > 0 {
    tampilkan(n)
    n = n - 1
}
```

### Keluar Loop — `henti`

```ky
ulang buat i = 0; i < 100; i = i + 1 {
    jika i == 10 {
        henti
    }
    tampilkan(i)
}
```

### Lewati Iterasi — `lanjut`

```ky
ulang buat i = 0; i < 10; i = i + 1 {
    jika i % 2 == 0 {
        lanjut
    }
    tampilkan(i)   // hanya mencetak bilangan ganjil
}
```

---

## Fungsi

### Definisi Fungsi

```ky
buat tambah = fungsi(a, b) {
    balik a + b
}
```

Fungsi adalah nilai kelas satu (first-class value) — bisa disimpan di variabel, dioper sebagai argumen, dan dikembalikan dari fungsi lain.

### Memanggil Fungsi

```ky
buat hasil = tambah(10, 20)
tampilkan(hasil)   // 30
```

### Fungsi Tanpa Nama (Anonim)

```ky
buat kali2 = fungsi(x) { balik x * 2 }
tampilkan(kali2(5))   // 10
```

### Rekursi

```ky
buat faktorial = fungsi(n) {
    jika n <= 1 {
        balik 1
    }
    balik n * faktorial(n - 1)
}

tampilkan(faktorial(10))   // 3628800
```

### Closure

```ky
buat buat_penghitung = fungsi() {
    buat n = 0
    balik fungsi() {
        n = n + 1
        balik n
    }
}

buat hitung = buat_penghitung()
tampilkan(hitung())   // 1
tampilkan(hitung())   // 2
tampilkan(hitung())   // 3
```

### Fungsi sebagai Argumen (Higher-Order)

```ky
buat terapkan = fungsi(nilai, fn) {
    balik fn(nilai)
}

buat kuadrat = fungsi(x) { balik x * x }
tampilkan(terapkan(5, kuadrat))   // 25
```

---

## Larik (Array)

### Membuat Larik

```ky
buat angka = [1, 2, 3, 4, 5]
buat nama  = ["Andi", "Budi", "Cici"]
buat campur = [1, "dua", benar, kosong]
```

### Mengakses Elemen

```ky
buat arr = [10, 20, 30, 40, 50]
arr[0]    // 10  (indeks pertama)
arr[4]    // 50  (indeks terakhir)
arr[-1]   // 50  (indeks negatif: dari belakang)
arr[-2]   // 40
```

### Operasi Larik

```ky
buat arr = [1, 2, 3]

masukkan(arr, 4)           // [1, 2, 3, 4]  — tambah di akhir
sisipkan(arr, 1, 99)       // [1, 99, 2, 3] — sisipkan di indeks 1
hapus(arr, 0)              // [2, 3]         — hapus indeks 0
irisan(arr, 1, 3)          // [2, 3]         — potong arr[1..2]
balik_larik(arr)           // [3, 2, 1]
urutkan(arr)               // [1, 2, 3]
unik([1,1,2,2,3])          // [1, 2, 3]
gabung([1,2], [3,4])       // [1, 2, 3, 4]
rentang(1, 6)              // [1, 2, 3, 4, 5]
pertama(arr)               // 1
terakhir(arr)              // 3
ekor(arr)                  // [2, 3]
```

> **Catatan penting:** Larik bersifat immutable — fungsi seperti `masukkan()` dan `hapus()` mengembalikan larik **baru**, bukan memodifikasi larik asli.

```ky
buat asli = [1, 2, 3]
buat baru = masukkan(asli, 4)
tampilkan(asli)   // [1, 2, 3]  — tidak berubah
tampilkan(baru)   // [1, 2, 3, 4]
```

---

## Kamus (Map)

### Membuat Kamus

```ky
buat siswa = {
    "nama":  "Budi",
    "kelas": "10A",
    "nilai": 95
}
```

### Mengakses Nilai

```ky
siswa["nama"]    // "Budi"
siswa["nilai"]   // 95
```

### Operasi Kamus

```ky
kunci(siswa)                    // ["nama", "kelas", "nilai"]
nilai_kamus(siswa)              // ["Budi", "10A", 95]
ada_kunci(siswa, "nama")        // benar
ada_kunci(siswa, "alamat")      // salah
hapus_kamus(siswa, "kelas")     // kamus baru tanpa kunci "kelas"
```

### Kunci yang Didukung

Kamus menerima `bilangan`, `desimal`, `teks`, dan `logika` sebagai kunci.

```ky
buat m = {1: "satu", benar: "ya", "teks": "nilai"}
```

---

## Fungsi Bawaan Lengkap

### I/O — Input & Output

| Fungsi | Tanda Tangan | Keterangan |
|--------|--------------|------------|
| `tampilkan` | `tampilkan(nilai...)` | Cetak ke layar dengan baris baru. Alias: `cetak`. |
| `cetak` | `cetak(nilai...)` | Alias `tampilkan`. |
| `tulis` | `tulis(nilai...)` | Cetak tanpa baris baru. |
| `format` | `format(template, nilai...)` | Format string gaya `sprintf`. Contoh: `format("Halo %s, umur %d", nama, umur)` |
| `baca` | `baca(prompt?)` | Baca satu baris dari input pengguna. |
| `baca_berkas` | `baca_berkas(path)` | Baca seluruh isi file teks. |
| `tulis_berkas` | `tulis_berkas(path, isi)` | Tulis string ke file. |
| `cetak_galat` | `cetak_galat(nilai...)` | Cetak ke stderr. |

```ky
tampilkan("Halo", "Dunia")              // Halo Dunia
tulis("Ketik nama: ")
buat nama = baca()
tampilkan(format("Halo, %s!", nama))

buat isi = baca_berkas("data.txt")
tulis_berkas("output.txt", "Hasil: " + isi)
```

---

### Tipe & Pemeriksaan

| Fungsi | Keterangan |
|--------|------------|
| `jenis(nilai)` | Kembalikan nama tipe: `"bilangan"`, `"desimal"`, `"teks"`, `"logika"`, `"larik"`, `"kamus"`, `"fungsi"`, `"kosong"` |
| `adalah_bilangan(nilai)` | Cek apakah integer |
| `adalah_desimal(nilai)` | Cek apakah float |
| `adalah_teks(nilai)` | Cek apakah string |
| `adalah_logika(nilai)` | Cek apakah boolean |
| `adalah_larik(nilai)` | Cek apakah array |
| `adalah_kamus(nilai)` | Cek apakah map |
| `adalah_fungsi(nilai)` | Cek apakah fungsi (termasuk fungsi bawaan) |
| `adalah_kosong(nilai)` | Cek apakah kosong/null |

```ky
jenis(42)         // "bilangan"
jenis("halo")     // "teks"
jenis([1, 2])     // "larik"

adalah_bilangan(3.14)   // salah
adalah_desimal(3.14)    // benar
adalah_kosong(kosong)   // benar
```

---

### Konversi Tipe

| Fungsi | Keterangan |
|--------|------------|
| `bilangan(nilai)` | Konversi ke integer. Mendukung desimal, teks, logika. |
| `desimal(nilai)` | Konversi ke float. Mendukung bilangan, teks, logika. |
| `teks(nilai)` | Konversi nilai apa pun ke representasi teks. |
| `logika(nilai)` | Konversi ke boolean berdasarkan aturan truthy/falsy. |

```ky
bilangan("42")      // 42
bilangan(3.9)       // 3   (potong, bukan bulatkan)
bilangan(benar)     // 1

desimal("3.14")     // 3.14
desimal(5)          // 5.0

teks(100)           // "100"
teks(benar)         // "benar"
teks([1, 2, 3])     // "[1, 2, 3]"

logika(0)           // salah
logika(1)           // benar
logika("")          // salah
logika("isi")       // benar
```

---

### String — Operasi Teks

| Fungsi | Tanda Tangan | Keterangan |
|--------|--------------|------------|
| `panjang` | `panjang(nilai)` | Panjang string, array, atau kamus. Mendukung UTF-8. |
| `huruf_besar` | `huruf_besar(teks)` | Ubah ke UPPERCASE. |
| `huruf_kecil` | `huruf_kecil(teks)` | Ubah ke lowercase. |
| `judul` | `judul(teks)` | Ubah ke Title Case. |
| `potong` | `potong(teks)` | Hapus spasi di kedua sisi. |
| `potong_kiri` | `potong_kiri(teks, karakter)` | Hapus karakter di sisi kiri. |
| `potong_kanan` | `potong_kanan(teks, karakter)` | Hapus karakter di sisi kanan. |
| `ganti` | `ganti(teks, lama, baru, batas?)` | Ganti semua kemunculan. Opsional: batas jumlah penggantian. |
| `ada` | `ada(teks, sub)` | Cek apakah sub ada di teks. Juga berlaku untuk larik. |
| `indeks` | `indeks(teks, sub)` | Posisi pertama sub di teks. `-1` jika tidak ditemukan. |
| `awalan` | `awalan(teks, prefix)` | Cek apakah teks diawali prefix. |
| `akhiran` | `akhiran(teks, suffix)` | Cek apakah teks diakhiri suffix. |
| `hitung` | `hitung(teks, sub)` | Hitung kemunculan sub dalam teks. |
| `ulang_teks` | `ulang_teks(teks, n)` | Ulangi teks sebanyak n kali. |
| `balik_teks` | `balik_teks(teks)` | Balik urutan karakter. |
| `pisah` | `pisah(teks, pemisah)` | Pecah teks menjadi larik. |
| `gabung_teks` | `gabung_teks(pemisah, larik)` | Gabung larik teks dengan pemisah. |
| `irisan_teks` | `irisan_teks(teks, awal, akhir)` | Potong bagian teks. |

```ky
panjang("halo")              // 4
panjang("こんにちは")         // 5  (Unicode-aware)
huruf_besar("halo dunia")    // "HALO DUNIA"
huruf_kecil("HALO")          // "halo"
judul("halo dunia")          // "Halo Dunia"
potong("  spasi  ")          // "spasi"
ganti("aababab", "ab", "X")  // "aXXX"
ada("halo dunia", "dunia")   // benar
indeks("halo", "lo")         // 2
awalan("halo", "ha")         // benar
akhiran("halo", "lo")        // benar
hitung("abcabc", "abc")      // 2
ulang_teks("ab", 3)          // "ababab"
balik_teks("halo")           // "olah"
pisah("a,b,c", ",")          // ["a", "b", "c"]
gabung_teks("-", ["a","b"])   // "a-b"
```

---

### Larik — Operasi Array

| Fungsi | Tanda Tangan | Keterangan |
|--------|--------------|------------|
| `panjang` | `panjang(larik)` | Jumlah elemen. |
| `masukkan` | `masukkan(larik, nilai)` | Tambah elemen di akhir. Kembalikan larik baru. |
| `hapus` | `hapus(larik, indeks)` | Hapus elemen di indeks. Kembalikan larik baru. |
| `sisipkan` | `sisipkan(larik, indeks, nilai)` | Sisipkan elemen di posisi tertentu. |
| `pertama` | `pertama(larik)` | Elemen pertama. |
| `terakhir` | `terakhir(larik)` | Elemen terakhir. |
| `ekor` | `ekor(larik)` | Semua elemen kecuali yang pertama. |
| `gabung` | `gabung(larik1, larik2)` | Gabungkan dua larik. |
| `irisan` | `irisan(larik, awal, akhir)` | Potong larik dari indeks awal hingga akhir-1. |
| `balik_larik` | `balik_larik(larik)` | Balik urutan elemen. |
| `urutkan` | `urutkan(larik)` | Urutkan elemen (ascending). |
| `unik` | `unik(larik)` | Hapus duplikat, pertahankan urutan pertama kemunculan. |
| `rentang` | `rentang(awal, akhir)` | Buat larik bilangan dari awal hingga akhir-1. |
| `kosongkan` | `kosongkan()` | Buat larik kosong baru. |
| `ada` | `ada(larik, nilai)` | Cek apakah nilai ada di larik. |
| `indeks` | `indeks(larik, nilai)` | Posisi pertama nilai di larik. `-1` jika tidak ada. |

```ky
buat angka = [3, 1, 4, 1, 5, 9]

masukkan(angka, 2)          // [3, 1, 4, 1, 5, 9, 2]
hapus(angka, 0)             // [1, 4, 1, 5, 9]
sisipkan(angka, 2, 99)      // [3, 1, 99, 4, 1, 5, 9]
pertama(angka)              // 3
terakhir(angka)             // 9
ekor(angka)                 // [1, 4, 1, 5, 9]
irisan(angka, 1, 4)         // [1, 4, 1]
balik_larik(angka)          // [9, 5, 1, 4, 1, 3]
urutkan(angka)              // [1, 1, 3, 4, 5, 9]
unik(angka)                 // [3, 1, 4, 5, 9]
rentang(1, 6)               // [1, 2, 3, 4, 5]
ada(angka, 4)               // benar
indeks(angka, 5)            // 4
```

---

### Fungsi Tingkat Tinggi

| Fungsi | Tanda Tangan | Keterangan |
|--------|--------------|------------|
| `petakan` | `petakan(larik, fn)` | Terapkan `fn` ke setiap elemen, kembalikan larik hasil. |
| `saring` | `saring(larik, fn)` | Kembalikan larik elemen yang lolos predikat `fn`. |
| `lipat` | `lipat(larik, awal, fn)` | Reduksi larik ke satu nilai. `fn(akumulasi, elemen)`. |
| `setiap` | `setiap(larik, fn)` | `benar` jika semua elemen lolos predikat. |
| `salah_satu` | `salah_satu(larik, fn)` | `benar` jika minimal satu elemen lolos predikat. |
| `zip` | `zip(larik1, larik2)` | Pasangkan elemen dari dua larik jadi larik pasangan. |

```ky
buat angka = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

// petakan — map
buat ganda = petakan(angka, fungsi(x) { balik x * 2 })
// [2, 4, 6, 8, 10, 12, 14, 16, 18, 20]

// saring — filter
buat genap_saja = saring(angka, fungsi(x) { balik genap(x) })
// [2, 4, 6, 8, 10]

// lipat — reduce/fold
buat total = lipat(angka, 0, fungsi(akum, x) { balik akum + x })
// 55

// setiap — all
buat semua_positif = setiap(angka, fungsi(x) { balik x > 0 })
// benar

// salah_satu — any
buat ada_besar = salah_satu(angka, fungsi(x) { balik x > 9 })
// benar

// zip
buat nama  = ["Andi", "Budi", "Cici"]
buat nilai = [90, 85, 92]
buat pasangan = zip(nama, nilai)
// [["Andi", 90], ["Budi", 85], ["Cici", 92]]

// Rantai fungsi tingkat tinggi
buat hasil = lipat(
    saring(
        petakan(rentang(1, 11), fungsi(x) { balik x * x }),
        fungsi(x) { balik x > 25 }
    ),
    0,
    fungsi(a, b) { balik a + b }
)
// Jumlah kuadrat yang > 25 dari [1..10]
```

---

### Matematika

#### Fungsi Dasar

| Fungsi | Keterangan | Contoh |
|--------|------------|--------|
| `akar(x)` | Akar kuadrat | `akar(25)` → `5.0` |
| `mutlak(x)` | Nilai absolut | `mutlak(-7)` → `7` |
| `lantai(x)` | Pembulatan ke bawah | `lantai(4.9)` → `4.0` |
| `langit(x)` | Pembulatan ke atas | `langit(4.1)` → `5.0` |
| `bulat(x)` | Pembulatan normal | `bulat(4.5)` → `5.0` |
| `pangkat(basis, eksponen)` | Perpangkatan | `pangkat(2, 10)` → `1024.0` |
| `eksponen(x)` | e^x | `eksponen(1)` → `2.718...` |

#### Trigonometri

| Fungsi | Keterangan |
|--------|------------|
| `sin(x)` | Sinus (radian) |
| `cos(x)` | Kosinus (radian) |
| `tan(x)` | Tangen (radian) |
| `asin(x)` | Arcsinus |
| `acos(x)` | Arkkosinus |
| `atan(x)` | Arktangen |
| `atan2(y, x)` | Arktangen dua argumen |

#### Logaritma

| Fungsi | Keterangan |
|--------|------------|
| `log(x)` | Logaritma natural (ln) |
| `log2(x)` | Logaritma basis 2 |
| `log10(x)` | Logaritma basis 10 |

#### Statistik

| Fungsi | Keterangan |
|--------|------------|
| `maks(a, b, ...)` atau `maks(larik)` | Nilai terbesar |
| `min(a, b, ...)` atau `min(larik)` | Nilai terkecil |
| `jumlah(larik)` | Total semua elemen |
| `rata(larik)` | Rata-rata semua elemen |
| `ganjil(n)` | Cek apakah bilangan ganjil |
| `genap(n)` | Cek apakah bilangan genap |

#### Konstanta

| Fungsi | Nilai | Keterangan |
|--------|-------|------------|
| `PI()` | `3.141592653589793` | Konstanta π |
| `E()` | `2.718281828459045` | Bilangan Euler |
| `MAKS_BILANGAN()` | `9223372036854775807` | Nilai maksimum integer 64-bit |

```ky
buat pi = PI()
tampilkan(sin(pi / 2))          // 1
tampilkan(cos(0))               // 1
tampilkan(log2(1024))           // 10
tampilkan(log10(1000))          // 3

buat data = [10, 20, 30, 40, 50]
tampilkan(jumlah(data))         // 150
tampilkan(rata(data))           // 30.0
tampilkan(maks(data))           // 50
tampilkan(min(data))            // 10
```

---

### Acak

| Fungsi | Keterangan |
|--------|------------|
| `acak()` | Float acak antara 0.0 dan 1.0 |
| `acak(n)` | Integer acak antara 0 dan n-1 |
| `acak_rentang(min, maks)` | Integer acak dalam rentang \[min, maks] (inklusif kedua ujung) |
| `acak_pilih(larik)` | Pilih satu elemen acak dari larik |
| `acak_campur(larik)` | Acak urutan elemen larik (Fisher-Yates shuffle), kembalikan larik baru |

```ky
acak()                          // misal: 0.6823...
acak(100)                       // misal: 42
acak_rentang(1, 6)              // simulasi dadu: 1–6
acak_pilih(["merah", "hijau", "biru"])   // salah satu warna
acak_campur([1, 2, 3, 4, 5])    // misal: [3, 1, 5, 2, 4]
```

---

### Waktu

| Fungsi | Keterangan |
|--------|------------|
| `waktu_sekarang()` | Unix timestamp dalam milidetik |
| `tidur(ms)` | Jeda eksekusi selama `ms` milidetik |

```ky
buat mulai = waktu_sekarang()
// ... kode yang diukur ...
buat selesai = waktu_sekarang()
tampilkan("Durasi:", teks(selesai - mulai), "ms")
```

---

### Kamus

| Fungsi | Keterangan |
|--------|------------|
| `kunci(kamus)` | Kembalikan semua kunci sebagai larik |
| `nilai_kamus(kamus)` | Kembalikan semua nilai sebagai larik |
| `hapus_kamus(kamus, kunci)` | Kembalikan kamus baru tanpa entri tersebut |
| `ada_kunci(kamus, kunci)` | `benar` jika kunci ada |

---

### Utilitas

| Fungsi | Keterangan |
|--------|------------|
| `galat(pesan)` | Lempar error runtime dengan pesan tertentu |
| `tegas(kondisi, pesan?)` | Assert — error jika kondisi salah |

```ky
buat bagi = fungsi(a, b) {
    tegas(b != 0, "Pembagi tidak boleh nol")
    balik a / b
}

galat("Sesuatu salah!")
```

---

## Benchmark Kecepatan

Benchmark dijalankan pada mesin dengan interpreter .ky yang dikompilasi dari Go. Hasil ini mencerminkan kinerja aktual *tree-walk interpreter* — bukan compiler atau JIT — dan dimaksudkan sebagai referensi yang jujur.

> **Lingkungan pengujian:** Go 1.25, Linux x86-64, build tanpa optimasi tambahan.  
> Setiap angka adalah hasil rata-rata dari beberapa run.

| Skenario Uji | Durasi | Keterangan |
|---|---|---|
| Fibonacci rekursif `fib(25)` | ~144 ms | 242.785 pemanggilan fungsi rekursif |
| Perulangan 100.000 iterasi | ~29 ms | Loop `ulang` sederhana tanpa komputasi berat |
| Faktorial rekursif `fak(20)` × 1.000 kali | ~18 ms | Rekursi kecil berulang |
| Parse + evaluasi ekspresi × 5.000 kali | ~23 ms | Termasuk tokenisasi dan parsing ulang |
| `petakan` + `saring` + `lipat` × 100 kali | ~11 ms | Rantai HOF pada larik 10 elemen |
| Konkatenasi string 2.000 kali | ~2 ms | String append berulang |
| `masukkan` 3.000 elemen ke larik | ~61 ms | Alokasi larik baru tiap operasi (immutable) |
| Closure & scope × 10 kali | ~7 ms | Pembuatan dan pemanggilan closure |

### Interpretasi Jujur

**.ky adalah tree-walk interpreter**, bukan compiler. Ini berarti:

- Kecepatan eksekusi sekitar **10–100× lebih lambat** dari bahasa yang dikompilasi seperti Go, C, atau Java.
- Fibonacci rekursif `fib(30)` akan terasa lambat — ini wajar dan konsisten dengan interpreter sejenis seperti Ruby MRI atau Python (CPython) tanpa JIT.
- Untuk **tujuan edukatif**, kecepatan ini lebih dari cukup. Hampir semua latihan pemrograman, algoritma dasar, dan proyek kelas berjalan mulus.
- **Bottleneck utama** adalah alokasi objek Go untuk setiap nilai .ky dan biaya rekursi tree-walk per node AST.

### Apa yang Berjalan Cepat

- Operasi string dan array dengan fungsi bawaan (diimplementasi langsung di Go)
- Fungsi matematika (meneruskan langsung ke `math` package Go)
- Loop dan kondisi sederhana
- Program skrip pendek hingga menengah

### Apa yang Relatif Lambat

- Rekursi dalam (seperti Fibonacci tanpa memoization)
- Operasi array berulang pada larik besar (karena immutability — setiap `masukkan` mengalokasi larik baru)
- Program yang membutuhkan komputasi numerik intensif

---

## Filosofi & Desain

### Tujuan Utama

**.ky** dibangun dengan satu premis sederhana: hambatan bahasa asing tidak seharusnya menghalangi seseorang untuk belajar logika pemrograman. Dengan kata kunci Indonesia, pelajar dapat fokus pada *konsep* — variabel, kondisi, perulangan, fungsi — tanpa perlu menghafal terminologi asing secara bersamaan.

### Keputusan Arsitektur

**Tree-Walk Interpreter**  
Sengaja dipilih karena mudah dibaca, mudah di-debug, dan mudah diperluas. Tidak ada bytecode, tidak ada JIT — apa yang Anda tulis langsung dieksekusi node per node. Ini menjadikan kode interpreter sendiri sebagai bahan belajar yang baik.

**Kata Kunci Terpusat**  
Semua kata kunci ada di satu tempat: `token/token.go`. Menambah atau mengubah kata kunci cukup edit satu map — lexer dan parser mengikuti otomatis. Desain ini sengaja untuk memudahkan kontribusi.

**Dukungan Alias**  
`balik`, `kembali`, dan `return` semuanya identik. Begitu pula `jika` dan `if`. Ini memungkinkan pelajar yang terbiasa Python/JavaScript tetap bisa menulis kode valid, sambil perlahan beralih ke kata kunci Indonesia.

**Fungsi Bawaan, Bukan Kata Kunci**  
`tampilkan`, `panjang`, dan fungsi lainnya adalah *identifier* biasa, bukan kata kunci. Ini berarti pengguna bisa mendefinisikan ulang mereka jika perlu — dan interpreter tidak menjadi kaku.

**Singleton untuk Nilai Tetap**  
`BENAR`, `SALAH`, dan `KOSONG` adalah pointer singleton. Setiap perbandingan `== benar` tidak mengalokasi objek baru — zero allocation untuk nilai paling sering digunakan.

**Immutable Arrays**  
Array dirancang tidak bisa dimodifikasi di tempat. Setiap operasi modifikasi mengembalikan array baru. Ini memperkenalkan pola pemrograman fungsional secara alami kepada pelajar.

---

*Dokumentasi ini mencerminkan kondisi interpreter versi saat ini.*  
*Language .ky — Created and Developed by Farhan Kertadiwangsa*
