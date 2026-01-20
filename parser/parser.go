// Package parser implements the parser for Cambridge Pseudocode
package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/andrinoff/cambridge-lang/ast"
	"github.com/andrinoff/cambridge-lang/lexer"
	"github.com/andrinoff/cambridge-lang/token"
)

// Operator precedence levels
const (
	_ int = iota
	LOWEST
	OR_PREC     // OR
	AND_PREC    // AND
	NOT_PREC    // NOT
	EQUALS      // = <>
	LESSGREATER // < > <= >=
	SUM         // + - &
	PRODUCT     // * / DIV MOD
	PREFIX      // -X NOT X
	CALL        // function(x)
	INDEX       // array[x]
	MEMBER      // object.field
)

var precedences = map[token.Type]int{
	token.OR:        OR_PREC,
	token.AND:       AND_PREC,
	token.EQ:        EQUALS,
	token.NOT_EQ:    EQUALS,
	token.LT:        LESSGREATER,
	token.GT:        LESSGREATER,
	token.LT_EQ:     LESSGREATER,
	token.GT_EQ:     LESSGREATER,
	token.PLUS:      SUM,
	token.MINUS:     SUM,
	token.AMPERSAND: SUM,
	token.ASTERISK:  PRODUCT,
	token.SLASH:     PRODUCT,
	token.DIV:       PRODUCT,
	token.MOD:       PRODUCT,
	token.LPAREN:    CALL,
	token.LBRACKET:  INDEX,
	token.DOT:       MEMBER,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser parses tokens into an AST
type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

// New creates a new parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INTEGER_LIT, p.parseIntegerLiteral)
	p.registerPrefix(token.REAL_LIT, p.parseRealLiteral)
	p.registerPrefix(token.STRING_LIT, p.parseStringLiteral)
	p.registerPrefix(token.CHAR_LIT, p.parseCharLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.NOT, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.NEW, p.parseNewExpression)

	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.DIV, p.parseInfixExpression)
	p.registerInfix(token.MOD, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LT_EQ, p.parseInfixExpression)
	p.registerInfix(token.GT_EQ, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.AMPERSAND, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseArrayAccess)
	p.registerInfix(token.DOT, p.parseMemberAccess)

	// Read two tokens to initialize curToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Errors returns parser errors
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, fmt.Sprintf("line %d, column %d: %s", p.curToken.Line, p.curToken.Column, msg))
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.addError(msg)
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekPrecedence() int {
	if prec, ok := precedences[p.peekToken.Type]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if prec, ok := precedences[p.curToken.Type]; ok {
		return prec
	}
	return LOWEST
}

// skipNewlines advances past any newline tokens
func (p *Parser) skipNewlines() {
	for p.curTokenIs(token.NEWLINE) {
		p.nextToken()
	}
}

// ParseProgram parses the entire program
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		p.skipNewlines()
		if p.curTokenIs(token.EOF) {
			break
		}
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.DECLARE:
		return p.parseDeclareStatement()
	case token.CONSTANT:
		return p.parseConstantStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.CASE:
		return p.parseCaseStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.REPEAT:
		return p.parseRepeatStatement()
	case token.PROCEDURE:
		return p.parseProcedureStatement()
	case token.FUNCTION:
		return p.parseFunctionStatement()
	case token.CALL:
		return p.parseCallStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.INPUT:
		return p.parseInputStatement()
	case token.OUTPUT:
		return p.parseOutputStatement()
	case token.OPENFILE:
		return p.parseOpenFileStatement()
	case token.CLOSEFILE:
		return p.parseCloseFileStatement()
	case token.READFILE:
		return p.parseReadFileStatement()
	case token.WRITEFILE:
		return p.parseWriteFileStatement()
	case token.TYPE:
		return p.parseTypeStatement()
	case token.CLASS:
		return p.parseClassStatement()
	case token.PUBLIC, token.PRIVATE:
		return p.parseAccessModifiedStatement()
	default:
		return p.parseAssignmentOrExpressionStatement()
	}
}

