package main

import (
	"crypto/tls"
	i "github.com/hearkat/burraw/interface"
	"github.com/hearkat/hearkat-go"
	"net/http"
)

type Lametric struct {
	config Config
}

func main() {}

func InitPlugin() i.Plugin {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
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
