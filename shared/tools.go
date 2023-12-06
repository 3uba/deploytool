package shared

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type ProjectConfig struct {
	RootDir string
	Name    string
	User    string
	Token   string
	GitURL  string
}

func ReadProjectConfigFile(filename string, projectName string) (ProjectConfig, error) {
	config := ProjectConfig{}

	file, err := os.Open(filename)
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

func WriteProjectConfigFile(filename string, config ProjectConfig) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	formatAndWrite := func(key, value string) error {
		_, err := file.WriteString(fmt.Sprintf("\n%s_%s=%s", config.Name, key, value))
		return err
	}

	_, err = file.WriteString(fmt.Sprintf("# %s configuration", config.Name))

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
