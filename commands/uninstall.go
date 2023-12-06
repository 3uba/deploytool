package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/3uba/deploytool/shared"
)

func UninstallDeploytool() error {
	fmt.Print("Are you sure you want to uninstall deploytool? (y/n): ")
	var answer string
	fmt.Scanln(&answer)

	if strings.ToLower(answer) != "y" {
		fmt.Println("Uninstall aborted.")
		return nil
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Error getting current directory: %v", err)
	}

	if err := os.Chdir(os.Getenv("HOME")); err != nil {
		return fmt.Errorf("Error changing directory to home: %v", err)
	}

	err = shared.RunCommand("sed", "-i", "/\\/usr\\/local\\/bin\\/deploytool\\/app/d", ".bashrc")
	if err != nil {
		return fmt.Errorf("Error removing deploytool from .bashrc: %v", err)
	}

	err = shared.RunCommand("sed", "-i", "/DT_PATH=\\/usr\\/local\\/bin\\/deploytool/d", ".bashrc")
	if err != nil {
		return fmt.Errorf("Error removing DT_PATH from .bashrc: %v", err)
	}

	if err := os.Chdir(currentDir); err != nil {
		return fmt.Errorf("Error changing back to the original directory: %v", err)
	}

	err = os.RemoveAll("/usr/local/bin/deploytool")
	if err != nil {
		return fmt.Errorf("Error removing deploytool directory: %v", err)
	}

	fmt.Println("Deploytool has been uninstalled.")
	return nil
}
