package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	command := os.Args[3]
	args := os.Args[4:len(os.Args)]

	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Println("Exit code:", exitError.ExitCode())
			os.Exit(exitError.ExitCode())
		}
		fmt.Println("Err: %v", err)
	}
}
