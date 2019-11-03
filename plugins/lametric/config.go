package main

type LametricIcon struct {
	Pattern string  `json:pattern`
	Icon    string  `json:"icon"`
	Sound   *string `json:"sound"`
}

type Config struct {
	IP     string         `json:"ip"`
	ApiKey string         `json:"apiKey"`
	Icons  []LametricIcon `json:"icons"`
}
