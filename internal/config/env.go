package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type EnvData struct {
	RawFlags []string
}

func LoadEnv(scriptDir string) (*EnvData, error) {
	path := filepath.Join(scriptDir, ".env")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("no .env found, copy .env.example to .env and fill it in")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := &EnvData{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "--") {
			data.RawFlags = append(data.RawFlags, line)
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if found {
			os.Setenv(key, strings.Trim(value, `"`))
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return data, nil
}
