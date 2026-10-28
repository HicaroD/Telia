package config

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// TODO: properly map these paths to user
var DEFAULT_ENV_FILE string = `T_STD_PATH=
`

//go:embed env
var DEFAULT_DEV_ENV_FILE string

var TELIA_CONFIG_DIR string

var ENVS *Envs

type Envs struct {
	STD string `env:"T_STD_PATH"`
}

// VendorPath follows the std installation layout unless explicitly overridden.
func VendorPath() string {
	if path := os.Getenv("T_VENDOR_PATH"); path != "" {
		return path
	}
	return filepath.Join(filepath.Dir(ENVS.STD), "vendor")
}

func (e *Envs) ShowAll() {
	v := reflect.ValueOf(e)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := range v.NumField() {
		field := v.Type().Field(i)
		fieldValue := v.Field(i)

		envTag := field.Tag.Get("env")
		if envTag != "" {
			fmt.Printf("%s='%s'\n", envTag, fieldValue.String())
		}
	}
}

func SetupConfigDir() error {
	teliaCfgDir, err := getConfigDir("telia")
	if err != nil {
		return err
	}
	TELIA_CONFIG_DIR = teliaCfgDir
	return nil
}

func SetupEnvFile() error {
	envFile := filepath.Join(TELIA_CONFIG_DIR, "env")
	envs, err := loadTeliaEnvFile(envFile)
	if err != nil {
		return err
	}

	parsedEnvs := Envs{}
	if value := os.Getenv("T_STD_PATH"); value != "" {
		envs["T_STD_PATH"] = value
	}
	err = MapEnvToStruct(envs, &parsedEnvs)
	if err != nil {
		return err
	}
	if parsedEnvs.STD == "" {
		executable, err := os.Executable()
		if err != nil {
			return err
		}
		if resolved, err := filepath.EvalSymlinks(executable); err == nil {
			executable = resolved
		}
		bundled := filepath.Join(filepath.Dir(executable), "base", "std")
		if info, err := os.Stat(bundled); err == nil && info.IsDir() {
			parsedEnvs.STD = bundled
		}
	}

	ENVS = &parsedEnvs
	return nil
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
	defer func() { _ = file.Close() }() // Read-only handle; no buffered writes to flush.

	// NOTE: if development mode and env file already exists, recreate
	// it because developers might have changed the environment variables
	// for debugging
	if DEV && !envFileCreated {
		envFileCreated = true
	}

	if envFileCreated {
		var err error
		if DEV {
			err = writeStringToFile(path, DEFAULT_DEV_ENV_FILE)
		} else {
			err = writeStringToFile(path, DEFAULT_ENV_FILE)
		}
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

func writeStringToFile(fileName, content string) error {
	// os.O_CREATE: Create the file if it doesn't exist.
	// os.O_WRONLY: Open the file for writing only.
	// os.O_TRUNC: Truncate the file if it already exists (overwrite it).
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	_, err = file.WriteString(content)
	closeErr := file.Close()
	if err != nil {
		return err
	}

	return closeErr
}

func MapEnvToStruct(data map[string]string, result any) error {
	v := reflect.ValueOf(result).Elem()
	t := v.Type()

	for i := range t.NumField() {
		field := t.Field(i)
		fieldValue := v.Field(i)

		envTag := field.Tag.Get("env")
		if envTag != "" {
			if value, ok := data[envTag]; ok {
				if fieldValue.CanSet() && fieldValue.Kind() == reflect.String {
					fieldValue.SetString(value)
				}

				if strings.HasSuffix(envTag, "PATH") && value != "" {
					fullPath := value

					if value == "~" || strings.HasPrefix(value, "~/") {
						homeDir, err := os.UserHomeDir()
						if err != nil {
							return err
						}
						fullPath = filepath.Join(homeDir, strings.TrimPrefix(value, "~/"))
						if value == "~" {
							fullPath = homeDir
						}
					}

					_, err := os.Stat(fullPath)
					if err != nil {
						return err
					}

					fieldValue.SetString(fullPath)
				}
			}
		}
	}

	return nil
}
