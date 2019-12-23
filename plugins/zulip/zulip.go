package main

import (
	"fmt"
	i "github.com/hearkat/burraw/interface"
	"github.com/hearkat/hearkat-go"
)

type Zulip struct{}

func main() {}

func InitPlugin() i.Plugin {
	return &Zulip{}
}

func (p *Zulip) Name() string {
	return "zulip"
}

func (p *Zulip) Start(b i.Burraw) {
	config := Config{}
	b.GetConfig(&config)

	z := NewClient(config.RootUrl, config.Email, config.ApiKey)

	b.OnMessage(func(container *hearkat.MessageContainer) {
		if container.Message.Notification != nil {
			notif := container.Message.Notification
			for _, c := range config.Channels {
				if c.Channel == container.Channel {
					z.SendMessage(c.Stream, c.Topic, fmt.Sprintf(`
**[%s](%s)**

%s
`, notif.Title, *notif.Link, notif.Message))
				}
			}
		}
	})
}

func (p *Zulip) Stop() {}
