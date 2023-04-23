package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	chrootDir := "chroot"

	err := os.Mkdir(chrootDir, 0755)
	if err != nil {
		fmt.Println("Err: %v", err)
		os.Exit(1)
	}

	err = os.MkdirAll(chrootDir+"/usr/local/bin", 0755)
	err = os.MkdirAll(chrootDir+"/dev/null", 0755)
	err = exec.Command("cp", "/usr/local/bin/docker-explorer", chrootDir+"/usr/local/bin").Run()
	//err = exec.Command("cp", "/usr/bin/ls", chrootDir+"/usr/local/bin").Run()

	command := os.Args[3]
	args := os.Args[4:len(os.Args)]

	// Change root directory to chroot directory
	if err := syscall.Chroot(chrootDir); err != nil {
		fmt.Printf("Failed to change root directory: %v\n", err)
		os.Exit(1)
	}
	if err := syscall.Chdir("/"); err != nil {
		fmt.Printf("Failed to change working directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Command:", command)
	fmt.Println("Args:", args)

	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println(cmd.Args)

	err = cmd.Run()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Println("Exit code:", exitError.ExitCode())
			os.Exit(exitError.ExitCode())
		}
		fmt.Println("Err: %v", err)
	}
}
