// Package lexer implements the lexical analyzer for Cambridge Pseudocode
package lexer

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/andrinoff/cambridge-lang/token"
)

// Lexer performs lexical analysis on source code
type Lexer struct {
	input   string
	pos     int  // current position in input
	readPos int  // current reading position (after current char)
	ch      byte // current char under examination
	line    int  // current line number
	column  int  // current column number
}

// New creates a new Lexer instance
func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar()
	return l
}

// readChar reads the next character and advances position
func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0 // ASCII NUL signifies EOF
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
	l.column++
}

// peekChar returns the next character without advancing
func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case '(':
		tok = l.newToken(token.LPAREN, l.ch)
	case ')':
		tok = l.newToken(token.RPAREN, l.ch)
	case '[':
		tok = l.newToken(token.LBRACKET, l.ch)
	case ']':
		tok = l.newToken(token.RBRACKET, l.ch)
	case ':':
		tok = l.newToken(token.COLON, l.ch)
	case ',':
		tok = l.newToken(token.COMMA, l.ch)
	case '.':
		tok = l.newToken(token.DOT, l.ch)
	case '+':
		tok = l.newToken(token.PLUS, l.ch)
	case '-':
		tok = l.newToken(token.MINUS, l.ch)
	case '*':
		tok = l.newToken(token.ASTERISK, l.ch)
	case '^':
		tok = l.newToken(token.CARET, l.ch)
	case '&':
		tok = l.newToken(token.AMPERSAND, l.ch)
	case '=':
		tok = l.newToken(token.EQ, l.ch)
	case '/':
		if l.peekChar() == '/' {
			// Comment - skip to end of line
			l.skipComment()
			return l.NextToken()
		}
		tok = l.newToken(token.SLASH, l.ch)
	case '<':
		if l.peekChar() == '>' {
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: "<>", Line: l.line, Column: l.column - 1}
		} else if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.LT_EQ, Literal: "<=", Line: l.line, Column: l.column - 1}
		} else if l.peekChar() == '-' {
			l.readChar()
			tok = token.Token{Type: token.ASSIGN, Literal: "<-", Line: l.line, Column: l.column - 1}
		} else {
			tok = l.newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.GT_EQ, Literal: ">=", Line: l.line, Column: l.column - 1}
		} else {
			tok = l.newToken(token.GT, l.ch)
		}
	case '"':
		tok.Type = token.STRING_LIT
		tok.Literal = l.readString()
		tok.Line = l.line
		tok.Column = l.column
		return tok
	case '\'':
		tok.Type = token.CHAR_LIT
		tok.Literal = l.readCharLiteral()
		tok.Line = l.line
		tok.Column = l.column
		return tok
	case '\n':
		tok = l.newToken(token.NEWLINE, l.ch)
		l.line++
		l.column = 0
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = l.line
		tok.Column = l.column
		return tok
	default:
		// Check for Unicode arrow ←
		if l.isArrow() {
			tok = token.Token{Type: token.ASSIGN, Literal: "←", Line: l.line, Column: l.column}
			l.readChar() // ← is multi-byte, need to skip remaining bytes
			l.readChar()
			return tok
		}
		if isLetter(l.ch) {
			tok.Column = l.column
			tok.Line = l.line
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(strings.ToUpper(tok.Literal))
			return tok
		} else if isDigit(l.ch) {
			tok.Column = l.column
			tok.Line = l.line
			literal, isReal := l.readNumber()
			tok.Literal = literal
			if isReal {
				tok.Type = token.REAL_LIT
			} else {
				tok.Type = token.INTEGER_LIT
			}
			return tok
		} else {
			tok = l.newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// isArrow checks if current position starts with Unicode arrow ←
func (l *Lexer) isArrow() bool {
	if l.pos+2 < len(l.input) {
		// ← is encoded as E2 86 90 in UTF-8
		return l.input[l.pos] == 0xE2 && l.input[l.pos+1] == 0x86 && l.input[l.pos+2] == 0x90
	}
	return false
}

// newToken creates a new token
func (l *Lexer) newToken(tokenType token.Type, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: l.line, Column: l.column}
}

// skipWhitespace skips spaces and tabs (but not newlines)
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

// skipComment skips from // to end of line
func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

// readIdentifier reads an identifier
func (l *Lexer) readIdentifier() string {
	start := l.pos
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[start:l.pos]
}

// readNumber reads a number (integer or real)
func (l *Lexer) readNumber() (string, bool) {
	start := l.pos
	isReal := false

	for isDigit(l.ch) {
		l.readChar()
	}

	// Check for decimal point
	if l.ch == '.' && isDigit(l.peekChar()) {
		isReal = true
		l.readChar() // consume '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[start:l.pos], isReal
}

// readString reads a string literal
func (l *Lexer) readString() string {
	l.readChar() // skip opening quote
	start := l.pos

	for l.ch != '"' && l.ch != 0 && l.ch != '\n' {
		l.readChar()
	}

	str := l.input[start:l.pos]
	if l.ch == '"' {
		l.readChar() // skip closing quote
	}
	return str
}

// readCharLiteral reads a character literal
func (l *Lexer) readCharLiteral() string {
	l.readChar() // skip opening quote
	ch := string(l.ch)
	l.readChar() // move past the character
	if l.ch == '\'' {
		l.readChar() // skip closing quote
	}
	return ch
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch))
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// Error represents a lexer error
type Error struct {
	Line    int
	Column  int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("Lexer error at line %d, column %d: %s", e.Line, e.Column, e.Message)
}
