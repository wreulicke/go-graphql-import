package imports

import (
	"bytes"
	"fmt"
)

type Parser struct {
	l      *Lexer
	errors []error

	curToken  Token
	peekToken Token
}

type ParseResult struct {
	Imports []Import
	From    string
}

type Import struct {
	Name   string
	Nested *Import
	isGlob bool
}

func NewParser(text string) *Parser {
	b := bytes.NewBufferString(text)
	l := NewLexer(b)
	p := &Parser{
		l: l,
	}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Parse() (*ParseResult, []error) {
	r := p.parseImportStatement()
	return r, p.errors
}

func (p *Parser) parseImportStatement() *ParseResult {
	if !p.curTokenIs(COMMENT_START) {
		p.curError(COMMENT_START)
		return nil
	}
	p.nextToken()

	if !p.curTokenIs(IMPORT) {
		p.curError(IMPORT)
		return nil
	}
	p.nextToken()

	imports := p.parseImports()
	if !p.curTokenIs(FROM) {
		p.curError(FROM)
		return nil
	}
	p.nextToken()
	fileName := p.parseFileName()
	return &ParseResult{
		Imports: imports,
		From:    fileName,
	}
}

func (p *Parser) parseImports() []Import {
	var imports []Import
	for !p.curTokenIs(FROM) && !p.curTokenIs(EOF) {
		switch p.curToken.Type {
		case GLOB:
			imports = append(imports, Import{
				isGlob: true,
			})
		case IDENTIFIER:
			t := p.curToken
			i := &Import{
				Name: t.Value,
			}

			if p.peekTokenIs(DOT) {
				p.nextToken()
				if p.peekTokenIs(IDENTIFIER) {
					p.nextToken()
					i.Nested = &Import{
						Name: p.curToken.Value,
					}
				} else if p.peekTokenIs(GLOB) {
					p.nextToken()
					i.Nested = &Import{
						isGlob: true,
					}
				}
			}
			imports = append(imports, *i)
		}
		p.nextToken()
	}
	return imports
}

func (p *Parser) parseFileName() string {
	if !p.curTokenIs(STRING) {
		p.curError(STRING)
		return ""
	}
	return p.curToken.Value
}

func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curError(t TokenType) {
	p.errors = append(p.errors, fmt.Errorf("expected next token to be %s, got %s instead", t, p.curToken.Type))
}

func (p *Parser) peekError(t TokenType) {
	p.errors = append(p.errors, fmt.Errorf("expected next token to be %s, got %s instead", t, p.peekToken.Type))
}
