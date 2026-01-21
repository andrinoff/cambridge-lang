// Package ast defines the Abstract Syntax Tree for Cambridge Pseudocode
package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/andrinoff/cambridge-lang/pkg/token"
)

// Node is the base interface for all AST nodes
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement represents a statement node
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression node
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of the AST
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
		out.WriteString("\n")
	}
	return out.String()
}

// ============ EXPRESSIONS ============

// Identifier represents a variable name
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// IntegerLiteral represents an integer value
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// RealLiteral represents a floating-point value
type RealLiteral struct {
	Token token.Token
	Value float64
}

func (rl *RealLiteral) expressionNode()      {}
func (rl *RealLiteral) TokenLiteral() string { return rl.Token.Literal }
func (rl *RealLiteral) String() string       { return rl.Token.Literal }

// StringLiteral represents a string value
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return "\"" + sl.Value + "\"" }

// CharLiteral represents a character value
type CharLiteral struct {
	Token token.Token
	Value string
}

func (cl *CharLiteral) expressionNode()      {}
func (cl *CharLiteral) TokenLiteral() string { return cl.Token.Literal }
func (cl *CharLiteral) String() string       { return "'" + cl.Value + "'" }

// BooleanLiteral represents a boolean value
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) String() string       { return bl.Token.Literal }

// PrefixExpression represents a prefix operation (e.g., NOT, -)
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	return "(" + pe.Operator + " " + pe.Right.String() + ")"
}

// InfixExpression represents a binary operation
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	return "(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")"
}

// ArrayAccess represents array indexing: arr[i] or arr[i,j]
type ArrayAccess struct {
	Token   token.Token
	Array   Expression
	Indices []Expression
}

func (aa *ArrayAccess) expressionNode()      {}
func (aa *ArrayAccess) TokenLiteral() string { return aa.Token.Literal }
func (aa *ArrayAccess) String() string {
	var indices []string
	for _, idx := range aa.Indices {
		indices = append(indices, idx.String())
	}
	return aa.Array.String() + "[" + strings.Join(indices, ", ") + "]"
}

// MemberAccess represents object member access: obj.field
type MemberAccess struct {
	Token  token.Token
	Object Expression
	Member string
}

func (ma *MemberAccess) expressionNode()      {}
func (ma *MemberAccess) TokenLiteral() string { return ma.Token.Literal }
func (ma *MemberAccess) String() string {
	return ma.Object.String() + "." + ma.Member
}

// CallExpression represents a function/procedure call
type CallExpression struct {
	Token     token.Token
	Function  Expression // Identifier or MemberAccess
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var args []string
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	return ce.Function.String() + "(" + strings.Join(args, ", ") + ")"
}

// NewExpression represents object instantiation: NEW ClassName(args)
type NewExpression struct {
	Token     token.Token
	ClassName string
	Arguments []Expression
}

func (ne *NewExpression) expressionNode()      {}
func (ne *NewExpression) TokenLiteral() string { return ne.Token.Literal }
func (ne *NewExpression) String() string {
	var args []string
	for _, a := range ne.Arguments {
		args = append(args, a.String())
	}
	return "NEW " + ne.ClassName + "(" + strings.Join(args, ", ") + ")"
}

// SuperExpression represents a reference to parent class: SUPER
type SuperExpression struct {
	Token token.Token
}

func (se *SuperExpression) expressionNode()      {}
func (se *SuperExpression) TokenLiteral() string { return se.Token.Literal }
func (se *SuperExpression) String() string       { return "SUPER" }

// ============ STATEMENTS ============

// DeclareStatement represents: DECLARE x : INTEGER
type DeclareStatement struct {
	Token    token.Token
	Name     *Identifier
	DataType DataType
	Access   string // "PUBLIC" or "PRIVATE" for class properties
}

func (ds *DeclareStatement) statementNode()       {}
func (ds *DeclareStatement) TokenLiteral() string { return ds.Token.Literal }
func (ds *DeclareStatement) String() string {
	var out bytes.Buffer
	if ds.Access != "" {
		out.WriteString(ds.Access + " ")
	}
	out.WriteString("DECLARE " + ds.Name.String() + " : " + ds.DataType.String())
	return out.String()
}

// ConstantStatement represents: CONSTANT PI = 3.14159
type ConstantStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (cs *ConstantStatement) statementNode()       {}
func (cs *ConstantStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ConstantStatement) String() string {
	return "CONSTANT " + cs.Name.String() + " = " + cs.Value.String()
}

