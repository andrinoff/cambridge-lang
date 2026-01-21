package parser

import (
	"testing"

	"github.com/andrinoff/cambridge-lang/pkg/ast"
	"github.com/andrinoff/cambridge-lang/pkg/lexer"
)

func TestParseDeclareStatement(t *testing.T) {
	tests := []struct {
		input        string
		expectedName string
		expectedType string
	}{
		{"DECLARE x : INTEGER", "x", "INTEGER"},
		{"DECLARE name : STRING", "name", "STRING"},
		{"DECLARE value : REAL", "value", "REAL"},
		{"DECLARE letter : CHAR", "letter", "CHAR"},
		{"DECLARE flag : BOOLEAN", "flag", "BOOLEAN"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.DeclareStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.DeclareStatement. got=%T",
				program.Statements[0])
		}

		if stmt.Name.Value != tt.expectedName {
			t.Errorf("stmt.Name.Value not '%s'. got=%s", tt.expectedName, stmt.Name.Value)
		}

		primitiveType, ok := stmt.DataType.(*ast.PrimitiveType)
		if !ok {
			t.Fatalf("stmt.DataType is not *ast.PrimitiveType. got=%T", stmt.DataType)
		}

		if primitiveType.Name != tt.expectedType {
			t.Errorf("stmt.DataType.Name not '%s'. got=%s", tt.expectedType, primitiveType.Name)
		}
	}
}

func TestParseArrayDeclaration(t *testing.T) {
	input := `DECLARE arr : ARRAY[1:10] OF INTEGER`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.DeclareStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.DeclareStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Name.Value != "arr" {
		t.Errorf("stmt.Name.Value not 'arr'. got=%s", stmt.Name.Value)
	}

	arrType, ok := stmt.DataType.(*ast.ArrayType)
	if !ok {
		t.Fatalf("stmt.DataType is not *ast.ArrayType. got=%T", stmt.DataType)
	}

	if len(arrType.Dimensions) != 1 {
		t.Fatalf("expected 1 dimension, got %d", len(arrType.Dimensions))
	}

	if arrType.Dimensions[0].Lower != 1 || arrType.Dimensions[0].Upper != 10 {
		t.Errorf("dimension bounds wrong. expected 1:10, got %d:%d",
			arrType.Dimensions[0].Lower, arrType.Dimensions[0].Upper)
	}
}

func TestParse2DArrayDeclaration(t *testing.T) {
	input := `DECLARE matrix : ARRAY[1:3,1:4] OF REAL`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.DeclareStatement)
	arrType := stmt.DataType.(*ast.ArrayType)

	if len(arrType.Dimensions) != 2 {
		t.Fatalf("expected 2 dimensions, got %d", len(arrType.Dimensions))
	}

	if arrType.Dimensions[0].Lower != 1 || arrType.Dimensions[0].Upper != 3 {
		t.Errorf("first dimension wrong. expected 1:3, got %d:%d",
			arrType.Dimensions[0].Lower, arrType.Dimensions[0].Upper)
	}

	if arrType.Dimensions[1].Lower != 1 || arrType.Dimensions[1].Upper != 4 {
		t.Errorf("second dimension wrong. expected 1:4, got %d:%d",
			arrType.Dimensions[1].Lower, arrType.Dimensions[1].Upper)
	}
}

func TestParseConstantStatement(t *testing.T) {
	tests := []struct {
		input         string
		expectedName  string
		expectedValue interface{}
	}{
		{"CONSTANT PI = 3.14159", "PI", 3.14159},
		{"CONSTANT MAX = 100", "MAX", int64(100)},
		{"CONSTANT NAME = \"Test\"", "NAME", "Test"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt, ok := program.Statements[0].(*ast.ConstantStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ConstantStatement. got=%T",
				program.Statements[0])
		}

		if stmt.Name.Value != tt.expectedName {
			t.Errorf("stmt.Name.Value not '%s'. got=%s", tt.expectedName, stmt.Name.Value)
		}
	}
}

