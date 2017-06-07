package json

import (
	"fmt"
	"unicode/utf8"
)

var colorMap = map[int]string{
	JSONIdentifier:         "b58900",
	JSONString:             "2aa198",
	JSONNumber:             "d33682",
	JSONOpenBrace:          "93a1a1",
	JSONCloseBrace:         "93a1a1",
	JSONColon:              "268bd2",
	JSONComma:              "859900",
	JSONOpenSquareBracket:  "6c71c4",
	JSONCloseSquareBracket: "6c71c4",
}

func spacePad(n int) string {
	var str string
	for i := 0; i < n; i++ {
		str += " "
	}
	return str
}

func printSpan(content, color string, spaces int) {
	fmt.Printf("%s<span style='color:#%s'>", spacePad(spaces), color)
	for _, r := range content {
		fmt.Printf("%s", getEscapedRune(r))
	}
	fmt.Printf("</span>")
}

func getEscapedRune(r rune) string {
	switch r {
	case '<':
		return "&lt;"
	case '>':
		return "&gt;"
	case '&':
		return "&amp;"
	case '"':
		return "&quot;"
	case '\'':
		return "&apos;"
	}
	return string(r)
}

func getObjectPadding(node ObjectNode) int {
	if len(node.properties) < 1 {
		return 0
	}
	var max int
	for _, property := range node.properties {
		if utf8.RuneCountInString(property.name) > max {
			max = len(property.name)
		}
		if max > 50 {
			return 0
		}
	}
	return max
}

func shouldSamelineArray(node ArrayNode) bool {
	if len(node.elements) > 10 {
		return false
	}
	for _, el := range node.elements {
		if _, ok := (*el).(ValueNode); !ok {
			return false
		}
	}
	return true
}

func printScalar(token Token) {
	switch token.TokenType {
	case JSONIdentifier:
		printSpan(token.Content, colorMap[JSONIdentifier], 0)
	case JSONString:
		printSpan(token.Content, colorMap[JSONString], 0)
	case JSONNumber:
		printSpan(token.Content, colorMap[JSONNumber], 0)
	default:
		print("I don't know what kind of value this is")
	}
}

func printObject(node ObjectNode, indent int) {
	padding := indent * 2

	printSpan("{", colorMap[JSONOpenBrace], 0)
	fmt.Println("")
	oPaddingNum := getObjectPadding(node)
	if len(node.properties) > 0 {
		property := node.properties[0]
		printSpan(property.name, colorMap[JSONString], padding+4)
		if oPaddingNum < 50 {
			propertyPad := spacePad(oPaddingNum - utf8.RuneCountInString(property.name))
			fmt.Printf("%s", propertyPad)
		}
		printSpan(":", colorMap[JSONColon], 0)
		fmt.Printf(" ")
		PrintTree(*property.value, indent+1)
		for _, property := range node.properties[1:] {
			fmt.Println("")
			printSpan(",", colorMap[JSONComma], padding+2)
			fmt.Printf(" ")
			printSpan(property.name, colorMap[JSONString], 0)
			if oPaddingNum < 50 {
				propertyPad := spacePad(oPaddingNum - utf8.RuneCountInString(property.name))
				fmt.Printf("%s", propertyPad)
			}
			printSpan(":", colorMap[JSONColon], 0)
			fmt.Printf(" ")
			PrintTree(*property.value, indent+1)
		}
	}
	fmt.Println("")
	printSpan("}", colorMap[JSONCloseBrace], padding)
}

func printArray(node ArrayNode, indent int) {
	padding := indent * 2

	printSpan("[", colorMap[JSONOpenSquareBracket], 0)
	if len(node.elements) > 0 {
		if shouldSamelineArray(node) {
			fmt.Printf(" ")
			for _, element := range node.elements[:len(node.elements)-1] {
				if node, ok := (*element).(ValueNode); ok {
					printScalar(node.token)
					printSpan(",", colorMap[JSONComma], 0)
					fmt.Printf(" ")
				} else {
					panic("Weird. This should have been a value node")
				}
			}
			if node, ok := (*node.elements[len(node.elements)-1]).(ValueNode); ok {
				printScalar(node.token)
			} else {
				panic("Weird. This should have been a value node")
			}
			fmt.Printf(" ")
			printSpan("]", colorMap[JSONCloseSquareBracket], 0)
		} else {
			fmt.Println("")
			fmt.Printf("%s", spacePad(padding+4))
			element := node.elements[0]
			PrintTree(*element, indent+1)
			for _, element := range node.elements[1:] {
				fmt.Println("")
				printSpan(",", colorMap[JSONComma], padding+2)
				fmt.Printf(" ")
				PrintTree(*element, indent+1)
			}
			fmt.Println("")
			printSpan("]", colorMap[JSONCloseSquareBracket], padding)
		}
	}
}

func PrintTree(tree Node, indent int) {
	if node, ok := tree.(ObjectNode); ok {
		printObject(node, indent)
	} else if node, ok := tree.(ArrayNode); ok {
		printArray(node, indent)
	} else if node, ok := tree.(ValueNode); ok {
		printScalar(node.token)
	} else {
		panic("I don't know what kind of a node this is")
	}
}
