package integration

import (
	"bytes"
	"strings"
	"testing"

	"github.com/andrinoff/cambridge-lang/pkg/builtins"
	"github.com/andrinoff/cambridge-lang/pkg/interpreter"
	"github.com/andrinoff/cambridge-lang/pkg/lexer"
	"github.com/andrinoff/cambridge-lang/pkg/parser"
)

// runProgram executes a Cambridge pseudocode program and returns the output
func runProgram(code string) (string, error) {
	var buf bytes.Buffer

	i := interpreter.New()
	i.SetBuiltins(builtins.GetBuiltins())
	i.SetOutput(&buf)

	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return "", &parserError{errors: p.Errors()}
	}

	result := i.Eval(program)
	if err, ok := result.(*interpreter.Error); ok {
		return buf.String(), &runtimeError{message: err.Message}
	}

	return buf.String(), nil
}

type parserError struct {
	errors []string
}

func (e *parserError) Error() string {
	return strings.Join(e.errors, "\n")
}

type runtimeError struct {
	message string
}

func (e *runtimeError) Error() string {
	return e.message
}

func TestIntegration_HelloWorld(t *testing.T) {
	code := `OUTPUT "Hello, World!"
OUTPUT "Welcome to Cambridge Pseudocode"`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Hello, World!\nWelcome to Cambridge Pseudocode\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestIntegration_Variables(t *testing.T) {
	code := `DECLARE Name : STRING
DECLARE Age : INTEGER
DECLARE Height : REAL
DECLARE IsStudent : BOOLEAN

CONSTANT PI = 3.14159
CONSTANT GREETING = "Hello"

Name <- "Alice"
Age <- 17
Height <- 1.65
IsStudent <- TRUE

OUTPUT GREETING, " ", Name, "!"
OUTPUT "Age: ", Age`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Hello Alice!") {
		t.Errorf("expected output to contain 'Hello Alice!', got %q", output)
	}
	if !strings.Contains(output, "Age: 17") {
		t.Errorf("expected output to contain 'Age: 17', got %q", output)
	}
}

func TestIntegration_Selection_If(t *testing.T) {
	code := `DECLARE Score : INTEGER
DECLARE Grade : STRING

Score <- 75

IF Score >= 90 THEN
    Grade <- "A"
ELSE
    IF Score >= 80 THEN
        Grade <- "B"
    ELSE
        IF Score >= 70 THEN
            Grade <- "C"
        ELSE
            IF Score >= 60 THEN
                Grade <- "D"
            ELSE
                Grade <- "F"
            ENDIF
        ENDIF
    ENDIF
ENDIF

OUTPUT "Score: ", Score, " Grade: ", Grade`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Score: 75 Grade: C") {
		t.Errorf("expected output to contain 'Score: 75 Grade: C', got %q", output)
	}
}

func TestIntegration_Selection_Case(t *testing.T) {
	code := `DECLARE DayNumber : INTEGER
DayNumber <- 3

CASE OF DayNumber
    1 : OUTPUT "Monday"
    2 : OUTPUT "Tuesday"
    3 : OUTPUT "Wednesday"
    4 : OUTPUT "Thursday"
    5 : OUTPUT "Friday"
    OTHERWISE : OUTPUT "Weekend"
ENDCASE`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Wednesday") {
		t.Errorf("expected output to contain 'Wednesday', got %q", output)
	}
}

func TestIntegration_ForLoop(t *testing.T) {
	code := `DECLARE i : INTEGER
FOR i <- 1 TO 5
    OUTPUT "i = ", i
NEXT i`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i := 1; i <= 5; i++ {
		expected := "i = " + string(rune('0'+i))
		if !strings.Contains(output, expected) {
			t.Errorf("expected output to contain %q, got %q", expected, output)
		}
	}
}

func TestIntegration_ForLoopWithStep(t *testing.T) {
	code := `DECLARE i : INTEGER
FOR i <- 0 TO 10 STEP 2
    OUTPUT i
NEXT i`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedValues := []string{"0", "2", "4", "6", "8", "10"}
	for _, v := range expectedValues {
		if !strings.Contains(output, v+"\n") {
			t.Errorf("expected output to contain %q, got %q", v, output)
		}
	}
}

