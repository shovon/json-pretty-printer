package main

import (
	"fmt"
	"strings"
	"testing"
	"text/scanner"
)

func TestParse(t *testing.T) {
	var s scanner.Scanner
	var reader *scanner.Scanner
	var tokenizer Tokenizer
	var tree interface{}
	var object ObjectNode
	var arr ArrayNode
	var ok bool

	reader = s.Init(strings.NewReader("{}"))
	tokenizer = NewTokenizer(reader)
	tree = Parse(&tokenizer)
	if _, ok := tree.(ObjectNode); !ok {
		panic("Should have been an object node")
	}

	reader = s.Init(strings.NewReader("{\"f\":1,\"g\":2}"))
	tokenizer = NewTokenizer(reader)
	tree = Parse(&tokenizer)
	object, ok = tree.(ObjectNode)
	if !ok {
		panic("Should have been an object node")
	} else if len(object.properties) != 2 {
		panic("Should have had at least one element as property")
	} else if object.properties[0].name != "\"f\"" {
		panic(fmt.Sprintf("The name of the property should have been \"f\", but was \"%s\"", object.properties[0].name))
	} else {
		v, ok := (*object.properties[0].value).(ValueNode)
		if !ok {
			panic("The type of the property should have been a value node")
		} else if v.token.Content != "1" {
			panic("The value of the property should have been 1")
		}
	}

	reader = s.Init(strings.NewReader("[1,2,3]"))
	tokenizer = NewTokenizer(reader)
	tree = Parse(&tokenizer)
	arr, ok = tree.(ArrayNode)
	if !ok {
		panic("Should have parsed as an array node")
	} else if len(arr.elements) != 3 {
		panic("Should have had 3 elements")
	}

	reader = s.Init(strings.NewReader("[{}, 2]"))
	tokenizer = NewTokenizer(reader)
	for {
		tokenizer.Scan()
		if tokenizer.Peek().TokenType == JSONEnd {
			break
		}
	}

	reader = s.Init(strings.NewReader("[{}, 2]"))
	tokenizer = NewTokenizer(reader)
	tree = Parse(&tokenizer)
	arr, ok = tree.(ArrayNode)
	if !ok {
		panic("Should have parsed as an array node")
	} else if len(arr.elements) != 2 {
		panic("Should have had 1 elements")
	} else {
		_, ok := (*arr.elements[0]).(ObjectNode)
		if !ok {
			panic("The first element should have been an object node")
		}
	}

	reader = s.Init(strings.NewReader("{} "))
	tokenizer = NewTokenizer(reader)
	tree = Parse(&tokenizer)
	_, ok = tree.(ObjectNode)
	if !ok {
		panic("Should have been parsed as an object")
	}

	reader = s.Init(strings.NewReader("{\"f\":{\"g\":\"1\"}}"))
	tokenizer = NewTokenizer(reader)
	tree = Parse(&tokenizer)
	object, ok = tree.(ObjectNode)
	if !ok {
		panic("Should have been an object node")
	} else if len(object.properties) != 1 {
		panic("Should have had exactly one element as property")
	} else if object.properties[0].name != "\"f\"" {
		panic(fmt.Sprintf("The name of the property should have been \"f\", but was \"%s\"", object.properties[0].name))
	} else {
		o, ok := (*object.properties[0].value).(ObjectNode)
		if !ok {
			panic("The type of the property should have been an object node")
		} else if len(o.properties) != 1 {
			panic("The sub object should have exactly one property")
		} else {
			v, ok := (*o.properties[0].value).(ValueNode)
			if !ok {
				panic("The value should have been a value node")
			} else if v.token.Content != "\"1\"" {
				panic(fmt.Sprintf("The value should have been the string \"1\", but instead got %s", v.token.Content))
			}
		}
	}
}
