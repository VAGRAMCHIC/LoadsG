package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Id            string `json:"id"`
	Key           string `json:"key"`
	JwtKey        string `json:"jwtKey"`
	MaxConcurrent int    `json:"MaxConcurrent"`
	PgConn        string `json:"pgConn"`
}

func ReadOSENV() (Config, error) {
	var config Config
	config.Id = os.Getenv("ID")
	config.Key = os.Getenv("KEY")
	config.JwtKey = os.Getenv("JWT_KEY")
	config.MaxConcurrent, _ = strconv.Atoi(os.Getenv("MAX_CONCURRENT"))
	config.PgConn = os.Getenv("PG_CONN")
	if config.JwtKey == "" || config.PgConn == "" || config.Key == "" {
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
