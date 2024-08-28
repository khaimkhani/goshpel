package goshpel

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// var SHELLPATH string = "/var/lib/goshpel/shell.go"

var SHELLPATH string = "./shell.go"

var ALLOWEDPARA []string = []string{"(", ")", "{", "}"}

func ReadStdin() {
	multiline := false
	scanner := bufio.NewScanner(os.Stdin)
	stack := NewStack()

	for {
		if multiline {
			fmt.Print("... ")
		} else {
			fmt.Print(">> ")
		}
		scanner.Scan()

		text := scanner.Text()
		fmt.Println(text)

		// need a stack here making sure all paranthesis are closed
		multiline = IsMulti(text)
		// if still ml then continue,
		// else Append to file and run
		// This retains last state effectively

	}
}

func IsMulti(text string) bool {

	// valid multilines
	para := strings.Count(text, "(") != strings.Count(text, ")")
	para = para || strings.Index(text, "(") > strings.Index(text, ")")
	curly := strings.Count(text, "{") != strings.Count(text, "}")
	curly = curly || strings.Index(text, "{") > strings.Index(text, "}")

	return para || curly

}

func AppendToFile(text string) error {

	fi, err := os.Open(SHELLPATH)
	if err != nil {
		return err
	}

	defer fi.Close()

	// Need checks here for what goes inside main {} and what is outside
	if _, err := fi.WriteString(text); err != nil {
		panic(err)
	}

	// is valid expression + shell returned error
	return bufio.ErrFinalToken
}

func ExecShell() error {
	// Exec shell.go and pass any errors to Stdout
	// Also need a way to revert to last state
	return nil
}

func main() {
	ReadStdin()
}
