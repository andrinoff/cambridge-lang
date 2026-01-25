package interpreter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/andrinoff/cambridge-lang/pkg/ast"
)

// Interpreter evaluates the AST
type Interpreter struct {
	env      *Environment
	builtins map[string]*Builtin
	files    map[string]*fileState
	input    io.Reader
	output   io.Writer
}

type fileState struct {
	file    *os.File
	mode    string
	scanner *bufio.Scanner
	atEOF   bool
}

// New creates a new interpreter
func New() *Interpreter {
	return &Interpreter{
		env:      NewEnvironment(),
		builtins: make(map[string]*Builtin),
		files:    make(map[string]*fileState),
		input:    os.Stdin,
		output:   os.Stdout,
	}
}

// SetBuiltins sets the built-in functions
func (i *Interpreter) SetBuiltins(builtins map[string]*Builtin) {
	i.builtins = builtins
}

// SetInput sets the input reader
func (i *Interpreter) SetInput(r io.Reader) {
	i.input = r
}

// SetOutput sets the output writer
func (i *Interpreter) SetOutput(w io.Writer) {
	i.output = w
}

// Eval evaluates a program
func (i *Interpreter) Eval(program *ast.Program) Object {
	var result Object

	for _, stmt := range program.Statements {
		result = i.evalStatement(stmt, i.env)

		switch result := result.(type) {
		case *ReturnValue:
			return result.Value
		case *Error:
			return result
		}
	}

	return result
}

func (i *Interpreter) evalStatement(stmt ast.Statement, env *Environment) Object {
	switch stmt := stmt.(type) {
	case *ast.DeclareStatement:
		return i.evalDeclareStatement(stmt, env)
	case *ast.ConstantStatement:
		return i.evalConstantStatement(stmt, env)
	case *ast.AssignmentStatement:
		return i.evalAssignmentStatement(stmt, env)
	case *ast.IfStatement:
		return i.evalIfStatement(stmt, env)
	case *ast.CaseStatement:
		return i.evalCaseStatement(stmt, env)
	case *ast.ForStatement:
		return i.evalForStatement(stmt, env)
	case *ast.WhileStatement:
		return i.evalWhileStatement(stmt, env)
	case *ast.RepeatStatement:
		return i.evalRepeatStatement(stmt, env)
	case *ast.ProcedureStatement:
		return i.evalProcedureStatement(stmt, env)
	case *ast.FunctionStatement:
		return i.evalFunctionStatement(stmt, env)
	case *ast.CallStatement:
		return i.evalCallStatement(stmt, env)
	case *ast.ReturnStatement:
		return i.evalReturnStatement(stmt, env)
	case *ast.InputStatement:
		return i.evalInputStatement(stmt, env)
	case *ast.OutputStatement:
		return i.evalOutputStatement(stmt, env)
	case *ast.OpenFileStatement:
		return i.evalOpenFileStatement(stmt, env)
	case *ast.CloseFileStatement:
		return i.evalCloseFileStatement(stmt, env)
	case *ast.ReadFileStatement:
		return i.evalReadFileStatement(stmt, env)
	case *ast.WriteFileStatement:
		return i.evalWriteFileStatement(stmt, env)
	case *ast.TypeStatement:
		return i.evalTypeStatement(stmt, env)
	case *ast.ClassStatement:
		return i.evalClassStatement(stmt, env)
	case *ast.ExpressionStatement:
		return i.evalExpression(stmt.Expression, env)
	default:
		return &Error{Message: fmt.Sprintf("unknown statement type: %T", stmt)}
	}
}

func (i *Interpreter) evalDeclareStatement(stmt *ast.DeclareStatement, env *Environment) Object {
	var value Object

	switch dt := stmt.DataType.(type) {
	case *ast.PrimitiveType:
		switch dt.Name {
		case "INTEGER":
			value = &Integer{Value: 0}
		case "REAL":
			value = &Real{Value: 0.0}
		case "STRING":
			value = &String{Value: ""}
		case "CHAR":
			value = &Char{Value: ' '}
		case "BOOLEAN":
			value = &Boolean{Value: false}
		case "DATE":
			value = &Date{Day: 1, Month: 1, Year: 1970}
		default:
			value = &Null{}
		}
	case *ast.ArrayType:
		arr := &Array{
			Elements:   make(map[string]Object),
			Dimensions: dt.Dimensions,
		}
		value = arr
	case *ast.CustomType:
		// Check if it's a defined type
		if typ, ok := env.GetType(dt.Name); ok {
			switch t := typ.(type) {
			case *Record:
				// Create a new record instance
				rec := &Record{
					TypeName: dt.Name,
					Fields:   make(map[string]Object),
				}
				for name := range t.Fields {
					rec.Fields[name] = &Null{}
				}
				value = rec
			default:
				value = &Null{}
			}
		} else {
			value = &Null{}
		}
	default:
		value = &Null{}
	}

	return env.Declare(stmt.Name.Value, value)
}

