package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func ReadFile(fileName string) *[]byte {
	// identify current environment and set file path according to that
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	var fullPath string

	if env == "local" {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		fullPath = filepath.Join(cwd, "ops", fileName)
	} else {
		// identify current file directory
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			panic(err)
		}
		fullPath = filepath.Join(dir, "ops", fileName)
	}
	content, err := os.ReadFile(fullPath)
	if err != nil {
		panic(fmt.Sprintf("error: %v", err))
	}

	return &content
}

func WriteFile(filePath string) bool {
	// attempt to create the file only if it does not exists
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			return true
		} else {
			panic(fmt.Sprintf("failed to create file. %v", err))
		}
	}
	defer file.Close()

	fmt.Println("File was created successfully!")

	return true
}

// check directory if exists
func checkIfFolderExists(filePath string) (bool, error) {
	if _, err := os.Stat(filePath); err == nil {
		// director exists
		return true, nil
	} else if os.IsNotExist(err) {
		// directory does not exists
		return false, nil
	} else {
		return false, err
	}
}

// create a directory if not exists
func CreateDirectoryIfNotExist(dirPath string) error {
	folderExist, err := checkIfFolderExists(dirPath)
	if err != nil {
		panic(fmt.Sprintf("Error checking existing folder. %v", err))
	}
	if !folderExist {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			panic(fmt.Sprintf("error creating directory. %v", err))
		}
		fmt.Println("Nested directories created: ", dirPath)
	}
	return nil
}

// handle permission for specific directory
func HandlePermissionDirectory(dirPath string) error {
	if err := os.Chmod(dirPath, 0755); err != nil {
		panic(fmt.Sprintf("error changing directory permissions. %v", err))
	}
	return nil
}
