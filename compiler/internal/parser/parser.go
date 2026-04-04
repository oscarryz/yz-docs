// Package parser implements the Yz recursive descent parser.
//
// It consumes the token stream produced by the lexer and builds an AST as
// defined in internal/ast. The grammar it implements is specified in
// spec/02-grammar.ebnf.
//
// Key grammar properties:
//   - All non-word method invocations have equal precedence, evaluated L→R.
//   - Semicolons (inserted by ASI or explicit) separate statements in a boc.
//   - Commas separate expressions inside () and [], and also serve as
//     multi-arg separators for non-word method calls (e.g. `? {a}, {b}`).
package parser

import (
	"fmt"
	"strings"

	"yz/internal/ast"
	"yz/internal/lexer"
	"yz/internal/token"
)

// ParseError is a syntax error with source location.
type ParseError struct {
	Msg  string
	Line int
	Col  int
	Len  int // byte length of the offending token (0 → 1 caret)
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error at L%d:C%d: %s", e.Line, e.Col, e.Msg)
}

// Parser holds the state for parsing one source file.
type Parser struct {
	tokens []token.Token
	pos    int // index into tokens
}

// New creates a Parser for the given source bytes.
func New(src []byte) *Parser {
	toks := lexer.Tokenize(src)
	return &Parser{tokens: toks}
}

// ParseFile parses an entire source file and returns the root SourceFile node.
func (p *Parser) ParseFile() (*ast.SourceFile, error) {
	sf := &ast.SourceFile{Pos: p.curPos()}
	for !p.at(token.EOF) {
		p.skipSeps()
		if p.at(token.EOF) {
			break
		}
		node, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		sf.Stmts = append(sf.Stmts, node)
		p.skipSeps()
	}
	return sf, nil
}

// ---------------------------------------------------------------------------
// Statement dispatch
// ---------------------------------------------------------------------------

// parseStatement parses one statement or expression-statement.
//
// The grammar's Statement = Declaration | Assignment | KeywordStmt | Expression.
// We disambiguate at the first token:
//
//   - BREAK / CONTINUE / RETURN / MIX → keyword stmt
//   - IDENT/TYPE_IDENT/GENERIC_IDENT followed by COLON → ShortDecl (possibly multi-name)
//   - IDENT/TYPE_IDENT/GENERIC_IDENT followed by HASH → BocWithSig
//   - IDENT/TYPE_IDENT followed by a type → TypedDecl
//   - IDENT list followed by ASSIGN → multi-assignment
//   - Otherwise → expression (which may itself be an assignment target)
func (p *Parser) parseStatement() (ast.Node, error) {
	tok := p.cur()

	switch tok.Type {
	case token.BREAK:
		p.advance()
		return &ast.BreakStmt{Pos: p.posOf(tok)}, nil

	case token.CONTINUE:
		p.advance()
		return &ast.ContinueStmt{Pos: p.posOf(tok)}, nil

	case token.RETURN:
		return p.parseReturn()

	case token.MIX:
		return p.parseMix()
	}

	// Check for info string: a string literal that immediately precedes a declaration.
	// We parse it as a node and attach it to the following decl if possible.
	if tok.Type == token.STRING_LIT && p.isInfoStringContext() {
		return p.parseInfoStringAndDecl()
	}

	// Check for multi-name decl/assignment: `a, b: ...` or `a, b = ...`
	if p.isMultiNameStart() {
		return p.parseMultiNameStmt()
	}

	// Check for BocWithSig: `name #(...)` or `name #(...) { ... }` or `name #(...) = { ... }`
	if p.isBocWithSigStart() {
		return p.parseBocWithSig()
	}

	// Check for TypedDecl: `name TypeName` or `name TypeName = expr`
	// (distinct from a call `name(args)` by the next token being a type identifier)
	if p.isTypedDeclStart() {
		return p.parseTypedDecl()
	}

	// Parse an expression (which may be a call, literal, etc.)
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	// After an expression, check for assignment `= rhs` or declaration `: rhs`.
	// This handles `name = "Bob"` and `name: "Alice"` as single-name forms.
	switch p.cur().Type {
	case token.COLON:
		// ShortDecl: expr was the LHS identifier
		return p.finishShortDecl(expr)
	case token.ASSIGN:
		// Assignment: expr is the target
		return p.finishAssignment(expr)
	}

	return expr, nil
}

// ---------------------------------------------------------------------------
// Keyword statements
// ---------------------------------------------------------------------------

