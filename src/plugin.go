package main

import (
	i "./interface"
	"github.com/hearkat/hearkat-go"
	"github.com/pkg/errors"
)

type Plugin struct {
	burraw *burraw
	file string
	plugin i.Plugin
	msg chan *hearkat.MessageContainer
}

func newPlugin(b *burraw, file string, pl i.Plugin) *Plugin {
	return &Plugin{
		b,
		file,
		pl,
		make(chan *hearkat.MessageContainer, 100),
	}
}

func (p *Plugin) Stream() chan *hearkat.MessageContainer {
	return p.msg
}

func (p *Plugin) Push(channel string, message *hearkat.Message) error {
	if p.burraw.hearkat == nil {
		return errors.New("Hearkat unavailable")
	}

	return p.burraw.hearkat.Push(channel, message)
}