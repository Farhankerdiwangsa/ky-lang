// Created and Developed by Farhan Kertadiwangsa
package main

import (
	"fmt"
	"time"

	"ky/evaluator"
	"ky/lexer"
	"ky/object"
	"ky/parser"
)

type hasil struct {
	nama    string
	durasi  time.Duration
	iterasi int
}

func jalankan(src string, env *object.Environment) object.Object {
	l := lexer.New(src)
	p := parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) > 0 {
		for _, e := range p.Errors() {
			fmt.Println("PARSE ERROR:", e)
		}
		return nil
	}
	return evaluator.Eval(prog, env)
}

func bench(nama, kode string, iterasi int) hasil {
	env := object.NewEnvironment()
	mulai := time.Now()
	for i := 0; i < iterasi; i++ {
		jalankan(kode, env)
	}
	return hasil{nama: nama, durasi: time.Since(mulai), iterasi: iterasi}
}

func fmtDurasi(d time.Duration) string {
	switch {
	case d < time.Millisecond:
		return fmt.Sprintf("%.2f µs", float64(d.Microseconds()))
	case d < time.Second:
		return fmt.Sprintf("%.2f ms", float64(d.Milliseconds()))
	default:
		return fmt.Sprintf("%.3f s", d.Seconds())
	}
}

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║    Benchmark    ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()

	tests := []struct {
		nama    string
		kode    string
		iterasi int
	}{
		{
			nama: "Fibonacci rekursif fib(25)",
			kode: `
buat fib = fungsi(n) {
    jika n <= 1 { balik n }
    balik fib(n-1) + fib(n-2)
}
fib(25)`,
			iterasi: 1,
		},
		{
			nama: "Perulangan 100.000 iterasi",
			kode: `
buat s = 0
ulang buat i = 0; i < 100000; i = i + 1 {
    s = s + i
}
s`,
			iterasi: 1,
		},
		{
			nama: "Faktorial rekursif fak(20) x1000",
			kode: `
buat fak = fungsi(n) {
    jika n <= 1 { balik 1 }
    balik n * fak(n-1)
}
fak(20)`,
			iterasi: 1000,
		},
		{
			nama: "Parse + Eval ekspresi (5000x)",
			kode:    `(1 + 2) * 3 - 4 / 2 % 3`,
			iterasi: 5000,
		},
		{
			nama: "petakan + saring + lipat (HOF)",
			kode: `
buat data = rentang(1, 101)
buat genap_saja = saring(data, fungsi(x) { balik genap(x) })
buat total = lipat(genap_saja, 0, fungsi(a, b) { balik a + b })
total`,
			iterasi: 100,
		},
		{
			nama: "String concat 2000x",
			kode: `
buat s = ""
ulang buat i = 0; i < 2000; i = i + 1 {
    s = s + "x"
}
panjang(s)`,
			iterasi: 1,
		},
		{
			nama: "Array masukkan 3000 elemen",
			kode: `
buat arr = []
ulang buat i = 0; i < 3000; i = i + 1 {
    arr = masukkan(arr, i)
}
panjang(arr)`,
			iterasi: 1,
		},
		{
			nama: "Closure & scope",
			kode: `
buat buat_penghitung = fungsi() {
    buat n = 0
    balik fungsi() {
        n = n + 1
        balik n
    }
}
buat hitung = buat_penghitung()
ulang buat i = 0; i < 1000; i = i + 1 {
    hitung()
}`,
			iterasi: 10,
		},
	}

	var hasilList []hasil
	for _, t := range tests {
		fmt.Printf("  ⏳ %s...\n", t.nama)
		hasilList = append(hasilList, bench(t.nama, t.kode, t.iterasi))
	}

	fmt.Println()
	fmt.Printf("┌──────────────────────────────────────────────────────────────────────────┐\n")
	fmt.Printf("│ %-46s │ %12s │ %8s │\n", "Benchmark", "Durasi", "Iterasi")
	fmt.Printf("├──────────────────────────────────────────────────────────────────────────┤\n")
	for _, r := range hasilList {
		fmt.Printf("│ %-46s │ %12s │ %8d │\n", r.nama, fmtDurasi(r.durasi), r.iterasi)
	}
	fmt.Printf("└──────────────────────────────────────────────────────────────────────────┘\n")
	fmt.Println()
	fmt.Println("Created and Developed by Farhan Kertadiwangsa")
}
