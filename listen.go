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
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print(">> ")
		eof := scanner.Scan()
		text := scanner.Text()
		fmt.Println(text)
		fmt.Println(eof)
	}
}

func SyntaxCheck(text string) (bool, error) {

	para := strings.Count(text, "(") == strings.Count(text, ")")
	curly := strings.Count(text, "{") == strings.Count(text, "}")

	// invalid
	square := strings.Count(text, "[") == strings.Count(text, "]")
	if square {
		return false, bufio.ErrBadReadCount
	}

	return !(para && curly), nil

}

func AppendToFile(text string) error {

	fi, err := os.Open(SHELLPATH)
	if err != nil {
		return err
	}

	defer fi.Close()

	// is valid expression + shell returned error
	return bufio.ErrFinalToken
}