func TestParseAssignmentStatement(t *testing.T) {
	input := `x <- 5`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.AssignmentStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Name.(*ast.Identifier)
	if !ok {
		t.Fatalf("stmt.Name is not *ast.Identifier. got=%T", stmt.Name)
	}

	if ident.Value != "x" {
		t.Errorf("ident.Value not 'x'. got=%s", ident.Value)
	}
}

func TestParseIntegerLiteral(t *testing.T) {
	input := `x <- 42`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.AssignmentStatement)
	lit, ok := stmt.Value.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("stmt.Value is not *ast.IntegerLiteral. got=%T", stmt.Value)
	}

	if lit.Value != 42 {
		t.Errorf("lit.Value not 42. got=%d", lit.Value)
	}
}

func TestParseRealLiteral(t *testing.T) {
	input := `x <- 3.14`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.AssignmentStatement)
	lit, ok := stmt.Value.(*ast.RealLiteral)
	if !ok {
		t.Fatalf("stmt.Value is not *ast.RealLiteral. got=%T", stmt.Value)
	}

	if lit.Value != 3.14 {
		t.Errorf("lit.Value not 3.14. got=%f", lit.Value)
	}
}

func TestParseStringLiteral(t *testing.T) {
	input := `x <- "hello"`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.AssignmentStatement)
	lit, ok := stmt.Value.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("stmt.Value is not *ast.StringLiteral. got=%T", stmt.Value)
	}

	if lit.Value != "hello" {
		t.Errorf("lit.Value not 'hello'. got=%s", lit.Value)
	}
}

func TestParseBooleanLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"x <- TRUE", true},
		{"x <- FALSE", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.AssignmentStatement)
		lit, ok := stmt.Value.(*ast.BooleanLiteral)
		if !ok {
			t.Fatalf("stmt.Value is not *ast.BooleanLiteral. got=%T", stmt.Value)
		}

		if lit.Value != tt.expected {
			t.Errorf("lit.Value not %v. got=%v", tt.expected, lit.Value)
		}
	}
}

func TestParsePrefixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
	}{
		{"x <- -5", "-"},
		{"x <- NOT TRUE", "NOT"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.AssignmentStatement)
		exp, ok := stmt.Value.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt.Value is not *ast.PrefixExpression. got=%T", stmt.Value)
		}

		if exp.Operator != tt.operator {
			t.Errorf("exp.Operator not '%s'. got=%s", tt.operator, exp.Operator)
		}
	}
}

func TestParseInfixExpressions(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"x <- 5 + 5", 5, "+", 5},
		{"x <- 5 - 5", 5, "-", 5},
		{"x <- 5 * 5", 5, "*", 5},
		{"x <- 5 / 5", 5, "/", 5},
		{"x <- 5 DIV 5", 5, "DIV", 5},
		{"x <- 5 MOD 5", 5, "MOD", 5},
		{"x <- 5 > 5", 5, ">", 5},
		{"x <- 5 < 5", 5, "<", 5},
		{"x <- 5 = 5", 5, "=", 5},
		{"x <- 5 <> 5", 5, "<>", 5},
		{"x <- 5 >= 5", 5, ">=", 5},
		{"x <- 5 <= 5", 5, "<=", 5},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.AssignmentStatement)
		exp, ok := stmt.Value.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt.Value is not *ast.InfixExpression. got=%T", stmt.Value)
		}

		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Errorf("exp.Operator not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func TestParseBooleanInfixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
	}{
		{"x <- TRUE AND FALSE", "AND"},
		{"x <- TRUE OR FALSE", "OR"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.AssignmentStatement)
		exp, ok := stmt.Value.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt.Value is not *ast.InfixExpression. got=%T", stmt.Value)
		}

		if exp.Operator != tt.operator {
			t.Errorf("exp.Operator not '%s'. got=%s", tt.operator, exp.Operator)
		}
	}
}

