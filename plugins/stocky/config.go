package main

import "./provider"

type SearchConfig struct {
	Cultura   provider.CulturaConfig   `json:"cultura"`
	Pikastore provider.PikastoreConfig `json:"pikastore"`
	Cardotaku provider.CardotakuConfig `json:"cardotaku"`
}

type Config struct {
	Timeout int          `json:"timeout"`
	Items   SearchConfig `json:"search"`
}
