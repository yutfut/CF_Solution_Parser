package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Files struct {
		Input  string `json:"input"`
		Output string `json:"output"`
	} `json:"files"`
	Workers struct {
		WorkerCount  int `json:"workerCount"`
		InputChanel  int `json:"inputChanel"`
		OutputChanel int `json:"outputChanel"`
	} `json:"workers"`
	Proxies []string `json:"proxies"`
}

func ReadConf(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	response := Config{}

	if err = json.Unmarshal(data, &response); err != nil {
		return Config{}, err
	}

	return response, nil
}
