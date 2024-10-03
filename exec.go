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

func ReadStdinAST() {
	// Read file to strings optionally
	content := fmt.Sprintf("package main\n%s\n %s\n func main() {\n%s\n%s\n } \n func UNUSED(x ...any) {\n}", IMPORTBREAK, FUNCDEFBREAK, MAINBREAK, UNUSEDBREAKS)

	staged := make(map[string]*ImportState)
	multiline := false
	textbuf := []string{}
	var err error
	var rollback string
	scanner := bufio.NewScanner(os.Stdin)
	stack := NewStack()

	root := NewRootAst()

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
						Inject(val.exec, "IMPORT", root)
						val.injected = true
					}
				}
				if stype == "MAIN" {
					// for now just do single line declarations
					decs := GetDeclarations(fulltext)
					for _, dec := range decs {
						// ideally prolly should have a queue with ops and routine execing those ops
						decl := fmt.Sprintf("UNUSED(%s)", dec)
						Inject(decl, "UNUSED", root)
					}
				}
			}

			// TODO for single var expr wrap in fmt.Println and write unwrapped to file
			Inject(fulltext, stype, root)
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

func Inject(expr string, stype string, root *ASTRoot) {
	if stype == "IMPORT" {
		root.AddImports(expr)
	} else if stype == "MAIN" {
		root.AddMain(expr)
	} else if stype == "FUNC_DEF" {
		root.AddDecls(expr)
	}
}

func ExecShell() ([]byte, error) {

	// needs external package support
	cmd := exec.Command("go", "run", SHELLPATH)

	output, err := cmd.CombinedOutput()
	return output, err

}

func main() {
	ReadStdinAST()
}
