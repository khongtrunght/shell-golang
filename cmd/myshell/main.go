package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
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
				continue
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

func changeDirCommand(args []string, writer io.Writer) {
	switch len(args) {
	case 1:
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(writer, "cd: error: ", err)
			return
		}
		err = os.Chdir(homeDir)
		if err != nil {
			fmt.Fprintf(writer, "%s: No such file or directory\n", homeDir)
		}
	case 2:
		err := os.Chdir(args[1])
		if err != nil {
			fmt.Fprintf(writer, "%s: No such file or directory\n", args[1])
		}
	default:
		fmt.Fprintln(writer, "cd: too many arguments")
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
	commands["cd"] = changeDirCommand
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

	if _, ok := excutableFiles[command]; ok {
		// run program with args
		execCmd := exec.Command(excutableFiles[command], splits[1:]...)
		execCmd.Stdout = s.writer
		execCmd.Stderr = s.writer
		err := execCmd.Run()
		if err != nil {
			log.Println("Error running command: ", err)
		}
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
