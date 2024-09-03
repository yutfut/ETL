package config

import (
	"encoding/json"
	"os"
)

type Secret struct {
	Postgres []struct {
		Host         string `json:"host"`
		Port         string `json:"port"`
		Username     string `json:"username"`
		Password     string `json:"password"`
		Database     string `json:"database"`
		PostgreSQLID string `json:"postgresql_id"`
	} `json:"postgres"`
	ClickHouse struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Database string `json:"database"`
		Username string `json:"username"`
		Password string `json:"password"`
		Debug    bool   `json:"debug"`
	} `json:"clickHouse"`
}

func LoadConfig(path string) (Secret, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Secret{}, err
	}

	response := Secret{}

	if err = json.Unmarshal(data, &response); err != nil {
		return Secret{}, err
	}

	return response, nil
}
