package main

import (
	"encoding/json"
	"io/ioutil"
)

type HearkatConfig struct {
	UserToken   string `json:"usertoken"`
	AccessToken string `json:"accesstoken"`
}

type Config struct {
	Hearkat HearkatConfig `json:"hearkat"`
}

func loadConfig(filename string) *Config {
	var config Config

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(file, &config); err != nil {
		panic(err)
	}

	return &config
}