func (p *Parser) parseReturn() (*ast.ReturnStmt, error) {
	pos := p.posOf(p.cur())
	p.advance() // consume 'return'

	// Bare return if next token ends the statement.
	if p.at(token.SEMICOLON) || p.at(token.EOF) || p.at(token.RBRACE) {
		return &ast.ReturnStmt{Pos: pos}, nil
	}
	val, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return &ast.ReturnStmt{Pos: pos, Value: val}, nil
}

func (p *Parser) parseMix() (*ast.MixStmt, error) {
	pos := p.posOf(p.cur())
	p.advance() // consume 'mix'
	if !p.atAnyIdent() {
		return nil, p.errorf("expected identifier after 'mix'")
	}
	name := p.parseIdent()
	return &ast.MixStmt{Pos: pos, Name: name}, nil
}

// ---------------------------------------------------------------------------
// Declarations
// ---------------------------------------------------------------------------

// isMultiNameStart returns true when we see `ident , ident` (possible multi-decl/assignment).
// Only lowercase IDENT is valid on the LHS of multi-assignment; TYPE_IDENT and GENERIC_IDENT
// (like T, E in a type boc body) are not multi-assignment targets.
func (p *Parser) isMultiNameStart() bool {
	if p.cur().Type != token.IDENT {
		return false
	}
	// Peek ahead: if next-next is COMMA then it's multi-name
	save := p.pos
	p.advance() // skip first ident
	result := p.at(token.COMMA)
	p.pos = save
	return result
}

// isBocWithSigStart returns true when the current token is an ident followed by HASH.
func (p *Parser) isBocWithSigStart() bool {
	if !p.atAnyIdent() {
		return false
	}
	save := p.pos
	p.advance()
	result := p.at(token.HASH)
	p.pos = save
	return result
}

// isTypedDeclStart returns true for `ident TypeIdent` or `ident GenericIdent`.
// We only trigger on TYPE_IDENT and GENERIC_IDENT as the following token —
// not on LBRACKET, because `array[0]` (index access) and `names [String]`
// (array-type decl) cannot be distinguished without deeper lookahead.
// Array-type declarations (`x [T]`) are handled via the expression path.
func (p *Parser) isTypedDeclStart() bool {
	if p.cur().Type != token.IDENT {
		return false
	}
	save := p.pos
	p.advance() // skip ident
	isType := p.at(token.TYPE_IDENT) || p.at(token.GENERIC_IDENT)
	p.pos = save
	return isType
}

// isInfoStringContext returns true when a string literal is in info-string position:
// immediately followed (after semicolons) by a declaration.
func (p *Parser) isInfoStringContext() bool {
	save := p.pos
	p.advance() // skip the string literal
	p.skipSemis()
	// It's an info string if the next token starts a declaration
	result := p.atAnyIdent() || p.at(token.MATCH)
	p.pos = save
	return result
}

// parseMultiNameStmt parses `a, b, c : exprs` or `a, b, c = exprs`.
func (p *Parser) parseMultiNameStmt() (ast.Node, error) {
	pos := p.curPos()
	names, err := p.parseIdentList()
	if err != nil {
		return nil, err
	}

	switch p.cur().Type {
	case token.COLON:
		p.advance()
		values, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		return &ast.ShortDecl{Pos: pos, Names: names, Values: values}, nil

	case token.ASSIGN:
		p.advance()
		values, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		return &ast.Assignment{Pos: pos, Names: names, Values: values}, nil

	default:
		return nil, p.errorf("expected ':' or '=' after identifier list")
	}
}

// finishShortDecl completes a ShortDecl when we already parsed the LHS expression.
func (p *Parser) finishShortDecl(lhs ast.Expr) (*ast.ShortDecl, error) {
	pos := p.curPos()
	id, ok := lhs.(*ast.Ident)
	if !ok {
		return nil, p.errorf("left side of ':' must be an identifier")
	}
	p.advance() // consume ':'
	values, err := p.parseExprList()
	if err != nil {
		return nil, err
	}
	return &ast.ShortDecl{Pos: pos, Names: []*ast.Ident{id}, Values: values}, nil
}

// finishAssignment completes an Assignment when we already parsed the LHS expression.
func (p *Parser) finishAssignment(lhs ast.Expr) (*ast.Assignment, error) {
	pos := p.curPos()
	p.advance() // consume '='
	values, err := p.parseExprList()
	if err != nil {
		return nil, err
	}
	return &ast.Assignment{Pos: pos, Target: lhs, Values: values}, nil
}