func (p *Parser) parseDeclareStatement() *ast.DeclareStatement {
	stmt := &ast.DeclareStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.COLON) {
		return nil
	}

	p.nextToken()
	stmt.DataType = p.parseDataType()

	return stmt
}

func (p *Parser) parseConstantStatement() *ast.ConstantStatement {
	stmt := &ast.ConstantStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.EQ) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.curToken}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.THEN) {
		return nil
	}

	p.nextToken()
	p.skipNewlines()

	stmt.Consequence = p.parseBlockStatements(token.ELSE, token.ENDIF)

	if p.curTokenIs(token.ELSE) {
		p.nextToken()
		p.skipNewlines()
		stmt.Alternative = p.parseBlockStatements(token.ENDIF)
	}

	return stmt
}

func (p *Parser) parseCaseStatement() *ast.CaseStatement {
	stmt := &ast.CaseStatement{Token: p.curToken}

	if !p.expectPeek(token.OF) {
		return nil
	}

	p.nextToken()
	stmt.Expr = p.parseExpression(LOWEST)

	p.nextToken()
	p.skipNewlines()

	for !p.curTokenIs(token.OTHERWISE) && !p.curTokenIs(token.ENDCASE) && !p.curTokenIs(token.EOF) {
		caseClause := p.parseCaseClause()
		stmt.Cases = append(stmt.Cases, caseClause)
		p.skipNewlines()
	}

	if p.curTokenIs(token.OTHERWISE) {
		if !p.expectPeek(token.COLON) {
			return nil
		}
		p.nextToken()
		p.skipNewlines()
		stmt.Otherwise = p.parseBlockStatements(token.ENDCASE)
	}

	return stmt
}

func (p *Parser) parseCaseClause() ast.CaseClause {
	clause := ast.CaseClause{}

	// Parse values (can be comma-separated or range)
	for {
		value := p.parseExpression(LOWEST)

		// Check for range (value TO value)
		if p.peekTokenIs(token.TO) {
			p.nextToken()
			p.nextToken()
			endValue := p.parseExpression(LOWEST)
			clause.Values = append(clause.Values, &ast.RangeExpression{
				Token: p.curToken,
				Start: value,
				End:   endValue,
			})
		} else {
			clause.Values = append(clause.Values, value)
		}

		if !p.peekTokenIs(token.COMMA) {
			break
		}
		p.nextToken() // consume comma
		p.nextToken()
	}

	if !p.expectPeek(token.COLON) {
		return clause
	}

	p.nextToken()
	p.skipNewlines()

	// Parse body until next case value, OTHERWISE, or ENDCASE
	for !p.curTokenIs(token.OTHERWISE) && !p.curTokenIs(token.ENDCASE) && !p.curTokenIs(token.EOF) {
		// Check if this looks like a new case value (number, string, or identifier followed by colon)
		if p.isStartOfCaseValue() {
			break
		}
		stmt := p.parseStatement()
		if stmt != nil {
			clause.Body = append(clause.Body, stmt)
		}
		p.nextToken()
		p.skipNewlines()
	}

	return clause
}

func (p *Parser) isStartOfCaseValue() bool {
	// Check if current token could be a case value
	switch p.curToken.Type {
	case token.INTEGER_LIT, token.REAL_LIT, token.STRING_LIT, token.CHAR_LIT, token.IDENT:
		// Look ahead to see if there's a colon or TO
		return p.peekTokenIs(token.COLON) || p.peekTokenIs(token.TO) || p.peekTokenIs(token.COMMA)
	}
	return false
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Variable = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Start = p.parseExpression(LOWEST)

	if !p.expectPeek(token.TO) {
		return nil
	}

	p.nextToken()
	stmt.End = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.STEP) {
		p.nextToken()
		p.nextToken()
		stmt.Step = p.parseExpression(LOWEST)
	}

	p.nextToken()
	p.skipNewlines()

	stmt.Body = p.parseBlockStatements(token.NEXT)

	// Expect NEXT variable
	if p.curTokenIs(token.NEXT) {
		p.nextToken() // skip variable name after NEXT
	}

	return stmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{Token: p.curToken}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	p.nextToken()
	p.skipNewlines()

	stmt.Body = p.parseBlockStatements(token.ENDWHILE)

	return stmt
}

