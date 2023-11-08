package main

import (
	"fmt"
	"io"
	"os"
)

type parseError struct {
	line    int
	message string
}

func (e parseError) Error() string {
	return fmt.Sprintf("line: %d, error: %s", e.line, e.message)
}

func main() {
	var err error
	if len(os.Args) > 2 {
		fmt.Println("error: invalid args")
		os.Exit(-1)
	}
	if len(os.Args) == 2 {
		err = runFile(os.Args[1])
	} else {
		err = runPrompt()
	}
	if err != nil {
		fmt.Println("error:", err)
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
	fmt.Println("running")
	return nil
}
