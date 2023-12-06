package shared

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ProjectConfig struct {
	RootDir string
	Name    string
	User    string
	Token   string
	GitURL  string
	Port    string
}

func ReadProjectConfigFile(projectName string) (ProjectConfig, error) {
	config := ProjectConfig{}

	dtLocation := os.Getenv("DT_PATH")
	if dtLocation == "" {
		fmt.Println("DT_PATH environment variable not set.")
		dtLocation = "."
		fmt.Println("Using current directory.")
	}

	configFilePath := filepath.Join(dtLocation, ".config")

	file, err := os.Open(configFilePath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "root_dir") {
			value := strings.TrimPrefix(line, "root_dir=")
			config.RootDir = strings.TrimSpace(value)
		}

		if strings.HasPrefix(line, projectName+"_name=") {
			value := strings.TrimPrefix(line, projectName+"_name=")
			config.Name = strings.TrimSpace(value)
		}

		if strings.HasPrefix(line, projectName+"_user=") {
			value := strings.TrimPrefix(line, projectName+"_user=")
			config.User = strings.TrimSpace(value)
		}

		if strings.HasPrefix(line, projectName+"_token=") {
			value := strings.TrimPrefix(line, projectName+"_token=")
			config.Token = strings.TrimSpace(value)
		}

		if strings.HasPrefix(line, projectName+"_git_url=") {
			value := strings.TrimPrefix(line, projectName+"_git_url=")
			config.GitURL = strings.TrimSpace(value)
		}
	}

	if err := scanner.Err(); err != nil {
		return config, err
	}

	return config, nil
}

func WriteProjectConfigFile(config ProjectConfig) error {
	dtLocation := os.Getenv("DT_PATH")
	if dtLocation == "" {
		fmt.Println("DT_PATH environment variable not set.")
		dtLocation = "."
		fmt.Println("Using current directory.")
	}

	configFilePath := filepath.Join(dtLocation, ".config")

	file, err := os.OpenFile(configFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("\n\n# %s configuration", config.Name))

	formatAndWrite := func(key, value string) error {
		_, err := file.WriteString(fmt.Sprintf("\n%s_%s=%s", config.Name, key, value))
		return err
	}

	if err := formatAndWrite("user", config.User); err != nil {
		return err
	}
	if err := formatAndWrite("token", config.Token); err != nil {
		return err
	}
	if err := formatAndWrite("git_url", config.GitURL); err != nil {
		return err
	}

	return nil
}

func RunCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
