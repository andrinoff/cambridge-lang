package builtins

import (
	"testing"

	"github.com/andrinoff/cambridge-lang/pkg/interpreter"
)

func TestLength(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"hello", 5},
		{"", 0},
		{"Hello World", 11},
		{"123", 3},
	}

	builtins := GetBuiltins()
	lengthFn := builtins["LENGTH"]

	for _, tt := range tests {
		result := lengthFn.Fn(&interpreter.String{Value: tt.input})

		intResult, ok := result.(*interpreter.Integer)
		if !ok {
			t.Fatalf("expected Integer, got %T", result)
		}

		if intResult.Value != tt.expected {
			t.Errorf("LENGTH(%q) = %d, want %d", tt.input, intResult.Value, tt.expected)
		}
	}
}

func TestLengthWrongArgCount(t *testing.T) {
	builtins := GetBuiltins()
	lengthFn := builtins["LENGTH"]

	result := lengthFn.Fn()

	if _, ok := result.(*interpreter.Error); !ok {
		t.Errorf("expected Error for wrong arg count, got %T", result)
	}
}

func TestLengthWrongArgType(t *testing.T) {
	builtins := GetBuiltins()
	lengthFn := builtins["LENGTH"]

	result := lengthFn.Fn(&interpreter.Integer{Value: 42})

	if _, ok := result.(*interpreter.Error); !ok {
		t.Errorf("expected Error for wrong arg type, got %T", result)
	}
}

func TestLeft(t *testing.T) {
	tests := []struct {
		input    string
		n        int64
		expected string
	}{
		{"Hello", 3, "Hel"},
		{"Hello", 0, ""},
		{"Hello", 5, "Hello"},
		{"Hello", 10, "Hello"},
	}

	builtins := GetBuiltins()
	leftFn := builtins["LEFT"]

	for _, tt := range tests {
		result := leftFn.Fn(&interpreter.String{Value: tt.input}, &interpreter.Integer{Value: tt.n})

		strResult, ok := result.(*interpreter.String)
		if !ok {
			t.Fatalf("expected String, got %T", result)
		}

		if strResult.Value != tt.expected {
			t.Errorf("LEFT(%q, %d) = %q, want %q", tt.input, tt.n, strResult.Value, tt.expected)
		}
	}
}

func TestLeftNegative(t *testing.T) {
	builtins := GetBuiltins()
	leftFn := builtins["LEFT"]

	result := leftFn.Fn(&interpreter.String{Value: "Hello"}, &interpreter.Integer{Value: -1})

	if _, ok := result.(*interpreter.Error); !ok {
		t.Errorf("expected Error for negative length, got %T", result)
	}
}

func TestRight(t *testing.T) {
	tests := []struct {
		input    string
		n        int64
		expected string
	}{
		{"Hello", 3, "llo"},
		{"Hello", 0, ""},
		{"Hello", 5, "Hello"},
		{"Hello", 10, "Hello"},
	}

	builtins := GetBuiltins()
	rightFn := builtins["RIGHT"]

	for _, tt := range tests {
		result := rightFn.Fn(&interpreter.String{Value: tt.input}, &interpreter.Integer{Value: tt.n})

		strResult, ok := result.(*interpreter.String)
		if !ok {
			t.Fatalf("expected String, got %T", result)
		}

		if strResult.Value != tt.expected {
			t.Errorf("RIGHT(%q, %d) = %q, want %q", tt.input, tt.n, strResult.Value, tt.expected)
		}
	}
}

func TestRightNegative(t *testing.T) {
	builtins := GetBuiltins()
	rightFn := builtins["RIGHT"]

	result := rightFn.Fn(&interpreter.String{Value: "Hello"}, &interpreter.Integer{Value: -1})

	if _, ok := result.(*interpreter.Error); !ok {
		t.Errorf("expected Error for negative length, got %T", result)
	}
}