func TestParseStringConcatenation(t *testing.T) {
	input := `x <- "Hello" & " " & "World"`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.AssignmentStatement)
	_, ok := stmt.Value.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("stmt.Value is not *ast.InfixExpression. got=%T", stmt.Value)
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"x <- 1 + 2 * 3", "x <- (1 + (2 * 3))"},
		{"x <- 1 * 2 + 3", "x <- ((1 * 2) + 3)"},
		{"x <- (1 + 2) * 3", "x <- ((1 + 2) * 3)"},
		{"x <- 1 + 2 + 3", "x <- ((1 + 2) + 3)"},
		{"x <- -5 + 3", "x <- ((- 5) + 3)"},
		{"x <- NOT TRUE AND FALSE", "x <- ((NOT TRUE) AND FALSE)"},
		{"x <- TRUE OR FALSE AND TRUE", "x <- (TRUE OR (FALSE AND TRUE))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.Statements[0].String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestParseIfStatement(t *testing.T) {
	input := `IF x > 5 THEN
    OUTPUT "Greater"
ENDIF`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.IfStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Condition == nil {
		t.Fatal("stmt.Condition is nil")
	}

	if len(stmt.Consequence) != 1 {
		t.Errorf("Consequence should have 1 statement. got=%d", len(stmt.Consequence))
	}

	if stmt.Alternative != nil {
		t.Errorf("Alternative should be nil. got=%+v", stmt.Alternative)
	}
}

func TestParseIfElseStatement(t *testing.T) {
	input := `IF x > 5 THEN
    OUTPUT "Greater"
ELSE
    OUTPUT "Less or equal"
ENDIF`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.IfStatement)

	if len(stmt.Consequence) != 1 {
		t.Errorf("Consequence should have 1 statement. got=%d", len(stmt.Consequence))
	}

	if stmt.Alternative == nil {
		t.Fatal("Alternative should not be nil")
	}

	if len(stmt.Alternative) != 1 {
		t.Errorf("Alternative should have 1 statement. got=%d", len(stmt.Alternative))
	}
}

func TestParseCaseStatement(t *testing.T) {
	input := `CASE OF grade
    'A' : OUTPUT "Excellent"
    'B', 'C' : OUTPUT "Good"
    OTHERWISE : OUTPUT "Needs improvement"
ENDCASE`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.CaseStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.CaseStatement. got=%T",
			program.Statements[0])
	}

	if len(stmt.Cases) != 2 {
		t.Errorf("expected 2 cases, got %d", len(stmt.Cases))
	}

	if stmt.Otherwise == nil {
		t.Error("Otherwise should not be nil")
	}
}

func TestParseCaseWithRange(t *testing.T) {
	input := `CASE OF score
    0 TO 49 : OUTPUT "Fail"
    50 TO 100 : OUTPUT "Pass"
ENDCASE`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.CaseStatement)

	if len(stmt.Cases) != 2 {
		t.Errorf("expected 2 cases, got %d", len(stmt.Cases))
	}

	// Check first case has range expression
	_, ok := stmt.Cases[0].Values[0].(*ast.RangeExpression)
	if !ok {
		t.Errorf("first case value should be RangeExpression, got %T", stmt.Cases[0].Values[0])
	}
}

func TestParseForStatement(t *testing.T) {
	input := `FOR i <- 1 TO 10
    OUTPUT i
NEXT i`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ForStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Variable.Value != "i" {
		t.Errorf("stmt.Variable.Value not 'i'. got=%s", stmt.Variable.Value)
	}

	if stmt.Step != nil {
		t.Error("Step should be nil when not specified")
	}

	if len(stmt.Body) != 1 {
		t.Errorf("Body should have 1 statement. got=%d", len(stmt.Body))
	}
}

func TestParseForStatementWithStep(t *testing.T) {
	input := `FOR i <- 10 TO 1 STEP -1
    OUTPUT i
NEXT i`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ForStatement)

	if stmt.Step == nil {
		t.Fatal("Step should not be nil")
	}

	prefix, ok := stmt.Step.(*ast.PrefixExpression)
	if !ok {
		t.Fatalf("Step should be PrefixExpression, got %T", stmt.Step)
	}

	if prefix.Operator != "-" {
		t.Errorf("Step operator should be '-', got '%s'", prefix.Operator)
	}
}

