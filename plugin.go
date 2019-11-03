package main

import (
	"bytes"
	"encoding/json"
	i "github.com/hearkat/burraw/interface"
	"github.com/hearkat/hearkat-go"
	"github.com/kandoo/beehive/gob"
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
	save     *i.Save
}

func newPlugin(b *burraw, file string, pl i.Plugin) *plugin {
	p := &plugin{
		b,
		file,
		pl,
		make(chan *hearkat.MessageContainer, 100),
		make([]func(*hearkat.MessageContainer), 0),
		nil,
	}

	p.save = p.loadSave()

	return p
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

func (p *plugin) GetSaveFile() string {
	return path.Join(p.burraw.getPluginFolder(p), "save.bin")
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

		out := bytes.NewBuffer(nil)
		enc := json.NewEncoder(out)
		enc.SetIndent("", "    ")
		if err := enc.Encode(config); err != nil {
			panic(err)
		}

		dat := out.Bytes()

		err = ioutil.WriteFile(cf, dat, os.ModePerm)
		if err != nil {
			return err
		}

		WARN("Generated default config for plugin", p.plugin.Name())
		WARN("Please make sure to edit it")
		os.Exit(0)

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

func (p *plugin) loadSave() *i.Save {
	sf := p.GetSaveFile()

	if _, err := os.Stat(sf); os.IsNotExist(err) {
		return &i.Save{}
	}

	file, err := ioutil.ReadFile(sf)
	if err != nil {
		panic(err)
	}

	sav := &i.Save{
		make(map[string]interface{}),
		make([]interface{}, 0),
		make(map[interface{}]interface{}),
	}

	err = gob.Decode(sav, file)
	if err != nil {
		panic(err)
	}

	return sav
}

func (p *plugin) writeSave() {
	dat, err := gob.Encode(p.save)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(p.GetSaveFile(), dat, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func (p *plugin) GetSave() *i.Save {
	return p.save
}