func TestMid(t *testing.T) {
	tests := []struct {
		input    string
		start    int64
		length   int64
		expected string
	}{
		{"Hello World", 1, 5, "Hello"},
		{"Hello World", 7, 5, "World"},
		{"Hello World", 1, 11, "Hello World"},
		{"Hello", 1, 10, "Hello"}, // Length exceeds string
		{"Hello", 3, 2, "ll"},     // Middle portion
		{"Hello", 10, 2, ""},      // Start beyond string
	}

	builtins := GetBuiltins()
	midFn := builtins["MID"]

	for _, tt := range tests {
		result := midFn.Fn(
			&interpreter.String{Value: tt.input},
			&interpreter.Integer{Value: tt.start},
			&interpreter.Integer{Value: tt.length},
		)

		strResult, ok := result.(*interpreter.String)
		if !ok {
			t.Fatalf("expected String, got %T", result)
		}

		if strResult.Value != tt.expected {
			t.Errorf("MID(%q, %d, %d) = %q, want %q",
				tt.input, tt.start, tt.length, strResult.Value, tt.expected)
		}
	}
}

func TestLcase(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected interface{}
	}{
		{&interpreter.Char{Value: 'A'}, 'a'},
		{&interpreter.Char{Value: 'z'}, 'z'},
		{&interpreter.Char{Value: '5'}, '5'},
		{&interpreter.String{Value: "HELLO"}, "hello"},
		{&interpreter.String{Value: "Hello World"}, "hello world"},
	}

	builtins := GetBuiltins()
	lcaseFn := builtins["LCASE"]

	for _, tt := range tests {
		result := lcaseFn.Fn(tt.input.(interpreter.Object))

		switch expected := tt.expected.(type) {
		case rune:
			charResult, ok := result.(*interpreter.Char)
			if !ok {
				t.Fatalf("expected Char, got %T", result)
			}
			if charResult.Value != expected {
				t.Errorf("LCASE(%v) = %c, want %c", tt.input, charResult.Value, expected)
			}
		case string:
			strResult, ok := result.(*interpreter.String)
			if !ok {
				t.Fatalf("expected String, got %T", result)
			}
			if strResult.Value != expected {
				t.Errorf("LCASE(%v) = %s, want %s", tt.input, strResult.Value, expected)
			}
		}
	}
}

func TestUcase(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected interface{}
	}{
		{&interpreter.Char{Value: 'a'}, 'A'},
		{&interpreter.Char{Value: 'Z'}, 'Z'},
		{&interpreter.Char{Value: '5'}, '5'},
		{&interpreter.String{Value: "hello"}, "HELLO"},
		{&interpreter.String{Value: "Hello World"}, "HELLO WORLD"},
	}

	builtins := GetBuiltins()
	ucaseFn := builtins["UCASE"]

	for _, tt := range tests {
		result := ucaseFn.Fn(tt.input.(interpreter.Object))

		switch expected := tt.expected.(type) {
		case rune:
			charResult, ok := result.(*interpreter.Char)
			if !ok {
				t.Fatalf("expected Char, got %T", result)
			}
			if charResult.Value != expected {
				t.Errorf("UCASE(%v) = %c, want %c", tt.input, charResult.Value, expected)
			}
		case string:
			strResult, ok := result.(*interpreter.String)
			if !ok {
				t.Fatalf("expected String, got %T", result)
			}
			if strResult.Value != expected {
				t.Errorf("UCASE(%v) = %s, want %s", tt.input, strResult.Value, expected)
			}
		}
	}
}

func TestToUpper(t *testing.T) {
	builtins := GetBuiltins()
	toUpperFn := builtins["TO_UPPER"]

	result := toUpperFn.Fn(&interpreter.String{Value: "hello world"})

	strResult, ok := result.(*interpreter.String)
	if !ok {
		t.Fatalf("expected String, got %T", result)
	}

	if strResult.Value != "HELLO WORLD" {
		t.Errorf("TO_UPPER(\"hello world\") = %q, want \"HELLO WORLD\"", strResult.Value)
	}
}

func TestToLower(t *testing.T) {
	builtins := GetBuiltins()
	toLowerFn := builtins["TO_LOWER"]

	result := toLowerFn.Fn(&interpreter.String{Value: "HELLO WORLD"})

	strResult, ok := result.(*interpreter.String)
	if !ok {
		t.Fatalf("expected String, got %T", result)
	}

	if strResult.Value != "hello world" {
		t.Errorf("TO_LOWER(\"HELLO WORLD\") = %q, want \"hello world\"", strResult.Value)
	}
}