func TestParseWhileStatement(t *testing.T) {
	input := `WHILE x < 10
    x <- x + 1
ENDWHILE`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.WhileStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Condition == nil {
		t.Fatal("Condition should not be nil")
	}

	if len(stmt.Body) != 1 {
		t.Errorf("Body should have 1 statement. got=%d", len(stmt.Body))
	}
}

func TestParseRepeatStatement(t *testing.T) {
	input := `REPEAT
    x <- x + 1
UNTIL x >= 10`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.RepeatStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.RepeatStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Condition == nil {
		t.Fatal("Condition should not be nil")
	}

	if len(stmt.Body) != 1 {
		t.Errorf("Body should have 1 statement. got=%d", len(stmt.Body))
	}
}

func TestParseProcedureStatement(t *testing.T) {
	input := `PROCEDURE Greet(name : STRING)
    OUTPUT "Hello, " & name
ENDPROCEDURE`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ProcedureStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ProcedureStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Name != "Greet" {
		t.Errorf("stmt.Name not 'Greet'. got=%s", stmt.Name)
	}

	if len(stmt.Parameters) != 1 {
		t.Fatalf("expected 1 parameter, got %d", len(stmt.Parameters))
	}

	if stmt.Parameters[0].Name != "name" {
		t.Errorf("parameter name not 'name'. got=%s", stmt.Parameters[0].Name)
	}
}

func TestParseProcedureWithByRef(t *testing.T) {
	input := `PROCEDURE Swap(BYREF a : INTEGER, BYREF b : INTEGER)
    DECLARE temp : INTEGER
    temp <- a
    a <- b
    b <- temp
ENDPROCEDURE`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ProcedureStatement)

	if len(stmt.Parameters) != 2 {
		t.Fatalf("expected 2 parameters, got %d", len(stmt.Parameters))
	}

	if !stmt.Parameters[0].ByRef {
		t.Error("first parameter should be BYREF")
	}

	if !stmt.Parameters[1].ByRef {
		t.Error("second parameter should be BYREF")
	}
}

func TestParseFunctionStatement(t *testing.T) {
	input := `FUNCTION Add(a : INTEGER, b : INTEGER) RETURNS INTEGER
    RETURN a + b
ENDFUNCTION`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.FunctionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.FunctionStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Name != "Add" {
		t.Errorf("stmt.Name not 'Add'. got=%s", stmt.Name)
	}

	if len(stmt.Parameters) != 2 {
		t.Fatalf("expected 2 parameters, got %d", len(stmt.Parameters))
	}

	returnType, ok := stmt.ReturnType.(*ast.PrimitiveType)
	if !ok {
		t.Fatalf("stmt.ReturnType is not *ast.PrimitiveType. got=%T", stmt.ReturnType)
	}

	if returnType.Name != "INTEGER" {
		t.Errorf("return type not 'INTEGER'. got=%s", returnType.Name)
	}
}

func TestParseCallStatement(t *testing.T) {
	input := `CALL Greet("World")`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.CallStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.CallStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Name.(*ast.Identifier)
	if !ok {
		t.Fatalf("stmt.Name is not *ast.Identifier. got=%T", stmt.Name)
	}

	if ident.Value != "Greet" {
		t.Errorf("stmt.Name not 'Greet'. got=%s", ident.Value)
	}

	if len(stmt.Arguments) != 1 {
		t.Errorf("expected 1 argument, got %d", len(stmt.Arguments))
	}
}

func TestParseReturnStatement(t *testing.T) {
	input := `RETURN 42`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ReturnStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Value == nil {
		t.Fatal("stmt.Value should not be nil")
	}
}

func TestParseInputStatement(t *testing.T) {
	input := `INPUT name`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.InputStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.InputStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Variable.(*ast.Identifier)
	if !ok {
		t.Fatalf("stmt.Variable is not *ast.Identifier. got=%T", stmt.Variable)
	}

	if ident.Value != "name" {
		t.Errorf("variable name not 'name'. got=%s", ident.Value)
	}
}

