package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	RootUser      string `json:"rootUser"`
	RootToken     string `json:"rootToken"`
	JwtKey        string `json:"jwtKey"`
	JwtRefreshKey string `json:"jwtRefreshKey"`
	AppName       string `json:"appName"`
	MaxConcurrent int    `json:"MaxConcurrent"`
	PgConn        string `json:"pgConn"`
}

func ReadOSENV() (Config, error) {
	var config Config
	config.RootUser = os.Getenv("ROOT_USER")
	config.RootToken = os.Getenv("ROOT_TOKEN")
	config.JwtKey = os.Getenv("JWT_KEY")
	config.JwtRefreshKey = os.Getenv("JWT_REFRESH_KEY")
	config.AppName = os.Getenv("APP_NAME")
	config.MaxConcurrent, _ = strconv.Atoi(os.Getenv("MAX_CONCURRENT"))
	config.PgConn = os.Getenv("PG_CONN")
	if config.JwtKey == "" || config.PgConn == "" || config.RootToken == "" {
		return config, errors.New("cant read config envs")
	}
	return config, nil
}

func ReadConfig(filename string) (Config, error) {
	var config Config

	file, err := os.Open(filename)
	if err != nil {
		return config, fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return config, fmt.Errorf("ошибка разбора JSON: %w", err)
	}

	return config, nil
}