func TestAsc(t *testing.T) {
	tests := []struct {
		input    interpreter.Object
		expected int64
	}{
		{&interpreter.Char{Value: 'A'}, 65},
		{&interpreter.Char{Value: 'a'}, 97},
		{&interpreter.Char{Value: '0'}, 48},
		{&interpreter.String{Value: "A"}, 65},
		{&interpreter.String{Value: "Hello"}, 72}, // First character
	}

	builtins := GetBuiltins()
	ascFn := builtins["ASC"]

	for _, tt := range tests {
		result := ascFn.Fn(tt.input)

		intResult, ok := result.(*interpreter.Integer)
		if !ok {
			t.Fatalf("expected Integer, got %T", result)
		}

		if intResult.Value != tt.expected {
			t.Errorf("ASC(%v) = %d, want %d", tt.input.Inspect(), intResult.Value, tt.expected)
		}
	}
}

func TestAscEmptyString(t *testing.T) {
	builtins := GetBuiltins()
	ascFn := builtins["ASC"]

	result := ascFn.Fn(&interpreter.String{Value: ""})

	if _, ok := result.(*interpreter.Error); !ok {
		t.Errorf("expected Error for empty string, got %T", result)
	}
}

func TestChr(t *testing.T) {
	tests := []struct {
		input    int64
		expected rune
	}{
		{65, 'A'},
		{97, 'a'},
		{48, '0'},
		{32, ' '},
	}

	builtins := GetBuiltins()
	chrFn := builtins["CHR"]

	for _, tt := range tests {
		result := chrFn.Fn(&interpreter.Integer{Value: tt.input})

		charResult, ok := result.(*interpreter.Char)
		if !ok {
			t.Fatalf("expected Char, got %T", result)
		}

		if charResult.Value != tt.expected {
			t.Errorf("CHR(%d) = %c, want %c", tt.input, charResult.Value, tt.expected)
		}
	}
}

func TestInt(t *testing.T) {
	tests := []struct {
		input    interpreter.Object
		expected int64
	}{
		{&interpreter.Real{Value: 3.7}, 3},
		{&interpreter.Real{Value: 3.2}, 3},
		{&interpreter.Real{Value: -3.7}, -3},
		{&interpreter.Integer{Value: 5}, 5},
	}

	builtins := GetBuiltins()
	intFn := builtins["INT"]

	for _, tt := range tests {
		result := intFn.Fn(tt.input)

		intResult, ok := result.(*interpreter.Integer)
		if !ok {
			t.Fatalf("expected Integer, got %T", result)
		}

		if intResult.Value != tt.expected {
			t.Errorf("INT(%v) = %d, want %d", tt.input.Inspect(), intResult.Value, tt.expected)
		}
	}
}

func TestRand(t *testing.T) {
	builtins := GetBuiltins()
	randFn := builtins["RAND"]

	// Test that result is within range [0, n)
	for i := 0; i < 100; i++ {
		result := randFn.Fn(&interpreter.Integer{Value: 10})

		realResult, ok := result.(*interpreter.Real)
		if !ok {
			t.Fatalf("expected Real, got %T", result)
		}

		if realResult.Value < 0 || realResult.Value >= 10 {
			t.Errorf("RAND(10) = %f, expected value in [0, 10)", realResult.Value)
		}
	}
}

func TestRandom(t *testing.T) {
	builtins := GetBuiltins()
	randomFn := builtins["RANDOM"]

	// Test that result is within range [0, 1]
	for i := 0; i < 100; i++ {
		result := randomFn.Fn()

		realResult, ok := result.(*interpreter.Real)
		if !ok {
			t.Fatalf("expected Real, got %T", result)
		}

		if realResult.Value < 0 || realResult.Value > 1 {
			t.Errorf("RANDOM() = %f, expected value in [0, 1]", realResult.Value)
		}
	}
}

