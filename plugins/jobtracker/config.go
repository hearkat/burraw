package main

type Search struct {
	Channel  string `json:"channel"`
	Provider string `json:"provider"`
}

type Config struct {
	Timeout  int      `json:"timeout"`
	Searches []Search `json:"searches"`
}
