package main

import "time"

func (p *Hearss) hasSeen(feedId string, date *time.Time) bool {
	if val, ok := p.burraw.GetSave().Store[feedId]; ok {
		d := val.(time.Time)
		return d.After(*date)
	}
	return false
}

func (p *Hearss) setSeen(feedId string, date *time.Time) error {
	if _, ok := p.burraw.GetSave().Store[feedId]; !ok {
		p.burraw.GetSave().Store[feedId] = *date
	}

	if p.burraw.GetSave().Store[feedId].(time.Time).Before(*date) {
		p.burraw.GetSave().Store[feedId] = *date
	}
	return nil
}
