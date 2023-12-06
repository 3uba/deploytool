package main

import (
	"fmt"
	"os"

	"github.com/3uba/deploytool/commands"
	_ "github.com/3uba/deploytool/shared"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "deploytool"
	app.Usage = "Tool for deploying and creating projects"

	app.Commands = []cli.Command{
		{
			Name:    "deploy",
			Aliases: []string{"d"},
			Usage:   "Deploy a project",
			Action: func(c *cli.Context) error {
				projectName := c.Args().First()
				if projectName == "" {
					return fmt.Errorf("Brak nazwy projektu. Użycie: ./deploytool deploy project_name")
				}
				commands.Deploy(projectName)
				return nil
			},
		},
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "Create something",
			Action: func(c *cli.Context) error {
				commands.Create()
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
