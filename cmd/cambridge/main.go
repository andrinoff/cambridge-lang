// Cambridge Pseudocode Interpreter and Compiler
// Based on Cambridge International AS & A Level Computer Science 9618 specification
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/andrinoff/cambridge-lang/pkg/builtins"
	"github.com/andrinoff/cambridge-lang/pkg/interpreter"
	"github.com/andrinoff/cambridge-lang/pkg/lexer"
	"github.com/andrinoff/cambridge-lang/pkg/parser"
)

const VERSION = "0.1.6"

func main() {
	if len(os.Args) < 2 {
		// Start REPL
		startREPL()
		return
	}

	switch os.Args[1] {
	case "run":
		if len(os.Args) < 3 {
			fmt.Println("Usage: cambridge run <filename>")
			os.Exit(1)
		}
		runFile(os.Args[2])
	case "repl":
		startREPL()
	case "version":
		fmt.Printf("Cambridge Pseudocode v%s\n", VERSION)
		fmt.Println("Based on Cambridge International AS & A Level Computer Science 9618")
	case "help":
		printHelp()
	default:
		// Assume it's a filename
		runFile(os.Args[1])
	}
}

func runFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			fmt.Fprintf(os.Stderr, "Parse error: %s\n", err)
		}
		os.Exit(1)
	}

	interp := interpreter.New()
	interp.SetBuiltins(builtins.GetBuiltins())

	result := interp.Eval(program)
	if result != nil {
		if err, ok := result.(*interpreter.Error); ok {
			fmt.Fprintf(os.Stderr, "%s\n", err.Inspect())
			os.Exit(1)
		}
	}
}

func startREPL() {
	fmt.Printf("Cambridge Pseudocode v%s\n", VERSION)
	fmt.Println("Based on Cambridge International AS & A Level Computer Science 9618")
	fmt.Printf("Type 'EXIT' to quit, 'HELP' for help\n")

	reader := bufio.NewReader(os.Stdin)
	interp := interpreter.New()
	interp.SetBuiltins(builtins.GetBuiltins())

	var multilineBuffer strings.Builder
	inMultiline := false

	for {
		if inMultiline {
			fmt.Print("... ")
		} else {
			fmt.Print(">>> ")
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nGoodbye!")
				return
			}
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			continue
		}

		line = strings.TrimRight(line, "\r\n")
		upperLine := strings.ToUpper(strings.TrimSpace(line))

		if upperLine == "EXIT" || upperLine == "QUIT" {
			fmt.Println("Goodbye!")
			return
		}

		if upperLine == "HELP" {
			printREPLHelp()
			continue
		}

		if upperLine == "CLEAR" {
			interp = interpreter.New()
			interp.SetBuiltins(builtins.GetBuiltins())
			fmt.Println("Environment cleared.")
			continue
		}

		// Check if this starts a multiline construct
		if startsMultiline(upperLine) {
			inMultiline = true
			multilineBuffer.WriteString(line)
			multilineBuffer.WriteString("\n")
			continue
		}

		// If in multiline mode, accumulate lines
		if inMultiline {
			multilineBuffer.WriteString(line)
			multilineBuffer.WriteString("\n")

			if endsMultiline(upperLine) {
				inMultiline = false
				line = multilineBuffer.String()
				multilineBuffer.Reset()
			} else {
				continue
			}
		}

		if strings.TrimSpace(line) == "" {
			continue
		}

		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			for _, err := range p.Errors() {
				fmt.Printf("Parse error: %s\n", err)
			}
			continue
		}

		result := interp.Eval(program)
		if result != nil {
			if _, ok := result.(*interpreter.Null); !ok {
				fmt.Println(result.Inspect())
			}
		}
	}
}

func startsMultiline(line string) bool {
	keywords := []string{
		"IF", "WHILE", "FOR", "REPEAT", "CASE",
		"PROCEDURE", "FUNCTION", "CLASS", "TYPE",
	}
	for _, kw := range keywords {
		if strings.HasPrefix(line, kw+" ") || line == kw {
			return true
		}
	}
	return false
}

func endsMultiline(line string) bool {
	endings := []string{
		"ENDIF", "ENDWHILE", "ENDFOR", "UNTIL",
		"ENDCASE", "ENDPROCEDURE", "ENDFUNCTION",
		"ENDCLASS", "ENDTYPE",
	}
	for _, end := range endings {
		if strings.HasPrefix(line, end) {
			return true
		}
	}
	return false
}

func printHelp() {
	fmt.Println(`Cambridge Pseudocode Interpreter

Usage:
  cambridge [command] [arguments]

Commands:
  run <file>    Run a pseudocode file
  repl          Start interactive REPL
  version       Show version information
  help          Show this help message

Examples:
  cambridge run program.pseudo
  cambridge repl

File Extensions:
  .pseudo, .cambridge, .cam, .txt

For more information, visit:
  https://github.com/andrinoff/cambridge-lang`)
}

func printREPLHelp() {
	fmt.Printf(`
REPL Commands:
  EXIT, QUIT    Exit the REPL
  HELP          Show this help
  CLEAR         Clear the environment

Syntax Reference:
  Variables:    DECLARE x : INTEGER
  Constants:    CONSTANT PI = 3.14159
  Assignment:   x <- 5  or  x ← 5

  Selection:    IF condition THEN ... ELSE ... ENDIF
                CASE OF x ... ENDCASE

  Iteration:    FOR i <- 1 TO 10 ... NEXT i
                WHILE condition ... ENDWHILE
                REPEAT ... UNTIL condition

  Procedures:   PROCEDURE Name(params) ... ENDPROCEDURE
  Functions:    FUNCTION Name(params) RETURNS type ... ENDFUNCTION

  I/O:          INPUT x
                OUTPUT "Hello", x

  Files:        OPENFILE "file.txt" FOR READ/WRITE/APPEND
                READFILE "file.txt", variable
                WRITEFILE "file.txt", data
                CLOSEFILE "file.txt"

Built-in Functions:
  String:       LENGTH, LEFT, RIGHT, MID, LCASE, UCASE
  Numeric:      INT, RAND, RANDOM, ROUND, ABS, SQRT, POW
  Conversion:   ASC, CHR, NUM_TO_STR, STR_TO_NUM
  File:         EOF
`)
}
