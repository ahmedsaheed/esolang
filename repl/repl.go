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
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/log"
)

var PROMPT = ">>"
var SHOULD_MULTILINE = false

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
	var inputBuffer strings.Builder
	for {
		PROMPT = map[bool]string{true: ">>>", false: ">>"}[SHOULD_MULTILINE]
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		if SHOULD_MULTILINE {

			//check the line for the multiline command
			if strings.HasPrefix(line, ".disable-multiline") {
				disableMultiline()
				continue
			}
			// Append the line to the input buffer
			inputBuffer.WriteString(line)
			inputBuffer.WriteString("\n")

			// If the line ends with a semicolon or is empty, evaluate the input
			if strings.HasSuffix(line, ";") {

				evaluteInput(inputBuffer.String(), logger, environmnet, stdLib)
				// Clear the input buffer
				inputBuffer.Reset()
				inputBuffer.WriteString("\n")
			}
		} else {
			if line[0] == '.' {
				evaluateReplCommand(line[1:])
			} else {
				evaluteInput(line, logger, environmnet, stdLib)
			}
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

	fmt.Println(syntaxHiglight(input))

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
			var DONT_PRINT = "flag=noshow"
			if !strings.HasSuffix(output, DONT_PRINT) {
				log.Info(output)
			}
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
		var DONT_PRINT = "flag=noshow"
		if strings.HasSuffix(output, DONT_PRINT) {
			output = output[:len(output)-len(DONT_PRINT)]
		}
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

//go:embed repl_help.md
var help string

func replHelp() {
	out, err := glamour.Render(help, "dark")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(out)
	fmt.Println("Feel free to type in commands")
}

func evaluateReplCommand(input string) {
	switch input {
	case "help":
		replHelp()
	case "exit":
		fmt.Println("Goodbye!")
		os.Exit(0)
	case "clear":
		fmt.Print("\033[H\033[2J")
	case "version":
		fmt.Println("esolang version 0.1.0")
	case "enable-multiline":
		SHOULD_MULTILINE = true
		fmt.Println("Multiline mode enabled")
	case "disable-multiline":
		disableMultiline()
	default:
		fmt.Println("Unknown command")
	}
}

func disableMultiline() {
	output := map[bool]string{true: "Multiline disabled", false: "Multiline not enabled, skipping."}[SHOULD_MULTILINE]
	if SHOULD_MULTILINE {
		SHOULD_MULTILINE = false
	}
	fmt.Println(output)
}

func syntaxHiglight(input string) string {
	md := "```js" + "\n" + input + "\n" + "```"
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(120),
	)
	styled, _ := r.Render(md)
	return strings.TrimSpace(styled)
}

func printParserErrors(errors []string, log *log.Logger) {
	for _, msg := range errors {
		log.Error(msg)
	}
}
