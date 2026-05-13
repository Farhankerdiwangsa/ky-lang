// Created and Developed by Farhan Kertadiwangsa
package parser

import (
	"fmt"
	"strconv"

	"ky/ast"
	"ky/lexer"
	"ky/token"
)

// Tingkat prioritas operator
const (
	_ int = iota
	LOWEST
	ATAU_PREC   // atau
	DAN_PREC    // dan
	EQUALS      // ==, !=
	LESSGREATER // <, >, <=, >=
	SUM         // +, -
	PRODUCT     // *, /, %
	PREFIX      // -x, !x, tidak x
	CALL        // fungsi()
	INDEX       // array[i]
)

var precedences = map[token.TokenType]int{
	token.ATAU:     ATAU_PREC,
	token.DAN:      DAN_PREC,
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.LT_EQ:    LESSGREATER,
	token.GT_EQ:    LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.MODULO:   PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser mengubah stream token menjadi AST menggunakan Pratt Parsing
type Parser struct {
	l      *lexer.Lexer
	cur    token.Token
	peek   token.Token
	errors []string

	prefixFns map[token.TokenType]prefixParseFn
	infixFns  map[token.TokenType]infixParseFn
}

// membuat Parser baru
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.prefixFns = map[token.TokenType]prefixParseFn{
		token.IDENT:    p.parseIdentifier,
		token.INT:      p.parseIntegerLiteral,
		token.FLOAT:    p.parseFloatLiteral,
		token.STRING:   p.parseStringLiteral,
		token.BENAR:    p.parseBooleanLiteral,
		token.SALAH:    p.parseBooleanLiteral,
		token.KOSONG:   p.parseKosongLiteral,
		token.BANG:     p.parsePrefixExpression,
		token.MINUS:    p.parsePrefixExpression,
		token.TIDAK:    p.parsePrefixExpression,
		token.LPAREN:   p.parseGroupedExpression,
		token.JIKA:     p.parseJikaExpression,
		token.SELAMA:   p.parseSelamaExpression,
		token.ULANG:    p.parseUlangExpression,
		token.FUNGSI:   p.parseFungsiLiteral,
		token.LBRACKET: p.parseArrayLiteral,
		token.LBRACE:   p.parseMapLiteral,
	}

	p.infixFns = map[token.TokenType]infixParseFn{
		token.PLUS:     p.parseInfixExpression,
		token.MINUS:    p.parseInfixExpression,
		token.ASTERISK: p.parseInfixExpression,
		token.SLASH:    p.parseInfixExpression,
		token.MODULO:   p.parseInfixExpression,
		token.EQ:       p.parseInfixExpression,
		token.NOT_EQ:   p.parseInfixExpression,
		token.LT:       p.parseInfixExpression,
		token.GT:       p.parseInfixExpression,
		token.LT_EQ:    p.parseInfixExpression,
		token.GT_EQ:    p.parseInfixExpression,
		token.DAN:      p.parseInfixExpression,
		token.ATAU:     p.parseInfixExpression,
		token.LPAREN:   p.parseCallExpression,
		token.LBRACKET: p.parseIndexExpression,
	}

	p.nextToken()
	p.nextToken()
	return p
}

