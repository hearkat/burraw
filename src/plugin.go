package main

import (
	i "./interface"
	"encoding/json"
	"github.com/hearkat/hearkat-go"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
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

func (p *Plugin) GetConfigFile() string {
	return path.Join(p.burraw.getPluginFolder(p), "config.json")
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

func (p *Plugin) GetConfig(config interface{}) error {
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