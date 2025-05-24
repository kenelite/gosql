package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Addr    string // Server address, e.g., ":3306"
	DataDir string // Directory for storing table files
}

func Load() *Config {
	file, err := os.Open("settings/server.conf")
	if err != nil {
		panic(fmt.Sprintf("failed to open config file: %v", err))
	}
	defer file.Close()

	port := "3306"
	dataDir := "data"

	scanner := bufio.NewScanner(file)
	inServerSection := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			inServerSection = line == "[server]"
			continue
		}
		if !inServerSection {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		switch key {
		case "port":
			port = val
		case "data_path":
			dataDir = val
		}
	}

	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("failed to read config file: %v", err))
	}

	return &Config{
		Addr:    ":" + port,
		DataDir: dataDir,
	}
}
