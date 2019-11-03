package main

import (
	i "github.com/hearkat/burraw/interface"
	"github.com/hearkat/hearkat-go"
)

type PushoverConfig struct {
	Pattern string `json:pattern`
	UserKey string `json:"userkey"`
	ApiKey  string `json:"apikey"`
}

type Config struct {
	Pushover []PushoverConfig `json:"pushover"`
}

type Plushover struct {
	config Config
}

func main() {}

func InitPlugin() i.Plugin {
	return &Plushover{}
}

func (p *Plushover) Name() string {
	return "plushover"
}

func (p *Plushover) Start(b i.Burraw) {
	err := b.GetConfig(&p.config)

	if err != nil {
		panic(err)
	}

	b.OnMessage(p.handler)
}

func (p *Plushover) Stop() {}

func (p *Plushover) handler(msg *hearkat.MessageContainer) {
	if msg.Message.Notification != nil {
		p.sendMatches(msg)
	}
}
