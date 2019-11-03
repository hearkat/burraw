package main

import (
	i "github.com/hearkat/burraw/interface"
	"github.com/hearkat/hearkat-go"
)

type Lametric struct {
	config Config
}

func main() {}

func InitPlugin() i.Plugin {
	return &Lametric{}
}

func (p *Lametric) Name() string {
	return "lametric"
}

func (p *Lametric) Start(b i.Burraw) {
	err := b.GetConfig(&p.config)

	if err != nil {
		panic(err)
	}

	b.OnMessage(p.handler)
}

func (p *Lametric) Stop() {}

func (p *Lametric) handler(msg *hearkat.MessageContainer) {
	n := msg.Message.Notification
	if n != nil {
		nn := p.toLametric(msg.Channel, n)
		p.sendNotification(nn)
	}
}
