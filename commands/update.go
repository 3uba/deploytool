package commands

import (
	"fmt"
	"os"

	"github.com/3uba/deploytool/shared"
)

func UpdateDeploytool() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Error getting current directory: %v", err)
	}

	fmt.Println("Updating deploytool...")

	err = shared.RunCommand("cd", "/usr/local/bin/deploytool", "&&", "git", "pull")

	fmt.Println("Test")
	if err != nil {
		return fmt.Errorf("Error updating deploytool: %v", err)
	}
	fmt.Println("Test2")
	err = shared.RunCommand("cd", currentDir)
	if err != nil {
		return fmt.Errorf("Error changing back to the original directory: %v", err)
	}
	fmt.Println("Test3")
	fmt.Println("Deploytool has been updated.")
	return nil
}