// parseTypedDecl parses `name TypeExpr [= expr]`.
func (p *Parser) parseTypedDecl() (*ast.TypedDecl, error) {
	pos := p.curPos()
	name := p.parseIdent()
	typ, err := p.parseTypeExpr()
	if err != nil {
		return nil, err
	}
	var val ast.Expr
	if p.at(token.ASSIGN) {
		p.advance()
		val, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}
	return &ast.TypedDecl{Pos: pos, Name: name, Type: typ, Value: val}, nil
}

// parseBocWithSig parses `name #(params) [body | = body]`.
func (p *Parser) parseBocWithSig() (*ast.BocWithSig, error) {
	pos := p.curPos()
	name := p.parseIdent()

	// consume '#'
	if !p.at(token.HASH) {
		return nil, p.errorf("expected '#'")
	}
	p.advance()

	sig, err := p.parseBocTypeExpr()
	if err != nil {
		return nil, err
	}

	var body *ast.BocLiteral
	var bodyOnly bool
	switch p.cur().Type {
	case token.LBRACE:
		// `name #(params) { body }` — params auto-scoped into body
		body, err = p.parseBocLiteral()
		if err != nil {
			return nil, err
		}
	case token.ASSIGN:
		p.advance()
		// `name #(params) = { body }` — body redeclares its own params
		body, err = p.parseBocLiteral()
		if err != nil {
			return nil, err
		}
		bodyOnly = true
	}

	return &ast.BocWithSig{Pos: pos, Name: name, Sig: sig, Body: body, BodyOnly: bodyOnly}, nil
}

// parseInfoStringAndDecl parses an info string and attaches it to the next decl.
func (p *Parser) parseInfoStringAndDecl() (ast.Node, error) {
	tok := p.cur()
	infoStr := &ast.InfoString{Pos: p.posOf(tok), Value: tok.Literal}
	p.advance()
	p.skipSemis()

	// Parse the following declaration
	node, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	// Attach info string to boc literals or short decls where applicable
	switch n := node.(type) {
	case *ast.ShortDecl:
		if len(n.Values) == 1 {
			if b, ok := n.Values[0].(*ast.BocLiteral); ok {
				b.InfoString = &ast.StringLit{Pos: infoStr.Pos, Value: infoStr.Value}
			}
		}
		return n, nil
	}
	// If we can't attach it, return the info string as a standalone node
	// followed by the declaration — but SourceFile only holds one node per call,
	// so we just return the declaration and lose the info string attachment.
	// The info string is still represented as InfoString node when standalone.
	_ = infoStr
	return node, nil
}

// ---------------------------------------------------------------------------
// Expressions
// ---------------------------------------------------------------------------

// parseExpr parses: UnaryExpr { non_word UnaryExpr }
// All non-word ops have equal precedence, evaluated left-to-right.
func (p *Parser) parseExpr() (ast.Expr, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for p.at(token.NON_WORD) {
		opTok := p.cur()
		p.advance()
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		bin := &ast.BinaryExpr{
			Pos:   p.posOf(opTok),
			Left:  left,
			Op:    opTok.Literal,
			Right: right,
		}
		left = bin

		// `cond ? {trueCase}, {falseCase}` — after parsing `cond ? {trueCase}`,
		// if we see a comma followed by a boc, consume it as the false branch.
		if opTok.Literal == "?" && p.at(token.COMMA) && p.peekIs(token.LBRACE) {
			p.advance() // consume ','
			falseCase, err := p.parseUnary()
			if err != nil {
				return nil, err
			}
			left = &ast.ConditionalExpr{
				Pos:       bin.Pos,
				Cond:      bin.Left,
				TrueCase:  bin.Right,
				FalseCase: falseCase,
			}
		}
	}

	return left, nil
}

func (p *Parser) parseUnary() (ast.Expr, error) {
	if p.at(token.NON_WORD) && p.cur().Literal == "-" {
		pos := p.posOf(p.cur())
		p.advance()
		operand, err := p.parsePostfix()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{Pos: pos, Op: "-", Operand: operand}, nil
	}
	return p.parsePostfix()
}

// parsePostfix parses: PrimaryExpr { Postfix }
// Postfix = "." Ident | "(" ArgList ")" | "[" Expr "]"
func (p *Parser) parsePostfix() (ast.Expr, error) {
	base, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	for {
		switch p.cur().Type {
		case token.DOT:
			p.advance()
			if !p.atAnyIdent() {
				return nil, p.errorf("expected identifier after '.'")
			}
			member := p.parseIdent()
			base = &ast.MemberExpr{Pos: member.Pos, Object: base, Member: member}

		case token.LPAREN:
			pos := p.curPos()
			args, err := p.parseArgList()
			if err != nil {
				return nil, err
			}
			base = &ast.CallExpr{Pos: pos, Callee: base, Args: args}

		case token.LBRACKET:
			pos := p.curPos()
			p.advance() // consume '['
			idx, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if err := p.expect(token.RBRACKET); err != nil {
				return nil, err
			}
			base = &ast.IndexExpr{Pos: pos, Object: base, Index: idx}

		default:
			return base, nil
		}
	}
}

