package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	fmt.Fprint(os.Stdout, "$ ")

	// Wait for user input
	commandName, err := bufio.NewReader(os.Stdin).ReadString('\n')
	commandName = strings.TrimSpace(commandName)
	if err != nil {
		log.Println("Error reading input: ", err)
	}

	// <command_name>: command not found
	fmt.Fprintf(os.Stdout, "%s: command not found\n", commandName)
}
