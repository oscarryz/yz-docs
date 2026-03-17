// Package lexer implements the Yz tokenizer.
//
// It scans UTF-8 source text and produces a stream of tokens as defined in
// the Yz spec §1. Features include:
//   - Unicode identifier scanning
//   - Multi-line strings with both ' and " delimiters
//   - String interpolation via backtick-delimited expressions
//   - Nested /* */ comments
//   - Automatic semicolon insertion (ASI)
//   - Open non-word identifier set
package lexer

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"yz/internal/token"
)

// Lexer tokenizes Yz source code.
type Lexer struct {
	src  []byte // source text
	pos  int    // current byte position
	line int    // current line (1-based)
	col  int    // current column (1-based, byte offset)
	prev token.Type // previous non-whitespace token type (for ASI)
}

// New creates a Lexer for the given source.
func New(src []byte) *Lexer {
	return &Lexer{
		src:  src,
		line: 1,
		col:  1,
	}
}

// Tokenize returns all tokens from the source, ending with EOF.
func Tokenize(src []byte) []token.Token {
	l := New(src)
	var tokens []token.Token
	for {
		tok := l.Next()
		tokens = append(tokens, tok)
		if tok.Type == token.EOF {
			break
		}
	}
	return tokens
}

// Next returns the next token from the input.
func (l *Lexer) Next() token.Token {
	for {
		l.skipWhitespaceAndComments()

		if l.atEnd() {
			// At EOF: maybe insert a trailing semicolon
			if l.shouldInsertSemicolon() {
				l.prev = token.SEMICOLON
				return token.Token{Type: token.SEMICOLON, Literal: "\n", Line: l.line, Col: l.col}
			}
			return token.Token{Type: token.EOF, Line: l.line, Col: l.col}
		}

		ch := l.peekRune()

		// ASI: skipWhitespaceAndComments stopped at a newline that should
		// trigger semicolon insertion. Consume the newline and produce
		// a SEMICOLON token.
		if ch == '\n' && l.shouldInsertSemicolon() {
			l.advanceNewline()
			l.prev = token.SEMICOLON
			return token.Token{Type: token.SEMICOLON, Literal: "\n", Line: l.line - 1, Col: l.col}
		}

		// Record position before consuming the token
		startLine, startCol := l.line, l.col

		var tok token.Token

		switch {
		case isLetter(ch):
			tok = l.scanIdentifier(startLine, startCol)
		case isDigit(ch):
			tok = l.scanNumber(startLine, startCol)
		case ch == '\'' || ch == '"':
			tok = l.scanString(startLine, startCol)
		default:
			tok = l.scanPunctOrNonWord(startLine, startCol)
		}

		l.prev = tok.Type
		return tok
	}
}

// ---------------------------------------------------------------------------
// Whitespace, newlines, and comments
// ---------------------------------------------------------------------------

// skipWhitespaceAndComments skips spaces, tabs, carriage returns, comments,
// and newlines. When a newline is encountered and ASI conditions are met, it
// does NOT skip it — the caller will see the inserted semicolon via the
// newline-triggered ASI path.
func (l *Lexer) skipWhitespaceAndComments() {
	for !l.atEnd() {
		ch := l.peekRune()

		switch {
		case ch == ' ' || ch == '\t' || ch == '\r':
			l.advance()

		case ch == '\n':
			if l.shouldInsertSemicolon() {
				// Don't skip this newline — Next() will produce a SEMICOLON
				// and the newline will be consumed on the next call.
				return
			}
			l.advanceNewline()

		case ch == '/' && l.peekRuneAt(1) == '/':
			l.skipLineComment()

		case ch == '/' && l.peekRuneAt(1) == '*':
			l.skipBlockComment()

		default:
			return
		}
	}
}

// shouldInsertSemicolon returns true if a semicolon should be inserted
// after the previous token when a newline or EOF is reached (spec §1.12).
func (l *Lexer) shouldInsertSemicolon() bool {
	switch l.prev {
	case token.IDENT, token.TYPE_IDENT, token.GENERIC_IDENT,
		token.INT_LIT, token.DECIMAL_LIT, token.STRING_LIT,
		token.NON_WORD,
		token.BREAK, token.CONTINUE, token.RETURN,
		token.RPAREN, token.RBRACKET, token.RBRACE:
		return true
	}
	return false
}

func (l *Lexer) skipLineComment() {
	// skip the '//'
	l.advance()
	l.advance()
	for !l.atEnd() {
		if l.peekRune() == '\n' {
			return // leave the newline for ASI
		}
		l.advance()
	}
}

func (l *Lexer) skipBlockComment() {
	// skip '/*'
	l.advance()
	l.advance()
	depth := 1
	for !l.atEnd() && depth > 0 {
		ch := l.peekRune()
		if ch == '/' && l.peekRuneAt(1) == '*' {
			depth++
			l.advance()
			l.advance()
		} else if ch == '*' && l.peekRuneAt(1) == '/' {
			depth--
			l.advance()
			l.advance()
		} else if ch == '\n' {
			l.advanceNewline()
		} else {
			l.advance()
		}
	}
}

