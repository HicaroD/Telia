package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var DEFAULT_ENV_FILE string = `T_STD=/usr/local/telia/std
T_RUNTIME=/usr/local/telia/runtime
`

func SetupConfigDir() (map[string]string, error) {
	teliaCfgDir, err := getConfigDir(APP_NAME)
	if err != nil {
		return nil, err
	}

	envFile := filepath.Join(teliaCfgDir, ENV_FILE)
	envs, err := loadTeliaEnvFile(envFile)
	if err != nil {
		return nil, err
	}

	return envs, nil
}

func getConfigDir(appName string) (string, error) {
	var configDir string

	if configHome := os.Getenv("XDG_CONFIG_HOME"); configHome != "" {
		configDir = filepath.Join(configHome, appName)
	} else if homeDir, err := os.UserHomeDir(); err == nil {
		if os.Getenv("OS") == "Windows_NT" {
			configDir = filepath.Join(os.Getenv("APPDATA"), appName)
		} else {
			configDir = filepath.Join(homeDir, ".config", appName)
		}
	} else {
		return "", fmt.Errorf("could not determine home directory")
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return configDir, nil
}

func loadTeliaEnvFile(path string) (map[string]string, error) {
	env := make(map[string]string)
	var envFileCreated bool

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			newFile, err := os.Create(path)
			if err != nil {
				return nil, err
			}
			file = newFile
			envFileCreated = true
		} else {
			return nil, err
		}
	}
	defer file.Close()

	if envFileCreated {
		err := writeStringToFile(path, DEFAULT_ENV_FILE)
		if err != nil {
			return nil, err
		}
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return env, nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func copyDir(srcDir, dstDir string) error {
	srcDirInfo, err := os.Stat(srcDir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dstDir, srcDirInfo.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		dstPath := filepath.Join(dstDir, entry.Name())

		if entry.IsDir() {
			err := copyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err := copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func writeStringToFile(fileName, content string) error {
	// os.O_CREATE: Create the file if it doesn't exist.
	// os.O_WRONLY: Open the file for writing only.
	// os.O_TRUNC: Truncate the file if it already exists (overwrite it).
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
