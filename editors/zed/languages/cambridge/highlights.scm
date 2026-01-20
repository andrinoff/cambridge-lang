; Keywords
[
  "DECLARE"
  "CONSTANT"
  "TYPE"
  "ENDTYPE"
  "IF"
  "THEN"
  "ELSE"
  "ENDIF"
  "CASE"
  "OF"
  "OTHERWISE"
  "ENDCASE"
  "FOR"
  "TO"
  "STEP"
  "NEXT"
  "WHILE"
  "ENDWHILE"
  "REPEAT"
  "UNTIL"
  "PROCEDURE"
  "ENDPROCEDURE"
  "FUNCTION"
  "ENDFUNCTION"
  "RETURNS"
  "RETURN"
  "CALL"
  "BYVAL"
  "BYREF"
  "INPUT"
  "OUTPUT"
  "OPENFILE"
  "CLOSEFILE"
  "READFILE"
  "WRITEFILE"
  "READ"
  "WRITE"
  "APPEND"
  "CLASS"
  "ENDCLASS"
  "INHERITS"
  "NEW"
  "SUPER"
  "ARRAY"
] @keyword

; Visibility modifiers
[
  "PUBLIC"
  "PRIVATE"
] @keyword.modifier

; Type keywords
(primitive_type) @type.builtin

; Operators
[
  "AND"
  "OR"
  "NOT"
  "MOD"
  "DIV"
] @keyword.operator

[
  "+"
  "-"
  "*"
  "/"
  "="
  "<>"
  "<"
  ">"
  "<="
  ">="
  "&"
  "<-"
  "←"
] @operator

; Punctuation
[
  "("
  ")"
  "["
  "]"
] @punctuation.bracket

[
  ":"
  ","
  "."
] @punctuation.delimiter

; Literals
(number) @number
(string) @string
(char) @character

(boolean) @boolean

; Comments
(comment) @comment

; Identifiers
(identifier) @variable

; Function/Procedure names
(function_declaration
  name: (identifier) @function)

(procedure_declaration
  name: (identifier) @function)

(function_call
  (identifier) @function.call)

(procedure_call
  (identifier) @function.call)

; Class names
(class_declaration
  name: (identifier) @type)

(class_declaration
  superclass: (identifier) @type)

(new_expression
  (identifier) @type)

; Parameter names
(parameter
  (identifier) @variable.parameter)

; Type annotations (custom types used as types)
(type
  (identifier) @type)

; Constant declarations
(constant_declaration
  name: (identifier) @constant)

; Field access
(member_access
  (identifier) @property)
