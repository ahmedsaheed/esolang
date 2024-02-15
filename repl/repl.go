/*
Package repl is a minimal Read-Eval-Print-Loop.
*/
package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lang-monkey/lexer"
	"monkey/lang-monkey/token"
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

		for tok := L.NextToken(); tok.Type != token.EOF; tok = L.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
