/*
Package repl is a minimal Read-Eval-Print-Loop.
*/
package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lang-monkey/lexer"
	"monkey/lang-monkey/parser"
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
		L := lexer.New(line)
		P := parser.New(L)

		program := P.ParseProgram()

		if len(P.Errors()) != 0 {
			printParserErrors(out, P.Errors())
			continue
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")

	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
