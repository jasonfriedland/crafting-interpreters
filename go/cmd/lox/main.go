package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("error: invalid args")
		os.Exit(-1)
	} else if len(os.Args) == 2 {
		runFile()
	} else {
		runPrompt()
	}
}

func runFile() {
	run()
}

func runPrompt() {
	run()
}

func run() {}
