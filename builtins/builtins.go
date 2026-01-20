// Package builtins provides built-in functions for Cambridge Pseudocode
// Based on Cambridge International AS & A Level Computer Science 9618 Insert
package builtins

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/andrinoff/cambridge-lang/interpreter"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetBuiltins returns all built-in functions
func GetBuiltins() map[string]*interpreter.Builtin {
	return map[string]*interpreter.Builtin{
		// String functions
		"LENGTH":   {Name: "LENGTH", Fn: length},
		"LEFT":     {Name: "LEFT", Fn: left},
		"RIGHT":    {Name: "RIGHT", Fn: right},
		"MID":      {Name: "MID", Fn: mid},
		"LCASE":    {Name: "LCASE", Fn: lcase},
		"UCASE":    {Name: "UCASE", Fn: ucase},
		"TO_UPPER": {Name: "TO_UPPER", Fn: toUpper},
		"TO_LOWER": {Name: "TO_LOWER", Fn: toLower},

		// Character/ASCII functions
		"ASC": {Name: "ASC", Fn: asc},
		"CHR": {Name: "CHR", Fn: chr},

		// Numeric functions
		"INT":    {Name: "INT", Fn: intFunc},
		"RAND":   {Name: "RAND", Fn: randFunc},
		"RANDOM": {Name: "RANDOM", Fn: random},
		"ROUND":  {Name: "ROUND", Fn: round},

		// Conversion functions
		"NUM_TO_STR": {Name: "NUM_TO_STR", Fn: numToStr},
		"STR_TO_NUM": {Name: "STR_TO_NUM", Fn: strToNum},

		// File function
		"EOF": {Name: "EOF", Fn: eof},

		// Math functions (additional)
		"ABS":  {Name: "ABS", Fn: abs},
		"SQRT": {Name: "SQRT", Fn: sqrt},
		"POW":  {Name: "POW", Fn: pow},
	}
}

// LENGTH(s) - returns the length of a string
func length(args ...interpreter.Object) interpreter.Object {
	if len(args) != 1 {
		return newError("LENGTH requires 1 argument, got %d", len(args))
	}

	switch arg := args[0].(type) {
	case *interpreter.String:
		return &interpreter.Integer{Value: int64(len(arg.Value))}
	default:
		return newError("LENGTH requires STRING argument, got %s", args[0].Type())
	}
}

// LEFT(s, n) - returns leftmost n characters
func left(args ...interpreter.Object) interpreter.Object {
	if len(args) != 2 {
		return newError("LEFT requires 2 arguments, got %d", len(args))
	}

	str, ok := args[0].(*interpreter.String)
	if !ok {
		return newError("LEFT requires STRING as first argument")
	}

	n, ok := args[1].(*interpreter.Integer)
	if !ok {
		return newError("LEFT requires INTEGER as second argument")
	}

	if n.Value < 0 {
		return newError("LEFT: length cannot be negative")
	}

	if int(n.Value) >= len(str.Value) {
		return &interpreter.String{Value: str.Value}
	}

	return &interpreter.String{Value: str.Value[:n.Value]}
}

// RIGHT(s, n) - returns rightmost n characters
func right(args ...interpreter.Object) interpreter.Object {
	if len(args) != 2 {
		return newError("RIGHT requires 2 arguments, got %d", len(args))
	}

	str, ok := args[0].(*interpreter.String)
	if !ok {
		return newError("RIGHT requires STRING as first argument")
	}

	n, ok := args[1].(*interpreter.Integer)
	if !ok {
		return newError("RIGHT requires INTEGER as second argument")
	}

	if n.Value < 0 {
		return newError("RIGHT: length cannot be negative")
	}

	length := len(str.Value)
	if int(n.Value) >= length {
		return &interpreter.String{Value: str.Value}
	}

	return &interpreter.String{Value: str.Value[length-int(n.Value):]}
}

// MID(s, start, length) - returns substring starting at position start with given length
// Note: Cambridge pseudocode uses 1-based indexing
func mid(args ...interpreter.Object) interpreter.Object {
	if len(args) != 3 {
		return newError("MID requires 3 arguments, got %d", len(args))
	}

	str, ok := args[0].(*interpreter.String)
	if !ok {
		return newError("MID requires STRING as first argument")
	}

	start, ok := args[1].(*interpreter.Integer)
	if !ok {
		return newError("MID requires INTEGER as second argument")
	}

	length, ok := args[2].(*interpreter.Integer)
	if !ok {
		return newError("MID requires INTEGER as third argument")
	}

	// Convert to 0-based indexing
	startIdx := int(start.Value) - 1
	if startIdx < 0 {
		startIdx = 0
	}

	strLen := len(str.Value)
	if startIdx >= strLen {
		return &interpreter.String{Value: ""}
	}

	endIdx := startIdx + int(length.Value)
	if endIdx > strLen {
		endIdx = strLen
	}

	return &interpreter.String{Value: str.Value[startIdx:endIdx]}
}

