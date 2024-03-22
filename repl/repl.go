/*
Package repl is a minimal Read-Eval-Print-Loop.
*/
package repl

import (
	"bufio"
	"esolang/lang-esolang/evaluator"
	"esolang/lang-esolang/lexer"
	"esolang/lang-esolang/object"
	"esolang/lang-esolang/parser"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/log"
)

const PROMPT = ">>"

// Start starts the Read-Eval-Print-Loop.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	environmnet := object.NewEnvironment()
	logger := log.New(os.Stderr)
	logger.SetReportTimestamp(false)
	logger.SetReportCaller(false)
	logger.SetLevel(log.DebugLevel)
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
			printParserErrors(parser.Errors(), logger)
			continue
		}
		evaluated := evaluator.Eval(program, environmnet)
		if evaluated != nil {
			message := evaluated.Inspect()
			if len(message) > 5 && message[:5] == "ERROR" {
				logger.Error(message[6:])
			} else {
				logger.Info(message)
				io.WriteString(out, "\n")
			}
		}
	}
}

func printParserErrors(errors []string, log *log.Logger) {
	for _, msg := range errors {
		log.Error(msg)
	}
}