func TestRound(t *testing.T) {
	tests := []struct {
		value    float64
		places   int64
		expected float64
	}{
		{3.14159, 2, 3.14},
		{3.14159, 0, 3.0},
		{3.5, 0, 4.0},
		{2.5, 0, 3.0}, // Go's math.Round rounds half away from zero
		{-2.5, 0, -3.0},
		{123.456, 1, 123.5},
	}

	builtins := GetBuiltins()
	roundFn := builtins["ROUND"]

	for _, tt := range tests {
		result := roundFn.Fn(
			&interpreter.Real{Value: tt.value},
			&interpreter.Integer{Value: tt.places},
		)

		realResult, ok := result.(*interpreter.Real)
		if !ok {
			t.Fatalf("expected Real, got %T", result)
		}

		if realResult.Value != tt.expected {
			t.Errorf("ROUND(%f, %d) = %f, want %f",
				tt.value, tt.places, realResult.Value, tt.expected)
		}
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		input    interpreter.Object
		expected interface{}
	}{
		{&interpreter.Integer{Value: 5}, int64(5)},
		{&interpreter.Integer{Value: -5}, int64(5)},
		{&interpreter.Integer{Value: 0}, int64(0)},
		{&interpreter.Real{Value: 3.14}, 3.14},
		{&interpreter.Real{Value: -3.14}, 3.14},
	}

	builtins := GetBuiltins()
	absFn := builtins["ABS"]

	for _, tt := range tests {
		result := absFn.Fn(tt.input)

		switch expected := tt.expected.(type) {
		case int64:
			intResult, ok := result.(*interpreter.Integer)
			if !ok {
				t.Fatalf("expected Integer, got %T", result)
			}
			if intResult.Value != expected {
				t.Errorf("ABS(%v) = %d, want %d", tt.input.Inspect(), intResult.Value, expected)
			}
		case float64:
			realResult, ok := result.(*interpreter.Real)
			if !ok {
				t.Fatalf("expected Real, got %T", result)
			}
			if realResult.Value != expected {
				t.Errorf("ABS(%v) = %f, want %f", tt.input.Inspect(), realResult.Value, expected)
			}
		}
	}
}

func TestSqrt(t *testing.T) {
	tests := []struct {
		input    interpreter.Object
		expected float64
	}{
		{&interpreter.Integer{Value: 4}, 2.0},
		{&interpreter.Integer{Value: 9}, 3.0},
		{&interpreter.Real{Value: 2.0}, 1.4142135623730951},
		{&interpreter.Integer{Value: 0}, 0.0},
	}

	builtins := GetBuiltins()
	sqrtFn := builtins["SQRT"]

	for _, tt := range tests {
		result := sqrtFn.Fn(tt.input)

		realResult, ok := result.(*interpreter.Real)
		if !ok {
			t.Fatalf("expected Real, got %T", result)
		}

		if realResult.Value != tt.expected {
			t.Errorf("SQRT(%v) = %f, want %f", tt.input.Inspect(), realResult.Value, tt.expected)
		}
	}
}

func TestSqrtNegative(t *testing.T) {
	builtins := GetBuiltins()
	sqrtFn := builtins["SQRT"]

	result := sqrtFn.Fn(&interpreter.Integer{Value: -1})

	if _, ok := result.(*interpreter.Error); !ok {
		t.Errorf("expected Error for negative sqrt, got %T", result)
	}
}

func TestPow(t *testing.T) {
	tests := []struct {
		base     interpreter.Object
		exp      interpreter.Object
		expected float64
	}{
		{&interpreter.Integer{Value: 2}, &interpreter.Integer{Value: 3}, 8.0},
		{&interpreter.Integer{Value: 10}, &interpreter.Integer{Value: 2}, 100.0},
		{&interpreter.Real{Value: 2.5}, &interpreter.Integer{Value: 2}, 6.25},
		{&interpreter.Integer{Value: 2}, &interpreter.Real{Value: 0.5}, 1.4142135623730951},
	}

	builtins := GetBuiltins()
	powFn := builtins["POW"]

	for _, tt := range tests {
		result := powFn.Fn(tt.base, tt.exp)

		realResult, ok := result.(*interpreter.Real)
		if !ok {
			t.Fatalf("expected Real, got %T", result)
		}

		if realResult.Value != tt.expected {
			t.Errorf("POW(%v, %v) = %f, want %f",
				tt.base.Inspect(), tt.exp.Inspect(), realResult.Value, tt.expected)
		}
	}
}

func TestNumToStr(t *testing.T) {
	tests := []struct {
		input    interpreter.Object
		expected string
	}{
		{&interpreter.Integer{Value: 42}, "42"},
		{&interpreter.Integer{Value: -10}, "-10"},
		{&interpreter.Real{Value: 3.14}, "3.14"},
		{&interpreter.Real{Value: 10.0}, "10"},
	}

	builtins := GetBuiltins()
	numToStrFn := builtins["NUM_TO_STR"]

	for _, tt := range tests {
		result := numToStrFn.Fn(tt.input)

		strResult, ok := result.(*interpreter.String)
		if !ok {
			t.Fatalf("expected String, got %T", result)
		}

		if strResult.Value != tt.expected {
			t.Errorf("NUM_TO_STR(%v) = %q, want %q",
				tt.input.Inspect(), strResult.Value, tt.expected)
		}
	}
}

