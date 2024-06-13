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

func (s *Shell) Run() {
	fmt.Fprint(s.writer, "$ ")

	// Wait for user input
	commandName, err := bufio.NewReader(s.reader).ReadString('\n')
	commandName = strings.TrimSpace(commandName)
	if err != nil {
		log.Println("Error reading input: ", err)
	}

	// <command_name>: command not found
	fmt.Fprintf(s.writer, "%s: command not found\n", commandName)
}

func main() {
	shell := &Shell{
		reader: os.Stdin,
		writer: os.Stdout,
	}

	for {
		shell.Run()
	}
}
