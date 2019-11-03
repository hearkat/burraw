package main

import (
	"./provider"
	"encoding/gob"
	i "github.com/hearkat/burraw/interface"
	"log"
	"time"
)

type Stocky struct {
	burraw i.Burraw
	config Config
	data   *provider.Data
}

func main() {}

func InitPlugin() i.Plugin {
	gob.Register(map[string]provider.StoreData{})
	gob.Register(map[string][]provider.DataItem{})
	return &Stocky{}
}

func (p *Stocky) Name() string {
	return "stocky"
}

func (p *Stocky) Start(b i.Burraw) {
	err := b.GetConfig(&p.config)
	if err != nil {
		panic(err)
	}

	p.data = provider.NewData(b)

	go func() {
		for {
			p.DoRequests()
			time.Sleep(time.Duration(p.config.Timeout) * time.Minute)
		}
	}()
}

func (p *Stocky) Stop() {
	p.data.Save()
}

func (p *Stocky) DoRequests() {
	data := p.data

	err := provider.FindCultura(&p.config.Items.Cultura, data)
	if err != nil {
		log.Println("Error:", err)
	}

	err = provider.FindPikastore(&p.config.Items.Pikastore, data)
	if err != nil {
		log.Println("Error:", err)
	}

	err = provider.FindCardotaku(&p.config.Items.Cardotaku, data)
	if err != nil {
		log.Println("Error:", err)
	}
}
