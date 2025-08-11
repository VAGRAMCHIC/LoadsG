package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Id     string `json:"id"`
	Key    string `json:"key"`
	JwtKey string `json:"jwtKey"`
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