func TestParseOutputStatement(t *testing.T) {
	input := `OUTPUT "Hello", name, 42`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.OutputStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.OutputStatement. got=%T",
			program.Statements[0])
	}

	if len(stmt.Values) != 3 {
		t.Errorf("expected 3 values, got %d", len(stmt.Values))
	}
}

func TestParseOpenFileStatement(t *testing.T) {
	tests := []struct {
		input        string
		expectedMode string
	}{
		{`OPENFILE "data.txt" FOR READ`, "READ"},
		{`OPENFILE "data.txt" FOR WRITE`, "WRITE"},
		{`OPENFILE "data.txt" FOR APPEND`, "APPEND"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt, ok := program.Statements[0].(*ast.OpenFileStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.OpenFileStatement. got=%T",
				program.Statements[0])
		}

		if stmt.Mode != tt.expectedMode {
			t.Errorf("stmt.Mode not '%s'. got=%s", tt.expectedMode, stmt.Mode)
		}
	}
}

func TestParseCloseFileStatement(t *testing.T) {
	input := `CLOSEFILE "data.txt"`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	_, ok := program.Statements[0].(*ast.CloseFileStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.CloseFileStatement. got=%T",
			program.Statements[0])
	}
}

func TestParseReadFileStatement(t *testing.T) {
	input := `READFILE "data.txt", line`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ReadFileStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ReadFileStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Filename == nil {
		t.Fatal("stmt.Filename should not be nil")
	}

	if stmt.Variable == nil {
		t.Fatal("stmt.Variable should not be nil")
	}
}

func TestParseWriteFileStatement(t *testing.T) {
	input := `WRITEFILE "data.txt", "Hello"`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.WriteFileStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.WriteFileStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Filename == nil {
		t.Fatal("stmt.Filename should not be nil")
	}

	if stmt.Data == nil {
		t.Fatal("stmt.Data should not be nil")
	}
}

func TestParseTypeStatement(t *testing.T) {
	input := `TYPE Person
    DECLARE name : STRING
    DECLARE age : INTEGER
ENDTYPE`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.TypeStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.TypeStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Name != "Person" {
		t.Errorf("stmt.Name not 'Person'. got=%s", stmt.Name)
	}

	recType, ok := stmt.Definition.(*ast.RecordType)
	if !ok {
		t.Fatalf("stmt.Definition is not *ast.RecordType. got=%T", stmt.Definition)
	}

	if len(recType.Fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(recType.Fields))
	}
}

func TestParseClassStatement(t *testing.T) {
	input := `CLASS Animal
    PRIVATE DECLARE name : STRING
    PUBLIC PROCEDURE Init(n : STRING)
        name <- n
    ENDPROCEDURE
ENDCLASS`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ClassStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ClassStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Name != "Animal" {
		t.Errorf("stmt.Name not 'Animal'. got=%s", stmt.Name)
	}

	if stmt.Parent != "" {
		t.Errorf("stmt.Parent should be empty. got=%s", stmt.Parent)
	}
}

func TestParseClassWithInheritance(t *testing.T) {
	input := `CLASS Dog INHERITS Animal
    PUBLIC PROCEDURE Bark()
        OUTPUT "Woof!"
    ENDPROCEDURE
ENDCLASS`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ClassStatement)

	if stmt.Name != "Dog" {
		t.Errorf("stmt.Name not 'Dog'. got=%s", stmt.Name)
	}

	if stmt.Parent != "Animal" {
		t.Errorf("stmt.Parent not 'Animal'. got=%s", stmt.Parent)
	}
}

func TestParseArrayAccess(t *testing.T) {
	input := `x <- arr[5]`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.AssignmentStatement)
	access, ok := stmt.Value.(*ast.ArrayAccess)
	if !ok {
		t.Fatalf("stmt.Value is not *ast.ArrayAccess. got=%T", stmt.Value)
	}

	if len(access.Indices) != 1 {
		t.Errorf("expected 1 index, got %d", len(access.Indices))
	}
}

