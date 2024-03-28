/*
Package repl is a minimal Read-Eval-Print-Loop.
*/
package repl

import (
	"bufio"
	_ "embed"
	"esolang/lang-esolang/evaluator"
	"esolang/lang-esolang/lexer"
	"esolang/lang-esolang/object"
	"esolang/lang-esolang/parser"
	"fmt"
	"io"
	"os"
	"os/user"

	"github.com/charmbracelet/log"
)

const PROMPT = ">>"

func Execute(input string, stdLib string) {
	environmnet := object.NewEnvironment()
	logger := generateLogger()
	evaluteInput(input, logger, environmnet, stdLib)
}

func Start(in io.Reader, out io.Writer, stdLib string) {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! Welcome to esolang's repl \n", user.Username)
	fmt.Printf("Feel free to type in commands \n")
	fmt.Printf("Type '.help' for assistance \n")
	scanner := bufio.NewScanner(in)
	environmnet := object.NewEnvironment()
	logger := generateLogger()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		if line == ".help" {
			replHelp()
		} else if line == ".exit" {
			fmt.Println("Goodbye!")
			return
		} else {
			evaluteInput(line, logger, environmnet, stdLib)
		}

	}
}

func evaluteInput(input string, log *log.Logger, environmnet *object.Environment, stdLib string) {
	initialLexer := lexer.New(input)
	initialParser := parser.New(initialLexer)

	program := initialParser.ParseProgram()

	if len(initialParser.Errors()) != 0 {
		printParserErrors(initialParser.Errors(), log)
		return
	}

	libLexer := lexer.New(stdLib)
	libParser := parser.New(libLexer)
	libProgram := libParser.ParseProgram()
	evaluator.Eval(libProgram, environmnet)

	evaluated := evaluator.Eval(program, environmnet)
	if evaluated != nil {
		output := evaluated.Inspect()
		if len(output) > 5 && output[:5] == "ERROR" {
			log.Error(output[6:])
		} else {
			log.Info(output)
		}
	}
}
func EvlauateFromPlayground(input string) string {
	initialLexer := lexer.New(input)
	initialParser := parser.New(initialLexer)
	environmnet := object.NewEnvironment()
	program := initialParser.ParseProgram()

	if len(initialParser.Errors()) != 0 {
	    return initialParser.Errors()[0]
    }

	evaluated := evaluator.Eval(program, environmnet)
	if evaluated != nil {
		output := evaluated.Inspect()
        return output
	} else {
        return "ERROR: No output found"
    }
}


func generateLogger() *log.Logger {
	logger := log.New(os.Stderr)
	logger.SetReportTimestamp(false)
	logger.SetReportCaller(false)
	logger.SetLevel(log.DebugLevel)

	return logger
}

func replHelp() {
	fmt.Println("Here are a few commands you can use:")
	fmt.Println(".help - Show this message")
	fmt.Println(".exit - Exit the repl")
	fmt.Println("Feel free to type in commands")
}

func printParserErrors(errors []string, log *log.Logger) {
	for _, msg := range errors {
		log.Error(msg)
	}
}
