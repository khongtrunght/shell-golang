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

type Command func(args []string, writer io.Writer)

func exitCommand(args []string, writer io.Writer) {
	if len(args) == 1 {
		if args[1] == "0" {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	} else {
		fmt.Fprintln(writer, "exit: invalid argument")
	}
}

func echoCommand(args []string, writer io.Writer) {
	if len(args) == 1 {
		fmt.Fprintln(writer, "")
	} else {
		fmt.Fprintln(writer, strings.Join(args[1:], " "))
	}
}

var commands = map[string]Command{
	"exit": exitCommand,
	"echo": echoCommand,
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

	if cmd, ok := commands[command]; ok {
		cmd(splits, s.writer)
		return
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
