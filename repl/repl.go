package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/mahiro72/monkey-lang/lexer"
	"github.com/mahiro72/monkey-lang/token"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		if line == "exit" {
			fmt.Printf("bye 🐵\n")
			os.Exit(0)
		}
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
