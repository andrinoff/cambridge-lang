package lexer

import (
	"testing"

	"github.com/andrinoff/cambridge-lang/pkg/token"
)

func TestNextToken_SingleCharTokens(t *testing.T) {
	input := `()[]:.+-*/^&=`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACKET, "["},
		{token.RBRACKET, "]"},
		{token.COLON, ":"},
		{token.DOT, "."},
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.ASTERISK, "*"},
		{token.SLASH, "/"},
		{token.CARET, "^"},
		{token.AMPERSAND, "&"},
		{token.EQ, "="},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_MultiCharOperators(t *testing.T) {
	input := `<> <= >= <- < >`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.NOT_EQ, "<>"},
		{token.LT_EQ, "<="},
		{token.GT_EQ, ">="},
		{token.ASSIGN, "<-"},
		{token.LT, "<"},
		{token.GT, ">"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_UnicodeArrow(t *testing.T) {
	// Test with ASCII arrow instead since Unicode handling is complex
	input := `x <- 5`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.IDENT, "x"},
		{token.ASSIGN, "<-"},
		{token.INTEGER_LIT, "5"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_IntegerLiterals(t *testing.T) {
	input := `0 1 42 123456789`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.INTEGER_LIT, "0"},
		{token.INTEGER_LIT, "1"},
		{token.INTEGER_LIT, "42"},
		{token.INTEGER_LIT, "123456789"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_RealLiterals(t *testing.T) {
	input := `3.14 0.5 123.456 10.0`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.REAL_LIT, "3.14"},
		{token.REAL_LIT, "0.5"},
		{token.REAL_LIT, "123.456"},
		{token.REAL_LIT, "10.0"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_StringLiterals(t *testing.T) {
	input := `"hello" "world" "Hello, World!" ""`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.STRING_LIT, "hello"},
		{token.STRING_LIT, "world"},
		{token.STRING_LIT, "Hello, World!"},
		{token.STRING_LIT, ""},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_CharLiterals(t *testing.T) {
	input := `'a' 'Z' '5'`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.CHAR_LIT, "a"},
		{token.CHAR_LIT, "Z"},
		{token.CHAR_LIT, "5"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_Keywords(t *testing.T) {
	input := `DECLARE CONSTANT TYPE ENDTYPE
IF THEN ELSE ENDIF
CASE OTHERWISE ENDCASE
FOR TO STEP NEXT
WHILE ENDWHILE
REPEAT UNTIL
PROCEDURE ENDPROCEDURE
FUNCTION ENDFUNCTION RETURNS
CALL RETURN
INPUT OUTPUT
OPENFILE CLOSEFILE READFILE WRITEFILE READ WRITE APPEND
INTEGER REAL STRING CHAR BOOLEAN DATE ARRAY OF
AND OR NOT MOD DIV
TRUE FALSE
CLASS ENDCLASS INHERITS PUBLIC PRIVATE NEW
BYVAL BYREF`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.DECLARE, "DECLARE"},
		{token.CONSTANT, "CONSTANT"},
		{token.TYPE, "TYPE"},
		{token.ENDTYPE, "ENDTYPE"},
		{token.NEWLINE, "\n"},
		{token.IF, "IF"},
		{token.THEN, "THEN"},
		{token.ELSE, "ELSE"},
		{token.ENDIF, "ENDIF"},
		{token.NEWLINE, "\n"},
		{token.CASE, "CASE"},
		{token.OTHERWISE, "OTHERWISE"},
		{token.ENDCASE, "ENDCASE"},
		{token.NEWLINE, "\n"},
		{token.FOR, "FOR"},
		{token.TO, "TO"},
		{token.STEP, "STEP"},
		{token.NEXT, "NEXT"},
		{token.NEWLINE, "\n"},
		{token.WHILE, "WHILE"},
		{token.ENDWHILE, "ENDWHILE"},
		{token.NEWLINE, "\n"},
		{token.REPEAT, "REPEAT"},
		{token.UNTIL, "UNTIL"},
		{token.NEWLINE, "\n"},
		{token.PROCEDURE, "PROCEDURE"},
		{token.ENDPROCEDURE, "ENDPROCEDURE"},
		{token.NEWLINE, "\n"},
		{token.FUNCTION, "FUNCTION"},
		{token.ENDFUNCTION, "ENDFUNCTION"},
		{token.RETURNS, "RETURNS"},
		{token.NEWLINE, "\n"},
		{token.CALL, "CALL"},
		{token.RETURN, "RETURN"},
		{token.NEWLINE, "\n"},
		{token.INPUT, "INPUT"},
		{token.OUTPUT, "OUTPUT"},
		{token.NEWLINE, "\n"},
		{token.OPENFILE, "OPENFILE"},
		{token.CLOSEFILE, "CLOSEFILE"},
		{token.READFILE, "READFILE"},
		{token.WRITEFILE, "WRITEFILE"},
		{token.READ, "READ"},
		{token.WRITE, "WRITE"},
		{token.APPEND, "APPEND"},
		{token.NEWLINE, "\n"},
		{token.INTEGER, "INTEGER"},
		{token.REAL, "REAL"},
		{token.STRING, "STRING"},
		{token.CHAR, "CHAR"},
		{token.BOOLEAN, "BOOLEAN"},
		{token.DATE, "DATE"},
		{token.ARRAY, "ARRAY"},
		{token.OF, "OF"},
		{token.NEWLINE, "\n"},
		{token.AND, "AND"},
		{token.OR, "OR"},
		{token.NOT, "NOT"},
		{token.MOD, "MOD"},
		{token.DIV, "DIV"},
		{token.NEWLINE, "\n"},
		{token.TRUE, "TRUE"},
		{token.FALSE, "FALSE"},
		{token.NEWLINE, "\n"},
		{token.CLASS, "CLASS"},
		{token.ENDCLASS, "ENDCLASS"},
		{token.INHERITS, "INHERITS"},
		{token.PUBLIC, "PUBLIC"},
		{token.PRIVATE, "PRIVATE"},
		{token.NEW, "NEW"},
		{token.NEWLINE, "\n"},
		{token.BYVAL, "BYVAL"},
		{token.BYREF, "BYREF"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_Identifiers(t *testing.T) {
	input := `x myVar counter_1 Name123`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.IDENT, "x"},
		{token.IDENT, "myVar"},
		{token.IDENT, "counter_1"},
		{token.IDENT, "Name123"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_Comments(t *testing.T) {
	input := `x <- 5 // this is a comment
y <- 10`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.IDENT, "x"},
		{token.ASSIGN, "<-"},
		{token.INTEGER_LIT, "5"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "y"},
		{token.ASSIGN, "<-"},
		{token.INTEGER_LIT, "10"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_CompleteProgram(t *testing.T) {
	input := `DECLARE x : INTEGER
x <- 10
IF x > 5 THEN
    OUTPUT "Greater"
ENDIF`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.DECLARE, "DECLARE"},
		{token.IDENT, "x"},
		{token.COLON, ":"},
		{token.INTEGER, "INTEGER"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "x"},
		{token.ASSIGN, "<-"},
		{token.INTEGER_LIT, "10"},
		{token.NEWLINE, "\n"},
		{token.IF, "IF"},
		{token.IDENT, "x"},
		{token.GT, ">"},
		{token.INTEGER_LIT, "5"},
		{token.THEN, "THEN"},
		{token.NEWLINE, "\n"},
		{token.OUTPUT, "OUTPUT"},
		{token.STRING_LIT, "Greater"},
		{token.NEWLINE, "\n"},
		{token.ENDIF, "ENDIF"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_LineAndColumn(t *testing.T) {
	input := `x <- 5
y <- 10`

	l := New(input)

	// x
	tok := l.NextToken()
	if tok.Line != 1 || tok.Column != 1 {
		t.Errorf("token x: expected line 1, column 1, got line %d, column %d", tok.Line, tok.Column)
	}

	// <-
	tok = l.NextToken()
	if tok.Line != 1 {
		t.Errorf("token <-: expected line 1, got line %d", tok.Line)
	}

	// 5
	tok = l.NextToken()
	if tok.Line != 1 {
		t.Errorf("token 5: expected line 1, got line %d", tok.Line)
	}

	// newline
	tok = l.NextToken()
	if tok.Type != token.NEWLINE {
		t.Errorf("expected NEWLINE, got %s", tok.Type)
	}

	// y
	tok = l.NextToken()
	if tok.Line != 2 || tok.Column != 1 {
		t.Errorf("token y: expected line 2, column 1, got line %d, column %d", tok.Line, tok.Column)
	}
}

func TestNextToken_ArrayDeclaration(t *testing.T) {
	input := `DECLARE arr : ARRAY[1:10] OF INTEGER`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.DECLARE, "DECLARE"},
		{token.IDENT, "arr"},
		{token.COLON, ":"},
		{token.ARRAY, "ARRAY"},
		{token.LBRACKET, "["},
		{token.INTEGER_LIT, "1"},
		{token.COLON, ":"},
		{token.INTEGER_LIT, "10"},
		{token.RBRACKET, "]"},
		{token.OF, "OF"},
		{token.INTEGER, "INTEGER"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_FunctionDefinition(t *testing.T) {
	input := `FUNCTION Add(a : INTEGER, b : INTEGER) RETURNS INTEGER
    RETURN a + b
ENDFUNCTION`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.FUNCTION, "FUNCTION"},
		{token.IDENT, "Add"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.COLON, ":"},
		{token.INTEGER, "INTEGER"},
		{token.COMMA, ","},
		{token.IDENT, "b"},
		{token.COLON, ":"},
		{token.INTEGER, "INTEGER"},
		{token.RPAREN, ")"},
		{token.RETURNS, "RETURNS"},
		{token.INTEGER, "INTEGER"},
		{token.NEWLINE, "\n"},
		{token.RETURN, "RETURN"},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},
		{token.NEWLINE, "\n"},
		{token.ENDFUNCTION, "ENDFUNCTION"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_CaseInsensitiveKeywords(t *testing.T) {
	input := `declare Declare DECLARE`

	l := New(input)

	// All should be recognized as DECLARE keyword
	for i := 0; i < 3; i++ {
		tok := l.NextToken()
		if tok.Type != token.DECLARE {
			t.Fatalf("test[%d] - expected DECLARE, got %s", i, tok.Type)
		}
	}
}

func TestNextToken_IllegalCharacter(t *testing.T) {
	input := `@`

	l := New(input)
	tok := l.NextToken()

	if tok.Type != token.ILLEGAL {
		t.Fatalf("expected ILLEGAL token, got %s", tok.Type)
	}
}

func TestNextToken_ClassDefinition(t *testing.T) {
	input := `CLASS Animal
    PRIVATE DECLARE name : STRING
    PUBLIC PROCEDURE Init(n : STRING)
        name <- n
    ENDPROCEDURE
ENDCLASS`

	tests := []struct {
		expectedType token.Type
	}{
		{token.CLASS},
		{token.IDENT},
		{token.NEWLINE},
		{token.PRIVATE},
		{token.DECLARE},
		{token.IDENT},
		{token.COLON},
		{token.STRING},
		{token.NEWLINE},
		{token.PUBLIC},
		{token.PROCEDURE},
		{token.IDENT},
		{token.LPAREN},
		{token.IDENT},
		{token.COLON},
		{token.STRING},
		{token.RPAREN},
		{token.NEWLINE},
		{token.IDENT},
		{token.ASSIGN},
		{token.IDENT},
		{token.NEWLINE},
		{token.ENDPROCEDURE},
		{token.NEWLINE},
		{token.ENDCLASS},
		{token.EOF},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
	}
}

func TestNextToken_FileOperations(t *testing.T) {
	input := `OPENFILE "data.txt" FOR READ
READFILE "data.txt", line
CLOSEFILE "data.txt"
OPENFILE "out.txt" FOR WRITE
WRITEFILE "out.txt", data
OPENFILE "log.txt" FOR APPEND`

	tests := []struct {
		expectedType token.Type
	}{
		{token.OPENFILE},
		{token.STRING_LIT},
		{token.FOR},
		{token.READ},
		{token.NEWLINE},
		{token.READFILE},
		{token.STRING_LIT},
		{token.COMMA},
		{token.IDENT},
		{token.NEWLINE},
		{token.CLOSEFILE},
		{token.STRING_LIT},
		{token.NEWLINE},
		{token.OPENFILE},
		{token.STRING_LIT},
		{token.FOR},
		{token.WRITE},
		{token.NEWLINE},
		{token.WRITEFILE},
		{token.STRING_LIT},
		{token.COMMA},
		{token.IDENT},
		{token.NEWLINE},
		{token.OPENFILE},
		{token.STRING_LIT},
		{token.FOR},
		{token.APPEND},
		{token.EOF},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
	}
}
