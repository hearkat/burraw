package main

import (
	"encoding/json"
	i "github.com/hearkat/burraw/interface"
	"github.com/hearkat/hearkat-go"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
)

type plugin struct {
	burraw   *burraw
	file     string
	plugin   i.Plugin
	msg      chan *hearkat.MessageContainer
	handlers []func(*hearkat.MessageContainer)
}

func newPlugin(b *burraw, file string, pl i.Plugin) *plugin {
	return &plugin{
		b,
		file,
		pl,
		make(chan *hearkat.MessageContainer, 100),
		make([]func(*hearkat.MessageContainer), 0),
	}
}

func (p *plugin) Handle(container *hearkat.MessageContainer) {
	p.msg <- container

	for _, h := range p.handlers {
		h(container)
	}
}

func (p *plugin) GetConfigFile() string {
	return path.Join(p.burraw.getPluginFolder(p), "config.json")
}

func (p *plugin) Stream() chan *hearkat.MessageContainer {
	return p.msg
}

func (p *plugin) OnMessage(f func(*hearkat.MessageContainer)) {
	p.handlers = append(p.handlers, f)
}

func (p *plugin) Push(channel string, message *hearkat.Message) error {
	if p.burraw.hearkat == nil {
		return errors.New("Hearkat unavailable")
	}

	return p.burraw.hearkat.Push(channel, message)
}

func (p *plugin) GetConfig(config interface{}) error {
	cf := p.GetConfigFile()

	if _, err := os.Stat(cf); os.IsNotExist(err) {
		dat, err := json.Marshal(config)

		if err != nil {
			return err
		}

		err = ioutil.WriteFile(cf, dat, os.ModePerm)
		if err != nil {
			return err
		}

		return nil
	}

	file, err := ioutil.ReadFile(cf)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(file, config); err != nil {
		panic(err)
	}

	return nil
}
