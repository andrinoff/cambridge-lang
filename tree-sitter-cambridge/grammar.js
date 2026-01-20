/// <reference types="tree-sitter-cli/dsl" />
// @ts-check

const KEYWORDS = [
  "DECLARE",
  "CONSTANT",
  "TYPE",
  "ENDTYPE",
  "INTEGER",
  "REAL",
  "STRING",
  "CHAR",
  "BOOLEAN",
  "DATE",
  "ARRAY",
  "OF",
  "TRUE",
  "FALSE",
  "AND",
  "OR",
  "NOT",
  "MOD",
  "DIV",
  "IF",
  "THEN",
  "ELSE",
  "ENDIF",
  "CASE",
  "OTHERWISE",
  "ENDCASE",
  "FOR",
  "TO",
  "STEP",
  "NEXT",
  "WHILE",
  "ENDWHILE",
  "REPEAT",
  "UNTIL",
  "PROCEDURE",
  "ENDPROCEDURE",
  "FUNCTION",
  "ENDFUNCTION",
  "RETURNS",
  "RETURN",
  "CALL",
  "BYVAL",
  "BYREF",
  "INPUT",
  "OUTPUT",
  "OPENFILE",
  "CLOSEFILE",
  "READFILE",
  "WRITEFILE",
  "READ",
  "WRITE",
  "APPEND",
  "CLASS",
  "ENDCLASS",
  "INHERITS",
  "PUBLIC",
  "PRIVATE",
  "NEW",
  "SUPER",
];