func (p *Parser) parseRepeatStatement() *ast.RepeatStatement {
	stmt := &ast.RepeatStatement{Token: p.curToken}

	p.nextToken()
	p.skipNewlines()

	stmt.Body = p.parseBlockStatements(token.UNTIL)

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseProcedureStatement() *ast.ProcedureStatement {
	stmt := &ast.ProcedureStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = p.curToken.Literal

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	stmt.Parameters = p.parseParameters()

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	p.nextToken()
	p.skipNewlines()

	stmt.Body = p.parseBlockStatements(token.ENDPROCEDURE)

	return stmt
}

func (p *Parser) parseFunctionStatement() *ast.FunctionStatement {
	stmt := &ast.FunctionStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = p.curToken.Literal

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	stmt.Parameters = p.parseParameters()

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.RETURNS) {
		return nil
	}

	p.nextToken()
	stmt.ReturnType = p.parseDataType()

	p.nextToken()
	p.skipNewlines()

	stmt.Body = p.parseBlockStatements(token.ENDFUNCTION)

	return stmt
}

func (p *Parser) parseParameters() []ast.Parameter {
	params := []ast.Parameter{}

	if p.peekTokenIs(token.RPAREN) {
		return params
	}

	p.nextToken()

	for {
		param := ast.Parameter{}

		// Check for BYREF or BYVAL
		if p.curTokenIs(token.BYREF) {
			param.ByRef = true
			p.nextToken()
		} else if p.curTokenIs(token.BYVAL) {
			p.nextToken()
		}

		if !p.curTokenIs(token.IDENT) {
			p.addError("expected parameter name")
			return params
		}

		param.Name = p.curToken.Literal

		if !p.expectPeek(token.COLON) {
			return params
		}

		p.nextToken()
		param.DataType = p.parseDataType()

		params = append(params, param)

		if !p.peekTokenIs(token.COMMA) {
			break
		}
		p.nextToken() // consume comma
		p.nextToken()
	}

	return params
}

