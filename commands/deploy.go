package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/3uba/deploytool/shared"
)

const configFile = ".config"

func prompt(question string) string {
	fmt.Print(question + ": ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func createFolders(projectPath string) error {
	folders := []string{
		filepath.Join(projectPath),
		filepath.Join(projectPath, "current"),
		filepath.Join(projectPath, "backup"),
		filepath.Join(projectPath, "shared"),
	}

	for _, folder := range folders {
		if err := os.MkdirAll(folder, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func createSymbolicLink(projectPath string) error {
	sharedPath := filepath.Join(projectPath, "shared")
	currentPath := filepath.Join(projectPath, "current")

	envFile := filepath.Join(sharedPath, ".env")
	if _, err := os.Stat(envFile); err == nil {
		linkName := filepath.Join(currentPath, ".env")
		if err := os.Symlink(envFile, linkName); err != nil {
			return fmt.Errorf("nie można utworzyć linku symbolicznego dla pliku .env: %v", err)
		}

		fmt.Printf("Utworzono link symboliczny dla pliku .env między %s a %s\n", envFile, linkName)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("nie można uzyskać dostępu do pliku .env: %v", err)
	}

	return nil
}

func gitClone(config shared.ProjectConfig, destination string) error {
	authURL := fmt.Sprintf("https://%s:%s@%s", config.User, config.Token, config.GitURL)

	cmd := exec.Command("git", "clone", authURL, destination)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func checkGitAccess(config shared.ProjectConfig, projectPath string) error {
	cmd := exec.Command("git", "ls-remote", fmt.Sprintf("https://%s:%s@%s", config.User, config.Token, config.GitURL))

	cmd.Env = append(os.Environ(), "GIT_ASKPASS=echo")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func moveCurrentToBackup(projectPath string) error {
	backupFolder := filepath.Join(projectPath, "backup")

	latestBackup, err := getLatestBackupNumber(backupFolder)
	if err != nil {
		return err
	}

	newBackupFolder := filepath.Join(backupFolder, fmt.Sprintf("%d", latestBackup+1))
	if err := os.MkdirAll(newBackupFolder, os.ModePerm); err != nil {
		return err
	}

	files, err := filepath.Glob(filepath.Join(projectPath, "current", "*"))
	if err != nil {
		return err
	}

	for _, file := range files {
		filename := filepath.Base(file)
		dest := filepath.Join(newBackupFolder, filename)

		if err := os.Rename(file, dest); err != nil {
			return err
		}
	}

	return nil
}

func getLatestBackupNumber(backupFolder string) (int, error) {
	files, err := filepath.Glob(filepath.Join(backupFolder, "*"))
	if err != nil {
		return 0, err
	}

	latestBackup := 0
	for _, file := range files {
		backupName := filepath.Base(file)
		backupNumber, err := strconv.Atoi(backupName)
		if err == nil && backupNumber > latestBackup {
			latestBackup = backupNumber
		}
	}

	return latestBackup, nil
}

func Deploy(projectName string) {
	config, err := shared.ReadProjectConfigFile(configFile, projectName)
	if err != nil {
		fmt.Printf("Error reading project configuration: %v\n", err)
		return
	}

	projectPath := filepath.Join(config.RootDir, projectName)

	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		if err := createFolders(projectPath); err != nil {
			fmt.Printf("Error creating project folders: %v\n", err)
			return
		}
	}

	if err := checkGitAccess(config, projectPath); err != nil {
		fmt.Printf("Error checking git access: %v\n", err)
		return
	}

	if err := moveCurrentToBackup(projectPath); err != nil {
		fmt.Printf("Error moving content from current to backup: %v\n", err)
		return
	}

	if err := gitClone(config, filepath.Join(projectPath, "current")); err != nil {
		fmt.Printf("Error cloning git repository: %v\n", err)
		return
	}

	err = createSymbolicLink(projectPath)
	if err != nil {
		fmt.Printf("Error creating symbolic links: %v\n", err)
		return
	}
}