// LCASE(c) - converts character to lowercase
func lcase(args ...interpreter.Object) interpreter.Object {
	if len(args) != 1 {
		return newError("LCASE requires 1 argument, got %d", len(args))
	}

	switch arg := args[0].(type) {
	case *interpreter.Char:
		return &interpreter.Char{Value: unicode.ToLower(arg.Value)}
	case *interpreter.String:
		return &interpreter.String{Value: strings.ToLower(arg.Value)}
	default:
		return newError("LCASE requires CHAR or STRING argument")
	}
}

// UCASE(c) - converts character to uppercase
func ucase(args ...interpreter.Object) interpreter.Object {
	if len(args) != 1 {
		return newError("UCASE requires 1 argument, got %d", len(args))
	}

	switch arg := args[0].(type) {
	case *interpreter.Char:
		return &interpreter.Char{Value: unicode.ToUpper(arg.Value)}
	case *interpreter.String:
		return &interpreter.String{Value: strings.ToUpper(arg.Value)}
	default:
		return newError("UCASE requires CHAR or STRING argument")
	}
}

// TO_UPPER(s) - converts string to uppercase
func toUpper(args ...interpreter.Object) interpreter.Object {
	if len(args) != 1 {
		return newError("TO_UPPER requires 1 argument, got %d", len(args))
	}

	str, ok := args[0].(*interpreter.String)
	if !ok {
		return newError("TO_UPPER requires STRING argument")
	}

	return &interpreter.String{Value: strings.ToUpper(str.Value)}
}

// TO_LOWER(s) - converts string to lowercase
func toLower(args ...interpreter.Object) interpreter.Object {
	if len(args) != 1 {
		return newError("TO_LOWER requires 1 argument, got %d", len(args))
	}

	str, ok := args[0].(*interpreter.String)
	if !ok {
		return newError("TO_LOWER requires STRING argument")
	}

	return &interpreter.String{Value: strings.ToLower(str.Value)}
}

// ASC(c) - returns ASCII value of character
func asc(args ...interpreter.Object) interpreter.Object {
	if len(args) != 1 {
		return newError("ASC requires 1 argument, got %d", len(args))
	}

	switch arg := args[0].(type) {
	case *interpreter.Char:
		return &interpreter.Integer{Value: int64(arg.Value)}
	case *interpreter.String:
		if len(arg.Value) == 0 {
			return newError("ASC: empty string")
		}
		return &interpreter.Integer{Value: int64(arg.Value[0])}
	default:
		return newError("ASC requires CHAR or STRING argument")
	}
}

// CHR(n) - returns character with given ASCII value
func chr(args ...interpreter.Object) interpreter.Object {
	if len(args) != 1 {
		return newError("CHR requires 1 argument, got %d", len(args))
	}

	n, ok := args[0].(*interpreter.Integer)
	if !ok {
		return newError("CHR requires INTEGER argument")
	}

	return &interpreter.Char{Value: rune(n.Value)}
}

// INT(x) - returns integer part of a real number
func intFunc(args ...interpreter.Object) interpreter.Object {
	if len(args) != 1 {
		return newError("INT requires 1 argument, got %d", len(args))
	}

	switch arg := args[0].(type) {
	case *interpreter.Real:
		return &interpreter.Integer{Value: int64(arg.Value)}
	case *interpreter.Integer:
		return arg
	default:
		return newError("INT requires REAL or INTEGER argument")
	}
}

// RAND(n) - returns random real number from 0 to n (exclusive)
func randFunc(args ...interpreter.Object) interpreter.Object {
	if len(args) != 1 {
		return newError("RAND requires 1 argument, got %d", len(args))
	}

	var max float64
	switch arg := args[0].(type) {
	case *interpreter.Integer:
		max = float64(arg.Value)
	case *interpreter.Real:
		max = arg.Value
	default:
		return newError("RAND requires numeric argument")
	}

	return &interpreter.Real{Value: rand.Float64() * max}
}

// RANDOM() - returns random real number from 0 to 1 (inclusive)
func random(args ...interpreter.Object) interpreter.Object {
	if len(args) != 0 {
		return newError("RANDOM requires 0 arguments, got %d", len(args))
	}
	return &interpreter.Real{Value: rand.Float64()}
}