func (p *Parser) parseCallStatement() *ast.CallStatement {
	stmt := &ast.CallStatement{Token: p.curToken}

	p.nextToken()
	stmt.Name = p.parseExpression(LOWEST)

	// Arguments should already be parsed as part of the call expression
	if call, ok := stmt.Name.(*ast.CallExpression); ok {
		stmt.Name = call.Function
		stmt.Arguments = call.Arguments
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	if p.peekTokenIs(token.NEWLINE) || p.peekTokenIs(token.EOF) {
		return stmt
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseInputStatement() *ast.InputStatement {
	stmt := &ast.InputStatement{Token: p.curToken}

	p.nextToken()
	stmt.Variable = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseOutputStatement() *ast.OutputStatement {
	stmt := &ast.OutputStatement{Token: p.curToken}

	p.nextToken()

	for {
		expr := p.parseExpression(LOWEST)
		stmt.Values = append(stmt.Values, expr)

		if !p.peekTokenIs(token.COMMA) {
			break
		}
		p.nextToken() // consume comma
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseOpenFileStatement() *ast.OpenFileStatement {
	stmt := &ast.OpenFileStatement{Token: p.curToken}

	p.nextToken()
	stmt.Filename = p.parseExpression(LOWEST)

	if !p.expectPeek(token.FOR) {
		return nil
	}

	p.nextToken()
	switch p.curToken.Type {
	case token.READ:
		stmt.Mode = "READ"
	case token.WRITE:
		stmt.Mode = "WRITE"
	case token.APPEND:
		stmt.Mode = "APPEND"
	default:
		p.addError("expected READ, WRITE, or APPEND after FOR")
		return nil
	}

	return stmt
}

func (p *Parser) parseCloseFileStatement() *ast.CloseFileStatement {
	stmt := &ast.CloseFileStatement{Token: p.curToken}

	p.nextToken()
	stmt.Filename = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseReadFileStatement() *ast.ReadFileStatement {
	stmt := &ast.ReadFileStatement{Token: p.curToken}

	p.nextToken()
	stmt.Filename = p.parseExpression(LOWEST)

	if !p.expectPeek(token.COMMA) {
		return nil
	}

	p.nextToken()
	stmt.Variable = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseWriteFileStatement() *ast.WriteFileStatement {
	stmt := &ast.WriteFileStatement{Token: p.curToken}

	p.nextToken()
	stmt.Filename = p.parseExpression(LOWEST)

	if !p.expectPeek(token.COMMA) {
		return nil
	}

	p.nextToken()
	stmt.Data = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseTypeStatement() *ast.TypeStatement {
	stmt := &ast.TypeStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = p.curToken.Literal

	// Check for = (enum or pointer) or newline (record)
	if p.peekTokenIs(token.EQ) {
		p.nextToken()
		p.nextToken()

		if p.curTokenIs(token.CARET) {
			// Pointer type
			p.nextToken()
			stmt.Definition = &ast.PointerType{TargetType: p.parseDataType()}
		} else if p.curTokenIs(token.LPAREN) {
			// Enum type
			stmt.Definition = p.parseEnumType()
		}
	} else {
		// Record type
		p.nextToken()
		p.skipNewlines()
		stmt.Definition = p.parseRecordType()
	}

	return stmt
}

func (p *Parser) parseRecordType() *ast.RecordType {
	record := &ast.RecordType{}

	for !p.curTokenIs(token.ENDTYPE) && !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.DECLARE) {
			p.nextToken()
			if !p.curTokenIs(token.IDENT) {
				p.addError("expected field name")
				return record
			}
			fieldName := p.curToken.Literal

			if !p.expectPeek(token.COLON) {
				return record
			}

			p.nextToken()
			fieldType := p.parseDataType()

			record.Fields = append(record.Fields, ast.RecordField{
				Name:     fieldName,
				DataType: fieldType,
			})
		}
		p.nextToken()
		p.skipNewlines()
	}

	return record
}

func (p *Parser) parseEnumType() *ast.EnumType {
	enum := &ast.EnumType{}

	if !p.curTokenIs(token.LPAREN) {
		return enum
	}

	p.nextToken()

	for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.IDENT) {
			enum.Values = append(enum.Values, p.curToken.Literal)
		}
		p.nextToken()
		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	return enum
}

func (p *Parser) parseClassStatement() *ast.ClassStatement {
	stmt := &ast.ClassStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = p.curToken.Literal

	if p.peekTokenIs(token.INHERITS) {
		p.nextToken()
		if !p.expectPeek(token.IDENT) {
			return nil
		}
		stmt.Parent = p.curToken.Literal
	}

	p.nextToken()
	p.skipNewlines()

	for !p.curTokenIs(token.ENDCLASS) && !p.curTokenIs(token.EOF) {
		member := p.parseStatement()
		if member != nil {
			stmt.Members = append(stmt.Members, member)
		}
		p.nextToken()
		p.skipNewlines()
	}

	return stmt
}

func (p *Parser) parseAccessModifiedStatement() ast.Statement {
	access := p.curToken.Literal
	p.nextToken()

	switch p.curToken.Type {
	case token.PROCEDURE:
		stmt := p.parseProcedureStatement()
		if stmt != nil {
			stmt.Access = access
		}
		return stmt
	case token.FUNCTION:
		stmt := p.parseFunctionStatement()
		if stmt != nil {
			stmt.Access = access
		}
		return stmt
	case token.DECLARE:
		// For class fields
		return p.parseDeclareStatement()
	default:
		p.addError("expected PROCEDURE, FUNCTION, or DECLARE after access modifier")
		return nil
	}
}

func (p *Parser) parseAssignmentOrExpressionStatement() ast.Statement {
	expr := p.parseExpression(LOWEST)

	if p.peekTokenIs(token.ASSIGN) {
		// This is an assignment
		p.nextToken()
		stmt := &ast.AssignmentStatement{Token: p.curToken, Name: expr}
		p.nextToken()
		stmt.Value = p.parseExpression(LOWEST)
		return stmt
	}

	return &ast.ExpressionStatement{Token: p.curToken, Expression: expr}
}

func (p *Parser) parseBlockStatements(endTokens ...token.Type) []ast.Statement {
	statements := []ast.Statement{}

	for !p.isEndToken(endTokens...) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}
		p.nextToken()
		p.skipNewlines()
	}

	return statements
}

func (p *Parser) isEndToken(tokens ...token.Type) bool {
	for _, t := range tokens {
		if p.curTokenIs(t) {
			return true
		}
	}
	return false
}

func (p *Parser) parseDataType() ast.DataType {
	switch p.curToken.Type {
	case token.INTEGER, token.REAL, token.STRING, token.CHAR, token.BOOLEAN, token.DATE:
		return &ast.PrimitiveType{Name: p.curToken.Literal}
	case token.ARRAY:
		return p.parseArrayType()
	case token.CARET:
		p.nextToken()
		return &ast.PointerType{TargetType: p.parseDataType()}
	case token.IDENT:
		return &ast.CustomType{Name: p.curToken.Literal}
	default:
		p.addError(fmt.Sprintf("unexpected data type: %s", p.curToken.Literal))
		return &ast.PrimitiveType{Name: "UNKNOWN"}
	}
}

func (p *Parser) parseArrayType() *ast.ArrayType {
	arrType := &ast.ArrayType{}

	if !p.expectPeek(token.LBRACKET) {
		return arrType
	}

	// Parse dimensions
	for {
		p.nextToken()
		lower, err := strconv.Atoi(p.curToken.Literal)
		if err != nil {
			p.addError("expected integer for array lower bound")
			return arrType
		}

		if !p.expectPeek(token.COLON) {
			return arrType
		}

		p.nextToken()
		upper, err := strconv.Atoi(p.curToken.Literal)
		if err != nil {
			p.addError("expected integer for array upper bound")
			return arrType
		}

		arrType.Dimensions = append(arrType.Dimensions, ast.ArrayDimension{Lower: lower, Upper: upper})

		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		} else {
			break
		}
	}

	if !p.expectPeek(token.RBRACKET) {
		return arrType
	}

	if !p.expectPeek(token.OF) {
		return arrType
	}

	p.nextToken()
	arrType.ElementType = p.parseDataType()

	return arrType
}

// ============ EXPRESSION PARSING ============

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.addError(fmt.Sprintf("no prefix parse function for %s", p.curToken.Type))
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.NEWLINE) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.addError(fmt.Sprintf("could not parse %q as integer", p.curToken.Literal))
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseRealLiteral() ast.Expression {
	lit := &ast.RealLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.addError(fmt.Sprintf("could not parse %q as real", p.curToken.Literal))
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseCharLiteral() ast.Expression {
	return &ast.CharLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{
		Token: p.curToken,
		Value: strings.ToUpper(p.curToken.Literal) == "TRUE",
	}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseExpressionList(end token.Type) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseArrayAccess(array ast.Expression) ast.Expression {
	exp := &ast.ArrayAccess{Token: p.curToken, Array: array}
	exp.Indices = p.parseExpressionList(token.RBRACKET)
	return exp
}

func (p *Parser) parseMemberAccess(object ast.Expression) ast.Expression {
	exp := &ast.MemberAccess{Token: p.curToken, Object: object}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	exp.Member = p.curToken.Literal
	return exp
}

func (p *Parser) parseNewExpression() ast.Expression {
	exp := &ast.NewExpression{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	exp.ClassName = p.curToken.Literal

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}
