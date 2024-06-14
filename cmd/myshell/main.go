package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Shell struct {
	reader io.Reader
	writer io.Writer
}

func (s *Shell) Repl() {
	// load path from env
	path := os.Getenv("PATH")

	if path != "" {
		pathList := strings.Split(path, ":")
		for _, p := range pathList {
			// get all executable files in the path
			files, err := os.ReadDir(p)
			if err != nil {
				log.Println("Error reading directory: ", err)
			}
			for _, file := range files {
				if !file.IsDir() {
					info, err := file.Info()
					if err != nil {
						log.Println("Error reading file info: ", err)
					}
					if info.Mode()&0111 != 0 {
						excutableFiles[file.Name()] = filepath.Join(p, file.Name())
					}
				}
			}
		}
	}
	for {
		s.Run()
	}
}

type Command func(args []string, writer io.Writer)

func exitCommand(args []string, writer io.Writer) {
	if len(args) == 2 {
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

var (
	commands       = map[string]Command{}
	excutableFiles = make(map[string]string)
)

func init() {
	commands["exit"] = exitCommand
	commands["echo"] = echoCommand
	commands["type"] = typeCommand
}

func typeCommand(args []string, writer io.Writer) {
	if len(args) == 1 {
		fmt.Fprintln(writer, "type: missing argument")
	} else if len(args) == 2 {
		if _, ok := commands[args[1]]; ok {
			fmt.Fprintf(writer, "%s is a shell builtin\n", args[1])
		} else if _, ok := excutableFiles[args[1]]; ok {
			fmt.Fprintf(writer, "%s is %s\n", args[1], excutableFiles[args[1]])
		} else {
			fmt.Fprintf(writer, "%s: not found\n", args[1])
		}
	} else {
		fmt.Fprintln(writer, "type: too many arguments")
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
