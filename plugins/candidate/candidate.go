package main

import (
	i "github.com/hearkat/burraw/interface"
	"github.com/hearkat/hearkat-go"
	"log"
)

type Candidate struct{}

func main() {}

func InitPlugin() i.Plugin {
	return &Candidate{}
}

func (p *Candidate) Name() string {
	return "candidate"
}

func (p *Candidate) Start(b i.Burraw) {
	log.Println("started")

	b.OnMessage(func(container *hearkat.MessageContainer) {
		log.Println(container)
	})
}

func (p *Candidate) Stop() {
	log.Println("stopped")
}
