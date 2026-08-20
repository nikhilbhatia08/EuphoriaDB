// Updated lexer with corrected date parsing logic
package parse

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type SyntaxError struct {
	Message string
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("syntax error: %s", e.Message)
}

// TokenType enumerates all recognized token types.
type TokenType int

const (
	TTDelimiter TokenType = iota
	TTNumber
	TTString
	TTWord
	TTBoolean
	TTDate
	TTOperator
	TTEOF
)

type Token struct {
	Type      TokenType
	StringVal string
	NumVal    int
	BoolVal   bool
	TimeVal   time.Time
	Rune      rune
}

type Lexer struct {
	input        string
	position     int
	currentToken Token
	keywords     map[string]struct{}
}

func NewLexer(s string) *Lexer {
	l := &Lexer{input: s}
	l.initKeywords()
	_ = l.nextToken()
	return l
}

func (l *Lexer) MatchDelim(d rune) bool {
	return l.currentToken.Type == TTDelimiter && l.currentToken.Rune == d
}

func (l *Lexer) MatchIntConstant() bool {
	return l.currentToken.Type == TTNumber
}

func (l *Lexer) MatchStringConstant() bool {
	return l.currentToken.Type == TTString
}

func (l *Lexer) MatchKeyword(w string) bool {
	return l.currentToken.Type == TTWord && l.currentToken.StringVal == strings.ToLower(w)
}

func (l *Lexer) MatchId() bool {
	if l.currentToken.Type != TTWord {
		return false
	}
	_, isKeyword := l.keywords[l.currentToken.StringVal]
	return !isKeyword
}

func (l *Lexer) MatchBooleanConstant() bool {
	return l.currentToken.Type == TTBoolean
}

func (l *Lexer) MatchDateConstant() bool {
	return l.currentToken.Type == TTDate
}

func (l *Lexer) MatchOperator(op string) bool {
	return l.currentToken.Type == TTOperator && l.currentToken.StringVal == op
}

func (l *Lexer) EatDelim(d rune) error {
	if !l.MatchDelim(d) {
		return &SyntaxError{Message: fmt.Sprintf("expected delimiter '%c'", d)}
	}
	return l.nextToken()
}

func (l *Lexer) EatIntConstant() (int, error) {
	if !l.MatchIntConstant() {
		return 0, &SyntaxError{Message: "expected integer constant"}
	}
	val := l.currentToken.NumVal
	if err := l.nextToken(); err != nil {
		return 0, err
	}
	return val, nil
}

func (l *Lexer) EatStringConstant() (string, error) {
	if !l.MatchStringConstant() {
		return "", &SyntaxError{Message: "expected string constant"}
	}
	val := l.currentToken.StringVal
	if err := l.nextToken(); err != nil {
		return "", err
	}
	return val, nil
}

func (l *Lexer) EatKeyword(w string) error {
	if !l.MatchKeyword(w) {
		return &SyntaxError{Message: fmt.Sprintf("expected keyword '%s'", w)}
	}
	return l.nextToken()
}

func (l *Lexer) EatId() (string, error) {
	if !l.MatchId() {
		return "", &SyntaxError{Message: "expected identifier"}
	}
	val := l.currentToken.StringVal
	if err := l.nextToken(); err != nil {
		return "", err
	}
	return val, nil
}

func (l *Lexer) EatBooleanConstant() (bool, error) {
	if !l.MatchBooleanConstant() {
		return false, &SyntaxError{Message: "expected boolean constant (true/false)"}
	}
	val := l.currentToken.BoolVal
	if err := l.nextToken(); err != nil {
		return false, err
	}
	return val, nil
}

func (l *Lexer) EatDateConstant() (time.Time, error) {
	if !l.MatchDateConstant() {
		return time.Time{}, &SyntaxError{Message: "expected date constant"}
	}
	val := l.currentToken.TimeVal
	if err := l.nextToken(); err != nil {
		return time.Time{}, err
	}
	return val, nil
}

func (l *Lexer) EatOperator(op string) error {
	if !l.MatchOperator(op) {
		return &SyntaxError{Message: fmt.Sprintf("expected operator '%s'", op)}
	}
	return l.nextToken()
}

