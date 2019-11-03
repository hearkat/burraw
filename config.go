package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
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

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		out := bytes.NewBuffer(nil)
		enc := json.NewEncoder(out)
		enc.SetIndent("", "    ")
		if err := enc.Encode(config); err != nil {
			panic(err)
		}

		dat := out.Bytes()

		err = ioutil.WriteFile(filename, dat, os.ModePerm)
		if err != nil {
			panic(err)
		}

		WARN("Generated default config for burraw")
		WARN("Please make sure to edit it")
		os.Exit(0)

		return nil
	}

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(file, &config); err != nil {
		panic(err)
	}

	return &config
}
