package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	URL string `yaml:"url"`
}

type AppConfig struct {
	LogFile  string         `yaml:"log_file"`
	Database DatabaseConfig `yaml:"database"`
	IsDebug  bool           `yaml:"is_debug"`
	BotKey   string         `yaml:"bot_key"`
}

const DefaultDatabaseURL = "postgres://postgres:root@localhost:5432/expense_db"

func workDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}

func configFile() string {
	return filepath.Join(workDir(), "config.yaml")
}

func logFile(name string) string {
	const constLogFolderName = "Logs"
	dir := filepath.Join(workDir(), constLogFolderName)
	_ = os.MkdirAll(dir, 0755)
	return filepath.Join(workDir(), constLogFolderName, name)
}

func defaultConfig() *AppConfig {
	return &AppConfig{
		LogFile: "expense-bot.log",
		Database: DatabaseConfig{
			URL: DefaultDatabaseURL,
		},
		IsDebug: true,
		BotKey:  "",
	}
}

func Load(fileName string) (*AppConfig, error) {

	//	fileName := configFile()

	data, err := os.ReadFile(fileName)
	if err != nil {
		cfg := defaultConfig()

		if err := Save(cfg); err != nil {
			return nil, err
		}

		cfg.LogFile = logFile(cfg.LogFile)

		return cfg, nil
	}

	cfg := defaultConfig()

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	cfg.LogFile = logFile(cfg.LogFile)

	return cfg, nil
}

func Save(cfg *AppConfig) error {

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(configFile(), data, 0644)
}