func (l *Lexer) nextToken() error {
	l.skipWhitespace()
	if l.position >= len(l.input) {
		l.currentToken = Token{Type: TTEOF}
		return nil
	}

	r, width := utf8.DecodeRuneInString(l.input[l.position:])

	if isOperatorStart(r) {
		op, err := l.scanOperator()
		if err != nil {
			return err
		}
		l.currentToken = Token{Type: TTOperator, StringVal: op}
		return nil
	}

	switch {
	case r == '\'':
		strVal, err := l.scanString()
		if err != nil {
			return err
		}
		l.currentToken = Token{Type: TTString, StringVal: strVal}
		return nil

	case isDelimiter(r):
		l.position += width
		l.currentToken = Token{Type: TTDelimiter, Rune: r}
		return nil

	case unicode.IsDigit(r):
		start := l.position
		for l.position < len(l.input) {
			r, width = utf8.DecodeRuneInString(l.input[l.position:])
			if !unicode.IsDigit(r) && r != '-' && r != ' ' && r != ':' {
				break
			}
			l.position += width
		}
		tokenStr := l.input[start:l.position]
		if t, err := parseDate(tokenStr); err == nil {
			l.currentToken = Token{Type: TTDate, TimeVal: t}
			return nil
		} else {
			tokenStr = strings.ReplaceAll(tokenStr, " ", "")
			if n, err := strconv.Atoi(tokenStr); err == nil {
				l.currentToken = Token{Type: TTNumber, NumVal: n}
				return nil
			} else {
				return &SyntaxError{Message: fmt.Sprintf("invalid token: '%s'", tokenStr)}
			}
		}

	case unicode.IsLetter(r) || r == '_':
		wordVal := l.scanWord()
		wordValLower := strings.ToLower(wordVal)

		if wordValLower == "true" || wordValLower == "false" {
			boolVal := wordValLower == "true"
			l.currentToken = Token{Type: TTBoolean, BoolVal: boolVal}
			return nil
		}

		l.currentToken = Token{Type: TTWord, StringVal: wordValLower}
		return nil
	}

	return &SyntaxError{
		Message: fmt.Sprintf("unexpected character '%c'", r),
	}
}

func (l *Lexer) scanOperator() (string, error) {
	r, width := utf8.DecodeRuneInString(l.input[l.position:])
	l.position += width

	if l.position < len(l.input) {
		r2, w2 := utf8.DecodeRuneInString(l.input[l.position:])

		if (r == '>' && r2 == '=') || (r == '<' && r2 == '=') ||
			(r == '!' && r2 == '=') || (r == '<' && r2 == '>') {
			l.position += w2
			return string([]rune{r, r2}), nil
		}
	}

	return string(r), nil
}

func (l *Lexer) scanString() (string, error) {
	l.position++
	var sb strings.Builder

	for l.position < len(l.input) {
		r, width := utf8.DecodeRuneInString(l.input[l.position:])
		if r == '\'' {
			l.position += width
			return sb.String(), nil
		}
		sb.WriteRune(r)
		l.position += width
	}
	return "", &SyntaxError{Message: "unterminated string constant"}
}

func (l *Lexer) scanWord() string {
	start := l.position
	for l.position < len(l.input) {
		r, width := utf8.DecodeRuneInString(l.input[l.position:])
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_') {
			break
		}
		l.position += width
	}
	return l.input[start:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.position < len(l.input) {
		r, width := utf8.DecodeRuneInString(l.input[l.position:])
		if !unicode.IsSpace(r) {
			break
		}
		l.position += width
	}
}

func (l *Lexer) initKeywords() {
	keywordList := []string{
		"select", "from", "where", "and",
		"insert", "into", "values", "delete", "update", "set",
		"create", "table", "int", "varchar", "view", "as", "index", "on",
		"group", "by", "having", "order", "asc", "desc",
		"max", "min", "count", "avg", "sum",
	}
	l.keywords = map[string]struct{}{}
	for _, keyword := range keywordList {
		l.keywords[strings.ToLower(keyword)] = struct{}{}
	}
}

func parseDate(s string) (time.Time, error) {
	layouts := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
	}
	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, &SyntaxError{
		Message: fmt.Sprintf("invalid date format: '%s'", s),
	}
}

func isOperatorStart(r rune) bool {
	switch r {
	case '<', '>', '=', '!':
		return true
	default:
		return false
	}
}

func isDelimiter(r rune) bool {
	delimiters := []rune{',', '(', ')', '.', ';', '+', '-'}
	for _, d := range delimiters {
		if r == d {
			return true
		}
	}
	return false
}