// AssignmentStatement represents: x ← 5
type AssignmentStatement struct {
	Token token.Token
	Name  Expression // Identifier, ArrayAccess, or MemberAccess
	Value Expression
}

func (as *AssignmentStatement) statementNode()       {}
func (as *AssignmentStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignmentStatement) String() string {
	return as.Name.String() + " <- " + as.Value.String()
}

// IfStatement represents: IF...THEN...ELSE...ENDIF
type IfStatement struct {
	Token       token.Token
	Condition   Expression
	Consequence []Statement
	Alternative []Statement // nil if no ELSE
}

func (is *IfStatement) statementNode()       {}
func (is *IfStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IfStatement) String() string {
	var out bytes.Buffer
	out.WriteString("IF " + is.Condition.String() + " THEN\n")
	for _, s := range is.Consequence {
		out.WriteString("  " + s.String() + "\n")
	}
	if is.Alternative != nil {
		out.WriteString("ELSE\n")
		for _, s := range is.Alternative {
			out.WriteString("  " + s.String() + "\n")
		}
	}
	out.WriteString("ENDIF")
	return out.String()
}

// CaseClause represents a single case in CASE statement
type CaseClause struct {
	Values []Expression // Can be single value or range
	Body   []Statement
}

// CaseStatement represents: CASE OF...ENDCASE
type CaseStatement struct {
	Token     token.Token
	Expr      Expression
	Cases     []CaseClause
	Otherwise []Statement // nil if no OTHERWISE
}

func (cs *CaseStatement) statementNode()       {}
func (cs *CaseStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *CaseStatement) String() string {
	var out bytes.Buffer
	out.WriteString("CASE OF " + cs.Expr.String() + "\n")
	for _, c := range cs.Cases {
		var vals []string
		for _, v := range c.Values {
			vals = append(vals, v.String())
		}
		out.WriteString("  " + strings.Join(vals, ", ") + " :\n")
		for _, s := range c.Body {
			out.WriteString("    " + s.String() + "\n")
		}
	}
	if cs.Otherwise != nil {
		out.WriteString("  OTHERWISE :\n")
		for _, s := range cs.Otherwise {
			out.WriteString("    " + s.String() + "\n")
		}
	}
	out.WriteString("ENDCASE")
	return out.String()
}

// ForStatement represents: FOR i ← 1 TO 10 STEP 1...NEXT i
type ForStatement struct {
	Token    token.Token
	Variable *Identifier
	Start    Expression
	End      Expression
	Step     Expression // nil if no STEP
	Body     []Statement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForStatement) String() string {
	var out bytes.Buffer
	out.WriteString("FOR " + fs.Variable.String() + " <- " + fs.Start.String() + " TO " + fs.End.String())
	if fs.Step != nil {
		out.WriteString(" STEP " + fs.Step.String())
	}
	out.WriteString("\n")
	for _, s := range fs.Body {
		out.WriteString("  " + s.String() + "\n")
	}
	out.WriteString("NEXT " + fs.Variable.String())
	return out.String()
}

// WhileStatement represents: WHILE...ENDWHILE
type WhileStatement struct {
	Token     token.Token
	Condition Expression
	Body      []Statement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WhileStatement) String() string {
	var out bytes.Buffer
	out.WriteString("WHILE " + ws.Condition.String() + "\n")
	for _, s := range ws.Body {
		out.WriteString("  " + s.String() + "\n")
	}
	out.WriteString("ENDWHILE")
	return out.String()
}

// RepeatStatement represents: REPEAT...UNTIL
type RepeatStatement struct {
	Token     token.Token
	Body      []Statement
	Condition Expression
}

func (rs *RepeatStatement) statementNode()       {}
func (rs *RepeatStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *RepeatStatement) String() string {
	var out bytes.Buffer
	out.WriteString("REPEAT\n")
	for _, s := range rs.Body {
		out.WriteString("  " + s.String() + "\n")
	}
	out.WriteString("UNTIL " + rs.Condition.String())
	return out.String()
}

// Parameter represents a procedure/function parameter
type Parameter struct {
	Name     string
	DataType DataType
	ByRef    bool
}

