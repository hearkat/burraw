package main

type Channel struct {
	Channel string `json:"channel"`
	Stream  string `json:"stream"`
	Topic   string `json:"topic"`
}

type Config struct {
	RootUrl  string    `json:"root_url"`
	Email    string    `json:"email"`
	ApiKey   string    `json:"apikey"`
	Channels []Channel `json:"channels"`
}