func TestStrToNum(t *testing.T) {
	tests := []struct {
		input       string
		expected    interface{}
		expectedTyp string
	}{
		{"42", int64(42), "INTEGER"},
		{"-10", int64(-10), "INTEGER"},
		{"3.14", 3.14, "REAL"},
		{"0", int64(0), "INTEGER"},
	}

	builtins := GetBuiltins()
	strToNumFn := builtins["STR_TO_NUM"]

	for _, tt := range tests {
		result := strToNumFn.Fn(&interpreter.String{Value: tt.input})

		switch expected := tt.expected.(type) {
		case int64:
			intResult, ok := result.(*interpreter.Integer)
			if !ok {
				t.Fatalf("expected Integer, got %T", result)
			}
			if intResult.Value != expected {
				t.Errorf("STR_TO_NUM(%q) = %d, want %d", tt.input, intResult.Value, expected)
			}
		case float64:
			realResult, ok := result.(*interpreter.Real)
			if !ok {
				t.Fatalf("expected Real, got %T", result)
			}
			if realResult.Value != expected {
				t.Errorf("STR_TO_NUM(%q) = %f, want %f", tt.input, realResult.Value, expected)
			}
		}
	}
}

func TestStrToNumInvalid(t *testing.T) {
	builtins := GetBuiltins()
	strToNumFn := builtins["STR_TO_NUM"]

	result := strToNumFn.Fn(&interpreter.String{Value: "abc"})

	if _, ok := result.(*interpreter.Error); !ok {
		t.Errorf("expected Error for invalid number string, got %T", result)
	}
}

func TestEOF(t *testing.T) {
	builtins := GetBuiltins()
	eofFn := builtins["EOF"]

	// EOF returns true by default (placeholder implementation)
	result := eofFn.Fn(&interpreter.String{Value: "test.txt"})

	boolResult, ok := result.(*interpreter.Boolean)
	if !ok {
		t.Fatalf("expected Boolean, got %T", result)
	}

	if !boolResult.Value {
		t.Error("EOF should return true by default")
	}
}

func TestGetBuiltins(t *testing.T) {
	builtins := GetBuiltins()

	expectedFunctions := []string{
		"LENGTH", "LEFT", "RIGHT", "MID",
		"LCASE", "UCASE", "TO_UPPER", "TO_LOWER",
		"ASC", "CHR",
		"INT", "RAND", "RANDOM", "ROUND",
		"NUM_TO_STR", "STR_TO_NUM",
		"EOF",
		"ABS", "SQRT", "POW",
	}

	for _, name := range expectedFunctions {
		if _, ok := builtins[name]; !ok {
			t.Errorf("builtin %s not found", name)
		}
	}
}

func TestBuiltinNames(t *testing.T) {
	builtins := GetBuiltins()

	for name, builtin := range builtins {
		if builtin.Name != name {
			t.Errorf("builtin name mismatch: map key=%q, Name=%q", name, builtin.Name)
		}
	}
}

// Date function tests

func TestDay(t *testing.T) {
	tests := []struct {
		day      int
		month    int
		year     int
		expected int64
	}{
		{4, 10, 2003, 4},
		{1, 1, 2000, 1},
		{31, 12, 2023, 31},
		{15, 6, 1990, 15},
	}

	builtins := GetBuiltins()
	dayFn := builtins["DAY"]

	for _, tt := range tests {
		result := dayFn.Fn(&interpreter.Date{Day: tt.day, Month: tt.month, Year: tt.year})

		intResult, ok := result.(*interpreter.Integer)
		if !ok {
			t.Fatalf("expected Integer, got %T", result)
		}

		if intResult.Value != tt.expected {
			t.Errorf("DAY(%02d/%02d/%04d) = %d, want %d",
				tt.day, tt.month, tt.year, intResult.Value, tt.expected)
		}
	}
}

func TestDayWrongArgType(t *testing.T) {
	builtins := GetBuiltins()
	dayFn := builtins["DAY"]

	result := dayFn.Fn(&interpreter.String{Value: "04/10/2003"})

	if _, ok := result.(*interpreter.Error); !ok {
		t.Errorf("expected Error for wrong arg type, got %T", result)
	}
}

