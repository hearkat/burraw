package main

import (
"bytes"
"encoding/base64"
"github.com/gregdel/pushover"
i "github.com/hearkat/burraw/interface"
"github.com/hearkat/hearkat-go"
"io/ioutil"
"net/http"
"regexp"
"strings"
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

func (p *Plushover) sendMatches(msg *hearkat.MessageContainer) {
	for _, v := range p.config.Pushover {
		m, err := regexp.MatchString(v.Pattern, msg.Channel)
		if err != nil {
			panic(err)
		}
		if !m {
			continue
		}

		po := pushover.New(v.ApiKey)
		r := pushover.NewRecipient(v.UserKey)
		n := msg.Message.Notification

		var notif *pushover.Message
		if n.Link == nil {
			notif = &pushover.Message{
				Message: n.Message,
				Title:   n.Title,
			}
		} else {
			notif = &pushover.Message{
				Message: n.Message,
				Title:   n.Title,
				URL:     *n.Link,
			}
		}

		p.Image(notif, msg)

		_, err = po.SendMessage(notif, r)
		if err != nil {
			panic(err)
		}
		return
	}
}

func (p *Plushover) Image(notif *pushover.Message, msg *hearkat.MessageContainer) {
	image := msg.Message.Notification.Image
	if image == nil {
		return
	}
	img := *image
	http, err := regexp.MatchString(`.*\.(jpg|jpeg|JPG|gif|GIF|png|PNG)`, img)
	if err != nil {
		panic(err)
	}

	if strings.HasPrefix(img, "data:image/") {
		p.ImageBase64(notif, msg)
	} else if http {
		p.ImageHTTP(notif, msg)
	}
}

func (p *Plushover) ImageBase64(notif *pushover.Message, msg *hearkat.MessageContainer) {
	img := *msg.Message.Notification.Image

	re := regexp.MustCompile(`data:image\/.*;base64,(.*)`)
	s := re.ReplaceAllString(img, "$1")
	s2, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}

	err = notif.AddAttachment(bytes.NewReader(s2))
	if err != nil {
		panic(err)
	}
}

func (p *Plushover) ImageHTTP(notif *pushover.Message, msg *hearkat.MessageContainer) {
	img := *msg.Message.Notification.Image

	resp, err := http.Get(img)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = notif.AddAttachment(bytes.NewReader(dat))
	if err != nil {
		panic(err)
	}
}

