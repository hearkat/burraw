package main

import (
	i "github.com/hearkat/burraw/interface"
	"github.com/skyforce77/jobtracker/providers"
	"time"
)

type Jobtracker struct {
	config Config
	burraw i.Burraw
}

func main() {}

func InitPlugin() i.Plugin {
	return &Jobtracker{}
}

func (p *Jobtracker) Name() string {
	return "jobtracker"
}

func (p *Jobtracker) Start(b i.Burraw) {
	p.burraw = b

	err := b.GetConfig(&p.config)
	if err != nil {
		panic(err)
	}

	go func() {
		for _, search := range p.config.Searches {
			pr := providers.ProviderFromName(search.Provider)
			pr.RetrieveJobs(func(job *providers.Job) {
				if p.hasSeen(job) {
					return
				}

				p.setSeen(job)
				msg := ToHearkat(job)
				b.Push(search.Channel, msg)
			})
			time.Sleep(time.Duration(p.config.Timeout) * time.Second)
		}
	}()
}

func (p *Jobtracker) Stop() {}