module.exports = grammar({
  name: "cambridge",

  extras: ($) => [/\s/, $.comment],

  conflicts: ($) => [[$.procedure_call]],

  rules: {
    source_file: ($) => repeat($._statement),

    _statement: ($) =>
      choice(
        $.declaration,
        $.constant_declaration,
        $.assignment,
        $.output_statement,
        $.input_statement,
        $.if_statement,
        $.case_statement,
        $.for_loop,
        $.while_loop,
        $.repeat_loop,
        $.procedure_declaration,
        $.function_declaration,
        $.procedure_call,
        $.return_statement,
        $.class_declaration,
        $.type_declaration,
        $.file_operation,
      ),

    // Comments
    comment: ($) => token(seq("//", /.*/)),

    // Declarations
    declaration: ($) =>
      seq(
        kw("DECLARE"),
        field("name", $.identifier),
        ":",
        field("type", $.type),
      ),

    constant_declaration: ($) =>
      seq(
        kw("CONSTANT"),
        field("name", $.identifier),
        "=",
        field("value", $._expression),
      ),

    type_declaration: ($) =>
      seq(
        kw("TYPE"),
        field("name", $.identifier),
        repeat($.type_field),
        kw("ENDTYPE"),
      ),

    type_field: ($) => seq(kw("DECLARE"), $.identifier, ":", $.type),

    // Types
    type: ($) => choice($.primitive_type, $.array_type, $.identifier),

    primitive_type: ($) =>
      choice(
        kw("INTEGER"),
        kw("REAL"),
        kw("STRING"),
        kw("CHAR"),
        kw("BOOLEAN"),
        kw("DATE"),
      ),

    array_type: ($) =>
      seq(
        kw("ARRAY"),
        "[",
        $.array_bounds,
        repeat(seq(",", $.array_bounds)),
        "]",
        kw("OF"),
        $.type,
      ),

    array_bounds: ($) => seq($._expression, ":", $._expression),

    // Assignment
    assignment: ($) =>
      seq(
        field("left", $.assignable),
        choice("<-", "←"),
        field("right", $._expression),
      ),

    assignable: ($) => choice($.identifier, $.array_access, $.member_access),

    // Expressions
    _expression: ($) =>
      choice(
        $.number,
        $.string,
        $.char,
        $.boolean,
        $.identifier,
        $.binary_expression,
        $.unary_expression,
        $.parenthesized_expression,
        $.function_call,
        $.array_access,
        $.member_access,
        $.new_expression,
      ),

    binary_expression: ($) =>
      choice(
        prec.left(1, seq($._expression, kw("OR"), $._expression)),
        prec.left(2, seq($._expression, kw("AND"), $._expression)),
        prec.left(3, seq($._expression, choice("=", "<>"), $._expression)),
        prec.left(
          4,
          seq($._expression, choice("<", ">", "<=", ">="), $._expression),
        ),
        prec.left(5, seq($._expression, choice("+", "-", "&"), $._expression)),
        prec.left(
          6,
          seq(
            $._expression,
            choice("*", "/", kw("MOD"), kw("DIV")),
            $._expression,
          ),
        ),
      ),

    unary_expression: ($) => prec(7, seq(kw("NOT"), $._expression)),

    parenthesized_expression: ($) => seq("(", $._expression, ")"),

    // Literals
    number: ($) => token(choice(/\d+\.\d+/, /\d+/)),

    string: ($) => /"[^"]*"/,

    char: ($) => /'[^']'/,

    boolean: ($) => choice(kw("TRUE"), kw("FALSE")),

    identifier: ($) => token(prec(-1, new RegExp(`[a-zA-Z_][a-zA-Z0-9_]*`))),

    // Array access
    array_access: ($) =>
      prec(
        8,
        seq(
          $.identifier,
          "[",
          $._expression,
          repeat(seq(",", $._expression)),
          "]",
        ),
      ),

    // Member access (for OOP)
    member_access: ($) =>
      prec.left(
        8,
        seq(
          choice($.identifier, $.member_access, kw("SUPER")),
          ".",
          $.identifier,
        ),
      ),

    // Function/method call
    function_call: ($) =>
      prec(
        8,
        seq(
          choice($.identifier, $.member_access),
          "(",
          optional(seq($._expression, repeat(seq(",", $._expression)))),
          ")",
        ),
      ),

    // NEW expression for OOP
    new_expression: ($) =>
      seq(
        kw("NEW"),
        $.identifier,
        "(",
        optional(seq($._expression, repeat(seq(",", $._expression)))),
        ")",
      ),

    // Output statement
    output_statement: ($) =>
      prec.right(
        seq(
          kw("OUTPUT"),
          optional(seq($._expression, repeat(seq(",", $._expression)))),
        ),
      ),

    // Input statement
    input_statement: ($) => seq(kw("INPUT"), $.identifier),

    // IF statement
    if_statement: ($) =>
      seq(
        kw("IF"),
        field("condition", $._expression),
        kw("THEN"),
        repeat($._statement),
        optional($.else_clause),
        kw("ENDIF"),
      ),

    else_clause: ($) => seq(kw("ELSE"), repeat($._statement)),

    // CASE statement
    case_statement: ($) =>
      seq(
        kw("CASE"),
        kw("OF"),
        $.identifier,
        repeat($.case_branch),
        optional($.otherwise_branch),
        kw("ENDCASE"),
      ),

    case_branch: ($) =>
      prec.right(
        seq(
          $._expression,
          repeat(seq(",", $._expression)),
          ":",
          repeat($._statement),
        ),
      ),

    otherwise_branch: ($) => seq(kw("OTHERWISE"), ":", repeat($._statement)),

    // FOR loop
    for_loop: ($) =>
      seq(
        kw("FOR"),
        field("variable", $.identifier),
        choice("<-", "←"),
        field("start", $._expression),
        kw("TO"),
        field("end", $._expression),
        optional(seq(kw("STEP"), field("step", $._expression))),
        repeat($._statement),
        kw("NEXT"),
        $.identifier,
      ),

    // WHILE loop
    while_loop: ($) =>
      seq(
        kw("WHILE"),
        field("condition", $._expression),
        repeat($._statement),
        kw("ENDWHILE"),
      ),

    // REPEAT loop
    repeat_loop: ($) =>
      seq(
        kw("REPEAT"),
        repeat($._statement),
        kw("UNTIL"),
        field("condition", $._expression),
      ),

    // Procedure declaration
    procedure_declaration: ($) =>
      seq(
        optional($.visibility),
        kw("PROCEDURE"),
        field("name", $.identifier),
        "(",
        optional($.parameter_list),
        ")",
        repeat($._statement),
        kw("ENDPROCEDURE"),
      ),

    // Function declaration
    function_declaration: ($) =>
      seq(
        optional($.visibility),
        kw("FUNCTION"),
        field("name", $.identifier),
        "(",
        optional($.parameter_list),
        ")",
        kw("RETURNS"),
        field("return_type", $.type),
        repeat($._statement),
        kw("ENDFUNCTION"),
      ),

    parameter_list: ($) => seq($.parameter, repeat(seq(",", $.parameter))),

    parameter: ($) =>
      seq(
        optional(choice(kw("BYVAL"), kw("BYREF"))),
        $.identifier,
        ":",
        $.type,
      ),

    // Procedure call
    procedure_call: ($) =>
      seq(
        kw("CALL"),
        choice($.identifier, $.member_access),
        optional(
          seq(
            "(",
            optional(seq($._expression, repeat(seq(",", $._expression)))),
            ")",
          ),
        ),
      ),

    // Return statement
    return_statement: ($) =>
      prec.right(seq(kw("RETURN"), optional($._expression))),

    // Class declaration
    class_declaration: ($) =>
      seq(
        kw("CLASS"),
        field("name", $.identifier),
        optional(seq(kw("INHERITS"), field("superclass", $.identifier))),
        repeat($._class_member),
        kw("ENDCLASS"),
      ),

    _class_member: ($) =>
      choice($.class_field, $.procedure_declaration, $.function_declaration),

    class_field: ($) => seq(optional($.visibility), $.identifier, ":", $.type),

    visibility: ($) => choice(kw("PUBLIC"), kw("PRIVATE")),

    // File operations
    file_operation: ($) =>
      choice($.openfile, $.closefile, $.readfile, $.writefile),

    openfile: ($) =>
      seq(
        kw("OPENFILE"),
        $._expression,
        kw("FOR"),
        choice(kw("READ"), kw("WRITE"), kw("APPEND")),
      ),

    closefile: ($) => seq(kw("CLOSEFILE"), $._expression),

    readfile: ($) => seq(kw("READFILE"), $._expression, ",", $.identifier),

    writefile: ($) => seq(kw("WRITEFILE"), $._expression, ",", $._expression),
  },
});

// Helper function for case-insensitive keywords
function kw(keyword) {
  return alias(
    token(
      prec(
        1,
        new RegExp(
          keyword
            .split("")
            .map((char) => `[${char.toLowerCase()}${char.toUpperCase()}]`)
            .join(""),
        ),
      ),
    ),
    keyword,
  );
}
