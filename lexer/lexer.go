// Created and Developed by Farhan Kertadiwangsa
package lexer

import (
	"ky/token"
)

// Lexer membaca karakter per karakter dan menghasilkan stream Token.
// Dioptimasi untuk alokasi memori minimal.
type Lexer struct {
	input   string
	pos     int  // posisi karakter saat ini
	readPos int  // posisi baca berikutnya
	ch      byte // karakter saat ini
	line    int
	col     int
}

// New membuat Lexer baru dari string kode sumber
func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, col: 0}
	l.advance()
	return l
}

// NextToken menghasilkan token berikutnya dari kode sumber
func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	line, col := l.line, l.col

	switch l.ch {
	case '+':
		l.advance()
		return token.New(token.PLUS, "+", line, col)
	case '-':
		l.advance()
		return token.New(token.MINUS, "-", line, col)
	case '*':
		l.advance()
		return token.New(token.ASTERISK, "*", line, col)
	case '/':
		if l.peek() == '/' {
			l.skipLineComment()
			return l.NextToken()
		}
		l.advance()
		return token.New(token.SLASH, "/", line, col)
	case '%':
		l.advance()
		return token.New(token.MODULO, "%", line, col)
	case '=':
		if l.peek() == '=' {
			l.advance()
			l.advance()
			return token.New(token.EQ, "==", line, col)
		}
		l.advance()
		return token.New(token.ASSIGN, "=", line, col)
	case '!':
		if l.peek() == '=' {
			l.advance()
			l.advance()
			return token.New(token.NOT_EQ, "!=", line, col)
		}
		l.advance()
		return token.New(token.BANG, "!", line, col)
	case '<':
		if l.peek() == '=' {
			l.advance()
			l.advance()
			return token.New(token.LT_EQ, "<=", line, col)
		}
		l.advance()
		return token.New(token.LT, "<", line, col)
	case '>':
		if l.peek() == '=' {
			l.advance()
			l.advance()
			return token.New(token.GT_EQ, ">=", line, col)
		}
		l.advance()
		return token.New(token.GT, ">", line, col)
	case ',':
		l.advance()
		return token.New(token.COMMA, ",", line, col)
	case ';':
		l.advance()
		return token.New(token.SEMICOLON, ";", line, col)
	case ':':
		l.advance()
		return token.New(token.COLON, ":", line, col)
	case '.':
		l.advance()
		return token.New(token.DOT, ".", line, col)
	case '(':
		l.advance()
		return token.New(token.LPAREN, "(", line, col)
	case ')':
		l.advance()
		return token.New(token.RPAREN, ")", line, col)
	case '{':
		l.advance()
		return token.New(token.LBRACE, "{", line, col)
	case '}':
		l.advance()
		return token.New(token.RBRACE, "}", line, col)
	case '[':
		l.advance()
		return token.New(token.LBRACKET, "[", line, col)
	case ']':
		l.advance()
		return token.New(token.RBRACKET, "]", line, col)
	case '\n':
		l.line++
		l.col = 0
		l.advance()
		return token.New(token.NEWLINE, "\n", line, col)
	case '"':
		str := l.readString()
		return token.New(token.STRING, str, line, col)
	case 0:
		return token.New(token.EOF, "", line, col)
	default:
		if isLetter(l.ch) {
			ident := l.readIdentifier()
			tt := token.LookupIdent(ident)
			return token.New(tt, ident, line, col)
		}
		if isDigit(l.ch) {
			num, isFloat := l.readNumber()
			if isFloat {
				return token.New(token.FLOAT, num, line, col)
			}
			return token.New(token.INT, num, line, col)
		}
		ch := l.ch
		l.advance()
		return token.New(token.ILLEGAL, string(ch), line, col)
	}
}

// Tokenize mengembalikan semua token sekaligus (berguna untuk debugging)
func (l *Lexer) Tokenize() []token.Token {
	var tokens []token.Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == token.EOF {
			break
		}
	}
	return tokens
}

// --- Fungsi bantu internal ---

func (l *Lexer) advance() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
	l.col++
}

func (l *Lexer) peek() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.advance()
	}
}

func (l *Lexer) skipLineComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.advance()
	}
}

func (l *Lexer) readIdentifier() string {
	start := l.pos
	for isLetter(l.ch) || isDigit(l.ch) {
		l.advance()
	}
	return l.input[start:l.pos]
}

func (l *Lexer) readNumber() (string, bool) {
	start := l.pos
	isFloat := false
	for isDigit(l.ch) {
		l.advance()
	}
	if l.ch == '.' && isDigit(l.peek()) {
		isFloat = true
		l.advance()
		for isDigit(l.ch) {
			l.advance()
		}
	}
	return l.input[start:l.pos], isFloat
}

func (l *Lexer) readString() string {
	l.advance() // skip pembuka "
	start := l.pos
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' && l.peek() == '"' {
			l.advance()
		}
		l.advance()
	}
	str := l.input[start:l.pos]
	l.advance() // skip penutup "
	return str
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
