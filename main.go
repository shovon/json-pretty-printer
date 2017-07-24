package main

import (
	"fmt"
	"os"
	"text/scanner"

	"./json"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Printf("Please provide a filename\n")
		os.Exit(1)
	}
	var s scanner.Scanner
	f, err := os.Open(args[0])
	if err != nil {
		panic(err)
	}
	scanner := s.Init(f)
	tokenizer := json.NewTokenizer(scanner)

	tree := json.Parse(&tokenizer)

	fmt.Printf("%s", `<!doctype html>
	<html lang='en'>
		<head>
			<meta charset='utf-8'>
			<title>Here ye some JSON!</title>
		</head>
		<body style="padding 0; margin: 0; background-color: #002b36">
			<div style="padding: 5px">
	<span style="font-family:monospace; white-space:pre">`)

	json.PrintTree(tree, 0)

	fmt.Printf("%s", `</span>
			</div>
		</body>
	</html>`)
}
