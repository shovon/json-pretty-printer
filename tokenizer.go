package main

import (
	"bytes"
	"fmt"
	"text/scanner"
	"unicode/utf8"
)

const (
	// JSONOpenBrace represents an open brace
	JSONOpenBrace = iota
	// JSONCloseBrace represents a close brace
	JSONCloseBrace
	// JSONOpenSquareBracket represents an open square bracket
	JSONOpenSquareBracket
	// JSONCloseSquareBracket represents a close square bracket
	JSONCloseSquareBracket
	// JSONColon represents a colon
	JSONColon
	// JSONComma represents a comma
	JSONComma
	// JSONIdentifier represents an identifier (there are only 3 legal identifiers
	// in go, techincally keywords)
	JSONIdentifier
	// JSONString represents a string
	JSONString
	// JSONNumber represents a number
	JSONNumber
	// JSONWhitespace represents whitespace in JSON
	JSONWhitespace
	// JSONEnd represents the end of the JSON stream
	JSONEnd
)

// Tokenizer represents a tokenizer for a CharStrema
type Tokenizer struct {
	scanner      *scanner.Scanner
	peekedTokens []Token
}

// Token represents a token
type Token struct {
	Content   string
	TokenType int
	Position  scanner.Position
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s", t.Content, t.Position.String())
}

// NewTokenizer initializes a new instance of a tokenizer.
func NewTokenizer(reader *scanner.Scanner) Tokenizer {
	return Tokenizer{reader, []Token{}}
}

func (t *Tokenizer) scanString() (string, error) {
	hasBackslash := false
	var buffer bytes.Buffer
	buffer.WriteRune('"')
	for {
		r := t.scanner.Next()
		switch {
		case r == '\\' && !hasBackslash:
			hasBackslash = true
		case r == '"' && !hasBackslash:
			buffer.WriteRune(r)
			return buffer.String(), nil // TODO: handle errors
		}
		hasBackslash = false
		buffer.WriteRune(r)
	}
}

// Finished represents whether or not the scanning is done.
func (t *Tokenizer) Finished() bool {
	return t.scanner.Peek() == scanner.EOF
}

func (t *Tokenizer) scanIdentifier() (string, error) {
	var buffer bytes.Buffer
	for {
		r := t.scanner.Peek()
		if !isIdentifierRune(r) {
			return buffer.String(), nil
		}
		buffer.WriteRune(r)
		t.scanner.Next()
	}
}

// TODO: perhaps move all the digit parsing code to a seperate file for the sake
//   of abstraction.

func isRuneDigit(r rune) bool {
	s := string(r)
	c := s[0]
	if utf8.RuneLen(r) > 1 {
		return false
	}
	return c >= 48 && c <= 57
}

func (t *Tokenizer) scanDigits() (string, error) {
	var buffer bytes.Buffer
	for {
		r := t.scanner.Peek()
		if !isRuneDigit(r) {
			return buffer.String(), nil
		}
		buffer.WriteRune(r)
		t.scanner.Next()
	}
}

func (t *Tokenizer) scanNumber(initial rune) (string, error) {
	var buffer bytes.Buffer
	buffer.WriteRune(initial)

	// Parse the second character (and onwards).
	//
	// Note: hopefully the next code seems like it may fail if the next character
	//  seems to be an illegal character.
	if initial != '0' && isRuneDigit(initial) {
		s, e := t.scanDigits()
		if e != nil {
			return "", e
		}
		buffer.WriteString(s)
	} else if initial == '-' {
		peek := t.scanner.Peek()
		if peek == '0' {
			buffer.WriteRune(peek)
			t.scanDigits()
		} else {
			s, e := t.scanDigits()
			if e != nil {
				return "", e
			}
			buffer.WriteString(s)
		}
	}

	if t.scanner.Peek() == '.' {
		buffer.WriteRune(t.scanner.Next())
		s, e := t.scanDigits()
		if e != nil {
			return "", e
		}
		buffer.WriteString(s)
	}

	if t.scanner.Peek() == 'e' || t.scanner.Peek() == 'E' {
		buffer.WriteRune(t.scanner.Next())
		next := t.scanner.Peek()
		if next == '-' || next == '+' {
			buffer.WriteRune(next)
			t.scanner.Next()
		}
		s, e := t.scanDigits()
		if e != nil {
			return "", e
		}
		buffer.WriteString(s)
	}

	return buffer.String(), nil
}

