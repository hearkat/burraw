package main

import (
	"fmt"
	i "github.com/hearkat/burraw/interface"
	"github.com/mmcdole/gofeed"
	"time"
)

type Hearss struct {
	config Config
	burraw i.Burraw
}

func main() {}

func InitPlugin() i.Plugin {
	return &Hearss{}
}

func (p *Hearss) Name() string {
	return "hearss"
}

func (p *Hearss) Start(b i.Burraw) {
	p.burraw = b

	err := b.GetConfig(&p.config)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			for _, rss := range p.config.RSS {
				p.parseRSS(&rss)
			}
			time.Sleep(time.Duration(p.config.Timeout) * time.Second)
		}
	}()
}

func (p *Hearss) Stop() {

}

func (p *Hearss) parseRSS(rss *RSSConfig) {
	d := findDriver(rss.Link)

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(rss.Link)
	if err != nil {
		panic(err)
	}

	for _, item := range feed.Items {
		if !p.hasSeen(feed.Link, item.PublishedParsed) {
			msg := d.ToHearkat(item)
			fmt.Println(item.Title)

			if rss.Title != nil {
				msg.Notification.Title = *rss.Title
			}

			err := p.burraw.Push(rss.Channel, msg)
			if err != nil {
				panic(err)
				continue
			}
			err = p.setSeen(feed.Link, item.PublishedParsed)
			if err != nil {
				panic(err)
				continue
			}
		}
	}
}