// ---------------------------------------------------------------------------
// Identifiers and keywords
// ---------------------------------------------------------------------------

func (l *Lexer) scanIdentifier(startLine, startCol int) token.Token {
	start := l.pos
	for !l.atEnd() {
		ch := l.peekRune()
		if isLetter(ch) || isDigit(ch) {
			l.advance()
		} else {
			break
		}
	}
	lit := string(l.src[start:l.pos])
	typ := token.LookupIdent(lit)
	return token.Token{Type: typ, Literal: lit, Line: startLine, Col: startCol}
}

// ---------------------------------------------------------------------------
// Number literals
// ---------------------------------------------------------------------------

func (l *Lexer) scanNumber(startLine, startCol int) token.Token {
	start := l.pos
	for !l.atEnd() && isDigit(l.peekRune()) {
		l.advance()
	}

	typ := token.INT_LIT

	// Check for decimal point
	if !l.atEnd() && l.peekRune() == '.' {
		// Look ahead to ensure there's a digit after the dot (not a method call)
		next := l.peekRuneAt(1)
		if next >= '0' && next <= '9' {
			typ = token.DECIMAL_LIT
			l.advance() // consume '.'
			for !l.atEnd() && isDigit(l.peekRune()) {
				l.advance()
			}
		}
	}

	lit := string(l.src[start:l.pos])
	return token.Token{Type: typ, Literal: lit, Line: startLine, Col: startCol}
}

// ---------------------------------------------------------------------------
// String literals (with interpolation)
// ---------------------------------------------------------------------------

func (l *Lexer) scanString(startLine, startCol int) token.Token {
	quote := l.peekRune() // ' or "
	l.advance()           // consume opening quote

	var lit []byte
	lit = append(lit, byte(quote))

	for !l.atEnd() {
		ch := l.peekRune()

		if ch == rune(quote) {
			l.advance()
			lit = append(lit, byte(quote))
			return token.Token{Type: token.STRING_LIT, Literal: string(lit), Line: startLine, Col: startCol}
		}

		if ch == '\\' {
			// Escape sequence: consume \ and the next char
			lit = append(lit, '\\')
			l.advance()
			if !l.atEnd() {
				esc := l.peekRune()
				l.advance()
				if esc == '\n' {
					l.line++
					l.col = 1
				}
				lit = appendRune(lit, esc)
			}
			continue
		}

		if ch == '`' {
			// String interpolation: consume ` expr `
			lit = append(lit, '`')
			l.advance()
			depth := 1
			for !l.atEnd() && depth > 0 {
				ic := l.peekRune()
				if ic == '`' {
					depth--
					lit = append(lit, '`')
					l.advance()
				} else if ic == '\n' {
					lit = append(lit, '\n')
					l.advanceNewline()
				} else {
					lit = appendRune(lit, ic)
					l.advance()
				}
			}
			continue
		}

		if ch == '\n' {
			// Multi-line strings allowed per decision #2
			lit = append(lit, '\n')
			l.advanceNewline()
			continue
		}

		lit = appendRune(lit, ch)
		l.advance()
	}

	// Unterminated string — return what we have
	return token.Token{Type: token.ILLEGAL, Literal: string(lit), Line: startLine, Col: startCol}
}

// ---------------------------------------------------------------------------
// Delimiters, assignment, fat arrow, and non-word identifiers
// ---------------------------------------------------------------------------

