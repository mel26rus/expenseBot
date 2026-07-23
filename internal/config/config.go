package config

import (
	"expense-bot/internal/filesystem"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	URL string `yaml:"url"`
}
type LogConfig struct {
	Folder  string `yaml:"folder"`
	Level   string `yaml:"level"`
	Console bool   `yaml:"console"`
	File    bool   `yaml:"file"`
}

type AppConfig struct {
	LogConfig LogConfig      `yaml:"log_config"`
	Database  DatabaseConfig `yaml:"database"`
	BotKey    string         `yaml:"bot_key"`
	AppName   string         `yaml:"app_name"`
}

const constConfigFileName = "config.yaml"

const DefaultDatabaseURL = "postgres://postgres:root@localhost:5432/expense_db"

func getConfigFile() string {
	return filepath.Join(
		filesystem.BaseDir(constConfigFileName),
		constConfigFileName,
	)
}

func defaultConfig() *AppConfig {
	return &AppConfig{
		LogConfig: LogConfig{
			Folder:  "Logs",
			Level:   "debug",
			Console: true,
			File:    true,
		},
		Database: DatabaseConfig{
			URL: DefaultDatabaseURL,
		},
		BotKey:  "",
		AppName: "ExpenseTGBot",
	}
}

func Load(fileName string) (*AppConfig, error) {

	configFile := getConfigFile()
	//slog.Info("LoadConfig", "configFile", configFile)
	data, err := os.ReadFile(configFile)
	if err != nil {
		defConfig := defaultConfig()
		Save(defConfig)
		return defConfig, err
	}

	var cfg AppConfig

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func Save(cfg *AppConfig) error {

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(getConfigFile(), data, 0644)
}
