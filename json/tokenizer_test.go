package json

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"text/scanner"
)

const (
	stringLitTestString  = "\"this is a \\\\test\""
	semiFullJSONTest     = "{\"hello\":\"world\",\"foo\":\"bar\",\"baz\":true,\"widgets\":[1,2,true,false]}"
	whitespaceJSONTest   = "{ \"foobar\": true }"
	noWhitespaceJSONTest = "{\"foobar\":true}"
	frivolousWhitespace  = "{} "
)

func getAllTokens(tokenizer Tokenizer) string {
	var buffer bytes.Buffer
	for token, _ := tokenizer.Scan(); token.TokenType != JSONEnd; token, _ = tokenizer.Scan() {
		buffer.WriteString(token.Content)
	}
	return buffer.String()
}

func testValidNumber(lit string) {
	var s scanner.Scanner
	reader := s.Init(strings.NewReader(lit))
	tokenizer := NewTokenizer(reader)
	token, _ := tokenizer.Scan()
	assert(token.TokenType == JSONNumber, fmt.Sprintf("Expected that the returned token be a number"))
	assert(token.Content == lit, fmt.Sprintf("Expected content to be '%s', but instead got '%s'", lit, token.Content))
}

func TestScan(t *testing.T) {
	var s scanner.Scanner
	reader := s.Init(strings.NewReader("{"))
	tokenizer := NewTokenizer(reader)
	token, done := tokenizer.Scan()
	assert(done, "Should have reached the last token")
	assert(token.Content == "{", fmt.Sprintf("%c should be {", token.Content))
	assert(token.TokenType == JSONOpenBrace, fmt.Sprintf("Token type should have been labeled as open brace"))

	reader = s.Init(strings.NewReader("{}"))
	tokenizer = NewTokenizer(reader)
	tokenizer.Scan()
	token, done = tokenizer.Scan()
	assert(done, "Should have reaced the last token")
	assert(token.Content == "}", fmt.Sprintf("%c should be }", token.Content))
	assert(token.TokenType == JSONCloseBrace, fmt.Sprintf("Token type should have been labeled as a close brace"))

	reader = s.Init(strings.NewReader(stringLitTestString))
	tokenizer = NewTokenizer(reader)
	token, done = tokenizer.Scan()
	assert(done, "Should have reached the last token")
	assert(token.Content == stringLitTestString, fmt.Sprintf("The string literal does not match"))
	assert(token.TokenType == JSONString, fmt.Sprintf("Token type should have been labeled as a string literal"))

	reader = s.Init(strings.NewReader(semiFullJSONTest))
	tokenizer = NewTokenizer(reader)
	str := getAllTokens(tokenizer)
	assert(str == semiFullJSONTest, "The reconstructed string should be equal")

	testValidNumber("0")
	testValidNumber("1")
	testValidNumber("20")
	testValidNumber("0.0")
	testValidNumber("20.0")
	testValidNumber("1.1")
	testValidNumber("1.10")
	testValidNumber("1.1")
	testValidNumber("-0")
	testValidNumber("-1")
	testValidNumber("-20")
	testValidNumber("-0.0")
	testValidNumber("-1.0")
	testValidNumber("0e1")
	testValidNumber("1e1")
	testValidNumber("1e+1")
	testValidNumber("1e+10")
	testValidNumber("10e+1")
	testValidNumber("10e+1")
	testValidNumber("10e-1")
	testValidNumber("0.10e1")
	testValidNumber("0.10e-10")
	testValidNumber("0.10e+10")

	reader = s.Init(strings.NewReader(whitespaceJSONTest))
	tokenizer = NewTokenizer(reader)
	str = getAllTokens(tokenizer)
	assert(str == noWhitespaceJSONTest, "The reconstructed string should have the whitespace removed")

	reader = s.Init(strings.NewReader("[{}, 2]"))
	tokenizer = NewTokenizer(reader)
	str = getAllTokens(tokenizer)
	assert(str == "[{},2]", "Should be able to reconstruct")
}
