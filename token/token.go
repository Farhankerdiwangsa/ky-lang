// Created and Developed by Farhan Kertadiwangsa
package token

// TokenType adalah jenis token dalam bahasa kylang
type TokenType string

// Token adalah unit terkecil hasil tokenisasi kode sumber
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Col     int
}

const (
	// Token khusus
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifier & literal
	IDENT  TokenType = "IDENT"
	INT    TokenType = "INT"
	FLOAT  TokenType = "FLOAT"
	STRING TokenType = "STRING"

	// Operator aritmatika
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	MODULO   TokenType = "%"

	// Operator perbandingan
	EQ     TokenType = "=="
	NOT_EQ TokenType = "!="
	LT     TokenType = "<"
	GT     TokenType = ">"
	LT_EQ  TokenType = "<="
	GT_EQ  TokenType = ">="

	// Operator prefix
	BANG TokenType = "!"

	// Pemisah
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	COLON     TokenType = ":"
	DOT       TokenType = "."
	NEWLINE   TokenType = "NEWLINE"

	// Tanda kurung & blok
	LPAREN   TokenType = "("
	RPAREN   TokenType = ")"
	LBRACE   TokenType = "{"
	RBRACE   TokenType = "}"
	LBRACKET TokenType = "["
	RBRACKET TokenType = "]"

	// =========================================================
	// KATA KUNCI BAHASA .ky
	//
	// Untuk menambah, mengganti, atau membuat alias kata kunci:
	// cukup edit map Keywords di bawah — Lexer & Parser
	// mengikuti secara otomatis tanpa perlu ubah file lain.
	// =========================================================

	FUNGSI TokenType = "fungsi"  // definisi fungsi
	BUAT   TokenType = "buat"    // deklarasi variabel
	BENAR  TokenType = "benar"   // nilai benar (true)
	SALAH  TokenType = "salah"   // nilai salah (false)
	JIKA   TokenType = "jika"    // percabangan if
	LAIN   TokenType = "lain"    // percabangan else
	BALIK  TokenType = "balik"   // kembalikan nilai (return)
	ULANG  TokenType = "ulang"   // perulangan for
	SELAMA TokenType = "selama"  // perulangan while
	DAN    TokenType = "dan"     // operator logika AND
	ATAU   TokenType = "atau"    // operator logika OR
	TIDAK  TokenType = "tidak"   // operator logika NOT
	KOSONG TokenType = "kosong"  // nilai null/nil
	HENTI  TokenType = "henti"   // keluar dari loop (break)
	LANJUT TokenType = "lanjut"  // lanjut ke iterasi berikutnya (continue)
)

// Keywords adalah SATU-SATUNYA tempat mendefinisikan kosakata bahasa kylang
//
// Cara menambah kata kunci baru:
//   1. Tambahkan konstanta TokenType di atas
//   2. Tambahkan entri di map ini
//   Lexer & Parser otomatis mengikuti — tidak perlu ubah file lain.
//
// Cara membuat alias (misal kata yang sama bisa ditulis dua cara):
//   Tambahkan baris: "alias": TOKEN_YANG_SAMA
var Keywords = map[string]TokenType{
	// Kata kunci utama
	"fungsi": FUNGSI,
	"buat":   BUAT,
	"benar":  BENAR,
	"salah":  SALAH,
	"jika":   JIKA,
	"lain":   LAIN,
	"balik":  BALIK,
	"ulang":  ULANG,
	"selama": SELAMA,
	"dan":    DAN,
	"atau":   ATAU,
	"tidak":  TIDAK,
	"kosong": KOSONG,
	"henti":  HENTI,
	"lanjut": LANJUT,

	// Alias alternatif (bisa dihapus/ditambah bebas)
	"bukan":      TIDAK,  // alias untuk tidak
	"nil":        KOSONG, // alias untuk kosong
	"stop":       HENTI,  // alias untuk henti
	"skip":       LANJUT, // alias untuk lanjut
	"kembali":    BALIK,  // alias untuk balik
	"else":       LAIN,   // alias untuk lain
	"if":         JIKA,   // alias untuk jika
	"for":        ULANG,  // alias untuk ulang
	"while":      SELAMA, // alias untuk selama
	"def":        FUNGSI, // alias untuk fungsi
	"fn":         FUNGSI, // alias singkat untuk fungsi
	"return":     BALIK,  // alias untuk balik
	"break":      HENTI,  // alias untuk henti
	"continue":   LANJUT, // alias untuk lanjut
	"true":       BENAR,  // alias untuk benar
	"false":      SALAH,  // alias untuk salah
}

// LookupIdent mengembalikan TokenType untuk identifier.
// Jika ada di Keywords, kembalikan jenis keyword-nya; jika tidak, IDENT biasa.
func LookupIdent(ident string) TokenType {
	if tok, ok := Keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// New membuat token baru
func New(t TokenType, lit string, line, col int) Token {
	return Token{Type: t, Literal: lit, Line: line, Col: col}
}
