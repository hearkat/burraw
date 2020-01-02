package main

import (
	"github.com/IDerr/go-dealabs"
	"time"
)

func (p *Dealabs) hasSeen(deal *dealabs.Data) bool {
	if val, ok := p.burraw.GetSave().Store["lastSeen"]; ok {
		d := val.(time.Time)
		return d.After(deal.PublishedAt())
	}
	return false
}

func (p *Dealabs) setSeen(deal *dealabs.Data) error {
	date := deal.PublishedAt()
	if _, ok := p.burraw.GetSave().Store["lastSeen"]; !ok {
		p.burraw.GetSave().Store["lastSeen"] = date
	}
	if p.burraw.GetSave().Store["lastSeen"].(time.Time).Before(date) {
		p.burraw.GetSave().Store["lastSeen"] = date
	}
	return nil
}
