package interpreter

// Environment stores variable bindings
type Environment struct {
	store     map[string]Object
	constants map[string]bool
	outer     *Environment
	types     map[string]Object // For TYPE declarations
	instance  *Instance         // For method execution context
}

// NewEnvironment creates a new environment
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	c := make(map[string]bool)
	t := make(map[string]Object)
	return &Environment{store: s, constants: c, outer: nil, types: t}
}

// NewEnclosedEnvironment creates a new environment with an outer scope
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get retrieves a variable from the environment
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.instance != nil {
		// Check instance fields
		if val, found := e.instance.Fields[name]; found {
			return val, true
		}
	}
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set sets a variable in the environment
func (e *Environment) Set(name string, val Object) Object {
	// Check if it's a constant
	if e.isConstant(name) {
		return &Error{Message: "cannot modify constant: " + name}
	}
	e.store[name] = val
	return val
}

// Declare declares a new variable in the current scope
func (e *Environment) Declare(name string, val Object) Object {
	e.store[name] = val
	return val
}

// DeclareConstant declares a constant
func (e *Environment) DeclareConstant(name string, val Object) Object {
	e.store[name] = val
	e.constants[name] = true
	return val
}

// isConstant checks if a name is a constant
func (e *Environment) isConstant(name string) bool {
	if e.constants[name] {
		return true
	}
	if e.outer != nil {
		return e.outer.isConstant(name)
	}
	return false
}

// SetInPlace updates a variable in its original scope
func (e *Environment) SetInPlace(name string, val Object) Object {
	if _, ok := e.store[name]; ok {
		if e.constants[name] {
			return &Error{Message: "cannot modify constant: " + name}
		}
		e.store[name] = val
		return val
	}
	// Check if it's an instance field
	if e.instance != nil {
		if _, isField := e.instance.Fields[name]; isField {
			e.instance.Fields[name] = val
			return val
		}
	}
	if e.outer != nil {
		return e.outer.SetInPlace(name, val)
	}
	// Variable not found, create it in current scope
	e.store[name] = val
	return val
}

// DefineType defines a type
func (e *Environment) DefineType(name string, typ Object) {
	e.types[name] = typ
}

// GetType retrieves a type definition
func (e *Environment) GetType(name string) (Object, bool) {
	typ, ok := e.types[name]
	if !ok && e.outer != nil {
		return e.outer.GetType(name)
	}
	return typ, ok
}