// parsePrimary parses a primary expression.
func (p *Parser) parsePrimary() (ast.Expr, error) {
	tok := p.cur()

	switch tok.Type {
	case token.INT_LIT:
		p.advance()
		return &ast.IntLit{Pos: p.posOf(tok), Value: tok.Literal}, nil

	case token.DECIMAL_LIT:
		p.advance()
		return &ast.DecimalLit{Pos: p.posOf(tok), Value: tok.Literal}, nil

	case token.STRING_LIT:
		p.advance()
		if segs := splitStringInterp(tok.Literal); segs != nil {
			return p.buildInterpExpr(p.posOf(tok), segs)
		}
		return &ast.StringLit{Pos: p.posOf(tok), Value: tok.Literal}, nil

	case token.IDENT, token.TYPE_IDENT, token.GENERIC_IDENT:
		return p.parseIdent(), nil

	case token.LBRACE:
		return p.parseBocLiteral()

	case token.LBRACKET:
		return p.parseArrayOrDict()

	case token.MATCH:
		return p.parseMatch()

	case token.LPAREN:
		return p.parseGroup()

	default:
		return nil, p.errorf("unexpected token %v in expression", tok)
	}
}

// parseGroup parses `( ExprList )`. Single element = GroupExpr.
func (p *Parser) parseGroup() (ast.Expr, error) {
	pos := p.curPos()
	p.advance() // consume '('
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if err := p.expect(token.RPAREN); err != nil {
		return nil, err
	}
	return &ast.GroupExpr{Pos: pos, Expr: expr}, nil
}

// ---------------------------------------------------------------------------
// Boc literals and bodies
// ---------------------------------------------------------------------------

// parseBocLiteral parses `{ BocBody }`.
func (p *Parser) parseBocLiteral() (*ast.BocLiteral, error) {
	pos := p.curPos()
	if err := p.expect(token.LBRACE); err != nil {
		return nil, err
	}
	boc := &ast.BocLiteral{Pos: pos}
	p.skipSeps()
	for !p.at(token.RBRACE) && !p.at(token.EOF) {
		node, err := p.parseBocElement()
		if err != nil {
			return nil, err
		}
		boc.Elements = append(boc.Elements, node)
		p.skipSeps()
	}
	if err := p.expect(token.RBRACE); err != nil {
		return nil, err
	}
	return boc, nil
}

// parseBocElement parses one element inside a boc: Declaration | Assignment |
// KeywordStmt | Expression | VariantDef.
func (p *Parser) parseBocElement() (ast.Node, error) {
	// VariantDef inside a type boc: TYPE_IDENT '(' params ')'
	if p.at(token.TYPE_IDENT) && p.peekIs(token.LPAREN) {
		return p.parseVariantDef()
	}
	return p.parseStatement()
}

// ---------------------------------------------------------------------------
// Arrays and dicts
// ---------------------------------------------------------------------------