// ProcedureStatement represents: PROCEDURE name(params)...ENDPROCEDURE
type ProcedureStatement struct {
	Token      token.Token
	Name       string
	Parameters []Parameter
	Body       []Statement
	Access     string // "PUBLIC" or "PRIVATE" for class methods
}

func (ps *ProcedureStatement) statementNode()       {}
func (ps *ProcedureStatement) TokenLiteral() string { return ps.Token.Literal }
func (ps *ProcedureStatement) String() string {
	var out bytes.Buffer
	if ps.Access != "" {
		out.WriteString(ps.Access + " ")
	}
	out.WriteString("PROCEDURE " + ps.Name + "(")
	var params []string
	for _, p := range ps.Parameters {
		pStr := ""
		if p.ByRef {
			pStr = "BYREF "
		}
		pStr += p.Name + " : " + p.DataType.String()
		params = append(params, pStr)
	}
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")\n")
	for _, s := range ps.Body {
		out.WriteString("  " + s.String() + "\n")
	}
	out.WriteString("ENDPROCEDURE")
	return out.String()
}

// FunctionStatement represents: FUNCTION name(params) RETURNS type...ENDFUNCTION
type FunctionStatement struct {
	Token      token.Token
	Name       string
	Parameters []Parameter
	ReturnType DataType
	Body       []Statement
	Access     string // "PUBLIC" or "PRIVATE" for class methods
}

func (fs *FunctionStatement) statementNode()       {}
func (fs *FunctionStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *FunctionStatement) String() string {
	var out bytes.Buffer
	if fs.Access != "" {
		out.WriteString(fs.Access + " ")
	}
	out.WriteString("FUNCTION " + fs.Name + "(")
	var params []string
	for _, p := range fs.Parameters {
		pStr := ""
		if p.ByRef {
			pStr = "BYREF "
		}
		pStr += p.Name + " : " + p.DataType.String()
		params = append(params, pStr)
	}
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") RETURNS " + fs.ReturnType.String() + "\n")
	for _, s := range fs.Body {
		out.WriteString("  " + s.String() + "\n")
	}
	out.WriteString("ENDFUNCTION")
	return out.String()
}

// CallStatement represents: CALL procedure(args)
type CallStatement struct {
	Token     token.Token
	Name      Expression // can be Identifier or MemberAccess
	Arguments []Expression
}

func (cs *CallStatement) statementNode()       {}
func (cs *CallStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *CallStatement) String() string {
	var args []string
	for _, a := range cs.Arguments {
		args = append(args, a.String())
	}
	return "CALL " + cs.Name.String() + "(" + strings.Join(args, ", ") + ")"
}

// ReturnStatement represents: RETURN value
type ReturnStatement struct {
	Token token.Token
	Value Expression // nil for procedures
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	if rs.Value != nil {
		return "RETURN " + rs.Value.String()
	}
	return "RETURN"
}

// InputStatement represents: INPUT variable
type InputStatement struct {
	Token    token.Token
	Variable Expression
}

func (is *InputStatement) statementNode()       {}
func (is *InputStatement) TokenLiteral() string { return is.Token.Literal }
func (is *InputStatement) String() string {
	return "INPUT " + is.Variable.String()
}

// OutputStatement represents: OUTPUT expr1, expr2, ...
type OutputStatement struct {
	Token  token.Token
	Values []Expression
}

func (os *OutputStatement) statementNode()       {}
func (os *OutputStatement) TokenLiteral() string { return os.Token.Literal }
func (os *OutputStatement) String() string {
	var vals []string
	for _, v := range os.Values {
		vals = append(vals, v.String())
	}
	return "OUTPUT " + strings.Join(vals, ", ")
}

// OpenFileStatement represents: OPENFILE filename FOR mode
type OpenFileStatement struct {
	Token    token.Token
	Filename Expression
	Mode     string // "READ", "WRITE", "APPEND"
}

func (of *OpenFileStatement) statementNode()       {}
func (of *OpenFileStatement) TokenLiteral() string { return of.Token.Literal }
func (of *OpenFileStatement) String() string {
	return "OPENFILE " + of.Filename.String() + " FOR " + of.Mode
}

// CloseFileStatement represents: CLOSEFILE filename
type CloseFileStatement struct {
	Token    token.Token
	Filename Expression
}

