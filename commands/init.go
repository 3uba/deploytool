package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func Init() {
	dtLocation := os.Getenv("DT_PATH")
	if dtLocation == "" {
		fmt.Println("DT_PATH environment variable not set. Using current directory.")
		dtLocation = "."
	}

	configFilePath := filepath.Join(dtLocation, ".config")

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		fmt.Println("The .config file doesn't exist.")
		fmt.Println("Please provide the root directory where all projects should be created:")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		rootDir := scanner.Text()

		err := os.MkdirAll(dtLocation, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating deploytool directory: %v\n", err)
			return
		}

		file, err := os.Create(configFilePath)
		if err != nil {
			fmt.Printf("Error creating .config file: %v\n", err)
			return
		}
		defer file.Close()

		_, err = file.WriteString(fmt.Sprintf("root_dir=%s\n", rootDir))
		if err != nil {
			fmt.Printf("Error writing to .config file: %v\n", err)
			return
		}

		fmt.Printf("Configuration saved. Projects will be created in %s\n", rootDir)
	} else if err != nil {
		fmt.Printf("Error checking .config file: %v\n", err)
		return
	} else {
		fmt.Println("The .config file already exists. No further action needed.")
	}
}