// parseArrayOrDict parses `[...]` which may be:
//   - Array literal: [1, 2, 3]
//   - Dict literal: ["a": 1, "b": 2]
//   - Empty array: [Int]()
//   - Empty dict: [K:V]()
func (p *Parser) parseArrayOrDict() (ast.Expr, error) {
	pos := p.curPos()
	p.advance() // consume '['

	// Empty array/dict with type: [Type]() or [K:V]()
	if p.at(token.TYPE_IDENT) || p.at(token.GENERIC_IDENT) {
		// Could be type or expression. Look ahead: if we see ] followed by ()
		// it's the empty form. Otherwise it's a dict/array literal.
		save := p.pos
		keyType, err := p.parseTypeExpr()
		if err == nil {
			if p.at(token.COLON) {
				// [K:V]()
				p.advance()
				valType, err := p.parseTypeExpr()
				if err != nil {
					return nil, err
				}
				if err := p.expect(token.RBRACKET); err != nil {
					return nil, err
				}
				// consume ()
				if err := p.expect(token.LPAREN); err != nil {
					return nil, err
				}
				if err := p.expect(token.RPAREN); err != nil {
					return nil, err
				}
				return &ast.DictLiteral{Pos: pos, KeyType: keyType, ValType: valType}, nil
			}
			if p.at(token.RBRACKET) {
				p.advance() // consume ]
				if p.at(token.LPAREN) {
					p.advance() // consume (
					if err := p.expect(token.RPAREN); err != nil {
						return nil, err
					}
					return &ast.ArrayLiteral{Pos: pos, ElemType: keyType}, nil
				}
				// Was not an empty form — backtrack and parse as expression.
				// The ] already consumed — this is tricky. Re-parse as expression.
				// Actually this case means [TypeIdent] which would be an array
				// type but not the empty form. Unusual. Fall back.
			}
		}
		// Backtrack and try as expression list.
		p.pos = save
	}

	// Try to detect empty `[]`
	if p.at(token.RBRACKET) {
		p.advance()
		return &ast.ArrayLiteral{Pos: pos}, nil
	}

	// Parse first element/key to determine array vs dict.
	first, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	if p.at(token.COLON) {
		// Dict literal: [key: val, ...]
		p.advance()
		val, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		entries := []*ast.DictEntry{{Pos: pos, Key: first, Value: val}}
		for p.at(token.COMMA) {
			p.advance()
			if p.at(token.RBRACKET) {
				break
			}
			k, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if err := p.expect(token.COLON); err != nil {
				return nil, err
			}
			v, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			entries = append(entries, &ast.DictEntry{Pos: p.curPos(), Key: k, Value: v})
		}
		if err := p.expect(token.RBRACKET); err != nil {
			return nil, err
		}
		return &ast.DictLiteral{Pos: pos, Entries: entries}, nil
	}

	// Array literal: [first, ...]
	elements := []ast.Expr{first}
	for p.at(token.COMMA) {
		p.advance()
		if p.at(token.RBRACKET) {
			break
		}
		el, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		elements = append(elements, el)
	}
	if err := p.expect(token.RBRACKET); err != nil {
		return nil, err
	}
	return &ast.ArrayLiteral{Pos: pos, Elements: elements}, nil
}

// ---------------------------------------------------------------------------
// Match expressions
// ---------------------------------------------------------------------------

// parseMatch parses `match [expr] { arm }, { arm }, ...`
func (p *Parser) parseMatch() (*ast.MatchExpr, error) {
	pos := p.curPos()
	p.advance() // consume 'match'

	var subject ast.Expr
	if !p.at(token.LBRACE) {
		// Has a subject expression
		var err error
		subject, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}

	// Skip any semicolons (ASI) between subject and the arm list.
	p.skipSemis()

	arms, err := p.parseConditionalBocList()
	if err != nil {
		return nil, err
	}

	return &ast.MatchExpr{Pos: pos, Subject: subject, Arms: arms}, nil
}

// parseConditionalBocList parses `{ arm }, { arm }, ...`
func (p *Parser) parseConditionalBocList() ([]*ast.ConditionalBoc, error) {
	var arms []*ast.ConditionalBoc
	for {
		if !p.at(token.LBRACE) {
			break
		}
		arm, err := p.parseConditionalBoc()
		if err != nil {
			return nil, err
		}
		arms = append(arms, arm)
		if p.at(token.COMMA) {
			p.advance()
		} else {
			break
		}
	}
	return arms, nil
}

// parseConditionalBoc parses `{ [condition =>] BocBody }`.
func (p *Parser) parseConditionalBoc() (*ast.ConditionalBoc, error) {
	pos := p.curPos()
	if err := p.expect(token.LBRACE); err != nil {
		return nil, err
	}
	p.skipSemis()

	arm := &ast.ConditionalBoc{Pos: pos}

	// Detect condition: if there's an expression followed by `=>`
	if !p.at(token.RBRACE) && p.hasConditionArrow() {
		cond, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		arm.Condition = cond
		if err := p.expect(token.FAT_ARROW); err != nil {
			return nil, err
		}
		p.skipSemis()
	}

	// Parse body elements
	for !p.at(token.RBRACE) && !p.at(token.EOF) {
		node, err := p.parseBocElement()
		if err != nil {
			return nil, err
		}
		arm.Body = append(arm.Body, node)
		p.skipSemis()
	}
	if err := p.expect(token.RBRACE); err != nil {
		return nil, err
	}
	return arm, nil
}

