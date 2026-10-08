package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/app/internal/ast"
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/parse"
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/scanner"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "tokenize":
		printTokens(os.Args[2])
	case "parse":
		parseFile(os.Args[2])
	case "evaluate":
		evaluateFile(os.Args[2])
	case "run":
		runFile(os.Args[2])
	}

	os.Exit(0)
}

func runFile(fileName string) {
	statements := parseExpr(fileName, false)
	interpreter := ast.NewInterpreter()
	resolver := ast.NewResolver(interpreter)
	resolver.ResolveStatements(statements)
	err := interpreter.InterpretStatements(statements)
	if err != 0 {
		os.Exit(err)
	}
}

func evaluateFile(fileName string) {
	expr := parseExpr(fileName, true)
	interpreter := ast.NewInterpreter()
	res, err := interpreter.InterpretExpr(expr[0])
	if err != 0 {
		os.Exit(err)
	}
	fmt.Println(res)
}

func parseFile(fileName string) {
	expr := parseExpr(fileName, true)
	astPrinter := ast.NewAstPrinter()
	fmt.Println(expr[0].Accept(astPrinter))
}

func parseExpr(fileName string, parseExpr bool) []parse.Statement {
	tokens, errorCode := tokenize(fileName)
	if errorCode != 0 {
		os.Exit(errorCode)
	}
	parser := parse.NewParser(tokens, parseExpr)
	expr, err := parser.Parse()
	if err != 0 {
		os.Exit(err)
	}
	return expr
}

func printTokens(fileName string) {
	tokens, errorCode := tokenize(fileName)

	for _, token := range tokens {
		fmt.Println(token)
	}
	os.Exit(errorCode)
}

func tokenize(fileName string) ([]*scanner.Token, int) {
	fileContents, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	tokenizer := scanner.NewTokenizer(string(fileContents))
	tokens, errorCode := tokenizer.Tokenize()
	return tokens, errorCode
}