func TestMonth(t *testing.T) {
	tests := []struct {
		day      int
		month    int
		year     int
		expected int64
	}{
		{4, 10, 2003, 10},
		{1, 1, 2000, 1},
		{31, 12, 2023, 12},
		{15, 6, 1990, 6},
	}

	builtins := GetBuiltins()
	monthFn := builtins["MONTH"]

	for _, tt := range tests {
		result := monthFn.Fn(&interpreter.Date{Day: tt.day, Month: tt.month, Year: tt.year})

		intResult, ok := result.(*interpreter.Integer)
		if !ok {
			t.Fatalf("expected Integer, got %T", result)
		}

		if intResult.Value != tt.expected {
			t.Errorf("MONTH(%02d/%02d/%04d) = %d, want %d",
				tt.day, tt.month, tt.year, intResult.Value, tt.expected)
		}
	}
}

func TestMonthWrongArgType(t *testing.T) {
	builtins := GetBuiltins()
	monthFn := builtins["MONTH"]

	result := monthFn.Fn(&interpreter.Integer{Value: 10})

	if _, ok := result.(*interpreter.Error); !ok {
		t.Errorf("expected Error for wrong arg type, got %T", result)
	}
}

func TestYear(t *testing.T) {
	tests := []struct {
		day      int
		month    int
		year     int
		expected int64
	}{
		{4, 10, 2003, 2003},
		{1, 1, 2000, 2000},
		{31, 12, 2023, 2023},
		{15, 6, 1990, 1990},
	}

	builtins := GetBuiltins()
	yearFn := builtins["YEAR"]

	for _, tt := range tests {
		result := yearFn.Fn(&interpreter.Date{Day: tt.day, Month: tt.month, Year: tt.year})

		intResult, ok := result.(*interpreter.Integer)
		if !ok {
			t.Fatalf("expected Integer, got %T", result)
		}

		if intResult.Value != tt.expected {
			t.Errorf("YEAR(%02d/%02d/%04d) = %d, want %d",
				tt.day, tt.month, tt.year, intResult.Value, tt.expected)
		}
	}
}

func TestYearWrongArgType(t *testing.T) {
	builtins := GetBuiltins()
	yearFn := builtins["YEAR"]

	result := yearFn.Fn(&interpreter.Real{Value: 2003.0})

	if _, ok := result.(*interpreter.Error); !ok {
		t.Errorf("expected Error for wrong arg type, got %T", result)
	}
}

func TestDayIndex(t *testing.T) {
	// Day index: Sunday = 1, Monday = 2, ..., Saturday = 7
	tests := []struct {
		day      int
		month    int
		year     int
		expected int64
		desc     string
	}{
		{9, 5, 2023, 3, "Tuesday"},   // 09/05/2023 is a Tuesday
		{7, 5, 2023, 1, "Sunday"},    // 07/05/2023 is a Sunday
		{8, 5, 2023, 2, "Monday"},    // 08/05/2023 is a Monday
		{13, 5, 2023, 7, "Saturday"}, // 13/05/2023 is a Saturday
		{1, 1, 2000, 7, "Saturday"},  // 01/01/2000 is a Saturday
		{25, 12, 2023, 2, "Monday"},  // Christmas 2023 is a Monday
	}

	builtins := GetBuiltins()
	dayIndexFn := builtins["DAYINDEX"]

	for _, tt := range tests {
		result := dayIndexFn.Fn(&interpreter.Date{Day: tt.day, Month: tt.month, Year: tt.year})

		intResult, ok := result.(*interpreter.Integer)
		if !ok {
			t.Fatalf("expected Integer, got %T", result)
		}

		if intResult.Value != tt.expected {
			t.Errorf("DAYINDEX(%02d/%02d/%04d) = %d, want %d (%s)",
				tt.day, tt.month, tt.year, intResult.Value, tt.expected, tt.desc)
		}
	}
}

func TestDayIndexWrongArgType(t *testing.T) {
	builtins := GetBuiltins()
	dayIndexFn := builtins["DAYINDEX"]

	result := dayIndexFn.Fn(&interpreter.String{Value: "09/05/2023"})

	if _, ok := result.(*interpreter.Error); !ok {
		t.Errorf("expected Error for wrong arg type, got %T", result)
	}
}

