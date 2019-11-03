package main

import (
	i "./interface"
	"github.com/hearkat/hearkat-go"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"
	"plugin"
	"syscall"
	"time"
)

type burraw struct {
	dir string
	config *Config
	plugins []Plugin
	hearkat hearkat.Hearkat
}

func newBurraw(dir string) *burraw {
	b := &burraw{
		dir,
		nil,
		nil,
		nil,
	}

	b.config = loadConfig(b.getConfigFile())

	b.hearkat = hearkat.NewClient(
		b.config.Hearkat.UserToken,
		b.config.Hearkat.AccessToken)

	return b
}

func (b *burraw) getPluginsFolder() string {
	return path.Join(b.dir, "plugins")
}

func (b *burraw) getPluginFolder(pl *Plugin) string {
	return path.Join(b.dir, "plugins", pl.plugin.Name())
}

func (b *burraw) getConfigFile() string {
	return path.Join(b.dir, "config.json")
}

func (b *burraw) init() {
	err := os.MkdirAll(b.getPluginsFolder(), os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func (b *burraw) loadGoPlugin(filename string) *Plugin {
	file := path.Join(b.getPluginsFolder(), filename)

	LOG("Initializing plugin", filename)

	plug, err := plugin.Open(file)
	if err != nil {
		WARN(err)
		return nil
	}

	pluginVar, err := plug.Lookup("InitPlugin")
	if err != nil {
		WARN(filename, "Please include the 'InitPlugin' function in your code")
		return nil
	}

	var plugin func() i.Plugin
	plugin, ok := pluginVar.(func() i.Plugin)
	if !ok {
		WARN(filename, "Type 'func() Plugin' expected as entry point")
		return nil
	}

	pl := plugin()

	LOG("Initialized", filename)
	return newPlugin(b, file, pl)
}

func (b *burraw) loadPlugins() []*Plugin {
	files, err := ioutil.ReadDir(b.getPluginsFolder())
	if err != nil {
		panic(err)
	}

	plugins := make([]*Plugin, 0)

	for _, f := range files {
		if ! f.IsDir() {
			p := b.loadGoPlugin(f.Name())
			if p != nil {
				plugins = append(plugins, p)

				err = os.MkdirAll(b.getPluginFolder(p), os.ModePerm)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	return plugins
}

func (b *burraw) handleMessage(container *hearkat.MessageContainer) {
	for _, p := range b.plugins {
		p.msg <- container
	}
}

func (b *burraw) listen() {
	for {
		stream, err := b.hearkat.StreamAll()

		if err != nil {
			panic(err)
		}

		go func() {
			for m := range stream.Messages {
				b.handleMessage(m)
			}
		}()

		err = <-stream.Errors
		ERR(err)
		time.Sleep(10 * time.Second)
	}
}

func (b *burraw) run() {
	b.init()
	plugins := b.loadPlugins()

	for _, plugin := range plugins {
		pl := plugin.plugin

		LOG("Starting plugin", pl.Name())
		pl.Start(plugin)
		LOG("Started plugin", pl.Name())
	}

	go b.listen()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	for _, plugin := range plugins {
		pl := plugin.plugin

		LOG("Sopping plugin", pl.Name())
		pl.Stop()
		LOG("Stopped plugin", pl.Name())
	}
}