// hasConditionArrow does a speculative scan to detect `expr =>` inside a boc arm.
// It looks ahead until it finds `=>` at depth 0 (not nested in () [] {}).
func (p *Parser) hasConditionArrow() bool {
	depth := 0
	for i := p.pos; i < len(p.tokens); i++ {
		switch p.tokens[i].Type {
		case token.LPAREN, token.LBRACKET, token.LBRACE:
			depth++
		case token.RPAREN, token.RBRACKET:
			depth--
		case token.RBRACE:
			if depth == 0 {
				return false // reached end of arm without finding =>
			}
			depth--
		case token.FAT_ARROW:
			return depth == 0
		case token.SEMICOLON:
			if depth == 0 {
				return false // new statement before =>
			}
		case token.EOF:
			return false
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// Type expressions
// ---------------------------------------------------------------------------

// parseTypeExpr parses a type annotation.
func (p *Parser) parseTypeExpr() (ast.TypeExpr, error) {
	switch p.cur().Type {
	case token.TYPE_IDENT, token.GENERIC_IDENT:
		return p.parseSimpleType()
	case token.LBRACKET:
		return p.parseArrayOrDictType()
	case token.HASH:
		p.advance()
		return p.parseBocTypeExpr()
	default:
		return nil, p.errorf("expected type expression, got %v", p.cur())
	}
}

func (p *Parser) parseSimpleType() (*ast.SimpleTypeExpr, error) {
	tok := p.cur()
	pos := p.posOf(tok)
	name := tok.Literal
	tokType := tok.Type
	p.advance()

	// Generic application: Type(T, U)
	var typeArgs []ast.TypeExpr
	if p.at(token.LPAREN) {
		p.advance()
		for !p.at(token.RPAREN) && !p.at(token.EOF) {
			arg, err := p.parseTypeExpr()
			if err != nil {
				return nil, err
			}
			typeArgs = append(typeArgs, arg)
			if p.at(token.COMMA) {
				p.advance()
			}
		}
		if err := p.expect(token.RPAREN); err != nil {
			return nil, err
		}
	}
	return &ast.SimpleTypeExpr{Pos: pos, Name: name, TokType: tokType, TypeArgs: typeArgs}, nil
}

func (p *Parser) parseArrayOrDictType() (ast.TypeExpr, error) {
	pos := p.curPos()
	p.advance() // consume '['
	keyType, err := p.parseTypeExpr()
	if err != nil {
		return nil, err
	}
	if p.at(token.COLON) {
		p.advance()
		valType, err := p.parseTypeExpr()
		if err != nil {
			return nil, err
		}
		if err := p.expect(token.RBRACKET); err != nil {
			return nil, err
		}
		return &ast.DictTypeExpr{Pos: pos, KeyType: keyType, ValType: valType}, nil
	}
	if err := p.expect(token.RBRACKET); err != nil {
		return nil, err
	}
	return &ast.ArrayTypeExpr{Pos: pos, ElemType: keyType}, nil
}

// parseBocTypeExpr parses `( params )` after `#`.
func (p *Parser) parseBocTypeExpr() (*ast.BocTypeExpr, error) {
	pos := p.curPos()
	if err := p.expect(token.LPAREN); err != nil {
		return nil, err
	}
	var params []*ast.BocParam
	for !p.at(token.RPAREN) && !p.at(token.EOF) {
		param, err := p.parseBocParam()
		if err != nil {
			return nil, err
		}
		params = append(params, param)
		if p.at(token.COMMA) {
			p.advance()
		}
	}
	if err := p.expect(token.RPAREN); err != nil {
		return nil, err
	}
	return &ast.BocTypeExpr{Pos: pos, Params: params}, nil
}

// parseBocParam parses one parameter in a boc signature.
// Forms:
//
//	ident TypeExpr [= expr]    — named param
//	ident : expr               — ShortDecl param (type inferred from default)
//	TypeExpr                   — anonymous param (return type)
//	TYPE_IDENT ( params )      — variant constructor param
func (p *Parser) parseBocParam() (*ast.BocParam, error) {
	pos := p.curPos()

	// Variant constructor: TYPE_IDENT '('
	if p.at(token.TYPE_IDENT) && p.peekIs(token.LPAREN) {
		v, err := p.parseVariantDef()
		if err != nil {
			return nil, err
		}
		return &ast.BocParam{Pos: pos, Variant: v}, nil
	}

	// Anonymous type param (return type): TYPE_IDENT or GENERIC_IDENT not followed by type
	if (p.at(token.TYPE_IDENT) || p.at(token.GENERIC_IDENT)) && !p.peekIsType() {
		typ, err := p.parseTypeExpr()
		if err != nil {
			return nil, err
		}
		return &ast.BocParam{Pos: pos, Type: typ}, nil
	}

	// Named param: ident TypeExpr [= expr]  OR  ShortDecl param: ident : expr
	if p.at(token.IDENT) {
		label := p.cur().Literal
		p.advance()

		// ShortDecl-style param: name : expr  (type inferred from default value)
		if p.at(token.COLON) {
			p.advance() // consume ':'
			def, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			return &ast.BocParam{Pos: pos, Label: label, Default: def}, nil
		}

		typ, err := p.parseTypeExpr()
		if err != nil {
			return nil, err
		}
		var def ast.Expr
		if p.at(token.ASSIGN) {
			p.advance()
			def, err = p.parseExpr()
			if err != nil {
				return nil, err
			}
		}
		return &ast.BocParam{Pos: pos, Label: label, Type: typ, Default: def}, nil
	}

	// Type-only param with complex type (array, dict, boc)
	if p.at(token.LBRACKET) || p.at(token.HASH) {
		typ, err := p.parseTypeExpr()
		if err != nil {
			return nil, err
		}
		return &ast.BocParam{Pos: pos, Type: typ}, nil
	}

	return nil, p.errorf("expected boc parameter, got %v", p.cur())
}

// peekIsType returns true when the token AFTER the current one is a type-starting token.
func (p *Parser) peekIsType() bool {
	if p.pos+1 >= len(p.tokens) {
		return false
	}
	next := p.tokens[p.pos+1]
	switch next.Type {
	case token.TYPE_IDENT, token.GENERIC_IDENT, token.LBRACKET, token.HASH:
		return true
	}
	return false
}

// ---------------------------------------------------------------------------
// Variant definitions
// ---------------------------------------------------------------------------

// parseVariantDef parses `TYPE_IDENT ( params )` inside a type boc body or
// boc signature.
func (p *Parser) parseVariantDef() (*ast.VariantDef, error) {
	pos := p.curPos()
	name := p.cur().Literal
	p.advance() // consume TYPE_IDENT
	if err := p.expect(token.LPAREN); err != nil {
		return nil, err
	}
	var params []*ast.BocParam
	for !p.at(token.RPAREN) && !p.at(token.EOF) {
		param, err := p.parseBocParam()
		if err != nil {
			return nil, err
		}
		params = append(params, param)
		if p.at(token.COMMA) {
			p.advance()
		}
	}
	if err := p.expect(token.RPAREN); err != nil {
		return nil, err
	}
	return &ast.VariantDef{Pos: pos, Name: name, Params: params}, nil
}

// ---------------------------------------------------------------------------
// Argument lists
// ---------------------------------------------------------------------------

// parseArgList parses `( [arg, arg, ...] )`.
func (p *Parser) parseArgList() ([]*ast.Argument, error) {
	p.advance() // consume '('
	var args []*ast.Argument
	for !p.at(token.RPAREN) && !p.at(token.EOF) {
		arg, err := p.parseArgument()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		if p.at(token.COMMA) {
			p.advance()
		}
	}
	if err := p.expect(token.RPAREN); err != nil {
		return nil, err
	}
	return args, nil
}

// parseArgument parses one argument: `[label:] expr`.
func (p *Parser) parseArgument() (*ast.Argument, error) {
	pos := p.curPos()

	// Named argument: ident ':'
	if p.at(token.IDENT) && p.peekIs(token.COLON) {
		label := p.cur().Literal
		p.advance() // consume ident
		p.advance() // consume ':'
		val, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return &ast.Argument{Pos: pos, Label: label, Value: val}, nil
	}

	val, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return &ast.Argument{Pos: pos, Value: val}, nil
}

// ---------------------------------------------------------------------------
// Identifier and list helpers
// ---------------------------------------------------------------------------

func (p *Parser) parseIdent() *ast.Ident {
	tok := p.cur()
	p.advance()
	return &ast.Ident{Pos: p.posOf(tok), Name: tok.Literal, TokType: tok.Type}
}

func (p *Parser) parseIdentList() ([]*ast.Ident, error) {
	if !p.atAnyIdent() {
		return nil, p.errorf("expected identifier")
	}
	ids := []*ast.Ident{p.parseIdent()}
	for p.at(token.COMMA) {
		p.advance()
		if !p.atAnyIdent() {
			return nil, p.errorf("expected identifier after ','")
		}
		ids = append(ids, p.parseIdent())
	}
	return ids, nil
}

func (p *Parser) parseExprList() ([]ast.Expr, error) {
	first, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	exprs := []ast.Expr{first}
	for p.at(token.COMMA) {
		p.advance()
		e, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, e)
	}
	return exprs, nil
}

// ---------------------------------------------------------------------------
// Token navigation helpers
// ---------------------------------------------------------------------------

func (p *Parser) cur() token.Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return token.Token{Type: token.EOF}
}

func (p *Parser) advance() {
	if p.pos < len(p.tokens) {
		p.pos++
	}
}

func (p *Parser) at(typ token.Type) bool {
	return p.cur().Type == typ
}

func (p *Parser) atAnyIdent() bool {
	t := p.cur().Type
	return t == token.IDENT || t == token.TYPE_IDENT || t == token.GENERIC_IDENT
}

func (p *Parser) peekIs(typ token.Type) bool {
	if p.pos+1 < len(p.tokens) {
		return p.tokens[p.pos+1].Type == typ
	}
	return false
}

func (p *Parser) skipSemis() {
	for p.at(token.SEMICOLON) {
		p.advance()
	}
}

// skipSeps skips any mix of semicolons and commas that act as statement
// separators at the top level of a SourceFile or BocLiteral body.
// Commas appear here in:
//   - `T, E,` generic param lists inside type boc bodies
//   - `{ arm }, { arm }` match arm / conditional boc lists
//   - `flag ? { a }, { b }` multi-arg non-word invocations (second boc is a
//     standalone expression following the first BinaryExpr)
func (p *Parser) skipSeps() {
	for p.at(token.SEMICOLON) || p.at(token.COMMA) {
		p.advance()
	}
}

func (p *Parser) expect(typ token.Type) error {
	if !p.at(typ) {
		return p.errorf("expected %v, got %v (%q)", typ, p.cur().Type, p.cur().Literal)
	}
	p.advance()
	return nil
}

func (p *Parser) curPos() ast.Pos {
	return p.posOf(p.cur())
}

func (p *Parser) posOf(tok token.Token) ast.Pos {
	return ast.Pos{Line: tok.Line, Col: tok.Col}
}

// ---------------------------------------------------------------------------
// String interpolation helpers
// ---------------------------------------------------------------------------

// interpSegment is one raw slice from splitting a STRING_LIT.
type interpSegment struct {
	isExpr  bool
	content string // raw text (escape sequences intact) or expression source
}

// splitStringInterp splits a raw STRING_LIT token value (including outer quotes)
// into alternating text and expression segments. Returns nil when there is no
// backtick interpolation in the literal.
func splitStringInterp(raw string) []interpSegment {
	if len(raw) < 2 {
		return nil
	}
	inner := raw[1 : len(raw)-1] // strip outer quotes

	var parts []interpSegment
	var cur strings.Builder
	hasInterp := false
	i := 0

	for i < len(inner) {
		ch := inner[i]

		// Preserve escape sequences verbatim.
		if ch == '\\' && i+1 < len(inner) {
			cur.WriteByte(ch)
			i++
			cur.WriteByte(inner[i])
			i++
			continue
		}

		if ch == '`' {
			hasInterp = true
			// Flush accumulated text.
			parts = append(parts, interpSegment{isExpr: false, content: cur.String()})
			cur.Reset()
			i++ // skip opening backtick
			// Collect expression content up to the closing backtick.
			for i < len(inner) && inner[i] != '`' {
				cur.WriteByte(inner[i])
				i++
			}
			parts = append(parts, interpSegment{isExpr: true, content: cur.String()})
			cur.Reset()
			if i < len(inner) {
				i++ // skip closing backtick
			}
			continue
		}

		cur.WriteByte(ch)
		i++
	}

	if !hasInterp {
		return nil
	}
	// Flush trailing text (may be empty — kept so we don't miss it).
	parts = append(parts, interpSegment{isExpr: false, content: cur.String()})
	return parts
}

// buildInterpExpr converts raw interpSegments into an InterpolatedStringExpr node.
// Expression segments are parsed using a sub-parser.
func (p *Parser) buildInterpExpr(pos ast.Pos, segs []interpSegment) (*ast.InterpolatedStringExpr, error) {
	node := &ast.InterpolatedStringExpr{Pos: pos}
	for _, seg := range segs {
		if seg.isExpr {
			sub := New([]byte(seg.content))
			e, err := sub.parseExpr()
			if err != nil {
				return nil, fmt.Errorf("string interpolation: %w", err)
			}
			node.Parts = append(node.Parts, ast.InterpPart{IsExpr: true, Expr: e})
		} else {
			// Skip empty text segments to keep the Part list minimal.
			if seg.content == "" {
				continue
			}
			node.Parts = append(node.Parts, ast.InterpPart{IsExpr: false, Text: seg.content})
		}
	}
	return node, nil
}

func (p *Parser) errorf(format string, args ...any) *ParseError {
	tok := p.cur()
	return &ParseError{
		Msg:  fmt.Sprintf(strings.TrimSpace(format), args...),
		Line: tok.Line,
		Col:  tok.Col,
		Len:  len(tok.Literal),
	}
}
