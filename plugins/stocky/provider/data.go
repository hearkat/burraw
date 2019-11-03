package provider

import (
	"fmt"
	i "github.com/hearkat/burraw/interface"
	"github.com/hearkat/hearkat-go"
	"log"
)

type DataItem struct {
	Id  string
	Url string
}

type StoreData map[string]int
type Data struct {
	store  map[string]StoreData
	search map[string][]DataItem
	burraw i.Burraw
}

func NewData(b i.Burraw) *Data {
	data := &Data{
		make(map[string]StoreData),
		make(map[string][]DataItem),
		b,
	}

	if st, ok := b.GetSave().Store["store"]; ok {
		stD, ok := st.(map[string]StoreData)
		if ok {
			data.store = stD
		}
	}

	if search, ok := b.GetSave().Store["search"]; ok {
		searchD, ok := search.(map[string][]DataItem)
		if ok {
			data.search = searchD
		}
	}

	return data
}

func (data *Data) Save() {
	data.burraw.GetSave().Store["store"] = data.store
	data.burraw.GetSave().Store["search"] = data.search
}

func (data *Data) setStock(store string, id string, stock int) (int, int) {
	st := data.store[store]

	if st == nil {
		st = make(StoreData)
		data.store[store] = st
	}

	if _, ok := st[id]; !ok {
		st[id] = stock
		return stock, 0
	}

	before := st[id]
	if before != stock {
		st[id] = stock
		return stock, stock - before
	}

	return stock, 0
}

func (data *Data) SetStock(channel string, store string, id string, stock int,
	url string) {
	stock, variation := data.setStock(store, id, stock)

	if variation == 0 {
		return
	}

	message := hearkat.Message{
		Notification: &hearkat.Notification{
			Title:   fmt.Sprintf("%s stock update", store),
			Message: fmt.Sprintf("%d %s available", stock, id),
			Link:    &url,
		},
		Tags: &hearkat.Tags{
			"stock",
		},
		Metadata: &hearkat.Metadata{
			"store": store,
			"item":  id,
			"from":  stock,
			"to":    stock + variation,
		},
	}

	err := data.burraw.Push(channel, &message)
	if err != nil {
		log.Println("Error:", err)
	}
}

func (data *Data) setSearch(store string, items []DataItem) []DataItem {
	st := data.search[store]

	if st == nil {
		st = make([]DataItem, 0)
		data.search[store] = st
		return make([]DataItem, 0)
	}

	m := make([]DataItem, 0)

	for _, a := range items {
		found := false
		for _, b := range st {
			if b == a {
				found = true
				break
			}
		}
		if !found {
			m = append(m, a)
		}
	}

	data.search[store] = items

	return m
}

func (data *Data) SetSearch(channel string, store string, items []DataItem) {
	added := data.setSearch(store, items)

	for _, item := range added {
		message := hearkat.Message{
			Notification: &hearkat.Notification{
				Title:   fmt.Sprintf("New item available at %s", store),
				Message: fmt.Sprintf("%s available !", item.Id),
				Link:    &item.Url,
			},
			Tags: &hearkat.Tags{
				"stock", "availability",
			},
			Metadata: &hearkat.Metadata{
				"store": store,
				"items": items,
			},
		}

		err := data.burraw.Push(channel, &message)
		if err != nil {
			log.Println("Error:", err)
		}
	}
}