func TestParse2DArrayAccess(t *testing.T) {
	input := `x <- matrix[1, 2]`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.AssignmentStatement)
	access, ok := stmt.Value.(*ast.ArrayAccess)
	if !ok {
		t.Fatalf("stmt.Value is not *ast.ArrayAccess. got=%T", stmt.Value)
	}

	if len(access.Indices) != 2 {
		t.Errorf("expected 2 indices, got %d", len(access.Indices))
	}
}

func TestParseMemberAccess(t *testing.T) {
	input := `x <- person.name`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.AssignmentStatement)
	access, ok := stmt.Value.(*ast.MemberAccess)
	if !ok {
		t.Fatalf("stmt.Value is not *ast.MemberAccess. got=%T", stmt.Value)
	}

	if access.Member != "name" {
		t.Errorf("access.Member not 'name'. got=%s", access.Member)
	}
}

func TestParseCallExpression(t *testing.T) {
	input := `x <- Add(1, 2)`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.AssignmentStatement)
	call, ok := stmt.Value.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Value is not *ast.CallExpression. got=%T", stmt.Value)
	}

	ident, ok := call.Function.(*ast.Identifier)
	if !ok {
		t.Fatalf("call.Function is not *ast.Identifier. got=%T", call.Function)
	}

	if ident.Value != "Add" {
		t.Errorf("function name not 'Add'. got=%s", ident.Value)
	}

	if len(call.Arguments) != 2 {
		t.Errorf("expected 2 arguments, got %d", len(call.Arguments))
	}
}

func TestParseNewExpression(t *testing.T) {
	input := `x <- NEW Person("John", 30)`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.AssignmentStatement)
	newExpr, ok := stmt.Value.(*ast.NewExpression)
	if !ok {
		t.Fatalf("stmt.Value is not *ast.NewExpression. got=%T", stmt.Value)
	}

	if newExpr.ClassName != "Person" {
		t.Errorf("newExpr.ClassName not 'Person'. got=%s", newExpr.ClassName)
	}

	if len(newExpr.Arguments) != 2 {
		t.Errorf("expected 2 arguments, got %d", len(newExpr.Arguments))
	}
}

func TestParseGroupedExpression(t *testing.T) {
	input := `x <- (1 + 2) * 3`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.AssignmentStatement)
	infix, ok := stmt.Value.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("stmt.Value is not *ast.InfixExpression. got=%T", stmt.Value)
	}

	if infix.Operator != "*" {
		t.Errorf("outer operator should be '*'. got=%s", infix.Operator)
	}

	// Left should be the grouped expression (1 + 2)
	leftInfix, ok := infix.Left.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("left is not *ast.InfixExpression. got=%T", infix.Left)
	}

	if leftInfix.Operator != "+" {
		t.Errorf("inner operator should be '+'. got=%s", leftInfix.Operator)
	}
}

func TestParseArrayAssignment(t *testing.T) {
	input := `arr[5] <- 100`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.AssignmentStatement. got=%T",
			program.Statements[0])
	}

	_, ok = stmt.Name.(*ast.ArrayAccess)
	if !ok {
		t.Fatalf("stmt.Name is not *ast.ArrayAccess. got=%T", stmt.Name)
	}
}

func TestParseMemberAssignment(t *testing.T) {
	input := `person.name <- "John"`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.AssignmentStatement. got=%T",
			program.Statements[0])
	}

	_, ok = stmt.Name.(*ast.MemberAccess)
	if !ok {
		t.Fatalf("stmt.Name is not *ast.MemberAccess. got=%T", stmt.Name)
	}
}

func TestParserErrors(t *testing.T) {
	input := `DECLARE x`

	l := lexer.New(input)
	p := New(l)
	p.ParseProgram()

	errors := p.Errors()
	if len(errors) == 0 {
		t.Error("expected parser errors, got none")
	}
}

// Helper functions

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int64) bool {
	integ, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("exp not *ast.IntegerLiteral. got=%T", exp)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	return true
}