func TestIntegration_WhileLoop(t *testing.T) {
	code := `DECLARE Sum : INTEGER
DECLARE Count : INTEGER

Sum <- 0
Count <- 0
WHILE Sum <= 50
    Count <- Count + 1
    Sum <- Sum + Count
ENDWHILE

OUTPUT "Final Count: ", Count
OUTPUT "Final Sum: ", Sum`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Final Count: 10") {
		t.Errorf("expected count to be 10, got %q", output)
	}
	if !strings.Contains(output, "Final Sum: 55") {
		t.Errorf("expected sum to be 55, got %q", output)
	}
}

func TestIntegration_RepeatLoop(t *testing.T) {
	code := `DECLARE i : INTEGER
i <- 5
REPEAT
    OUTPUT i
    i <- i - 1
UNTIL i = 0
OUTPUT "Done!"`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i := 5; i >= 1; i-- {
		if !strings.Contains(output, string(rune('0'+i))+"\n") {
			t.Errorf("expected output to contain %d, got %q", i, output)
		}
	}
	if !strings.Contains(output, "Done!") {
		t.Errorf("expected output to contain 'Done!', got %q", output)
	}
}

func TestIntegration_ProcedureCall(t *testing.T) {
	code := `PROCEDURE Greet(Name : STRING)
    OUTPUT "Hello, ", Name, "!"
ENDPROCEDURE

CALL Greet("World")`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Hello, World!") {
		t.Errorf("expected output to contain 'Hello, World!', got %q", output)
	}
}

func TestIntegration_Factorial(t *testing.T) {
	code := `FUNCTION Factorial(N : INTEGER) RETURNS INTEGER
    DECLARE Result : INTEGER
    DECLARE I : INTEGER
    Result <- 1
    FOR I <- 1 TO N
        Result <- Result * I
    NEXT I
    RETURN Result
ENDFUNCTION

OUTPUT "Factorial of 5: ", Factorial(5)
OUTPUT "Factorial of 7: ", Factorial(7)`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Factorial of 5: 120") {
		t.Errorf("expected 5! = 120, got %q", output)
	}
	if !strings.Contains(output, "Factorial of 7: 5040") {
		t.Errorf("expected 7! = 5040, got %q", output)
	}
}

func TestIntegration_IsPrime(t *testing.T) {
	code := `FUNCTION IsPrime(N : INTEGER) RETURNS BOOLEAN
    DECLARE I : INTEGER
    IF N < 2 THEN
        RETURN FALSE
    ENDIF
    FOR I <- 2 TO N - 1
        IF N MOD I = 0 THEN
            RETURN FALSE
        ENDIF
    NEXT I
    RETURN TRUE
ENDFUNCTION

OUTPUT "Is 7 prime? ", IsPrime(7)
OUTPUT "Is 10 prime? ", IsPrime(10)
OUTPUT "Is 13 prime? ", IsPrime(13)`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Is 7 prime? TRUE") {
		t.Errorf("expected 7 to be prime, got %q", output)
	}
	if !strings.Contains(output, "Is 10 prime? FALSE") {
		t.Errorf("expected 10 not to be prime, got %q", output)
	}
	if !strings.Contains(output, "Is 13 prime? TRUE") {
		t.Errorf("expected 13 to be prime, got %q", output)
	}
}

func TestIntegration_MaxFunction(t *testing.T) {
	code := `FUNCTION Max(A : INTEGER, B : INTEGER) RETURNS INTEGER
    IF A > B THEN
        RETURN A
    ELSE
        RETURN B
    ENDIF
ENDFUNCTION

OUTPUT "Max of 10 and 25: ", Max(10, 25)
OUTPUT "Max of 100 and 50: ", Max(100, 50)`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Max of 10 and 25: 25") {
		t.Errorf("expected max(10,25)=25, got %q", output)
	}
	if !strings.Contains(output, "Max of 100 and 50: 100") {
		t.Errorf("expected max(100,50)=100, got %q", output)
	}
}

