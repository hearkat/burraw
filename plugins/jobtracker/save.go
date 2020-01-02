package main

import (
	"github.com/skyforce77/jobtracker/providers"
)

func (p *Jobtracker) hasSeen(job *providers.Job) bool {
	if _, ok := p.burraw.GetSave().Store[job.Link]; ok {
		return true
	}
	return false
}

func (p *Jobtracker) setSeen(job *providers.Job) error {
	if _, ok := p.burraw.GetSave().Store[job.Link]; !ok {
		p.burraw.GetSave().Store[job.Link] = true
	}
	return nil
}
