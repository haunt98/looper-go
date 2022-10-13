package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	commandsFile  = flag.String("commands-file", "", "commands file path, which contains a list of command to run")
	isClearScreen = flag.Bool("clear-screen", false, "clear screen between each command")
	waitDuration  = flag.Duration("wait-duration", 5*time.Second, "wait duration between each command")
)

func main() {
	flag.Parse()

	if *commandsFile == "" {
		fmt.Println("commands-file is required")
		fmt.Println("Example: go run main.go --commands-file=commands.txt")
		return
	}

	commandsFileBytes, err := os.ReadFile(*commandsFile)
	if err != nil {
		fmt.Printf("os: failed to read file: %s\n", err)
		return
	}

	rawCommands := strings.Split(strings.TrimSpace(string(commandsFileBytes)), "\n")
	commands := make([][2]string, 0, len(rawCommands))

	for _, rawCommand := range rawCommands {
		rawCommand = strings.TrimSpace(rawCommand)
		if rawCommand == "" {
			continue
		}

		// Command expect to be in the format of "command arg1 arg2 ..."
		before, after, found := strings.Cut(rawCommand, " ")
		if !found {
			continue
		}

		commands = append(commands, [2]string{before, after})
	}

	// For is forever
	i := 0
	for {
		i = i % len(commands)
		before, after := commands[i][0], commands[i][1]
		i += 1

		cmd := exec.Command(before, after)
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			fmt.Printf("exec: failed to run command %s %s: %s\n", before, after, err)
			continue
		}

		time.Sleep(*waitDuration)

		if *isClearScreen {
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			if err := cmd.Run(); err != nil {
				fmt.Printf("exec: failed to run command clear: %s\n", err)
				continue
			}
		}
	}
}
