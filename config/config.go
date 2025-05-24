package config

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type ServerConfig struct {
	Port     string
	DataPath string
	Users    map[string]string // Multiple users
}

func LoadConfig(path string) (*ServerConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	cfg := &ServerConfig{
		Port:     "3306", // default port
		DataPath: "data", // default data path
	}
	inServerSection := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section := strings.ToLower(strings.TrimSpace(line[1 : len(line)-1]))
			inServerSection = (section == "server")
			continue
		}
		if !inServerSection {
			continue
		}
		kv := strings.SplitN(line, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		switch key {
		case "port":
			cfg.Port = value
		case "data_path":
			cfg.DataPath = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if cfg.Port == "" || cfg.DataPath == "" {
		return nil, errors.New("missing required server config")
	}

	users, err := LoadUsers("settings/users.conf")
	if err != nil {
		return nil, err
	}
	cfg.Users = users

	return cfg, nil
}

func LoadUsers(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	users := make(map[string]string)
	scanner := bufio.NewScanner(file)
	inUsers := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section := strings.ToLower(strings.TrimSpace(line[1 : len(line)-1]))
			inUsers = (section == "users")
			continue
		}
		if !inUsers {
			continue
		}
		kv := strings.SplitN(line, "=", 2)
		if len(kv) != 2 {
			continue
		}
		user := strings.TrimSpace(kv[0])
		pass := strings.TrimSpace(kv[1])
		users[user] = pass
	}

	return users, scanner.Err()
}