// ParseProgram mem-parse seluruh program dan mengembalikan root AST
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	for p.cur.Type != token.EOF {
		p.skipNewlines()
		if p.cur.Type == token.EOF {
			break
		}
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// Errors mengembalikan daftar error parsing
func (p *Parser) Errors() []string { return p.errors }

// helper

func (p *Parser) nextToken() {
	p.cur = p.peek
	p.peek = p.l.NextToken()
}

func (p *Parser) skipNewlines() {
	for p.cur.Type == token.NEWLINE {
		p.nextToken()
	}
}

func (p *Parser) curTokenIs(t token.TokenType) bool  { return p.cur.Type == t }
func (p *Parser) peekTokenIs(t token.TokenType) bool { return p.peek.Type == t }

func (p *Parser) expect(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.errors = append(p.errors,
		fmt.Sprintf("baris %d: diharapkan '%s', tapi dapat '%s'",
			p.peek.Line, t, p.peek.Type))
	return false
}

func (p *Parser) curPrecedence() int {
	if pr, ok := precedences[p.cur.Type]; ok {
		return pr
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if pr, ok := precedences[p.peek.Type]; ok {
		return pr
	}
	return LOWEST
}

// Statements

func (p *Parser) parseStatement() ast.Statement {
	switch p.cur.Type {
	case token.BUAT:
		return p.parseBuatStatement()
	case token.BALIK:
		return p.parseBalikStatement()
	case token.HENTI:
		return &ast.HentiStatement{Token: p.cur}
	case token.LANJUT:
		return &ast.LanjutStatement{Token: p.cur}
	case token.IDENT:
		if p.peekTokenIs(token.ASSIGN) {
			return p.parseAssignStatement()
		}
		return p.parseExpressionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseBuatStatement() *ast.BuatStatement {
	stmt := &ast.BuatStatement{Token: p.cur}
	if !p.expect(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.cur, Value: p.cur.Literal}
	if !p.expect(token.ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if fl, ok := stmt.Value.(*ast.FungsiLiteral); ok {
		fl.Name = stmt.Name.Value
	}
	p.skipSemicolonOrNewline()
	return stmt
}

func (p *Parser) parseAssignStatement() *ast.AssignStatement {
	stmt := &ast.AssignStatement{
		Token: p.cur,
		Name:  &ast.Identifier{Token: p.cur, Value: p.cur.Literal},
	}
	p.nextToken() // skip IDENT
	p.nextToken() // skip =
	stmt.Value = p.parseExpression(LOWEST)
	p.skipSemicolonOrNewline()
	return stmt
}

func (p *Parser) parseBalikStatement() *ast.BalikStatement {
	stmt := &ast.BalikStatement{Token: p.cur}
	p.nextToken()
	if p.cur.Type != token.NEWLINE && p.cur.Type != token.SEMICOLON && p.cur.Type != token.RBRACE {
		stmt.ReturnValue = p.parseExpression(LOWEST)
	}
	p.skipSemicolonOrNewline()
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.cur}
	stmt.Expression = p.parseExpression(LOWEST)
	p.skipSemicolonOrNewline()
	return stmt
}

func (p *Parser) skipSemicolonOrNewline() {
	if p.peekTokenIs(token.SEMICOLON) || p.peekTokenIs(token.NEWLINE) {
		p.nextToken()
	}
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.cur}
	p.nextToken() // skip {
	p.skipNewlines()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
		p.skipNewlines()
	}
	return block
}

// Expressions

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixFns[p.cur.Type]
	if prefix == nil {
		p.errors = append(p.errors,
			fmt.Sprintf("baris %d: tidak bisa parse '%s'", p.cur.Line, p.cur.Type))
		return nil
	}
	left := prefix()

	for !p.peekTokenIs(token.NEWLINE) &&
		!p.peekTokenIs(token.SEMICOLON) &&
		!p.peekTokenIs(token.EOF) &&
		precedence < p.peekPrecedence() {
		infix := p.infixFns[p.peek.Type]
		if infix == nil {
			return left
		}
		p.nextToken()
		left = infix(left)
	}
	return left
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.cur, Value: p.cur.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	val, err := strconv.ParseInt(p.cur.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("baris %d: '%s' bukan bilangan bulat", p.cur.Line, p.cur.Literal))
		return nil
	}
	return &ast.IntegerLiteral{Token: p.cur, Value: val}
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	val, err := strconv.ParseFloat(p.cur.Literal, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("baris %d: '%s' bukan bilangan desimal", p.cur.Line, p.cur.Literal))
		return nil
	}
	return &ast.FloatLiteral{Token: p.cur, Value: val}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.cur, Value: p.cur.Literal}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.cur, Value: p.curTokenIs(token.BENAR)}
}