func TestIntegration_Array1D(t *testing.T) {
	code := `DECLARE Numbers : ARRAY[1:5] OF INTEGER
DECLARE I : INTEGER
DECLARE Sum : INTEGER

Numbers[1] <- 10
Numbers[2] <- 20
Numbers[3] <- 30
Numbers[4] <- 40
Numbers[5] <- 50

Sum <- 0
FOR I <- 1 TO 5
    Sum <- Sum + Numbers[I]
NEXT I
OUTPUT "Sum: ", Sum
OUTPUT "Average: ", Sum / 5`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Sum: 150") {
		t.Errorf("expected sum=150, got %q", output)
	}
	if !strings.Contains(output, "Average: 30") {
		t.Errorf("expected average=30, got %q", output)
	}
}

func TestIntegration_Array2D(t *testing.T) {
	code := `DECLARE Matrix : ARRAY[1:3, 1:3] OF INTEGER
DECLARE I : INTEGER
DECLARE J : INTEGER

FOR I <- 1 TO 3
    FOR J <- 1 TO 3
        Matrix[I, J] <- I * J
    NEXT J
NEXT I

OUTPUT Matrix[1, 1], " ", Matrix[1, 2], " ", Matrix[1, 3]
OUTPUT Matrix[2, 1], " ", Matrix[2, 2], " ", Matrix[2, 3]
OUTPUT Matrix[3, 1], " ", Matrix[3, 2], " ", Matrix[3, 3]`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "1 2 3") {
		t.Errorf("expected row 1 to be '1 2 3', got %q", output)
	}
	if !strings.Contains(output, "2 4 6") {
		t.Errorf("expected row 2 to be '2 4 6', got %q", output)
	}
	if !strings.Contains(output, "3 6 9") {
		t.Errorf("expected row 3 to be '3 6 9', got %q", output)
	}
}

func TestIntegration_StringFunctions(t *testing.T) {
	code := `DECLARE Text : STRING

Text <- "Hello, World!"

OUTPUT "Original: ", Text
OUTPUT "Length: ", LENGTH(Text)
OUTPUT "LEFT(Text, 5): ", LEFT(Text, 5)
OUTPUT "RIGHT(Text, 6): ", RIGHT(Text, 6)
OUTPUT "MID(Text, 8, 5): ", MID(Text, 8, 5)
OUTPUT "UCASE: ", UCASE(Text)
OUTPUT "LCASE: ", LCASE(Text)`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		expected string
	}{
		{"Original: Hello, World!"},
		{"Length: 13"},
		{"LEFT(Text, 5): Hello"},
		{"RIGHT(Text, 6): World!"},
		{"MID(Text, 8, 5): World"},
		{"UCASE: HELLO, WORLD!"},
		{"LCASE: hello, world!"},
	}

	for _, tt := range tests {
		if !strings.Contains(output, tt.expected) {
			t.Errorf("expected output to contain %q, got %q", tt.expected, output)
		}
	}
}

func TestIntegration_ASCandCHR(t *testing.T) {
	code := `DECLARE Ch : CHAR

Ch <- 'A'
OUTPUT "Character: ", Ch
OUTPUT "ASC('A'): ", ASC(Ch)
OUTPUT "CHR(66): ", CHR(66)`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Character: A") {
		t.Errorf("expected 'Character: A', got %q", output)
	}
	if !strings.Contains(output, "ASC('A'): 65") {
		t.Errorf("expected 'ASC('A'): 65', got %q", output)
	}
	if !strings.Contains(output, "CHR(66): B") {
		t.Errorf("expected 'CHR(66): B', got %q", output)
	}
}

func TestIntegration_StringConcatenation(t *testing.T) {
	code := `DECLARE FirstName : STRING
DECLARE LastName : STRING
DECLARE FullName : STRING

FirstName <- "John"
LastName <- "Smith"
FullName <- FirstName & " " & LastName
OUTPUT "Full Name: ", FullName`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Full Name: John Smith") {
		t.Errorf("expected 'Full Name: John Smith', got %q", output)
	}
}

