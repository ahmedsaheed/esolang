/*
Package repl is a minimal Read-Eval-Print-Loop.
*/
package repl

import (
	"bufio"
	"esolang/lang-esolang/evaluator"
	"esolang/lang-esolang/lexer"
	"esolang/lang-esolang/parser"
	"fmt"
	"io"
)

const PROMPT = ">>"

// Start starts the Read-Eval-Print-Loop.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		lexer := lexer.New(line)
		parser := parser.New(lexer)

		program := parser.ParseProgram()

		if len(parser.Errors()) != 0 {
			printParserErrors(out, parser.Errors())
			continue
		}
		evaluated := evaluator.Eval(program)
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")

	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