func (p *Parser) parseKosongLiteral() ast.Expression {
	return &ast.KosongLiteral{Token: p.cur}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{Token: p.cur, Operator: p.cur.Literal}
	p.nextToken()
	expr.Right = p.parseExpression(PREFIX)
	return expr
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token:    p.cur,
		Operator: p.cur.Literal,
		Left:     left,
	}
	prec := p.curPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(prec)
	return expr
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expect(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseJikaExpression() ast.Expression {
	expr := &ast.JikaExpression{Token: p.cur}
	p.nextToken()
	expr.Condition = p.parseExpression(LOWEST)
	if !p.expect(token.LBRACE) {
		return nil
	}
	expr.Consequence = p.parseBlockStatement()

	for p.peekTokenIs(token.NEWLINE) {
		p.nextToken()
	}
	if p.peekTokenIs(token.LAIN) {
		p.nextToken()
		if !p.expect(token.LBRACE) {
			return nil
		}
		expr.Alternative = p.parseBlockStatement()
	}
	return expr
}

func (p *Parser) parseSelamaExpression() ast.Expression {
	expr := &ast.SelamaExpression{Token: p.cur}
	p.nextToken()
	expr.Condition = p.parseExpression(LOWEST)
	if !p.expect(token.LBRACE) {
		return nil
	}
	expr.Body = p.parseBlockStatement()
	return expr
}

func (p *Parser) parseUlangExpression() ast.Expression {
	expr := &ast.UlangExpression{Token: p.cur}
	p.nextToken() // lewati 'ulang'

	// Init: buat i = 0
	if p.curTokenIs(token.BUAT) {
		expr.Init = p.parseBuatStatement()
		if p.curTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
	} else if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.ASSIGN) {
		expr.Init = p.parseAssignStatement()
		if p.curTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
	}
	p.skipNewlines()

	// Condition: i < 10
	if !p.curTokenIs(token.SEMICOLON) && !p.curTokenIs(token.LBRACE) && !p.curTokenIs(token.EOF) {
		expr.Condition = p.parseExpression(LOWEST)
	}
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	if p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	p.skipNewlines()

	//  Post: i = i + 1
	if !p.curTokenIs(token.LBRACE) && !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.ASSIGN) {
			expr.Post = p.parseAssignStatement()
		}
		p.skipNewlines()
	}

	//  Body
	if !p.curTokenIs(token.LBRACE) {
		if p.peekTokenIs(token.LBRACE) {
			p.nextToken()
		}
	}
	if !p.curTokenIs(token.LBRACE) {
		p.errors = append(p.errors, fmt.Sprintf("baris %d: diharapkan '{' pada ulang", p.cur.Line))
		return expr
	}
	expr.Body = p.parseBlockStatement()
	return expr
}

func (p *Parser) parseFungsiLiteral() ast.Expression {
	lit := &ast.FungsiLiteral{Token: p.cur}
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.Name = p.cur.Literal
	}
	if !p.expect(token.LPAREN) {
		return nil
	}
	lit.Parameters = p.parseFungsiParameters()
	if !p.expect(token.LBRACE) {
		return nil
	}
	lit.Body = p.parseBlockStatement()
	return lit
}

func (p *Parser) parseFungsiParameters() []*ast.Identifier {
	var params []*ast.Identifier
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return params
	}
	p.nextToken()
	params = append(params, &ast.Identifier{Token: p.cur, Value: p.cur.Literal})
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		params = append(params, &ast.Identifier{Token: p.cur, Value: p.cur.Literal})
	}
	if !p.expect(token.RPAREN) {
		return nil
	}
	return params
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	return &ast.CallExpression{
		Token:     p.cur,
		Function:  fn,
		Arguments: p.parseExpressionList(token.RPAREN),
	}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	return &ast.ArrayLiteral{
		Token:    p.cur,
		Elements: p.parseExpressionList(token.RBRACKET),
	}
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	var list []ast.Expression
	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}
	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}
	if !p.expect(end) {
		return nil
	}
	return list
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expr := &ast.IndexExpression{Token: p.cur, Left: left}
	p.nextToken()
	expr.Index = p.parseExpression(LOWEST)
	if !p.expect(token.RBRACKET) {
		return nil
	}
	return expr
}

func (p *Parser) parseMapLiteral() ast.Expression {
	ml := &ast.MapLiteral{Token: p.cur, Pairs: make(map[ast.Expression]ast.Expression)}
	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		p.skipNewlines()
		if p.curTokenIs(token.RBRACE) {
			break
		}
		key := p.parseExpression(LOWEST)
		if !p.expect(token.COLON) {
			return nil
		}
		p.nextToken()
		val := p.parseExpression(LOWEST)
		ml.Pairs[key] = val
		if !p.peekTokenIs(token.RBRACE) {
			if !p.expect(token.COMMA) {
				if !p.peekTokenIs(token.NEWLINE) {
					return nil
				}
			}
		}
		p.skipNewlines()
	}
	if !p.expect(token.RBRACE) {
		return nil
	}
	return ml
}
