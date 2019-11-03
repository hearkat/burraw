package main

func (p *Hearss) hasItem(feedId string, guid string) bool {
	if val, ok := p.burraw.GetSave().Store[feedId]; ok {
		arr := val.([]string)
		for _, v := range arr {
			if v == guid {
				return true
			}
		}
	}
	return false
}

func (p *Hearss) addItem(feedId string, guid string) error {
	if _, ok := p.burraw.GetSave().Store[feedId]; !ok {
		p.burraw.GetSave().Store[feedId] = make([]string, 0)
	}

	p.burraw.GetSave().Store[feedId] = append(
		p.burraw.GetSave().Store[feedId].([]string), guid)

	return nil
}