func TestIntegration_BuildAlphabet(t *testing.T) {
	code := `DECLARE Result : STRING
DECLARE I : INTEGER

Result <- ""
FOR I <- 65 TO 90
    Result <- Result & CHR(I)
NEXT I
OUTPUT Result`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if !strings.Contains(output, expected) {
		t.Errorf("expected alphabet, got %q", output)
	}
}

func TestIntegration_RecordType(t *testing.T) {
	code := `TYPE Person
    DECLARE name : STRING
    DECLARE age : INTEGER
ENDTYPE

DECLARE p : Person
p.name <- "Alice"
p.age <- 25

OUTPUT "Name: ", p.name
OUTPUT "Age: ", p.age`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Name: Alice") {
		t.Errorf("expected 'Name: Alice', got %q", output)
	}
	if !strings.Contains(output, "Age: 25") {
		t.Errorf("expected 'Age: 25', got %q", output)
	}
}

func TestIntegration_NestedLoops(t *testing.T) {
	code := `DECLARE i : INTEGER
DECLARE j : INTEGER
DECLARE count : INTEGER

count <- 0
FOR i <- 1 TO 3
    FOR j <- 1 TO 4
        count <- count + 1
    NEXT j
NEXT i

OUTPUT "Total iterations: ", count`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Total iterations: 12") {
		t.Errorf("expected 'Total iterations: 12', got %q", output)
	}
}

func TestIntegration_RecursiveFunction(t *testing.T) {
	code := `FUNCTION Fibonacci(N : INTEGER) RETURNS INTEGER
    IF N <= 1 THEN
        RETURN N
    ENDIF
    RETURN Fibonacci(N - 1) + Fibonacci(N - 2)
ENDFUNCTION

OUTPUT "Fibonacci(0) = ", Fibonacci(0)
OUTPUT "Fibonacci(1) = ", Fibonacci(1)
OUTPUT "Fibonacci(5) = ", Fibonacci(5)
OUTPUT "Fibonacci(10) = ", Fibonacci(10)`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		expected string
	}{
		{"Fibonacci(0) = 0"},
		{"Fibonacci(1) = 1"},
		{"Fibonacci(5) = 5"},
		{"Fibonacci(10) = 55"},
	}

	for _, tt := range tests {
		if !strings.Contains(output, tt.expected) {
			t.Errorf("expected output to contain %q, got %q", tt.expected, output)
		}
	}
}

func TestIntegration_MathFunctions(t *testing.T) {
	code := `OUTPUT "ABS(-5): ", ABS(-5)
OUTPUT "ABS(3.14): ", ABS(3.14)
OUTPUT "SQRT(16): ", SQRT(16)
OUTPUT "SQRT(2): ", SQRT(2)
OUTPUT "POW(2, 10): ", POW(2, 10)
OUTPUT "INT(3.7): ", INT(3.7)
OUTPUT "ROUND(3.14159, 2): ", ROUND(3.14159, 2)`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		expected string
	}{
		{"ABS(-5): 5"},
		{"ABS(3.14): 3.14"},
		{"SQRT(16): 4"},
		{"POW(2, 10): 1024"},
		{"INT(3.7): 3"},
		{"ROUND(3.14159, 2): 3.14"},
	}

	for _, tt := range tests {
		if !strings.Contains(output, tt.expected) {
			t.Errorf("expected output to contain %q, got %q", tt.expected, output)
		}
	}
}

func TestIntegration_NumToStrAndStrToNum(t *testing.T) {
	code := `DECLARE numStr : STRING
DECLARE num : INTEGER

numStr <- NUM_TO_STR(42)
OUTPUT "NUM_TO_STR(42): ", numStr

num <- STR_TO_NUM("123")
OUTPUT "STR_TO_NUM: ", num`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "NUM_TO_STR(42): 42") {
		t.Errorf("expected NUM_TO_STR result, got %q", output)
	}
	if !strings.Contains(output, "STR_TO_NUM: 123") {
		t.Errorf("expected STR_TO_NUM result, got %q", output)
	}
}

