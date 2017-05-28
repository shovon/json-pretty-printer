package main

// Node the base node type
type Node interface {
	GetType() string
}

// PropertyNode represents a property of an object
type PropertyNode struct {
	name  string
	value *Node
}

// GetType returns the string PropertyNode
func (p PropertyNode) GetType() string {
	return "PropertyNode"
}

// ObjectNode represents an object
type ObjectNode struct {
	properties []*PropertyNode
}

// GetType returns the string ObjectNode
func (o ObjectNode) GetType() string {
	return "ObjectNode"
}

// ArrayNode represents an array
type ArrayNode struct {
	elements []*Node
}

// GetType returns the string ArrayNode
func (a ArrayNode) GetType() string {
	return "ArrayNode"
}

// ValueNode represents a leaf value node
type ValueNode struct {
	token Token
}

// GetType returns the string ValueNode
func (a ValueNode) GetType() string {
	return "ValueNode"
}

func parseObject(tokenizer *Tokenizer) ObjectNode {
	var node ObjectNode

	for {
		token := tokenizer.Peek()
		if token.TokenType == JSONCloseBrace {
			tokenizer.Scan()
			return node
		}
		tokenizer.Scan() // Key
		key := token
		tokenizer.Scan() // Colon
		value := Parse(tokenizer)
		node.properties = append(node.properties, &PropertyNode{key.Content, &value})
		token, _ = tokenizer.Scan()
		if token.TokenType == JSONCloseBrace {
			return node
		}
		// Otherwise, we'll just assume that we have a comma
	}
}

func parseArray(tokenizer *Tokenizer) ArrayNode {
	var node ArrayNode
	for {
		if tokenizer.Peek().TokenType == JSONCloseSquareBracket {
			tokenizer.Scan()
			return node
		}
		value := Parse(tokenizer)
		node.elements = append(node.elements, &value)
		token, _ := tokenizer.Scan()
		if token.TokenType == JSONCloseSquareBracket {
			return node
		}
		assert(token.TokenType == JSONComma, "Was expecting a comma, at least!")
		// Otherwise, we'll juse assume that we have a comma
	}
}

// Parse parses the input tokens
func Parse(tokenizer *Tokenizer) Node {
	var node Node
	token, _ := tokenizer.Scan()

	switch token.TokenType {

	case JSONOpenBrace:
		node = parseObject(tokenizer)
	case JSONOpenSquareBracket:
		node = parseArray(tokenizer)
	case JSONIdentifier:
		node = ValueNode{token}
	case JSONString:
		node = ValueNode{token}
	case JSONNumber:
		node = ValueNode{token}
	default:
		panic("I wrote a bad parser")
	}

	return node
}