func (cf *CloseFileStatement) statementNode()       {}
func (cf *CloseFileStatement) TokenLiteral() string { return cf.Token.Literal }
func (cf *CloseFileStatement) String() string {
	return "CLOSEFILE " + cf.Filename.String()
}

// ReadFileStatement represents: READFILE filename, variable
type ReadFileStatement struct {
	Token    token.Token
	Filename Expression
	Variable Expression
}

func (rf *ReadFileStatement) statementNode()       {}
func (rf *ReadFileStatement) TokenLiteral() string { return rf.Token.Literal }
func (rf *ReadFileStatement) String() string {
	return "READFILE " + rf.Filename.String() + ", " + rf.Variable.String()
}

// WriteFileStatement represents: WRITEFILE filename, data
type WriteFileStatement struct {
	Token    token.Token
	Filename Expression
	Data     Expression
}

func (wf *WriteFileStatement) statementNode()       {}
func (wf *WriteFileStatement) TokenLiteral() string { return wf.Token.Literal }
func (wf *WriteFileStatement) String() string {
	return "WRITEFILE " + wf.Filename.String() + ", " + wf.Data.String()
}

// TypeStatement represents: TYPE name...ENDTYPE (for records, enums, etc.)
type TypeStatement struct {
	Token      token.Token
	Name       string
	Definition DataType
}

func (ts *TypeStatement) statementNode()       {}
func (ts *TypeStatement) TokenLiteral() string { return ts.Token.Literal }
func (ts *TypeStatement) String() string {
	return "TYPE " + ts.Name + "\n" + ts.Definition.String() + "\nENDTYPE"
}

// ClassStatement represents: CLASS name INHERITS parent...ENDCLASS
type ClassStatement struct {
	Token   token.Token
	Name    string
	Parent  string // empty if no inheritance
	Members []Statement
}

func (cs *ClassStatement) statementNode()       {}
func (cs *ClassStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ClassStatement) String() string {
	var out bytes.Buffer
	out.WriteString("CLASS " + cs.Name)
	if cs.Parent != "" {
		out.WriteString(" INHERITS " + cs.Parent)
	}
	out.WriteString("\n")
	for _, m := range cs.Members {
		out.WriteString("  " + m.String() + "\n")
	}
	out.WriteString("ENDCLASS")
	return out.String()
}

// ExpressionStatement wraps an expression as a statement
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// ============ DATA TYPES ============

// DataType represents a data type
type DataType interface {
	String() string
}

// PrimitiveType represents basic types: INTEGER, REAL, STRING, CHAR, BOOLEAN, DATE
type PrimitiveType struct {
	Name string
}

func (pt *PrimitiveType) String() string { return pt.Name }

// ArrayType represents: ARRAY[lower:upper] OF type
type ArrayType struct {
	Dimensions  []ArrayDimension
	ElementType DataType
}

type ArrayDimension struct {
	Lower int
	Upper int
}

func (at *ArrayType) String() string {
	var dims []string
	for _, d := range at.Dimensions {
		dims = append(dims, fmt.Sprintf("%d:%d", d.Lower, d.Upper))
	}
	return "ARRAY[" + strings.Join(dims, ",") + "] OF " + at.ElementType.String()
}

// RecordType represents a composite record type
type RecordType struct {
	Fields []RecordField
}

type RecordField struct {
	Name     string
	DataType DataType
}

func (rt *RecordType) String() string {
	var out bytes.Buffer
	for _, f := range rt.Fields {
		out.WriteString("  DECLARE " + f.Name + " : " + f.DataType.String() + "\n")
	}
	return out.String()
}

// EnumType represents an enumerated type
type EnumType struct {
	Values []string
}

func (et *EnumType) String() string {
	return "(" + strings.Join(et.Values, ", ") + ")"
}

// PointerType represents: ^type
type PointerType struct {
	TargetType DataType
}

func (pt *PointerType) String() string {
	return "^" + pt.TargetType.String()
}

// CustomType represents a user-defined type reference
type CustomType struct {
	Name string
}

func (ct *CustomType) String() string { return ct.Name }

// RangeExpression represents a range in CASE: value1 TO value2
type RangeExpression struct {
	Token token.Token
	Start Expression
	End   Expression
}

func (re *RangeExpression) expressionNode()      {}
func (re *RangeExpression) TokenLiteral() string { return re.Token.Literal }
func (re *RangeExpression) String() string {
	return re.Start.String() + " TO " + re.End.String()
}
