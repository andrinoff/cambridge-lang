// Package token defines token types for Cambridge Pseudocode
// Based on Cambridge International AS & A Level Computer Science 9618 specification
package token

// Type represents the type of a token
type Type string

const (
	// Special tokens
	ILLEGAL Type = "ILLEGAL"
	EOF     Type = "EOF"
	NEWLINE Type = "NEWLINE"

	// Literals
	INTEGER_LIT Type = "INTEGER_LIT"
	REAL_LIT    Type = "REAL_LIT"
	STRING_LIT  Type = "STRING_LIT"
	CHAR_LIT    Type = "CHAR_LIT"
	TRUE        Type = "TRUE"
	FALSE       Type = "FALSE"

	// Identifier
	IDENT Type = "IDENT"

	// Data Type Keywords
	INTEGER Type = "INTEGER"
	REAL    Type = "REAL"
	STRING  Type = "STRING"
	CHAR    Type = "CHAR"
	BOOLEAN Type = "BOOLEAN"
	DATE    Type = "DATE"
	ARRAY   Type = "ARRAY"
	OF      Type = "OF"
	SET     Type = "SET"

	// Declaration Keywords
	DECLARE  Type = "DECLARE"
	CONSTANT Type = "CONSTANT"
	TYPE     Type = "TYPE"
	ENDTYPE  Type = "ENDTYPE"
	DEFINE   Type = "DEFINE"

	// Assignment
	ASSIGN Type = "ASSIGN" // ← or <-

	// Arithmetic Operators
	PLUS     Type = "PLUS"
	MINUS    Type = "MINUS"
	ASTERISK Type = "ASTERISK"
	SLASH    Type = "SLASH"
	MOD      Type = "MOD"
	DIV      Type = "DIV"

	// Comparison Operators
	EQ     Type = "EQ"     // =
	NOT_EQ Type = "NOT_EQ" // <>
	LT     Type = "LT"     // <
	GT     Type = "GT"     // >
	LT_EQ  Type = "LT_EQ"  // <=
	GT_EQ  Type = "GT_EQ"  // >=

	// Logical Operators
	AND Type = "AND"
	OR  Type = "OR"
	NOT Type = "NOT"

	// String Concatenation
	AMPERSAND Type = "AMPERSAND" // &

	// Selection
	IF        Type = "IF"
	THEN      Type = "THEN"
	ELSE      Type = "ELSE"
	ENDIF     Type = "ENDIF"
	CASE      Type = "CASE"
	OTHERWISE Type = "OTHERWISE"
	ENDCASE   Type = "ENDCASE"

	// Iteration
	FOR      Type = "FOR"
	TO       Type = "TO"
	STEP     Type = "STEP"
	NEXT     Type = "NEXT"
	WHILE    Type = "WHILE"
	ENDWHILE Type = "ENDWHILE"
	REPEAT   Type = "REPEAT"
	UNTIL    Type = "UNTIL"

	// Procedures and Functions
	PROCEDURE    Type = "PROCEDURE"
	ENDPROCEDURE Type = "ENDPROCEDURE"
	FUNCTION     Type = "FUNCTION"
	ENDFUNCTION  Type = "ENDFUNCTION"
	CALL         Type = "CALL"
	RETURN       Type = "RETURN"
	RETURNS      Type = "RETURNS"
	BYVAL        Type = "BYVAL"
	BYREF        Type = "BYREF"

	// Input/Output
	INPUT  Type = "INPUT"
	OUTPUT Type = "OUTPUT"

	// File Handling
	OPENFILE  Type = "OPENFILE"
	CLOSEFILE Type = "CLOSEFILE"
	READFILE  Type = "READFILE"
	WRITEFILE Type = "WRITEFILE"
	READ      Type = "READ"
	WRITE     Type = "WRITE"
	APPEND    Type = "APPEND"

	// OOP Keywords
	CLASS    Type = "CLASS"
	ENDCLASS Type = "ENDCLASS"
	INHERITS Type = "INHERITS"
	PUBLIC   Type = "PUBLIC"
	PRIVATE  Type = "PRIVATE"
	NEW      Type = "NEW"
	SUPER    Type = "SUPER"

	// Punctuation
	COLON    Type = "COLON"
	COMMA    Type = "COMMA"
	DOT      Type = "DOT"
	LPAREN   Type = "LPAREN"
	RPAREN   Type = "RPAREN"
	LBRACKET Type = "LBRACKET"
	RBRACKET Type = "RBRACKET"
	CARET    Type = "CARET" // ^ for pointers
)

// Token represents a lexical token
type Token struct {
	Type    Type
	Literal string
	Line    int
	Column  int
}

// keywords maps keyword strings to their token types
var Keywords = map[string]Type{
	// Data types
	"INTEGER": INTEGER,
	"REAL":    REAL,
	"STRING":  STRING,
	"CHAR":    CHAR,
	"BOOLEAN": BOOLEAN,
	"DATE":    DATE,
	"ARRAY":   ARRAY,
	"OF":      OF,
	"SET":     SET,

	// Boolean literals
	"TRUE":  TRUE,
	"FALSE": FALSE,

	// Declaration
	"DECLARE":  DECLARE,
	"CONSTANT": CONSTANT,
	"TYPE":     TYPE,
	"ENDTYPE":  ENDTYPE,
	"DEFINE":   DEFINE,

	// Arithmetic
	"MOD": MOD,
	"DIV": DIV,

	// Logical
	"AND": AND,
	"OR":  OR,
	"NOT": NOT,

	// Selection
	"IF":        IF,
	"THEN":      THEN,
	"ELSE":      ELSE,
	"ENDIF":     ENDIF,
	"CASE":      CASE,
	"OTHERWISE": OTHERWISE,
	"ENDCASE":   ENDCASE,

	// Iteration
	"FOR":      FOR,
	"TO":       TO,
	"STEP":     STEP,
	"NEXT":     NEXT,
	"WHILE":    WHILE,
	"ENDWHILE": ENDWHILE,
	"REPEAT":   REPEAT,
	"UNTIL":    UNTIL,

	// Procedures/Functions
	"PROCEDURE":    PROCEDURE,
	"ENDPROCEDURE": ENDPROCEDURE,
	"FUNCTION":     FUNCTION,
	"ENDFUNCTION":  ENDFUNCTION,
	"CALL":         CALL,
	"RETURN":       RETURN,
	"RETURNS":      RETURNS,
	"BYVAL":        BYVAL,
	"BYREF":        BYREF,

	// I/O
	"INPUT":  INPUT,
	"OUTPUT": OUTPUT,

	// File handling
	"OPENFILE":  OPENFILE,
	"CLOSEFILE": CLOSEFILE,
	"READFILE":  READFILE,
	"WRITEFILE": WRITEFILE,
	"READ":      READ,
	"WRITE":     WRITE,
	"APPEND":    APPEND,

	// OOP
	"CLASS":    CLASS,
	"ENDCLASS": ENDCLASS,
	"INHERITS": INHERITS,
	"PUBLIC":   PUBLIC,
	"PRIVATE":  PRIVATE,
	"NEW":      NEW,
	"SUPER":    SUPER,
}

// LookupIdent checks if an identifier is a keyword
func LookupIdent(ident string) Type {
	if tok, ok := Keywords[ident]; ok {
		return tok
	}
	return IDENT
}
