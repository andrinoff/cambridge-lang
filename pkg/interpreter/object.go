// Package interpreter implements the interpreter for Cambridge Pseudocode
package interpreter

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/andrinoff/cambridge-lang/pkg/ast"
)

// ObjectType represents the type of an object
type ObjectType string

const (
	INTEGER_OBJ      ObjectType = "INTEGER"
	REAL_OBJ         ObjectType = "REAL"
	STRING_OBJ       ObjectType = "STRING"
	CHAR_OBJ         ObjectType = "CHAR"
	BOOLEAN_OBJ      ObjectType = "BOOLEAN"
	NULL_OBJ         ObjectType = "NULL"
	RETURN_VALUE_OBJ ObjectType = "RETURN_VALUE"
	ERROR_OBJ        ObjectType = "ERROR"
	FUNCTION_OBJ     ObjectType = "FUNCTION"
	PROCEDURE_OBJ    ObjectType = "PROCEDURE"
	BUILTIN_OBJ      ObjectType = "BUILTIN"
	ARRAY_OBJ        ObjectType = "ARRAY"
	RECORD_OBJ       ObjectType = "RECORD"
	CLASS_OBJ        ObjectType = "CLASS"
	INSTANCE_OBJ     ObjectType = "INSTANCE"
	FILE_OBJ         ObjectType = "FILE"
	BOUND_METHOD_OBJ ObjectType = "BOUND_METHOD"
	SUPER_OBJ        ObjectType = "SUPER"
)

// Object is the interface all values implement
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer represents an integer value
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

// Real represents a floating-point value
type Real struct {
	Value float64
}

func (r *Real) Type() ObjectType { return REAL_OBJ }
func (r *Real) Inspect() string  { return fmt.Sprintf("%g", r.Value) }

// String represents a string value
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

// Char represents a character value
type Char struct {
	Value rune
}

func (c *Char) Type() ObjectType { return CHAR_OBJ }
func (c *Char) Inspect() string  { return string(c.Value) }

// Boolean represents a boolean value
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string {
	if b.Value {
		return "TRUE"
	}
	return "FALSE"
}

// Null represents a null/nil value
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "NULL" }

// ReturnValue wraps a return value
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// Error represents an error
type Error struct {
	Message string
	Line    int
	Column  int
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string {
	if e.Line > 0 {
		return fmt.Sprintf("ERROR at line %d, column %d: %s", e.Line, e.Column, e.Message)
	}
	return fmt.Sprintf("ERROR: %s", e.Message)
}

// Function represents a user-defined function
type Function struct {
	Name       string
	Parameters []ast.Parameter
	ReturnType ast.DataType
	Body       []ast.Statement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.Name)
	}
	out.WriteString("FUNCTION ")
	out.WriteString(f.Name)
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	return out.String()
}

// Procedure represents a user-defined procedure
type Procedure struct {
	Name       string
	Parameters []ast.Parameter
	Body       []ast.Statement
	Env        *Environment
}

func (p *Procedure) Type() ObjectType { return PROCEDURE_OBJ }
func (p *Procedure) Inspect() string {
	var out bytes.Buffer
	var params []string
	for _, param := range p.Parameters {
		params = append(params, param.Name)
	}
	out.WriteString("PROCEDURE ")
	out.WriteString(p.Name)
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	return out.String()
}

// BuiltinFunction represents a built-in function
type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Name string
	Fn   BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function: " + b.Name }

// Array represents an array
type Array struct {
	Elements   map[string]Object // key is index as string, e.g., "1" or "1,2"
	Dimensions []ast.ArrayDimension
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	return fmt.Sprintf("ARRAY[%d elements]", len(a.Elements))
}

func (a *Array) GetIndex(indices ...int64) string {
	parts := make([]string, len(indices))
	for i, idx := range indices {
		parts[i] = fmt.Sprintf("%d", idx)
	}
	return strings.Join(parts, ",")
}

// Record represents a record instance
type Record struct {
	TypeName string
	Fields   map[string]Object
}

func (r *Record) Type() ObjectType { return RECORD_OBJ }
func (r *Record) Inspect() string {
	return fmt.Sprintf("RECORD %s", r.TypeName)
}

// Class represents a class definition
type Class struct {
	Name    string
	Parent  *Class
	Methods map[string]Object // Function or Procedure
	Fields  map[string]ast.DataType
}

func (c *Class) Type() ObjectType { return CLASS_OBJ }
func (c *Class) Inspect() string  { return fmt.Sprintf("CLASS %s", c.Name) }

// Instance represents an instance of a class
type Instance struct {
	Class  *Class
	Fields map[string]Object
}

func (i *Instance) Type() ObjectType { return INSTANCE_OBJ }
func (i *Instance) Inspect() string  { return fmt.Sprintf("<%s instance>", i.Class.Name) }

// BoundMethod represents a method bound to an instance
type BoundMethod struct {
	Instance *Instance
	Method   Object // Function or Procedure
}

func (bm *BoundMethod) Type() ObjectType { return BOUND_METHOD_OBJ }
func (bm *BoundMethod) Inspect() string {
	return fmt.Sprintf("<bound method of %s>", bm.Instance.Class.Name)
}

// Super represents a reference to the parent class from within an instance
type Super struct {
	Instance *Instance
	Class    *Class // The parent class to use for method lookup
}

func (s *Super) Type() ObjectType { return SUPER_OBJ }
func (s *Super) Inspect() string  { return "SUPER" }

// File represents an open file
type File struct {
	Name   string
	Mode   string
	Handle interface{} // *os.File in actual implementation
}

func (f *File) Type() ObjectType { return FILE_OBJ }
func (f *File) Inspect() string  { return fmt.Sprintf("FILE(%s, %s)", f.Name, f.Mode) }

// Reference represents a reference to a variable (for BYREF)
type Reference struct {
	Name string
	Env  *Environment
}

func (r *Reference) Type() ObjectType { return "REFERENCE" }
func (r *Reference) Inspect() string  { return fmt.Sprintf("&%s", r.Name) }

func (r *Reference) Get() Object {
	obj, _ := r.Env.Get(r.Name)
	return obj
}

func (r *Reference) Set(val Object) {
	r.Env.Set(r.Name, val)
}
