package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/3uba/deploytool/shared"
)

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

		if _, err := os.Stat(linkName); err == nil {
			if err := os.Remove(linkName); err != nil {
				return fmt.Errorf("Unable to remove .env file: %v", err)
			}
		}

		if err := os.Symlink(envFile, linkName); err != nil {
			return fmt.Errorf("Unable to create symbolic link for .env file: %v", err)
		}

		fmt.Printf("Created symbolic link for .env file between %s and %s\n", envFile, linkName)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("Unable to access .env file: %v", err)
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

func runByDocker(config shared.ProjectConfig, projectPath string) error {
	composeFile := filepath.Join(projectPath, "docker-compose.yaml")
	if _, err := os.Stat(composeFile); err == nil {

		fmt.Print("Do you want to deploy using docker-compose? (yes/no): ")
		var answer string
		fmt.Scanln(&answer)

		if strings.ToLower(answer) != "yes" {
			return nil
		}

		cmd := exec.Command("docker", "compose", "down")
		cmd.Dir = projectPath
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error running docker-compose: %v\n", err)
			return err
		}

		cmd = exec.Command("docker", "compose", "up", "-d", "--build")
		cmd.Dir = projectPath
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Error running docker-compose: %v\n", err)
			return err
		}
		fmt.Println("Deployment successful using docker-compose.")
		return nil
	}

	dockerfile := filepath.Join(projectPath, "Dockerfile")
	if _, err := os.Stat(dockerfile); err == nil {

		fmt.Print("Do you want to deploy using dockerfile? (yes/no): ")
		var answer string
		fmt.Scanln(&answer)

		if strings.ToLower(answer) != "yes" {
			return nil
		}

		imageName := fmt.Sprintf("%s-image", config.Name)
		if imageExists(imageName) {
			fmt.Printf("Removing existing image: %s\n", imageName)
			removeDockerImage(imageName)
		}

		projectPort := config.Port
		cmd := exec.Command("docker", "build", "-t", imageName, ".")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = projectPath
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error building Docker image: %v\n", err)
			return err
		}

		cmd = exec.Command("docker", "run", "--name", config.Name, "-dp", fmt.Sprintf("%s:%s", projectPort, projectPort), imageName)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = projectPath
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Error running Docker container: %v\n", err)
			return err
		}
		fmt.Println("Deployment successful using Dockerfile.")
		return nil
	}

	fmt.Println("No docker-compose.yaml or Dockerfile found. Deployment aborted.")
	return nil
}

func imageExists(imageName string) bool {
	cmd := exec.Command("docker", "image", "inspect", imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err == nil
}

func removeDockerImage(imageName string) {
	cmd := exec.Command("docker", "image", "rm", imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error removing Docker image: %v\n", err)
	}
}

func Deploy(projectName string) {
	config, err := shared.ReadProjectConfigFile(projectName)
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

	err = runByDocker(config, projectPath+"/current")
	if err != nil {
		fmt.Printf("Error running by Docker: %v\n", err)
		return
	}
}
