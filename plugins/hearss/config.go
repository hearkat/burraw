package main

type RSSConfig struct {
	Channel string  `json:"channel"`
	Link    string  `json:"link"`
	Title   *string `json:"title"`
}

type Config struct {
	Timeout int         `json:"timeout"`
	RSS     []RSSConfig `json:"rss"`
}
