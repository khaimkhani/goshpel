package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
)

func GetUsedPkgs(text string) []string {

	re := regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*)\s*\.`)
	match := re.FindAllStringSubmatch(text, -1)

	pkgs := []string{}

	for _, m := range match {
		pkg := m[1]
		if strings.Contains(pkg, "/") {
			splitstr := strings.Split(pkg, "/")
			pkg = splitstr[len(splitstr)-1]
		}
		pkgs = append(pkgs, pkg)
	}

	return pkgs

}

func GetPkgNames(text string) []string {

	re := regexp.MustCompile(`"(.*)"`)
	match := re.FindAllStringSubmatch(text, -1)

	pkgs := []string{}

	for _, m := range match {
		pkg := m[1]
		if strings.Contains(pkg, "/") {
			splitstr := strings.Split(pkg, "/")
			pkg = splitstr[len(splitstr)-1]
		}
		pkgs = append(pkgs, pkg)
	}

	return pkgs
}

func GetDeclarations(text string) []string {

	re := regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_]*(?:\s*,\s*[a-zA-Z_][a-zA-Z0-9_]*)*)\s*:=\s*.+?`)
	match := re.FindAllStringSubmatch(text, -1)

	decs := []string{}

	for _, m := range match {
		dec := m[1]
		if strings.Contains(dec, "/") {
			splitstr := strings.Split(dec, "/")
			dec = splitstr[len(splitstr)-1]
		}
		decs = append(decs, dec)
	}

	return decs

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

func GetStatementType(text string) (string, error) {

	// TODO: empty strings should be ignored.

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
