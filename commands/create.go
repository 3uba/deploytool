package commands

import (
	"fmt"
	"github.com/3uba/deploytool/shared"
)

func Create() {
	var name, user, token, gitURL string

	fmt.Print("Enter project name: ")
	fmt.Scanln(&name)

	fmt.Print("Enter user (press Enter to skip): ")
	fmt.Scanln(&user)

	fmt.Print("Enter token (press Enter to skip): ")
	fmt.Scanln(&token)

	fmt.Print("Enter Git URL: ")
	fmt.Scanln(&gitURL)

	newConfig := shared.ProjectConfig{
		Name:   name,
		User:   user,
		Token:  token,
		GitURL: gitURL,
	}

    err := shared.WriteProjectConfigFile(".config", newConfig)
    if err != nil {
        fmt.Println("Error writing project config:", err)
    } else {
        fmt.Println("Project config written successfully.")
    }
}
