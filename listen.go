package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// var SHELLPATH string = "/var/lib/goshpel/shell.go"

var SHELLPATH string = "./dev/shell.go"

var RESTORE string = "./dev/shell-old.txt"

// TODO:
// StdoutPipe errors should trigger a reroll (after piping error to terminal)

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
			fmt.Println(stype, err)
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

func GetStatementType(txt []string) (string, error) {

	text := strings.Join(txt, " ")

	stype := "MAIN"

	if imprt := strings.HasPrefix(text, "import "); imprt {
		stype = "IMPORT"
	} else if funcdef := strings.HasPrefix(text, "func "); funcdef {
		stype = "FUNC_DEF"
	}

	return stype, nil
}

func Inject(expr string) {

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

	// needs external package support
	cmd := exec.Command("go", "run", SHELLPATH)

	output, err := cmd.CombinedOutput()

	fmt.Println(string(output))
	fmt.Println(err)

	return nil

}

func main() {
	// TEMP
	ExecShell()
	//ReadStdin()
}