func (i *Interpreter) evalConstantStatement(stmt *ast.ConstantStatement, env *Environment) Object {
	value := i.evalExpression(stmt.Value, env)
	if isError(value) {
		return value
	}
	return env.DeclareConstant(stmt.Name.Value, value)
}

func (i *Interpreter) evalAssignmentStatement(stmt *ast.AssignmentStatement, env *Environment) Object {
	value := i.evalExpression(stmt.Value, env)
	if isError(value) {
		return value
	}

	switch target := stmt.Name.(type) {
	case *ast.Identifier:
		return env.SetInPlace(target.Value, value)
	case *ast.ArrayAccess:
		return i.evalArrayAssignment(target, value, env)
	case *ast.MemberAccess:
		return i.evalMemberAssignment(target, value, env)
	default:
		return &Error{Message: "invalid assignment target"}
	}
}

func (i *Interpreter) evalArrayAssignment(access *ast.ArrayAccess, value Object, env *Environment) Object {
	arr := i.evalExpression(access.Array, env)
	if isError(arr) {
		return arr
	}

	array, ok := arr.(*Array)
	if !ok {
		return &Error{Message: "not an array"}
	}

	indices := []int64{}
	for _, idx := range access.Indices {
		idxVal := i.evalExpression(idx, env)
		if isError(idxVal) {
			return idxVal
		}
		intVal, ok := idxVal.(*Integer)
		if !ok {
			return &Error{Message: "array index must be an integer"}
		}
		indices = append(indices, intVal.Value)
	}

	key := array.GetIndex(indices...)
	array.Elements[key] = value
	return value
}

func (i *Interpreter) evalMemberAssignment(access *ast.MemberAccess, value Object, env *Environment) Object {
	obj := i.evalExpression(access.Object, env)
	if isError(obj) {
		return obj
	}

	switch o := obj.(type) {
	case *Record:
		o.Fields[access.Member] = value
		return value
	case *Instance:
		o.Fields[access.Member] = value
		return value
	default:
		return &Error{Message: "cannot access member of non-record/instance"}
	}
}

func (i *Interpreter) evalIfStatement(stmt *ast.IfStatement, env *Environment) Object {
	condition := i.evalExpression(stmt.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return i.evalStatements(stmt.Consequence, env)
	} else if stmt.Alternative != nil {
		return i.evalStatements(stmt.Alternative, env)
	}

	return &Null{}
}

func (i *Interpreter) evalCaseStatement(stmt *ast.CaseStatement, env *Environment) Object {
	value := i.evalExpression(stmt.Expr, env)
	if isError(value) {
		return value
	}

	for _, caseClause := range stmt.Cases {
		for _, caseValue := range caseClause.Values {
			if i.matchesCaseValue(value, caseValue, env) {
				return i.evalStatements(caseClause.Body, env)
			}
		}
	}

	if stmt.Otherwise != nil {
		return i.evalStatements(stmt.Otherwise, env)
	}

	return &Null{}
}

func (i *Interpreter) matchesCaseValue(value Object, caseValue ast.Expression, env *Environment) bool {
	switch cv := caseValue.(type) {
	case *ast.RangeExpression:
		start := i.evalExpression(cv.Start, env)
		end := i.evalExpression(cv.End, env)
		return i.valueInRange(value, start, end)
	default:
		evalValue := i.evalExpression(caseValue, env)
		return i.objectsEqual(value, evalValue)
	}
}

func (i *Interpreter) valueInRange(value, start, end Object) bool {
	switch v := value.(type) {
	case *Integer:
		s, sok := start.(*Integer)
		e, eok := end.(*Integer)
		if sok && eok {
			return v.Value >= s.Value && v.Value <= e.Value
		}
	case *Char:
		s, sok := start.(*Char)
		e, eok := end.(*Char)
		if sok && eok {
			return v.Value >= s.Value && v.Value <= e.Value
		}
	}
	return false
}

func (i *Interpreter) objectsEqual(a, b Object) bool {
	switch av := a.(type) {
	case *Integer:
		if bv, ok := b.(*Integer); ok {
			return av.Value == bv.Value
		}
	case *Real:
		if bv, ok := b.(*Real); ok {
			return av.Value == bv.Value
		}
	case *String:
		if bv, ok := b.(*String); ok {
			return av.Value == bv.Value
		}
	case *Char:
		if bv, ok := b.(*Char); ok {
			return av.Value == bv.Value
		}
	case *Boolean:
		if bv, ok := b.(*Boolean); ok {
			return av.Value == bv.Value
		}
	}
	return false
}