// ROUND(x, places) - rounds to specified decimal places
func round(args ...interpreter.Object) interpreter.Object {
	if len(args) != 2 {
		return newError("ROUND requires 2 arguments, got %d", len(args))
	}

	var value float64
	switch arg := args[0].(type) {
	case *interpreter.Real:
		value = arg.Value
	case *interpreter.Integer:
		value = float64(arg.Value)
	default:
		return newError("ROUND requires numeric first argument")
	}

	places, ok := args[1].(*interpreter.Integer)
	if !ok {
		return newError("ROUND requires INTEGER as second argument")
	}

	multiplier := math.Pow(10, float64(places.Value))
	rounded := math.Round(value*multiplier) / multiplier

	return &interpreter.Real{Value: rounded}
}

// NUM_TO_STR(n) - converts number to string
func numToStr(args ...interpreter.Object) interpreter.Object {
	if len(args) != 1 {
		return newError("NUM_TO_STR requires 1 argument, got %d", len(args))
	}

	switch arg := args[0].(type) {
	case *interpreter.Integer:
		return &interpreter.String{Value: strconv.FormatInt(arg.Value, 10)}
	case *interpreter.Real:
		return &interpreter.String{Value: strconv.FormatFloat(arg.Value, 'f', -1, 64)}
	default:
		return newError("NUM_TO_STR requires numeric argument")
	}
}

// STR_TO_NUM(s) - converts string to number
func strToNum(args ...interpreter.Object) interpreter.Object {
	if len(args) != 1 {
		return newError("STR_TO_NUM requires 1 argument, got %d", len(args))
	}

	str, ok := args[0].(*interpreter.String)
	if !ok {
		return newError("STR_TO_NUM requires STRING argument")
	}

	// Try to parse as integer first
	if i, err := strconv.ParseInt(str.Value, 10, 64); err == nil {
		return &interpreter.Integer{Value: i}
	}

	// Try to parse as float
	if f, err := strconv.ParseFloat(str.Value, 64); err == nil {
		return &interpreter.Real{Value: f}
	}

	return newError("STR_TO_NUM: cannot convert '%s' to number", str.Value)
}

// EOF(filename) - checks if at end of file
// This is a placeholder - actual implementation depends on file handling
func eof(args ...interpreter.Object) interpreter.Object {
	if len(args) != 1 {
		return newError("EOF requires 1 argument, got %d", len(args))
	}
	// This will be handled by the interpreter with access to file state
	return &interpreter.Boolean{Value: true}
}

// ABS(n) - returns absolute value
func abs(args ...interpreter.Object) interpreter.Object {
	if len(args) != 1 {
		return newError("ABS requires 1 argument, got %d", len(args))
	}

	switch arg := args[0].(type) {
	case *interpreter.Integer:
		if arg.Value < 0 {
			return &interpreter.Integer{Value: -arg.Value}
		}
		return arg
	case *interpreter.Real:
		return &interpreter.Real{Value: math.Abs(arg.Value)}
	default:
		return newError("ABS requires numeric argument")
	}
}

// SQRT(n) - returns square root
func sqrt(args ...interpreter.Object) interpreter.Object {
	if len(args) != 1 {
		return newError("SQRT requires 1 argument, got %d", len(args))
	}

	var value float64
	switch arg := args[0].(type) {
	case *interpreter.Integer:
		value = float64(arg.Value)
	case *interpreter.Real:
		value = arg.Value
	default:
		return newError("SQRT requires numeric argument")
	}

	if value < 0 {
		return newError("SQRT: cannot take square root of negative number")
	}

	return &interpreter.Real{Value: math.Sqrt(value)}
}

// POW(base, exp) - returns base raised to power exp
func pow(args ...interpreter.Object) interpreter.Object {
	if len(args) != 2 {
		return newError("POW requires 2 arguments, got %d", len(args))
	}

	var base, exp float64

	switch arg := args[0].(type) {
	case *interpreter.Integer:
		base = float64(arg.Value)
	case *interpreter.Real:
		base = arg.Value
	default:
		return newError("POW requires numeric first argument")
	}

	switch arg := args[1].(type) {
	case *interpreter.Integer:
		exp = float64(arg.Value)
	case *interpreter.Real:
		exp = arg.Value
	default:
		return newError("POW requires numeric second argument")
	}

	return &interpreter.Real{Value: math.Pow(base, exp)}
}

func newError(format string, a ...interface{}) *interpreter.Error {
	return &interpreter.Error{Message: fmt.Sprintf(format, a...)}
}
