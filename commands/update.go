package commands

import (
	"fmt"
	"os"
	"os/exec"
)

func UpdateDeploytool() error {
	fmt.Println("Updating deploytool...")

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Error getting current directory: %v", err)
	}

	if err := os.Chdir("/usr/local/bin/deploytool"); err != nil {
		return fmt.Errorf("Error changing directory: %v", err)
	}

	cmd := exec.Command("git", "pull")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error updating deploytool: %v", err)
	}

	if err := os.Chdir(currentDir); err != nil {
		return fmt.Errorf("Error changing back to the original directory: %v", err)
	}

	return nil
}