func (l *Lexer) scanPunctOrNonWord(startLine, startCol int) token.Token {
	ch := l.peekRune()

	// Single-character delimiters
	switch ch {
	case '{':
		l.advance()
		return token.Token{Type: token.LBRACE, Literal: "{", Line: startLine, Col: startCol}
	case '}':
		l.advance()
		return token.Token{Type: token.RBRACE, Literal: "}", Line: startLine, Col: startCol}
	case '(':
		l.advance()
		return token.Token{Type: token.LPAREN, Literal: "(", Line: startLine, Col: startCol}
	case ')':
		l.advance()
		return token.Token{Type: token.RPAREN, Literal: ")", Line: startLine, Col: startCol}
	case '[':
		l.advance()
		return token.Token{Type: token.LBRACKET, Literal: "[", Line: startLine, Col: startCol}
	case ']':
		l.advance()
		return token.Token{Type: token.RBRACKET, Literal: "]", Line: startLine, Col: startCol}
	case ':':
		l.advance()
		return token.Token{Type: token.COLON, Literal: ":", Line: startLine, Col: startCol}
	case ',':
		l.advance()
		return token.Token{Type: token.COMMA, Literal: ",", Line: startLine, Col: startCol}
	case ';':
		l.advance()
		return token.Token{Type: token.SEMICOLON, Literal: ";", Line: startLine, Col: startCol}
	case '.':
		l.advance()
		return token.Token{Type: token.DOT, Literal: ".", Line: startLine, Col: startCol}
	case '#':
		l.advance()
		return token.Token{Type: token.HASH, Literal: "#", Line: startLine, Col: startCol}
	}

	// '=' can be:
	//   '='  alone → ASSIGN
	//   '=>' → FAT_ARROW
	//   '==', '!=', etc. are non-word (handled by falling through below)
	if ch == '=' {
		next := l.peekRuneAt(1)
		if next == '>' {
			l.advance()
			l.advance()
			return token.Token{Type: token.FAT_ARROW, Literal: "=>", Line: startLine, Col: startCol}
		}
		if !isNonWordChar(next) || l.pos+1 >= len(l.src) {
			// Lone '=' → assignment
			l.advance()
			return token.Token{Type: token.ASSIGN, Literal: "=", Line: startLine, Col: startCol}
		}
		// Otherwise fall through to scan as non-word (e.g. ==)
	}

	// Non-word identifier: sequence of non-word characters
	if isNonWordChar(ch) {
		return l.scanNonWord(startLine, startCol)
	}

	// Unknown character
	l.advance()
	return token.Token{Type: token.ILLEGAL, Literal: string(ch), Line: startLine, Col: startCol}
}

func (l *Lexer) scanNonWord(startLine, startCol int) token.Token {
	start := l.pos
	for !l.atEnd() {
		ch := l.peekRune()
		// Stop at '=>' to avoid consuming the fat arrow as part of a non-word
		if ch == '=' && l.peekRuneAt(1) == '>' {
			// If we haven't consumed anything yet, this is handled by scanPunctOrNonWord
			if l.pos == start {
				break
			}
			break
		}
		if !isNonWordChar(ch) {
			break
		}
		l.advance()
	}
	lit := string(l.src[start:l.pos])
	if lit == "" {
		l.advance()
		return token.Token{Type: token.ILLEGAL, Literal: string(l.src[start:l.pos]), Line: startLine, Col: startCol}
	}
	return token.Token{Type: token.NON_WORD, Literal: lit, Line: startLine, Col: startCol}
}

// ---------------------------------------------------------------------------
// Character classification
// ---------------------------------------------------------------------------

// isLetter returns true for Unicode letters and underscore.
func isLetter(ch rune) bool {
	return ch == '_' || unicode.IsLetter(ch)
}

// isDigit returns true for ASCII digits 0-9.
func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

// isNonWordChar returns true if ch is a valid non-word identifier character.
// Per spec §1.9, a non-word char is anything that is NOT:
//   - a letter or digit
//   - a delimiter: { } ( ) [ ] : , ; . #
//   - whitespace
//   - a quote: ' " `
//   - the lone character = (handled by the caller)
//   - the character > when preceded by = (=> fat arrow, handled by caller)
func isNonWordChar(ch rune) bool {
	if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
		return false
	}
	if unicode.IsSpace(ch) {
		return false
	}
	switch ch {
	case '{', '}', '(', ')', '[', ']',
		':', ',', ';', '.', '#',
		'\'', '"', '`',
		0: // EOF sentinel
		return false
	}
	return true
}

// ---------------------------------------------------------------------------
// Low-level input helpers
// ---------------------------------------------------------------------------

func (l *Lexer) atEnd() bool {
	return l.pos >= len(l.src)
}

// peekRune returns the current rune without consuming it.
func (l *Lexer) peekRune() rune {
	if l.atEnd() {
		return 0
	}
	r, _ := utf8.DecodeRune(l.src[l.pos:])
	return r
}

// peekRuneAt returns the rune at offset positions ahead (in runes, not bytes).
func (l *Lexer) peekRuneAt(offset int) rune {
	pos := l.pos
	for i := 0; i < offset; i++ {
		if pos >= len(l.src) {
			return 0
		}
		_, size := utf8.DecodeRune(l.src[pos:])
		pos += size
	}
	if pos >= len(l.src) {
		return 0
	}
	r, _ := utf8.DecodeRune(l.src[pos:])
	return r
}

// advance moves past the current rune (non-newline).
func (l *Lexer) advance() {
	if l.atEnd() {
		return
	}
	_, size := utf8.DecodeRune(l.src[l.pos:])
	l.pos += size
	l.col += size
}

// advanceNewline moves past a '\n' and updates line/col tracking.
func (l *Lexer) advanceNewline() {
	l.pos++
	l.line++
	l.col = 1
}

// appendRune appends a rune to a byte slice.
func appendRune(buf []byte, r rune) []byte {
	var tmp [utf8.UTFMax]byte
	n := utf8.EncodeRune(tmp[:], r)
	return append(buf, tmp[:n]...)
}

// Ensure Lexer satisfies a useful debug stringer.
var _ fmt.Stringer = token.Token{}
