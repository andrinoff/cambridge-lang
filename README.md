# Cambridge Pseudocode Interpreter

A fully-featured interpreter for the Cambridge International AS & A Level Computer Science (9618) pseudocode language.

## Extensions

Extensions are available for:

* **Visual Studio Code**
    * [Marketplace](https://marketplace.visualstudio.com/items?itemName=andrinoff.cambridge-pseudo)
    * [Pre-built](https://github.com/andrinoff/cambridge-lang/releases)
* **Zed**
    * To install, download the .zip file from [releases](https://github.com/andrinoff/cambridge-lang/releases) and extract it. Go into Zed -> Cmd(Ctrl)+Shift+P -> Install dev extension -> Choose the folder

## Overview

This project implements an interpreter for the pseudocode specification used in Cambridge International AS & A Level Computer Science examinations. It supports all major language features including:

- Variables and constants
- All primitive data types (INTEGER, REAL, STRING, CHAR, BOOLEAN, DATE)
- Arrays (1D and 2D)
- Records (composite data types)
- Selection (IF/THEN/ELSE, CASE)
- Iteration (FOR, WHILE, REPEAT)
- Procedures and functions
- Parameter passing (BYVAL, BYREF)
- File handling
- Object-oriented programming (classes, inheritance)
- Built-in functions

## Installation

### 🍺 Homebrew

```bash
brew tap andrinoff/cambridge
brew install cambridge
brew install cambridge-lsp # Optinal, Language Server Protocol
```

### Snapcraft

```bash
snap install cambridge
```

### Already built binaries

I build all binaries, and they are available [here](https://github.com/andrinoff/cambridge-lang/releases)

### Build From Source

```bash
git clone https://github.com/andrinoff/cambridge-lang.git
cd cambridge-lang
go build -o cambridge ./cmd/cambridge
```

### Usage

```bash
# Run a pseudocode file
./cambridge run program.pseudo

# Start interactive REPL
./cambridge repl

# Show version
./cambridge version

# Show help
./cambridge help
```

## Language Reference

### Variables and Constants

```
DECLARE Name : STRING
DECLARE Age : INTEGER
DECLARE Height : REAL
DECLARE IsStudent : BOOLEAN

CONSTANT PI = 3.14159
CONSTANT GREETING = "Hello"

Name <- "Alice"
Age <- 17
```

### Data Types

| Type | Description | Example |
|------|-------------|---------|
| INTEGER | Whole numbers | `42`, `-17` |
| REAL | Floating-point numbers | `3.14`, `-0.5` |
| STRING | Text strings | `"Hello"` |
| CHAR | Single character | `'A'` |
| BOOLEAN | True/False | `TRUE`, `FALSE` |
| DATE | Date values | - |

### Arrays

```
// 1D Array
DECLARE Numbers : ARRAY[1:10] OF INTEGER
Numbers[1] <- 100

// 2D Array
DECLARE Matrix : ARRAY[1:3, 1:3] OF INTEGER
Matrix[1, 2] <- 5
```

### Selection

```
// IF statement
IF Score >= 90 THEN
    Grade <- "A"
ELSE
    Grade <- "B"
ENDIF

// CASE statement
CASE OF DayNumber
    1 : OUTPUT "Monday"
    2 : OUTPUT "Tuesday"
    6, 7 : OUTPUT "Weekend"
    OTHERWISE : OUTPUT "Midweek"
ENDCASE
```

### Iteration

```
// FOR loop
FOR i <- 1 TO 10
    OUTPUT i
NEXT i

// FOR with STEP
FOR i <- 10 TO 0 STEP -2
    OUTPUT i
NEXT i

// WHILE loop
WHILE Count < 10
    Count <- Count + 1
ENDWHILE

// REPEAT loop
REPEAT
    INPUT Value
UNTIL Value > 0
```

### Procedures and Functions

```
// Procedure
PROCEDURE Greet(Name : STRING)
    OUTPUT "Hello, ", Name
ENDPROCEDURE

CALL Greet("World")

// Function
FUNCTION Square(N : INTEGER) RETURNS INTEGER
    RETURN N * N
ENDFUNCTION

Result <- Square(5)

// BYREF parameter
PROCEDURE Swap(BYREF A : INTEGER, BYREF B : INTEGER)
    DECLARE Temp : INTEGER
    Temp <- A
    A <- B
    B <- Temp
ENDPROCEDURE
```

### Records

```
TYPE Student
    DECLARE Name : STRING
    DECLARE Age : INTEGER
    DECLARE Grade : CHAR
ENDTYPE

DECLARE MyStudent : Student
MyStudent.Name <- "Alice"
MyStudent.Age <- 17
```

### File Handling

```
// Writing to a file
OPENFILE "data.txt" FOR WRITE
WRITEFILE "data.txt", "Hello, World!"
CLOSEFILE "data.txt"

// Reading from a file
OPENFILE "data.txt" FOR READ
WHILE NOT EOF("data.txt")
    READFILE "data.txt", Line
    OUTPUT Line
ENDWHILE
CLOSEFILE "data.txt"

// Appending to a file
OPENFILE "data.txt" FOR APPEND
WRITEFILE "data.txt", "New line"
CLOSEFILE "data.txt"
```

### Object-Oriented Programming

```
CLASS Animal
    PRIVATE Name : STRING
    
    PUBLIC PROCEDURE NEW(GivenName : STRING)
        Name <- GivenName
    ENDPROCEDURE
    
    PUBLIC FUNCTION GetName() RETURNS STRING
        RETURN Name
    ENDFUNCTION
ENDCLASS

CLASS Dog INHERITS Animal
    PUBLIC PROCEDURE Speak()
        OUTPUT GetName(), " says Woof!"
    ENDPROCEDURE
ENDCLASS

DECLARE MyDog : Dog
MyDog <- NEW Dog("Buddy")
CALL MyDog.Speak()
```

### Built-in Functions

#### String Functions
| Function | Description | Example |
|----------|-------------|---------|
| `LENGTH(s)` | Returns string length | `LENGTH("Hello")` → `5` |
| `LEFT(s, n)` | Returns leftmost n characters | `LEFT("Hello", 2)` → `"He"` |
| `RIGHT(s, n)` | Returns rightmost n characters | `RIGHT("Hello", 2)` → `"lo"` |
| `MID(s, start, len)` | Returns substring | `MID("Hello", 2, 3)` → `"ell"` |
| `UCASE(c)` | Converts to uppercase | `UCASE('a')` → `'A'` |
| `LCASE(c)` | Converts to lowercase | `LCASE('A')` → `'a'` |

#### Character/ASCII Functions
| Function | Description | Example |
|----------|-------------|---------|
| `ASC(c)` | Returns ASCII value | `ASC('A')` → `65` |
| `CHR(n)` | Returns character for ASCII value | `CHR(65)` → `'A'` |

#### Numeric Functions
| Function | Description | Example |
|----------|-------------|---------|
| `INT(x)` | Returns integer part | `INT(3.7)` → `3` |
| `RAND(n)` | Random real 0 to n | `RAND(10)` → `7.23` |
| `ROUND(x, p)` | Rounds to p decimal places | `ROUND(3.456, 2)` → `3.46` |
| `ABS(n)` | Absolute value | `ABS(-5)` → `5` |
| `SQRT(n)` | Square root | `SQRT(16)` → `4` |
| `POW(b, e)` | Power (b^e) | `POW(2, 3)` → `8` |

#### Conversion Functions
| Function | Description | Example |
|----------|-------------|---------|
| `NUM_TO_STR(n)` | Number to string | `NUM_TO_STR(42)` → `"42"` |
| `STR_TO_NUM(s)` | String to number | `STR_TO_NUM("42")` → `42` |

#### File Functions
| Function | Description |
|----------|-------------|
| `EOF(filename)` | Returns TRUE if at end of file |

### Operators

#### Arithmetic
| Operator | Description |
|----------|-------------|
| `+` | Addition |
| `-` | Subtraction |
| `*` | Multiplication |
| `/` | Division (returns REAL) |
| `DIV` | Integer division |
| `MOD` | Modulus (remainder) |

#### Comparison
| Operator | Description |
|----------|-------------|
| `=` | Equal to |
| `<>` | Not equal to |
| `<` | Less than |
| `>` | Greater than |
| `<=` | Less than or equal |
| `>=` | Greater than or equal |

#### Logical
| Operator | Description |
|----------|-------------|
| `AND` | Logical AND |
| `OR` | Logical OR |
| `NOT` | Logical NOT |

#### String
| Operator | Description |
|----------|-------------|
| `&` | Concatenation |

## Examples

See the `examples/` directory for complete example programs:

- `hello.pseudo` - Hello World
- `variables.pseudo` - Variables and constants
- `selection.pseudo` - IF and CASE statements
- `loops.pseudo` - FOR, WHILE, REPEAT loops
- `functions.pseudo` - Procedures and functions
- `arrays.pseudo` - Array operations
- `strings.pseudo` - String manipulation
- `records.pseudo` - Record types
- `oop.pseudo` - Object-oriented programming
- `fileio.pseudo` - File handling

## Specification

This interpreter is based on the Cambridge International AS & A Level Computer Science (9618) pseudocode specification. For the official specification, see:

- [Cambridge 9618 Pseudocode Guide for Teachers (2026)](https://www.cambridgeinternational.org/Images/697401-2026-pseudocode-guide-for-teachers.pdf)

## License

MIT License
