package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Shell struct {
	reader io.Reader
	writer io.Writer
}

func (s *Shell) Repl() {
	for {
		s.Run()
	}
}

func (s *Shell) Run() {
	fmt.Fprint(s.writer, "$ ")
	var command string
	// Wait for user input
	userInput, err := bufio.NewReader(s.reader).ReadString('\n')
	userInput = strings.TrimSpace(userInput)
	if err != nil {
		log.Println("Error reading input: ", err)
	}

	splits := strings.Split(userInput, " ")
	if len(splits) == 0 {
		return
	} else {
		command = splits[0]
	}

	if command == "exit" {
		if len(splits) == 2 {
			if splits[1] == "0" {
				os.Exit(0)
			} else {
				os.Exit(1)
			}
		} else {
			fmt.Fprintln(s.writer, "exit: invalid argument")
		}
	}

	// <command_name>: command not found
	fmt.Fprintf(s.writer, "%s: command not found\n", userInput)
}

func main() {
	shell := &Shell{
		reader: os.Stdin,
		writer: os.Stdout,
	}

	shell.Repl()
}