func TestSetDate(t *testing.T) {
	tests := []struct {
		day   int64
		month int64
		year  int64
	}{
		{26, 10, 2003},
		{1, 1, 2000},
		{31, 12, 2023},
		{15, 6, 1990},
	}

	builtins := GetBuiltins()
	setDateFn := builtins["SETDATE"]

	for _, tt := range tests {
		result := setDateFn.Fn(
			&interpreter.Integer{Value: tt.day},
			&interpreter.Integer{Value: tt.month},
			&interpreter.Integer{Value: tt.year},
		)

		dateResult, ok := result.(*interpreter.Date)
		if !ok {
			t.Fatalf("expected Date, got %T", result)
		}

		if dateResult.Day != int(tt.day) || dateResult.Month != int(tt.month) || dateResult.Year != int(tt.year) {
			t.Errorf("SETDATE(%d, %d, %d) = %02d/%02d/%04d, want %02d/%02d/%04d",
				tt.day, tt.month, tt.year,
				dateResult.Day, dateResult.Month, dateResult.Year,
				tt.day, tt.month, tt.year)
		}
	}
}

func TestSetDateWrongArgCount(t *testing.T) {
	builtins := GetBuiltins()
	setDateFn := builtins["SETDATE"]

	result := setDateFn.Fn(
		&interpreter.Integer{Value: 26},
		&interpreter.Integer{Value: 10},
	)

	if _, ok := result.(*interpreter.Error); !ok {
		t.Errorf("expected Error for wrong arg count, got %T", result)
	}
}

func TestSetDateWrongArgType(t *testing.T) {
	builtins := GetBuiltins()
	setDateFn := builtins["SETDATE"]

	// Wrong type for day
	result := setDateFn.Fn(
		&interpreter.String{Value: "26"},
		&interpreter.Integer{Value: 10},
		&interpreter.Integer{Value: 2003},
	)

	if _, ok := result.(*interpreter.Error); !ok {
		t.Errorf("expected Error for wrong arg type, got %T", result)
	}

	// Wrong type for month
	result = setDateFn.Fn(
		&interpreter.Integer{Value: 26},
		&interpreter.String{Value: "10"},
		&interpreter.Integer{Value: 2003},
	)

	if _, ok := result.(*interpreter.Error); !ok {
		t.Errorf("expected Error for wrong arg type, got %T", result)
	}

	// Wrong type for year
	result = setDateFn.Fn(
		&interpreter.Integer{Value: 26},
		&interpreter.Integer{Value: 10},
		&interpreter.Real{Value: 2003.0},
	)

	if _, ok := result.(*interpreter.Error); !ok {
		t.Errorf("expected Error for wrong arg type, got %T", result)
	}
}

func TestToday(t *testing.T) {
	builtins := GetBuiltins()
	todayFn := builtins["TODAY"]

	result := todayFn.Fn()

	dateResult, ok := result.(*interpreter.Date)
	if !ok {
		t.Fatalf("expected Date, got %T", result)
	}

	// Just verify that it returns a reasonable date (day 1-31, month 1-12, year > 2000)
	if dateResult.Day < 1 || dateResult.Day > 31 {
		t.Errorf("TODAY() returned invalid day: %d", dateResult.Day)
	}
	if dateResult.Month < 1 || dateResult.Month > 12 {
		t.Errorf("TODAY() returned invalid month: %d", dateResult.Month)
	}
	if dateResult.Year < 2000 {
		t.Errorf("TODAY() returned invalid year: %d", dateResult.Year)
	}
}

func TestTodayWrongArgCount(t *testing.T) {
	builtins := GetBuiltins()
	todayFn := builtins["TODAY"]

	result := todayFn.Fn(&interpreter.Integer{Value: 1})

	if _, ok := result.(*interpreter.Error); !ok {
		t.Errorf("expected Error for wrong arg count, got %T", result)
	}
}

func TestDateBuiltinsRegistered(t *testing.T) {
	builtins := GetBuiltins()

	dateFunctions := []string{"DAY", "MONTH", "YEAR", "DAYINDEX", "SETDATE", "TODAY"}

	for _, name := range dateFunctions {
		if _, ok := builtins[name]; !ok {
			t.Errorf("date builtin %s not found", name)
		}
	}
}
