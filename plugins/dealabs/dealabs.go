package main

import (
	"github.com/IDerr/go-dealabs"
	i "github.com/hearkat/burraw/interface"
	"time"
)

type Dealabs struct {
	config Config
	burraw i.Burraw
}

func main() {}

func InitPlugin() i.Plugin {
	return &Dealabs{}
}

func (p *Dealabs) Name() string {
	return "dealabs"
}

func (p *Dealabs) Start(b i.Burraw) {
	p.burraw = b

	err := b.GetConfig(&p.config)
	if err != nil {
		panic(err)
	}

	d := dealabs.New()

	go func() {
		for {
			for _, deal := range d.GetNewDeals(nil).Data {
				if p.hasSeen(&deal) {
					continue
				}

				p.setSeen(&deal)
				msg := ToHearkat(&deal)
				b.Push(p.config.Channel, msg)
			}
			time.Sleep(time.Duration(p.config.Timeout) * time.Second)
		}
	}()
}

func (p *Dealabs) Stop() {}
