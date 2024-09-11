package main

import (
	"bufio"
	"fmt"
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
const UNUSEDBREAKS string = "//UB"
const FUNCDEFBREAK string = "//FDB"

// TODO:
// Clean up/refactor
// StdoutPipe errors should trigger a reroll (after piping error to terminal)

type ImportState struct {
	exec     string
	injected bool
}

func ReadStdin() {
	// Read file to strings optionally
	content := fmt.Sprintf("package main\n%s\n %s\n func main() {\n%s\n%s\n } \n func UNUSED(x ...any) {\n}", IMPORTBREAK, FUNCDEFBREAK, MAINBREAK, UNUSEDBREAKS)

	staged := make(map[string]*ImportState)
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
			textbuf = nil
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
		textbuf = append(textbuf, "\n")
		if !multiline {
			fulltext := strings.Join(textbuf, " ")
			rollback = strings.Clone(content)

			// determine type (import, package, inside main())
			stype, err := GetStatementType(fulltext)
			if err != nil {
				break
			}

			// if stype import, track it and only inject if a given expr has the import
			if stype == "IMPORT" {
				pkgs := GetPkgNames(fulltext)
				for _, pkg := range pkgs {
					if _, ok := staged[pkg]; !ok {
						staged[pkg] = &ImportState{fmt.Sprintf("import \"%s\"", pkg), false}
					}
				}
				// read next
				continue
			} else {
				imprts := GetUsedPkgs(fulltext)
				for _, imp := range imprts {
					if val, ok := staged[imp]; ok && !val.injected {
						Inject(val.exec, "IMPORT", &content)
						val.injected = true
					}
				}
				if stype == "MAIN" {
					// for now just do single line declarations
					decs := GetDeclarations(fulltext)
					for _, dec := range decs {
						// ideally prolly should have a queue with ops and routine execing those ops
						decl := fmt.Sprintf("UNUSED(%s)", dec)
						Inject(decl, "UNUSED", &content)
					}
				}
			}

			// TODO for single var expr wrap in fmt.Println and write unwrapped to file

			Inject(fulltext, stype, &content)
			AppendToFile(content)

			out, err := ExecShell()
			fmt.Println(string(out))
			if err != nil {
				content = rollback
			}
			// remove all fmt Stdout from main

		}
	}
}

func Inject(expr string, stype string, content *string) {

	var sb strings.Builder
	var before, after, breaker string

	switch stype {
	case "MAIN":
		breaker = MAINBREAK
	case "FUNC_DEF":
		// include var/const decs here?
		breaker = FUNCDEFBREAK
	case "IMPORT":
		breaker = IMPORTBREAK
	case "UNUSED":
		breaker = UNUSEDBREAKS
	}

	before, after, _ = strings.Cut(*content, breaker)

	sb.WriteString(before)
	sb.WriteString(expr)
	sb.WriteString("\n")
	sb.WriteString(" ")
	sb.WriteString(breaker)
	sb.WriteString(after)
	*content = sb.String()

}

func ExecShell() ([]byte, error) {

	// needs external package support
	cmd := exec.Command("go", "run", SHELLPATH)

	output, err := cmd.CombinedOutput()
	return output, err

}

func main() {
	ReadStdin()
}
