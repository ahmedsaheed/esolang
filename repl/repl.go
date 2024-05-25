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

var (
	PROMPT                = ">>"
	SHOULD_MULTILINE      = false
	MULTILINE_CURR_NUMBER = 1
	REPL_MODE             = false
	REPL_VERSION          = "0.0.1" // todo: get esolang version on the machine
)

func Execute(sourceName, input string) {
	environmnet := object.NewEnvironment()
	logger := generateLogger()
	evaluteInput(sourceName, input, logger, environmnet)
}

func Start(in io.Reader, out io.Writer) {
	REPL_MODE = true
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Welcome to Esolang version (%s) ~ %s.\n", REPL_VERSION, user.Username)
	fmt.Printf("Type ':help' for assistance. \n")
	scanner := bufio.NewScanner(in)
	environmnet := object.NewEnvironment()
	logger := generateLogger()
	var inputBuffer strings.Builder
	for {
		PROMPT = map[bool]string{true: fmt.Sprint(MULTILINE_CURR_NUMBER) + ">", false: ">>"}[SHOULD_MULTILINE && REPL_MODE]
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		if SHOULD_MULTILINE {

			//check the line for the multiline command
			if strings.HasPrefix(line, ":disable-multiline") {
				disableMultiline()
				continue
			}
			// Append the line to the input buffer
			inputBuffer.WriteString(line)
			inputBuffer.WriteString("\n")
			MULTILINE_CURR_NUMBER += 1

			// If the line ends with a semicolon or is empty, evaluate the input
			if strings.HasSuffix(line, ";") {

				evaluteInput("repl.eso", inputBuffer.String(), logger, environmnet)
				// Clear the input buffer
				inputBuffer.Reset()
				inputBuffer.WriteString("\n")
			}
		} else {
			line = strings.Trim(line, " ")
			if len(line) != 0 && line[0] == ':' {
				evaluateReplCommand(line[1:])
			} else {
				evaluteInput("repl.eso", line, logger, environmnet)
			}
		}
	}
}

func evaluteInput(sourceName, input string, log *log.Logger, environmnet *object.Environment) {
	initialLexer := lexer.New(sourceName, input)
	initialParser := parser.New(initialLexer)
	program := initialParser.ParseProgram()

	if len(initialParser.Errors()) != 0 {
		printParserErrors(initialParser.Errors(), log)
		return
	}
	evaluated := evaluator.Eval(program, environmnet)
	if evaluated != nil {
		output := evaluated.Inspect()

		if len(output) > 5 && output[:5] == "ERROR" {
			log.Error(output[6:])
		} else {
			var DONT_PRINT = "flag=noshow"
			if !strings.HasSuffix(output, DONT_PRINT) {
				if REPL_MODE {
					cyan := "\033[36m" //Cyan color
					reset := "\033[0m" // Reset color
					// Print cyan-colored text
					fmt.Println(cyan + output + reset)
				} else {
					fmt.Println(output)
				}

			}
		}
	}
}
func EvlauateFromPlayground(input string) string {
	initialLexer := lexer.New("plaground.eso", input)
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
		fmt.Println("Esolang Version(" + REPL_VERSION + ")")
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
		MULTILINE_CURR_NUMBER = 1
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
