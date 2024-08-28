package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

// var SHELLPATH string = "/var/lib/goshpel/shell.go"

var SHELLPATH string = "./dev/shell.go"

var RESTORE string = "./dev/shell-old.txt"

// TODO:
// Need a versioning system so breaking inputs do not crash the system.
// StdoutPipe errors should trigger a reroll (after piping error to terminal)
// Theres a smart and dumb way to do this

func ReadStdin() {
	multiline := false
	textbuf := []string{}
	var err error
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

		multiline, err = CheckMultiline(stack, text)
		if err != nil {
			if !multiline {
				// revert
				CopyFile(RESTORE, SHELLPATH)
			}
			textbuf = nil
			break
		}

		textbuf = append(textbuf, text)
		if !multiline {
			// determine type (import, package, inside main())
			stype, err := GetStatementType(textbuf)
			// exec code
			textbuf = nil
		}
	}
}

func CheckMultiline(s *stack, line string) (bool, error) {

	for _, i := range line {
		switch i {
		case '{':
			s.Push('}')
		case '(':
			s.Push(')')
		case ')', '}':
			char, err := s.Pop()
			if err != nil {
				return false, err
			}
			if char, ok := char.(rune); ok && char != i {
				return false, errors.New("Paranthesis not closed")
			}
		}
	}

	return len(s.s) > 0, nil
}

func GetStatementType(text []string) (string, error) {
	return "yur", nil
}

func CopyFile(src string, dest string) error {
	srcf, err := os.OpenFile(src, os.O_CREATE|os.O_RDONLY, 0644)
	defer srcf.Close()

	if err != nil {
		return err
	}

	destf, err := os.OpenFile(dest, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	defer destf.Close()

	if err != nil {
		return err
	}
	_, err = io.Copy(destf, srcf)

	return err
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
	// TEMP
	err := CopyFile(RESTORE, SHELLPATH)
	fmt.Println(err)
	//ReadStdin()
}