func (i *Interpreter) evalForStatement(stmt *ast.ForStatement, env *Environment) Object {
	start := i.evalExpression(stmt.Start, env)
	if isError(start) {
		return start
	}

	end := i.evalExpression(stmt.End, env)
	if isError(end) {
		return end
	}

	step := int64(1)
	if stmt.Step != nil {
		stepVal := i.evalExpression(stmt.Step, env)
		if isError(stepVal) {
			return stepVal
		}
		if s, ok := stepVal.(*Integer); ok {
			step = s.Value
		}
	}

	startInt, ok := start.(*Integer)
	if !ok {
		return &Error{Message: "FOR loop start must be an integer"}
	}

	endInt, ok := end.(*Integer)
	if !ok {
		return &Error{Message: "FOR loop end must be an integer"}
	}

	loopEnv := NewEnclosedEnvironment(env)
	loopEnv.Declare(stmt.Variable.Value, startInt)

	var result Object
	for current := startInt.Value; ; {
		// Check loop condition
		if step > 0 && current > endInt.Value {
			break
		}
		if step < 0 && current < endInt.Value {
			break
		}

		loopEnv.SetInPlace(stmt.Variable.Value, &Integer{Value: current})
		result = i.evalStatements(stmt.Body, loopEnv)

		if isError(result) {
			return result
		}
		if _, ok := result.(*ReturnValue); ok {
			return result
		}

		current += step
	}

	return result
}

func (i *Interpreter) evalWhileStatement(stmt *ast.WhileStatement, env *Environment) Object {
	var result Object

	for {
		condition := i.evalExpression(stmt.Condition, env)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result = i.evalStatements(stmt.Body, env)
		if isError(result) {
			return result
		}
		if _, ok := result.(*ReturnValue); ok {
			return result
		}
	}

	return result
}

func (i *Interpreter) evalRepeatStatement(stmt *ast.RepeatStatement, env *Environment) Object {
	var result Object

	for {
		result = i.evalStatements(stmt.Body, env)
		if isError(result) {
			return result
		}
		if _, ok := result.(*ReturnValue); ok {
			return result
		}

		condition := i.evalExpression(stmt.Condition, env)
		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			break
		}
	}

	return result
}

func (i *Interpreter) evalProcedureStatement(stmt *ast.ProcedureStatement, env *Environment) Object {
	proc := &Procedure{
		Name:       stmt.Name,
		Parameters: stmt.Parameters,
		Body:       stmt.Body,
		Env:        env,
	}
	return env.Declare(stmt.Name, proc)
}

func (i *Interpreter) evalFunctionStatement(stmt *ast.FunctionStatement, env *Environment) Object {
	fn := &Function{
		Name:       stmt.Name,
		Parameters: stmt.Parameters,
		ReturnType: stmt.ReturnType,
		Body:       stmt.Body,
		Env:        env,
	}
	return env.Declare(stmt.Name, fn)
}

func (i *Interpreter) evalCallStatement(stmt *ast.CallStatement, env *Environment) Object {
	// Evaluate the call
	call := &ast.CallExpression{
		Token:     stmt.Token,
		Function:  stmt.Name,
		Arguments: stmt.Arguments,
	}
	return i.evalCallExpression(call, env)
}

func (i *Interpreter) evalReturnStatement(stmt *ast.ReturnStatement, env *Environment) Object {
	if stmt.Value == nil {
		return &ReturnValue{Value: &Null{}}
	}

	value := i.evalExpression(stmt.Value, env)
	if isError(value) {
		return value
	}

	return &ReturnValue{Value: value}
}

func (i *Interpreter) evalInputStatement(stmt *ast.InputStatement, env *Environment) Object {
	reader := bufio.NewReader(i.input)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return &Error{Message: fmt.Sprintf("input error: %v", err)}
	}

	line = strings.TrimRight(line, "\r\n")

	switch target := stmt.Variable.(type) {
	case *ast.Identifier:
		env.SetInPlace(target.Value, &String{Value: line})
	case *ast.ArrayAccess:
		return i.evalArrayAssignment(target, &String{Value: line}, env)
	}

	return &Null{}
}

func (i *Interpreter) evalOutputStatement(stmt *ast.OutputStatement, env *Environment) Object {
	var parts []string

	for _, expr := range stmt.Values {
		value := i.evalExpression(expr, env)
		if isError(value) {
			return value
		}
		parts = append(parts, value.Inspect())
	}

	fmt.Fprintln(i.output, strings.Join(parts, ""))
	return &Null{}
}

