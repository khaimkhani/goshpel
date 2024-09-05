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

// dont need this
// load from old or empty
var SHELLPATH string = "./dev/shell.go"

// Be able to load "old" using a flag
var RESTORE string = "./dev/shell-old.txt"

const IMPORTBREAK string = "//IB"
const MAINBREAK string = "//MB"
const FUNCDEFBREAK string = "//FDB"

// TODO:
// Clean up/refactor
// StdoutPipe errors should trigger a reroll (after piping error to terminal)

func ReadStdin() {
	// Read file to strings optionally
	content := fmt.Sprintf("package main\n%s\n %s\n func main() {\n%s\n }", IMPORTBREAK, FUNCDEFBREAK, MAINBREAK)
	t := NewTracker()
	multiline := false
	textbuf := []string{}
	var err error
	var rollback string
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
			stype, err := GetStatementType(textbuf, t)
			if err != nil {
				break
			}

			fmt.Println(content)
			rollback = strings.Clone(content)
			Inject(textbuf, stype, &content)
			AppendToFile(content)

			// dont exec import statements
			// wait till func def
			if out, err := ExecShell(); err != nil {
				fmt.Println(string(out))
				content = rollback
			}

			// reset buffer
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

func GetStatementType(txt []string, t *Tracker) (string, error) {

	text := strings.Join(txt, " ")

	stype := "MAIN"

	if imprt := strings.HasPrefix(text, "import "); imprt {
		stype = "IMPORT"
	} else if funcdef := strings.HasPrefix(text, "func "); funcdef {
		stype = "FUNC_DEF"
	} else if constdef := strings.HasPrefix(text, "const "); constdef {
		stype = "FUNC_DEF"
	} else if vardef := strings.HasPrefix(text, "var "); vardef {
		stype = "FUNC_DEF"
	}

	if stype == "MAIN" {
		// check if var declaration/reassignment
		// check if tracker has var
		// do appropriate things
		// regex match for "=" and ":="
	}

	return stype, nil
}

func Inject(text []string, stype string, content *string) {

	expr := strings.Join(text, " ")
	var sb strings.Builder
	var before, after, breaker string
	var ok bool

	switch stype {
	case "MAIN":
		// inject at bottom of main func
		breaker = MAINBREAK

	case "FUNC_DEF":
		// before funcdefbreak
		breaker = FUNCDEFBREAK

	case "IMPORT":
		// before importbreak
		breaker = IMPORTBREAK

	case "REPLACE":
		// replace existing var
	}

	before, after, ok = strings.Cut(*content, breaker)

	if !ok {
		// error
		fmt.Println(ok)
		return
	}

	sb.WriteString(before)
	sb.WriteString(expr)
	sb.WriteString("\n")
	sb.WriteString(" ")
	sb.WriteString(breaker)
	sb.WriteString(after)
	*content = sb.String()
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

	fi, err := os.OpenFile(SHELLPATH, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
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

func ExecShell() ([]byte, error) {

	// needs external package support
	cmd := exec.Command("go", "run", SHELLPATH)

	output, err := cmd.CombinedOutput()

	return output, err

}

func main() {
	// TODO run init
	// TEMP
	// ExecShell()
	ReadStdin()
}
