package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var SHELLPATH string = "/var/lib/goshpel/shell.go"

func main() {
	ReadStdin()
}

func ReadStdin() {
	multiline := false
	for {
		scanner := bufio.NewScanner(os.Stdin)
		if multiline {
			fmt.Print("...")
		} else {
			fmt.Print(">> ")
		}
		eof := scanner.Scan()
		text := scanner.Text()
		fmt.Println(text)
		fmt.Println(eof)

		multiline = IsMulti(text)

	}
}

func IsMulti(text string) bool {

	para := strings.Count(text, "(") == strings.Count(text, ")")
	para = para && strings.Index(text, "(") < strings.Index(text, ")")
	curly := strings.Count(text, "{") == strings.Count(text, "}")
	curly = curly && strings.Index(text, "{") < strings.Index(text, "}")

	// invalid multi line
	square := strings.Count(text, "[") == strings.Count(text, "]")
	if square {
		return false
	}

	return !(para && curly)

}

func AppendToFile(text string) error {

	fi, err := os.Open(SHELLPATH)
	if err != nil {
		return err
	}

	defer fi.Close()

	if _, err := fi.WriteString(text); err != nil {
		panic(err)
	}

	// is valid expression + shell returned error
	return bufio.ErrFinalToken
}
