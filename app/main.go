package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"syscall"
)

func main() {
	chrootDir := "chroot"

	err := os.MkdirAll(chrootDir, 0755)
	if err != nil {
		fmt.Println("Err:", err)
		os.Exit(1)
	}

	err = os.MkdirAll(chrootDir+"/dev/null", 0755)
	if err != nil {
		fmt.Println("Err:", err)
		os.Exit(1)
	}

	command := os.Args[3]
	args := os.Args[4:len(os.Args)]

	err = copyExecutableIntoDir(command, chrootDir)
	if err != nil {
		fmt.Println("Err:", err)
		os.Exit(1)
	}

	// Change root directory to chroot directory
	if err := syscall.Chroot(chrootDir); err != nil {
		fmt.Printf("Failed to change root directory: %v\n", err)
		os.Exit(1)
	}
	if err := syscall.Chdir("/"); err != nil {
		fmt.Printf("Failed to change working directory: %v\n", err)
		os.Exit(1)
	}

	if err := syscall.Unshare(syscall.CLONE_NEWPID); err != nil {
		fmt.Printf("Failed to create new PID namespace: %s\n", err)
		os.Exit(1)
	}

	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Println("Exit code:", exitError.ExitCode())
			os.Exit(exitError.ExitCode())
		}
		fmt.Println("Err:", err)
	}
}

// eg. /usr/lib/docker-explorer -> chroot/usr/lib/docker-explorer
func copyExecutableIntoDir(chrootDir string, executablePath string) error {
	executablePathInChrootDir := path.Join(chrootDir, executablePath)

	if err := os.MkdirAll(path.Dir(executablePathInChrootDir), 0750); err != nil {
		return err
	}

	return copyFile(executablePath, executablePathInChrootDir)
}

func copyFile(sourceFilePath, destinationFilePath string) error {
	sourceFileStat, err := os.Stat(sourceFilePath)
	if err != nil {
		return err
	}

	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.OpenFile(destinationFilePath, os.O_RDWR|os.O_CREATE, sourceFileStat.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	return err
}
