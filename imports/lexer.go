package imports

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

const eof = -1

type Token struct {
	Type  TokenType
	Value string
}

type TokenType string

const (
	COMMENT_START = "#"
	IMPORT        = "import"
	IDENTIFIER    = "identifier"
	DOT           = "dot"
	FROM          = "from"
	GLOB          = "glob"
	STRING        = "string"
	COMMA         = "comma"
	ILLEGAL       = "illegal"
	EOF           = "eof"
)

type Lexer struct {
	input  *bufio.Reader
	buffer bytes.Buffer
	offset int
	error  error
}

func NewLexer(input io.Reader) *Lexer {
	l := &Lexer{input: bufio.NewReader(input)}
	return l
}

func (l *Lexer) Next() rune {
	r, w, err := l.input.ReadRune()
	if err == io.EOF {
		return eof
	}
	l.offset = w
	l.buffer.WriteRune(r)
	return r
}

func (l *Lexer) Skip() rune {
	r, w, err := l.input.ReadRune()
	if err == io.EOF {
		return eof
	}
	l.offset += w
	return r
}

func (l *Lexer) Peek() rune {
	lead, err := l.input.Peek(1)
	if err == io.EOF {
		return eof
	} else if err != nil {
		l.Error(err.Error())
		return 0
	}

	p, err := l.input.Peek(runeLen(lead[0]))

	if err == io.EOF {
		return eof
	} else if err != nil {
		l.Error("unexpected input error")
		return 0
	}

	ruNe, _ := utf8.DecodeRune(p)
	return ruNe
}

func (l *Lexer) skipWhitespace() {
	ruNe := l.Peek()
	for unicode.IsSpace(ruNe) {
		l.Next()
		ruNe = l.Peek()
	}
	l.buffer.Reset()
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	switch l.Peek() {
	case '"':
		l.Skip()
		l.readString('"')
		return l.token(STRING)
	case '\'':
		l.Skip()
		l.readString('\'')
		return l.token(STRING)
	}
	next := l.Next()
	switch next {
	case '#':
		return l.token(COMMENT_START)
	case '.':
		return l.token(DOT)
	case '*':
		return l.token(GLOB)
	case ',':
		return l.token(COMMA)
	case eof:
		return l.token(EOF)

	default:
		if isLetter(next) {
			l.readIdentifier()
			text := l.TokenText()
			if text == "import" {
				return l.token(IMPORT)
			} else if text == "from" {
				return l.token(FROM)
			}
			return l.token(IDENTIFIER)
		}
		return l.token(ILLEGAL)
	}
}

func (l *Lexer) token(t TokenType) Token {
	return Token{
		Type:  t,
		Value: l.TokenText(),
	}
}

func (l *Lexer) TokenText() string {
	return l.buffer.String()
}

func runeLen(lead byte) int {
	if lead < 0xC0 {
		return 1
	} else if lead < 0xE0 {
		return 2
	} else if lead < 0xF0 {
		return 3
	}
	return 4
}

func (l *Lexer) readIdentifier() {
	next := l.Peek()
	for unicode.IsLetter(next) {
		l.Next()
		next = l.Peek()
	}
}

func (l *Lexer) readString(start rune) {
	for {
		next := l.Peek()
		if next == start {
			l.Skip()
			return
		}
		switch {
		case next == '\\':
			l.Skip()
			next := l.Peek()
			if next == start {
				l.Next()
			} else if next == 'b' {
				l.Skip()
				l.buffer.WriteRune('\b')
			} else if next == 'f' {
				l.Skip()
				l.buffer.WriteRune('\f')
			} else if next == 'n' {
				l.Skip()
				l.buffer.WriteRune('\n')
			} else if next == 'r' {
				l.Skip()
				l.buffer.WriteRune('\r')
			} else if next == 't' {
				l.Skip()
				l.buffer.WriteRune('\t')
			} else {
				l.Error("unsupported escape character")
				return
			}
		case unicode.IsControl(next):
			l.Error("cannot contain control characters in strings")
			return
		case next == eof:
			l.Error("unclosed string")
			return
		default:
			l.Next()
		}
	}
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) Error(e string) {
	err := fmt.Errorf("%s in %d", e, l.offset)
	l.error = err
}