func TestIntegration_ComplexExpression(t *testing.T) {
	code := `DECLARE result : REAL

result <- (5 + 3) * 2 - 10 / 2
OUTPUT "Result: ", result

result <- 2 * 3 + 4 * 5
OUTPUT "2 * 3 + 4 * 5 = ", result

result <- (2 + 3) * (4 + 5)
OUTPUT "(2 + 3) * (4 + 5) = ", result`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Result: 11") {
		t.Errorf("expected 'Result: 11', got %q", output)
	}
	if !strings.Contains(output, "2 * 3 + 4 * 5 = 26") {
		t.Errorf("expected '2 * 3 + 4 * 5 = 26', got %q", output)
	}
	if !strings.Contains(output, "(2 + 3) * (4 + 5) = 45") {
		t.Errorf("expected '(2 + 3) * (4 + 5) = 45', got %q", output)
	}
}

func TestIntegration_BubbleSort(t *testing.T) {
	code := `DECLARE arr : ARRAY[1:5] OF INTEGER
DECLARE i : INTEGER
DECLARE j : INTEGER
DECLARE temp : INTEGER

arr[1] <- 5
arr[2] <- 2
arr[3] <- 8
arr[4] <- 1
arr[5] <- 9

FOR i <- 1 TO 4
    FOR j <- 1 TO 5 - i
        IF arr[j] > arr[j + 1] THEN
            temp <- arr[j]
            arr[j] <- arr[j + 1]
            arr[j + 1] <- temp
        ENDIF
    NEXT j
NEXT i

OUTPUT "Sorted array:"
FOR i <- 1 TO 5
    OUTPUT arr[i]
NEXT i`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that the output contains sorted values in order
	lines := strings.Split(output, "\n")
	var numbers []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "1" || line == "2" || line == "5" || line == "8" || line == "9" {
			numbers = append(numbers, line)
		}
	}

	expected := []string{"1", "2", "5", "8", "9"}
	for i, v := range expected {
		if i >= len(numbers) || numbers[i] != v {
			t.Errorf("sorting failed, expected %v, got %v", expected, numbers)
			break
		}
	}
}

func TestIntegration_BinarySearch(t *testing.T) {
	code := `FUNCTION BinarySearch(arr : ARRAY[1:10] OF INTEGER, target : INTEGER, low : INTEGER, high : INTEGER) RETURNS INTEGER
    DECLARE mid : INTEGER

    IF low > high THEN
        RETURN -1
    ENDIF

    mid <- (low + high) DIV 2

    IF arr[mid] = target THEN
        RETURN mid
    ENDIF

    IF arr[mid] > target THEN
        RETURN BinarySearch(arr, target, low, mid - 1)
    ELSE
        RETURN BinarySearch(arr, target, mid + 1, high)
    ENDIF
ENDFUNCTION

DECLARE numbers : ARRAY[1:10] OF INTEGER
DECLARE i : INTEGER

// Initialize sorted array
FOR i <- 1 TO 10
    numbers[i] <- i * 10
NEXT i

OUTPUT "Searching for 50: position ", BinarySearch(numbers, 50, 1, 10)
OUTPUT "Searching for 10: position ", BinarySearch(numbers, 10, 1, 10)
OUTPUT "Searching for 100: position ", BinarySearch(numbers, 100, 1, 10)
OUTPUT "Searching for 35: position ", BinarySearch(numbers, 35, 1, 10)`

	output, err := runProgram(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "Searching for 50: position 5") {
		t.Errorf("expected 50 at position 5, got %q", output)
	}
	if !strings.Contains(output, "Searching for 10: position 1") {
		t.Errorf("expected 10 at position 1, got %q", output)
	}
	if !strings.Contains(output, "Searching for 100: position 10") {
		t.Errorf("expected 100 at position 10, got %q", output)
	}
	if !strings.Contains(output, "Searching for 35: position -1") {
		t.Errorf("expected 35 not found (-1), got %q", output)
	}
}
