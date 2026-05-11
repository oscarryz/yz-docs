// Package token defines the token types produced by the Yz lexer.
package token

import "fmt"

// Type represents the category of a lexical token.
type Type int

const (
	// Special
	ILLEGAL Type = iota
	EOF
	SEMICOLON // ; (explicit or inserted by ASI)

	// Identifiers & literals
	IDENT         // lowercase identifier: name, count, fetch_user
	TYPE_IDENT    // uppercase multi-char: Person, NetworkResponse
	GENERIC_IDENT // single uppercase letter: T, E, K, V
	INT_LIT       // 42
	DECIMAL_LIT   // 3.14
	STRING_LIT    // "hello" or 'hello'
	NON_WORD      // +, -, ==, !=, &&, ||, ?, <<, etc.

	// Keywords
	BREAK    // break
	CONTINUE // continue
	RETURN   // return
	MATCH    // match

	// Delimiters
	LBRACE   // {
	RBRACE   // }
	LPAREN   // (
	RPAREN   // )
	LBRACKET // [
	RBRACKET // ]
	COLON    // :
	ASSIGN   // =
	COMMA    // ,
	DOT      // .
	HASH     // #
	FAT_ARROW // =>
)

var typeNames = [...]string{
	ILLEGAL:       "ILLEGAL",
	EOF:           "EOF",
	SEMICOLON:     "SEMICOLON",
	IDENT:         "IDENT",
	TYPE_IDENT:    "TYPE_IDENT",
	GENERIC_IDENT: "GENERIC_IDENT",
	INT_LIT:       "INT_LIT",
	DECIMAL_LIT:   "DECIMAL_LIT",
	STRING_LIT:    "STRING_LIT",
	NON_WORD:      "NON_WORD",
	BREAK:         "BREAK",
	CONTINUE:      "CONTINUE",
	RETURN:        "RETURN",
	MATCH:         "MATCH",
	LBRACE:        "LBRACE",
	RBRACE:        "RBRACE",
	LPAREN:        "LPAREN",
	RPAREN:        "RPAREN",
	LBRACKET:      "LBRACKET",
	RBRACKET:      "RBRACKET",
	COLON:         "COLON",
	ASSIGN:        "ASSIGN",
	COMMA:         "COMMA",
	DOT:           "DOT",
	HASH:          "HASH",
	FAT_ARROW:     "FAT_ARROW",
}

func (t Type) String() string {
	if int(t) < len(typeNames) {
		return typeNames[t]
	}
	return fmt.Sprintf("Type(%d)", int(t))
}

// keywords maps reserved words to their token types.
var keywords = map[string]Type{
	"break":    BREAK,
	"continue": CONTINUE,
	"return":   RETURN,
	"match":    MATCH,
}

// LookupIdent classifies an identifier string. If the string is a keyword
// it returns the keyword token type. Otherwise it classifies based on casing:
//   - single uppercase letter → GENERIC_IDENT
//   - starts with uppercase   → TYPE_IDENT
//   - otherwise               → IDENT
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	if len(ident) == 1 && ident[0] >= 'A' && ident[0] <= 'Z' {
		return GENERIC_IDENT
	}
	if ident[0] >= 'A' && ident[0] <= 'Z' {
		return TYPE_IDENT
	}
	return IDENT
}

// Token represents a single lexical token.
type Token struct {
	Type    Type
	Literal string // the raw text of the token
	Line    int    // 1-based line number
	Col     int    // 1-based column (byte offset from start of line)
}

func (t Token) String() string {
	if t.Literal != "" {
		return fmt.Sprintf("%s(%q L%d:C%d)", t.Type, t.Literal, t.Line, t.Col)
	}
	return fmt.Sprintf("%s(L%d:C%d)", t.Type, t.Line, t.Col)
}
