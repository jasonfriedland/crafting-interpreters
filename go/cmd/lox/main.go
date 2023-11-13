package main

import (
	"fmt"
	"io"
	"os"

	"github.com/jasonfriedland/crafting-interpreters/pkg/scanner"
)

type parseError struct {
	line    int
	message string
}

func (e parseError) Error() string {
	return fmt.Sprintf("%s, line: %d", e.message, e.line)
}

func main() {
	var err error
	if len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, "error: invalid args\n")
		os.Exit(-1)
	}
	if len(os.Args) == 2 {
		err = runFile(os.Args[1])
	} else {
		err = runPrompt()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(-1)
	}
}

func runFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	return run(f)
}

func runPrompt() error {
	return run(nil)
}

func run(r io.Reader) error {
	s, err := scanner.New(r)
	if err != nil {
		return err
	}
	err = s.Scan()
	if err != nil {
		return parseError{message: err.Error(), line: s.Line()}
	}
	fmt.Println(s.Tokens())
	return nil
}