func (i *Interpreter) evalOpenFileStatement(stmt *ast.OpenFileStatement, env *Environment) Object {
	filename := i.evalExpression(stmt.Filename, env)
	if isError(filename) {
		return filename
	}

	filenameStr, ok := filename.(*String)
	if !ok {
		return &Error{Message: "filename must be a string"}
	}

	var file *os.File
	var err error

	switch stmt.Mode {
	case "READ":
		file, err = os.Open(filenameStr.Value)
	case "WRITE":
		file, err = os.Create(filenameStr.Value)
	case "APPEND":
		file, err = os.OpenFile(filenameStr.Value, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}

	if err != nil {
		return &Error{Message: fmt.Sprintf("cannot open file: %v", err)}
	}

	fs := &fileState{
		file: file,
		mode: stmt.Mode,
	}
	if stmt.Mode == "READ" {
		fs.scanner = bufio.NewScanner(file)
	}

	i.files[filenameStr.Value] = fs
	return &Null{}
}

func (i *Interpreter) evalCloseFileStatement(stmt *ast.CloseFileStatement, env *Environment) Object {
	filename := i.evalExpression(stmt.Filename, env)
	if isError(filename) {
		return filename
	}

	filenameStr, ok := filename.(*String)
	if !ok {
		return &Error{Message: "filename must be a string"}
	}

	fs, ok := i.files[filenameStr.Value]
	if !ok {
		return &Error{Message: "file not open"}
	}

	fs.file.Close()
	delete(i.files, filenameStr.Value)
	return &Null{}
}

func (i *Interpreter) evalReadFileStatement(stmt *ast.ReadFileStatement, env *Environment) Object {
	filename := i.evalExpression(stmt.Filename, env)
	if isError(filename) {
		return filename
	}

	filenameStr, ok := filename.(*String)
	if !ok {
		return &Error{Message: "filename must be a string"}
	}

	fs, ok := i.files[filenameStr.Value]
	if !ok {
		return &Error{Message: "file not open"}
	}

	if fs.mode != "READ" {
		return &Error{Message: "file not open for reading"}
	}

	if fs.scanner.Scan() {
		line := fs.scanner.Text()
		switch target := stmt.Variable.(type) {
		case *ast.Identifier:
			env.SetInPlace(target.Value, &String{Value: line})
		case *ast.ArrayAccess:
			return i.evalArrayAssignment(target, &String{Value: line}, env)
		}
	} else {
		fs.atEOF = true
	}

	return &Null{}
}

func (i *Interpreter) evalWriteFileStatement(stmt *ast.WriteFileStatement, env *Environment) Object {
	filename := i.evalExpression(stmt.Filename, env)
	if isError(filename) {
		return filename
	}

	filenameStr, ok := filename.(*String)
	if !ok {
		return &Error{Message: "filename must be a string"}
	}

	fs, ok := i.files[filenameStr.Value]
	if !ok {
		return &Error{Message: "file not open"}
	}

	if fs.mode != "WRITE" && fs.mode != "APPEND" {
		return &Error{Message: "file not open for writing"}
	}

	data := i.evalExpression(stmt.Data, env)
	if isError(data) {
		return data
	}

	_, err := fmt.Fprintln(fs.file, data.Inspect())
	if err != nil {
		return &Error{Message: fmt.Sprintf("write error: %v", err)}
	}

	return &Null{}
}

func (i *Interpreter) evalTypeStatement(stmt *ast.TypeStatement, env *Environment) Object {
	switch def := stmt.Definition.(type) {
	case *ast.RecordType:
		rec := &Record{
			TypeName: stmt.Name,
			Fields:   make(map[string]Object),
		}
		for _, field := range def.Fields {
			rec.Fields[field.Name] = &Null{}
		}
		env.DefineType(stmt.Name, rec)
	case *ast.EnumType:
		// Store enum values
		for idx, val := range def.Values {
			env.Declare(val, &Integer{Value: int64(idx)})
		}
	}
	return &Null{}
}

func (i *Interpreter) evalClassStatement(stmt *ast.ClassStatement, env *Environment) Object {
	class := &Class{
		Name:    stmt.Name,
		Methods: make(map[string]Object),
		Fields:  make(map[string]ast.DataType),
	}

	if stmt.Parent != "" {
		if parentObj, ok := env.Get(stmt.Parent); ok {
			if parent, ok := parentObj.(*Class); ok {
				class.Parent = parent
			}
		}
	}

	classEnv := NewEnclosedEnvironment(env)

	for _, member := range stmt.Members {
		switch m := member.(type) {
		case *ast.DeclareStatement:
			class.Fields[m.Name.Value] = m.DataType
		case *ast.ProcedureStatement:
			proc := &Procedure{
				Name:       m.Name,
				Parameters: m.Parameters,
				Body:       m.Body,
				Env:        classEnv,
			}
			class.Methods[m.Name] = proc
		case *ast.FunctionStatement:
			fn := &Function{
				Name:       m.Name,
				Parameters: m.Parameters,
				ReturnType: m.ReturnType,
				Body:       m.Body,
				Env:        classEnv,
			}
			class.Methods[m.Name] = fn
		}
	}

	return env.Declare(stmt.Name, class)
}

func (i *Interpreter) evalStatements(stmts []ast.Statement, env *Environment) Object {
	var result Object

	for _, stmt := range stmts {
		result = i.evalStatement(stmt, env)

		if result != nil {
			if result.Type() == RETURN_VALUE_OBJ || result.Type() == ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func (i *Interpreter) evalExpression(expr ast.Expression, env *Environment) Object {
	switch expr := expr.(type) {
	case *ast.IntegerLiteral:
		return &Integer{Value: expr.Value}
	case *ast.RealLiteral:
		return &Real{Value: expr.Value}
	case *ast.StringLiteral:
		return &String{Value: expr.Value}
	case *ast.CharLiteral:
		if len(expr.Value) > 0 {
			return &Char{Value: rune(expr.Value[0])}
		}
		return &Char{Value: ' '}
	case *ast.BooleanLiteral:
		return &Boolean{Value: expr.Value}
	case *ast.Identifier:
		return i.evalIdentifier(expr, env)
	case *ast.PrefixExpression:
		return i.evalPrefixExpression(expr, env)
	case *ast.InfixExpression:
		return i.evalInfixExpression(expr, env)
	case *ast.ArrayAccess:
		return i.evalArrayAccess(expr, env)
	case *ast.MemberAccess:
		return i.evalMemberAccess(expr, env)
	case *ast.CallExpression:
		return i.evalCallExpression(expr, env)
	case *ast.NewExpression:
		return i.evalNewExpression(expr, env)
	case *ast.SuperExpression:
		return i.evalSuperExpression(expr, env)
	default:
		return &Error{Message: fmt.Sprintf("unknown expression type: %T", expr)}
	}
}

func (i *Interpreter) evalIdentifier(node *ast.Identifier, env *Environment) Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := i.builtins[node.Value]; ok {
		return builtin
	}

	return &Error{Message: fmt.Sprintf("identifier not found: %s", node.Value)}
}

func (i *Interpreter) evalPrefixExpression(expr *ast.PrefixExpression, env *Environment) Object {
	right := i.evalExpression(expr.Right, env)
	if isError(right) {
		return right
	}

	switch expr.Operator {
	case "-":
		return i.evalMinusPrefixOperator(right)
	case "NOT":
		return i.evalNotOperator(right)
	default:
		return &Error{Message: fmt.Sprintf("unknown operator: %s", expr.Operator)}
	}
}

func (i *Interpreter) evalMinusPrefixOperator(right Object) Object {
	switch obj := right.(type) {
	case *Integer:
		return &Integer{Value: -obj.Value}
	case *Real:
		return &Real{Value: -obj.Value}
	default:
		return &Error{Message: fmt.Sprintf("unknown operator: -%s", right.Type())}
	}
}

func (i *Interpreter) evalNotOperator(right Object) Object {
	switch obj := right.(type) {
	case *Boolean:
		return &Boolean{Value: !obj.Value}
	default:
		return &Error{Message: fmt.Sprintf("unknown operator: NOT %s", right.Type())}
	}
}

func (i *Interpreter) evalInfixExpression(expr *ast.InfixExpression, env *Environment) Object {
	left := i.evalExpression(expr.Left, env)
	if isError(left) {
		return left
	}

	right := i.evalExpression(expr.Right, env)
	if isError(right) {
		return right
	}

	switch {
	case left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
		return i.evalIntegerInfixExpression(expr.Operator, left, right)
	case left.Type() == REAL_OBJ || right.Type() == REAL_OBJ:
		return i.evalRealInfixExpression(expr.Operator, left, right)
	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		return i.evalStringInfixExpression(expr.Operator, left, right)
	case left.Type() == BOOLEAN_OBJ && right.Type() == BOOLEAN_OBJ:
		return i.evalBooleanInfixExpression(expr.Operator, left, right)
	case expr.Operator == "&":
		// String concatenation - convert operands to strings
		return i.evalConcatenation(left, right)
	case expr.Operator == "=":
		return &Boolean{Value: i.objectsEqual(left, right)}
	case expr.Operator == "<>":
		return &Boolean{Value: !i.objectsEqual(left, right)}
	default:
		return &Error{Message: fmt.Sprintf("type mismatch: %s %s %s", left.Type(), expr.Operator, right.Type())}
	}
}

func (i *Interpreter) evalIntegerInfixExpression(op string, left, right Object) Object {
	leftVal := left.(*Integer).Value
	rightVal := right.(*Integer).Value

	switch op {
	case "+":
		return &Integer{Value: leftVal + rightVal}
	case "-":
		return &Integer{Value: leftVal - rightVal}
	case "*":
		return &Integer{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return &Error{Message: "division by zero"}
		}
		return &Real{Value: float64(leftVal) / float64(rightVal)}
	case "DIV":
		if rightVal == 0 {
			return &Error{Message: "division by zero"}
		}
		return &Integer{Value: leftVal / rightVal}
	case "MOD":
		if rightVal == 0 {
			return &Error{Message: "division by zero"}
		}
		return &Integer{Value: leftVal % rightVal}
	case "<":
		return &Boolean{Value: leftVal < rightVal}
	case ">":
		return &Boolean{Value: leftVal > rightVal}
	case "<=":
		return &Boolean{Value: leftVal <= rightVal}
	case ">=":
		return &Boolean{Value: leftVal >= rightVal}
	case "=":
		return &Boolean{Value: leftVal == rightVal}
	case "<>":
		return &Boolean{Value: leftVal != rightVal}
	default:
		return &Error{Message: fmt.Sprintf("unknown operator: %s %s %s", left.Type(), op, right.Type())}
	}
}

func (i *Interpreter) evalRealInfixExpression(op string, left, right Object) Object {
	var leftVal, rightVal float64

	switch l := left.(type) {
	case *Real:
		leftVal = l.Value
	case *Integer:
		leftVal = float64(l.Value)
	}

	switch r := right.(type) {
	case *Real:
		rightVal = r.Value
	case *Integer:
		rightVal = float64(r.Value)
	}

	switch op {
	case "+":
		return &Real{Value: leftVal + rightVal}
	case "-":
		return &Real{Value: leftVal - rightVal}
	case "*":
		return &Real{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return &Error{Message: "division by zero"}
		}
		return &Real{Value: leftVal / rightVal}
	case "<":
		return &Boolean{Value: leftVal < rightVal}
	case ">":
		return &Boolean{Value: leftVal > rightVal}
	case "<=":
		return &Boolean{Value: leftVal <= rightVal}
	case ">=":
		return &Boolean{Value: leftVal >= rightVal}
	case "=":
		return &Boolean{Value: leftVal == rightVal}
	case "<>":
		return &Boolean{Value: leftVal != rightVal}
	default:
		return &Error{Message: fmt.Sprintf("unknown operator: %s %s %s", left.Type(), op, right.Type())}
	}
}

func (i *Interpreter) evalStringInfixExpression(op string, left, right Object) Object {
	leftVal := left.(*String).Value
	rightVal := right.(*String).Value

	switch op {
	case "&":
		return &String{Value: leftVal + rightVal}
	case "=":
		return &Boolean{Value: leftVal == rightVal}
	case "<>":
		return &Boolean{Value: leftVal != rightVal}
	case "<":
		return &Boolean{Value: leftVal < rightVal}
	case ">":
		return &Boolean{Value: leftVal > rightVal}
	case "<=":
		return &Boolean{Value: leftVal <= rightVal}
	case ">=":
		return &Boolean{Value: leftVal >= rightVal}
	default:
		return &Error{Message: fmt.Sprintf("unknown operator: %s %s %s", left.Type(), op, right.Type())}
	}
}

func (i *Interpreter) evalBooleanInfixExpression(op string, left, right Object) Object {
	leftVal := left.(*Boolean).Value
	rightVal := right.(*Boolean).Value

	switch op {
	case "AND":
		return &Boolean{Value: leftVal && rightVal}
	case "OR":
		return &Boolean{Value: leftVal || rightVal}
	case "=":
		return &Boolean{Value: leftVal == rightVal}
	case "<>":
		return &Boolean{Value: leftVal != rightVal}
	default:
		return &Error{Message: fmt.Sprintf("unknown operator: %s %s %s", left.Type(), op, right.Type())}
	}
}

func (i *Interpreter) evalConcatenation(left, right Object) Object {
	leftStr := i.objectToString(left)
	rightStr := i.objectToString(right)
	return &String{Value: leftStr + rightStr}
}

func (i *Interpreter) objectToString(obj Object) string {
	switch o := obj.(type) {
	case *String:
		return o.Value
	case *Char:
		return string(o.Value)
	case *Integer:
		return fmt.Sprintf("%d", o.Value)
	case *Real:
		return fmt.Sprintf("%g", o.Value)
	case *Boolean:
		if o.Value {
			return "TRUE"
		}
		return "FALSE"
	default:
		return obj.Inspect()
	}
}

func (i *Interpreter) evalArrayAccess(expr *ast.ArrayAccess, env *Environment) Object {
	arr := i.evalExpression(expr.Array, env)
	if isError(arr) {
		return arr
	}

	array, ok := arr.(*Array)
	if !ok {
		return &Error{Message: "not an array"}
	}

	indices := []int64{}
	for _, idx := range expr.Indices {
		idxVal := i.evalExpression(idx, env)
		if isError(idxVal) {
			return idxVal
		}
		intVal, ok := idxVal.(*Integer)
		if !ok {
			return &Error{Message: "array index must be an integer"}
		}
		indices = append(indices, intVal.Value)
	}

	key := array.GetIndex(indices...)
	if val, ok := array.Elements[key]; ok {
		return val
	}

	return &Null{}
}

func (i *Interpreter) evalMemberAccess(expr *ast.MemberAccess, env *Environment) Object {
	obj := i.evalExpression(expr.Object, env)
	if isError(obj) {
		return obj
	}

	switch o := obj.(type) {
	case *Record:
		if val, ok := o.Fields[expr.Member]; ok {
			return val
		}
		return &Error{Message: fmt.Sprintf("field not found: %s", expr.Member)}
	case *Instance:
		if val, ok := o.Fields[expr.Member]; ok {
			return val
		}
		// Look up method in class hierarchy
		if method := i.lookupMethod(o.Class, expr.Member); method != nil {
			return &BoundMethod{Instance: o, Method: method}
		}
		return &Error{Message: fmt.Sprintf("member not found: %s", expr.Member)}
	case *Super:
		// Look up method in parent class
		if o.Class == nil {
			return &Error{Message: "no parent class"}
		}
		if method := i.lookupMethod(o.Class, expr.Member); method != nil {
			return &BoundMethod{Instance: o.Instance, Method: method}
		}
		return &Error{Message: fmt.Sprintf("method not found in parent class: %s", expr.Member)}
	default:
		return &Error{Message: "cannot access member of non-record/instance"}
	}
}

// lookupMethod searches for a method in the class hierarchy
func (i *Interpreter) lookupMethod(class *Class, name string) Object {
	for c := class; c != nil; c = c.Parent {
		if method, ok := c.Methods[name]; ok {
			return method
		}
	}
	return nil
}

func (i *Interpreter) evalCallExpression(expr *ast.CallExpression, env *Environment) Object {
	fn := i.evalExpression(expr.Function, env)
	if isError(fn) {
		return fn
	}

	args := i.evalExpressions(expr.Arguments, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	return i.applyFunction(fn, args, env)
}

func (i *Interpreter) evalExpressions(exprs []ast.Expression, env *Environment) []Object {
	var result []Object

	for _, e := range exprs {
		evaluated := i.evalExpression(e, env)
		if isError(evaluated) {
			return []Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func (i *Interpreter) applyFunction(fn Object, args []Object, callerEnv *Environment) Object {
	switch fn := fn.(type) {
	case *Function:
		extendedEnv := i.extendFunctionEnv(fn, args, fn.Parameters, callerEnv)
		evaluated := i.evalStatements(fn.Body, extendedEnv)
		return i.unwrapReturnValue(evaluated)

	case *Procedure:
		extendedEnv := i.extendFunctionEnv(&Function{Env: fn.Env}, args, fn.Parameters, callerEnv)
		evaluated := i.evalStatements(fn.Body, extendedEnv)
		return i.unwrapReturnValue(evaluated)

	case *BoundMethod:
		return i.applyBoundMethod(fn, args, callerEnv)

	case *Builtin:
		return fn.Fn(args...)

	default:
		return &Error{Message: fmt.Sprintf("not a function: %s", fn.Type())}
	}
}

func (i *Interpreter) applyBoundMethod(bm *BoundMethod, args []Object, callerEnv *Environment) Object {
	// Create a method environment that has access to instance fields and methods
	methodEnv := i.createMethodEnv(bm.Instance, callerEnv)

	switch method := bm.Method.(type) {
	case *Function:
		// Add parameters to the environment
		for idx, param := range method.Parameters {
			if idx < len(args) {
				methodEnv.Declare(param.Name, args[idx])
			}
		}
		evaluated := i.evalStatements(method.Body, methodEnv)
		return i.unwrapReturnValue(evaluated)

	case *Procedure:
		// Add parameters to the environment
		for idx, param := range method.Parameters {
			if idx < len(args) {
				methodEnv.Declare(param.Name, args[idx])
			}
		}
		evaluated := i.evalStatements(method.Body, methodEnv)
		return i.unwrapReturnValue(evaluated)

	default:
		return &Error{Message: "invalid method type"}
	}
}

// createMethodEnv creates an environment for method execution with access to instance fields and class methods
func (i *Interpreter) createMethodEnv(instance *Instance, callerEnv *Environment) *Environment {
	// Create a new environment enclosed by the caller's environment
	env := NewEnclosedEnvironment(callerEnv)

	// Set instance reference so field access/assignment goes through the instance
	env.instance = instance

	// Bind "this" to the instance for explicit self-reference
	env.Declare("this", instance)

	// Bind SUPER if there's a parent class
	if instance.Class.Parent != nil {
		env.Declare("SUPER", &Super{Instance: instance, Class: instance.Class.Parent})
	}

	// Bind class methods as bound methods (so GetName() works without this. prefix)
	for name, method := range instance.Class.Methods {
		env.Declare(name, &BoundMethod{Instance: instance, Method: method})
	}

	// Also bind inherited methods
	for parent := instance.Class.Parent; parent != nil; parent = parent.Parent {
		for name, method := range parent.Methods {
			// Don't override if already defined (child methods take precedence)
			if _, exists := env.store[name]; !exists {
				env.Declare(name, &BoundMethod{Instance: instance, Method: method})
			}
		}
	}

	return env
}

func (i *Interpreter) extendFunctionEnv(fn *Function, args []Object, params []ast.Parameter, callerEnv *Environment) *Environment {
	env := NewEnclosedEnvironment(fn.Env)

	for idx, param := range params {
		if idx < len(args) {
			if param.ByRef {
				// For BYREF, we need to create a reference
				// This is a simplified implementation
				env.Declare(param.Name, args[idx])
			} else {
				env.Declare(param.Name, args[idx])
			}
		}
	}

	return env
}

func (i *Interpreter) unwrapReturnValue(obj Object) Object {
	if rv, ok := obj.(*ReturnValue); ok {
		return rv.Value
	}
	return obj
}

func (i *Interpreter) evalNewExpression(expr *ast.NewExpression, env *Environment) Object {
	classObj, ok := env.Get(expr.ClassName)
	if !ok {
		return &Error{Message: fmt.Sprintf("class not found: %s", expr.ClassName)}
	}

	class, ok := classObj.(*Class)
	if !ok {
		return &Error{Message: fmt.Sprintf("%s is not a class", expr.ClassName)}
	}

	instance := &Instance{
		Class:  class,
		Fields: make(map[string]Object),
	}

	// Initialize fields from entire class hierarchy
	i.initializeInstanceFields(instance, class)

	// Call constructor (NEW procedure) if exists
	if constructor, ok := class.Methods["NEW"]; ok {
		args := i.evalExpressions(expr.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		// Create method environment with proper instance context
		ctorEnv := i.createMethodEnv(instance, env)

		if proc, ok := constructor.(*Procedure); ok {
			for idx, param := range proc.Parameters {
				if idx < len(args) {
					ctorEnv.Declare(param.Name, args[idx])
				}
			}
			i.evalStatements(proc.Body, ctorEnv)
		}
	}

	return instance
}

// initializeInstanceFields initializes fields from the class and all parent classes
func (i *Interpreter) initializeInstanceFields(instance *Instance, class *Class) {
	// First initialize parent fields
	if class.Parent != nil {
		i.initializeInstanceFields(instance, class.Parent)
	}
	// Then initialize this class's fields (may override parent fields with same name)
	for name := range class.Fields {
		instance.Fields[name] = &Null{}
	}
}

func (i *Interpreter) evalSuperExpression(expr *ast.SuperExpression, env *Environment) Object {
	// SUPER should be available in method context
	if superObj, ok := env.Get("SUPER"); ok {
		return superObj
	}
	return &Error{Message: "SUPER can only be used within a class method"}
}

// IsEOF checks if file is at EOF
func (i *Interpreter) IsEOF(filename string) bool {
	fs, ok := i.files[filename]
	if !ok {
		return true
	}
	return fs.atEOF
}

func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}

func isTruthy(obj Object) bool {
	switch obj := obj.(type) {
	case *Null:
		return false
	case *Boolean:
		return obj.Value
	default:
		return true
	}
}