func isIdentifierRune(r rune) bool {
	// TODO: not sure if this should be made a little bit more sophisticated.
	s := string(r)
	c := s[0]
	if len(s) > 1 {
		return false
	}
	// From the ascii table: http://www.asciitable.com/index/asciifull.gif
	return utf8.RuneLen(r) == 1 && (c >= 65 && c <= 89 || c >= 97 && c <= 122)
}

func isNumberStart(r rune) bool {
	return isRuneDigit(r) || r == '-'
}

func isInsignificantWhitespace(r rune) bool {
	s := string(r)
	c := s[0]
	if len(s) > 1 {
		return false
	}
	// The list of insignificant whitespaces: https://www.ietf.org/rfc/rfc4627.txt
	return c == 0x20 || c == 0x09 || c == 0x0A || c == 0x0D
}

func (t *Tokenizer) scanWhitespaces() (string, error) {
	var buffer bytes.Buffer
	for {
		r := t.scanner.Peek()
		if !isInsignificantWhitespace(r) {
			return buffer.String(), nil
		}
		buffer.WriteRune(r)
		t.scanner.Next()
	}
}

// Peek peers at the latest token
func (t *Tokenizer) Peek() Token {
	token, _ := t.Scan()
	t.peekedTokens = append(t.peekedTokens, token)
	return token
}

// Scan scans the next token.
func (t *Tokenizer) Scan() (Token, bool) {
	// TODO: remove the boolean.

	if len(t.peekedTokens) > 0 {
		p := t.peekedTokens
		token, p := p[len(p)-1], p[:len(p)-1]
		t.peekedTokens = p
		return token, false
	}

	position := t.scanner.Pos()
	var r rune
	for {
		r = t.scanner.Next()
		if !isInsignificantWhitespace(r) {
			break
		}
	}
	var token Token
	switch {
	case r == '"':
		s, e := t.scanString()
		if e != nil {
			panic("My scanner does not handle errors yet")
		}
		token = Token{s, JSONString, position}
	case r == '{':
		token = Token{string(r), JSONOpenBrace, position}
	case r == '}':
		token = Token{string(r), JSONCloseBrace, position}
	case r == '[':
		token = Token{string(r), JSONOpenSquareBracket, position}
	case r == ']':
		token = Token{string(r), JSONCloseSquareBracket, position}
	case r == ':':
		token = Token{string(r), JSONColon, position}
	case r == ',':
		token = Token{string(r), JSONComma, position}
	case isIdentifierRune(r):
		var toCat bytes.Buffer
		toCat.WriteRune(r)
		s, e := t.scanIdentifier()
		if e != nil {
			panic("My scanner does not handle errors yet")
		}
		toCat.WriteString(s)
		token = Token{toCat.String(), JSONIdentifier, position}
	case isNumberStart(r):
		var toCat bytes.Buffer
		s, e := t.scanNumber(r)
		if e != nil {
			panic("My scanner does not handle errors yet")
		}
		toCat.WriteString(s)
		token = Token{toCat.String(), JSONNumber, position}
	case isInsignificantWhitespace(r):
		// The assignment does not require that we reconstruct the white spaces
		// and therefore, we will just ignore them. So fuck it.
		// var toCat bytes.Buffer
		// toCat.WriteRune(r)
		// s, e := t.scanWhitespaces()
		// if e != nil {
		// 	panic("My scanner does not handle errors yet")
		// }
		// toCat.WriteString(s)
		// token = Token{toCat.String(), JSONWhitespace, position}
		panic("We should not be here.")
		// Maybe in the future we will consider
	case r == scanner.EOF:
		return Token{string(r), JSONEnd, position}, true
	default:
		panic("I wrote a bad scanner")
	}
	return token, t.scanner.Peek() == scanner.EOF
}
