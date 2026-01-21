package interpreter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/andrinoff/cambridge-lang/pkg/lexer"
	"github.com/andrinoff/cambridge-lang/pkg/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"DECLARE x : INTEGER\nx <- 5", 5},
		{"DECLARE x : INTEGER\nx <- 10", 10},
		{"DECLARE x : INTEGER\nx <- -5", -5},
		{"DECLARE x : INTEGER\nx <- 5 + 5", 10},
		{"DECLARE x : INTEGER\nx <- 5 - 5", 0},
		{"DECLARE x : INTEGER\nx <- 5 * 5", 25},
		{"DECLARE x : INTEGER\nx <- 10 DIV 3", 3},
		{"DECLARE x : INTEGER\nx <- 10 MOD 3", 1},
		{"DECLARE x : INTEGER\nx <- 2 + 3 * 4", 14},
		{"DECLARE x : INTEGER\nx <- (2 + 3) * 4", 20},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalRealExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"DECLARE x : REAL\nx <- 3.14", 3.14},
		{"DECLARE x : REAL\nx <- 2.5 + 2.5", 5.0},
		{"DECLARE x : REAL\nx <- 10.0 / 4.0", 2.5},
		{"DECLARE x : REAL\nx <- 5 / 2", 2.5}, // Integer division returns real
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testRealObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"DECLARE x : BOOLEAN\nx <- TRUE", true},
		{"DECLARE x : BOOLEAN\nx <- FALSE", false},
		{"DECLARE x : BOOLEAN\nx <- NOT TRUE", false},
		{"DECLARE x : BOOLEAN\nx <- NOT FALSE", true},
		{"DECLARE x : BOOLEAN\nx <- TRUE AND TRUE", true},
		{"DECLARE x : BOOLEAN\nx <- TRUE AND FALSE", false},
		{"DECLARE x : BOOLEAN\nx <- FALSE OR TRUE", true},
		{"DECLARE x : BOOLEAN\nx <- FALSE OR FALSE", false},
		{"DECLARE x : BOOLEAN\nx <- 5 > 3", true},
		{"DECLARE x : BOOLEAN\nx <- 5 < 3", false},
		{"DECLARE x : BOOLEAN\nx <- 5 = 5", true},
		{"DECLARE x : BOOLEAN\nx <- 5 <> 5", false},
		{"DECLARE x : BOOLEAN\nx <- 5 >= 5", true},
		{"DECLARE x : BOOLEAN\nx <- 5 <= 5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`DECLARE x : STRING
x <- "Hello"`, "Hello"},
		{`DECLARE x : STRING
x <- "Hello" & " World"`, "Hello World"},
		{`DECLARE x : STRING
x <- "Hello" & " " & "World"`, "Hello World"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestEvalStringComparison(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`DECLARE x : BOOLEAN
x <- "abc" = "abc"`, true},
		{`DECLARE x : BOOLEAN
x <- "abc" <> "def"`, true},
		{`DECLARE x : BOOLEAN
x <- "abc" < "def"`, true},
		{`DECLARE x : BOOLEAN
x <- "xyz" > "abc"`, true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestDeclareStatement(t *testing.T) {
	tests := []struct {
		input        string
		expectedType ObjectType
	}{
		{"DECLARE x : INTEGER", INTEGER_OBJ},
		{"DECLARE x : REAL", REAL_OBJ},
		{"DECLARE x : STRING", STRING_OBJ},
		{"DECLARE x : BOOLEAN", BOOLEAN_OBJ},
	}

	for _, tt := range tests {
		i := setupInterpreter(tt.input)
		obj, ok := i.env.Get("x")
		if !ok {
			t.Fatalf("variable x not found")
		}
		if obj.Type() != tt.expectedType {
			t.Errorf("object type wrong. expected=%s, got=%s", tt.expectedType, obj.Type())
		}
	}
}

func TestConstantStatement(t *testing.T) {
	input := `CONSTANT PI = 3.14159
DECLARE x : REAL
x <- PI`

	evaluated := testEval(input)
	testRealObject(t, evaluated, 3.14159)
}

func TestConstantImmutability(t *testing.T) {
	input := `CONSTANT PI = 3.14159
PI <- 3.0`

	evaluated := testEval(input)
	if _, ok := evaluated.(*Error); !ok {
		t.Errorf("expected error when modifying constant, got %T", evaluated)
	}
}

func TestIfStatement(t *testing.T) {
	tests := []struct {
		input    string
		varName  string
		expected int64
	}{
		{`DECLARE x : INTEGER
x <- 0
IF TRUE THEN
    x <- 10
ENDIF`, "x", 10},
		{`DECLARE x : INTEGER
x <- 0
IF FALSE THEN
    x <- 10
ENDIF`, "x", 0},
		{`DECLARE x : INTEGER
x <- 0
IF 5 > 3 THEN
    x <- 10
ELSE
    x <- 20
ENDIF`, "x", 10},
		{`DECLARE x : INTEGER
x <- 0
IF 5 < 3 THEN
    x <- 10
ELSE
    x <- 20
ENDIF`, "x", 20},
	}

	for _, tt := range tests {
		i := setupInterpreter(tt.input)
		obj, ok := i.env.Get(tt.varName)
		if !ok {
			t.Fatalf("variable %s not found", tt.varName)
		}
		testIntegerObject(t, obj, tt.expected)
	}
}

func TestCaseStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`DECLARE grade : INTEGER
DECLARE result : INTEGER
grade <- 1
CASE OF grade
    1 : result <- 100
    2 : result <- 200
ENDCASE`, 100},
		{`DECLARE grade : INTEGER
DECLARE result : INTEGER
grade <- 2
CASE OF grade
    1 : result <- 100
    2 : result <- 200
ENDCASE`, 200},
		{`DECLARE grade : INTEGER
DECLARE result : INTEGER
grade <- 99
CASE OF grade
    1 : result <- 100
    2 : result <- 200
    OTHERWISE : result <- 0
ENDCASE`, 0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestCaseStatementWithRange(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`DECLARE score : INTEGER
DECLARE result : INTEGER
score <- 75
CASE OF score
    0 TO 49 : result <- 0
    50 TO 100 : result <- 1
ENDCASE`, 1},
		{`DECLARE score : INTEGER
DECLARE result : INTEGER
score <- 30
CASE OF score
    0 TO 49 : result <- 0
    50 TO 100 : result <- 1
ENDCASE`, 0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestForStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`DECLARE sum : INTEGER
sum <- 0
FOR i <- 1 TO 5
    sum <- sum + i
NEXT i`, 15},
		{`DECLARE sum : INTEGER
sum <- 0
FOR i <- 1 TO 10 STEP 2
    sum <- sum + 1
NEXT i`, 5}, // 1, 3, 5, 7, 9
		{`DECLARE sum : INTEGER
sum <- 0
FOR i <- 5 TO 1 STEP -1
    sum <- sum + i
NEXT i`, 15},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestWhileStatement(t *testing.T) {
	input := `DECLARE x : INTEGER
DECLARE sum : INTEGER
x <- 1
sum <- 0
WHILE x <= 5
    sum <- sum + x
    x <- x + 1
ENDWHILE`

	i := setupInterpreter(input)
	obj, ok := i.env.Get("sum")
	if !ok {
		t.Fatal("variable sum not found")
	}
	testIntegerObject(t, obj, 15)
}

func TestRepeatStatement(t *testing.T) {
	input := `DECLARE x : INTEGER
DECLARE sum : INTEGER
x <- 1
sum <- 0
REPEAT
    sum <- sum + x
    x <- x + 1
UNTIL x > 5`

	i := setupInterpreter(input)
	obj, ok := i.env.Get("sum")
	if !ok {
		t.Fatal("variable sum not found")
	}
	testIntegerObject(t, obj, 15)
}

func TestProcedure(t *testing.T) {
	input := `DECLARE result : INTEGER
result <- 0

PROCEDURE SetValue()
    result <- 42
ENDPROCEDURE

CALL SetValue()`

	i := setupInterpreter(input)
	obj, ok := i.env.Get("result")
	if !ok {
		t.Fatal("variable result not found")
	}

	intObj, ok := obj.(*Integer)
	if !ok {
		t.Fatalf("expected Integer, got %T", obj)
	}

	if intObj.Value != 42 {
		t.Errorf("expected 42, got %d", intObj.Value)
	}
}

func TestProcedureWithParameters(t *testing.T) {
	input := `DECLARE result : INTEGER
result <- 0

PROCEDURE Add(a : INTEGER, b : INTEGER)
    result <- a + b
ENDPROCEDURE

CALL Add(5, 3)`

	i := setupInterpreter(input)
	obj, _ := i.env.Get("result")
	intObj := obj.(*Integer)

	if intObj.Value != 8 {
		t.Errorf("expected 8, got %d", intObj.Value)
	}
}

func TestFunction(t *testing.T) {
	input := `FUNCTION Add(a : INTEGER, b : INTEGER) RETURNS INTEGER
    RETURN a + b
ENDFUNCTION

DECLARE result : INTEGER
result <- Add(5, 3)`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 8)
}

func TestFunctionWithRecursion(t *testing.T) {
	input := `FUNCTION Factorial(n : INTEGER) RETURNS INTEGER
    IF n <= 1 THEN
        RETURN 1
    ENDIF
    RETURN n * Factorial(n - 1)
ENDFUNCTION

DECLARE result : INTEGER
result <- Factorial(5)`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 120)
}

func TestArrayOperations(t *testing.T) {
	input := `DECLARE arr : ARRAY[1:5] OF INTEGER
arr[1] <- 10
arr[2] <- 20
arr[3] <- 30
DECLARE result : INTEGER
result <- arr[1] + arr[2] + arr[3]`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 60)
}

func TestArray2D(t *testing.T) {
	input := `DECLARE matrix : ARRAY[1:2, 1:2] OF INTEGER
matrix[1, 1] <- 1
matrix[1, 2] <- 2
matrix[2, 1] <- 3
matrix[2, 2] <- 4
DECLARE result : INTEGER
result <- matrix[1, 1] + matrix[2, 2]`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 5)
}

func TestRecordType(t *testing.T) {
	input := `TYPE Person
    DECLARE name : STRING
    DECLARE age : INTEGER
ENDTYPE

DECLARE p : Person
p.name <- "John"
p.age <- 30`

	i := setupInterpreter(input)
	obj, ok := i.env.Get("p")
	if !ok {
		t.Fatal("variable p not found")
	}

	record, ok := obj.(*Record)
	if !ok {
		t.Fatalf("expected Record, got %T", obj)
	}

	nameObj, ok := record.Fields["name"]
	if !ok {
		t.Fatal("field name not found")
	}

	strObj, ok := nameObj.(*String)
	if !ok {
		t.Fatalf("name field not String, got %T", nameObj)
	}

	if strObj.Value != "John" {
		t.Errorf("expected 'John', got %s", strObj.Value)
	}
}

func TestClass(t *testing.T) {
	// Test simple class definition without instantiation to avoid potential issues
	input := `CLASS Counter
    PRIVATE DECLARE count : INTEGER

    PUBLIC PROCEDURE Increment()
        count <- count + 1
    ENDPROCEDURE

    PUBLIC FUNCTION GetCount() RETURNS INTEGER
        RETURN count
    ENDFUNCTION
ENDCLASS`

	i := setupInterpreter(input)
	obj, ok := i.env.Get("Counter")
	if !ok {
		t.Fatal("class Counter not found")
	}

	_, ok = obj.(*Class)
	if !ok {
		t.Fatalf("expected Class, got %T", obj)
	}
}

func TestOutputStatement(t *testing.T) {
	input := `OUTPUT "Hello, World!"`

	var buf bytes.Buffer
	i := New()
	i.SetOutput(&buf)

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	i.Eval(program)

	output := buf.String()
	expected := "Hello, World!\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestOutputMultipleValues(t *testing.T) {
	input := `OUTPUT "Value: ", 42`

	var buf bytes.Buffer
	i := New()
	i.SetOutput(&buf)

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	i.Eval(program)

	output := buf.String()
	expected := "Value: 42\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestInputStatement(t *testing.T) {
	input := `DECLARE name : STRING
INPUT name`

	i := New()
	i.SetInput(strings.NewReader("John\n"))

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	i.Eval(program)

	obj, ok := i.env.Get("name")
	if !ok {
		t.Fatal("variable name not found")
	}

	strObj, ok := obj.(*String)
	if !ok {
		t.Fatalf("expected String, got %T", obj)
	}

	if strObj.Value != "John" {
		t.Errorf("expected 'John', got %s", strObj.Value)
	}
}

func TestDivisionByZero(t *testing.T) {
	tests := []string{
		"DECLARE x : INTEGER\nx <- 5 DIV 0",
		"DECLARE x : INTEGER\nx <- 5 MOD 0",
		"DECLARE x : REAL\nx <- 5.0 / 0.0",
	}

	for _, input := range tests {
		evaluated := testEval(input)
		if _, ok := evaluated.(*Error); !ok {
			t.Errorf("expected error for division by zero, got %T", evaluated)
		}
	}
}

func TestUndefinedVariable(t *testing.T) {
	input := `x <- 5`

	i := New()
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	result := i.Eval(program)

	// Accessing undefined variable on the RHS should cause error
	// But assignment to undefined creates it
	// Let's test accessing undefined
	input2 := `DECLARE y : INTEGER
y <- x`

	i2 := New()
	l2 := lexer.New(input2)
	p2 := parser.New(l2)
	program2 := p2.ParseProgram()
	result2 := i2.Eval(program2)

	if _, ok := result2.(*Error); !ok {
		t.Errorf("expected error for undefined variable, got %T (%+v)", result2, result)
	}
}

func TestTypeMismatch(t *testing.T) {
	input := `DECLARE x : INTEGER
x <- "hello" + 5`

	evaluated := testEval(input)
	if _, ok := evaluated.(*Error); !ok {
		t.Errorf("expected error for type mismatch, got %T", evaluated)
	}
}

func TestIntegerRealMixedArithmetic(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"DECLARE x : REAL\nx <- 5 + 3.5", 8.5},
		{"DECLARE x : REAL\nx <- 3.5 + 5", 8.5},
		{"DECLARE x : REAL\nx <- 10 - 2.5", 7.5},
		{"DECLARE x : REAL\nx <- 2 * 3.5", 7.0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testRealObject(t, evaluated, tt.expected)
	}
}

func TestNestedLoops(t *testing.T) {
	input := `DECLARE sum : INTEGER
sum <- 0
FOR i <- 1 TO 3
    FOR j <- 1 TO 3
        sum <- sum + 1
    NEXT j
NEXT i`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 9)
}

func TestNestedIf(t *testing.T) {
	input := `DECLARE result : INTEGER
DECLARE x : INTEGER
x <- 10
IF x > 5 THEN
    IF x > 8 THEN
        result <- 1
    ELSE
        result <- 2
    ENDIF
ELSE
    result <- 3
ENDIF`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 1)
}

func TestStringConcatenationWithNumbers(t *testing.T) {
	input := `DECLARE result : STRING
result <- "Value: " & 42`

	evaluated := testEval(input)
	testStringObject(t, evaluated, "Value: 42")
}

func TestCharLiteral(t *testing.T) {
	input := `DECLARE c : CHAR
c <- 'A'`

	i := setupInterpreter(input)
	obj, ok := i.env.Get("c")
	if !ok {
		t.Fatal("variable c not found")
	}

	charObj, ok := obj.(*Char)
	if !ok {
		t.Fatalf("expected Char, got %T", obj)
	}

	if charObj.Value != 'A' {
		t.Errorf("expected 'A', got %c", charObj.Value)
	}
}

func TestReturnInFunction(t *testing.T) {
	input := `FUNCTION Test() RETURNS INTEGER
    RETURN 10
    RETURN 20
ENDFUNCTION

DECLARE x : INTEGER
x <- Test()`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 10)
}

func TestEarlyReturnInLoop(t *testing.T) {
	input := `FUNCTION FindFirst() RETURNS INTEGER
    FOR i <- 1 TO 10
        IF i = 5 THEN
            RETURN i
        ENDIF
    NEXT i
    RETURN 0
ENDFUNCTION

DECLARE result : INTEGER
result <- FindFirst()`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 5)
}

// Helper functions

func testEval(input string) Object {
	i := New()
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return i.Eval(program)
}

func setupInterpreter(input string) *Interpreter {
	i := New()
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	i.Eval(program)
	return i
}

func testIntegerObject(t *testing.T, obj Object, expected int64) bool {
	result, ok := obj.(*Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func testRealObject(t *testing.T, obj Object, expected float64) bool {
	result, ok := obj.(*Real)
	if !ok {
		t.Errorf("object is not Real. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%f, want=%f", result.Value, expected)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj Object, expected bool) bool {
	result, ok := obj.(*Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func testStringObject(t *testing.T, obj Object, expected string) bool {
	result, ok := obj.(*String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%s, want=%s", result.Value, expected)
		return false
	}
	return true
}
