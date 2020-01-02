package main

import (
	"github.com/IDerr/go-dealabs"
	"github.com/hearkat/hearkat-go"
)

func ToHearkat(deal *dealabs.Data) *hearkat.Message {
	meta := hearkat.Metadata{}
	meta["deal"] = *deal

	msg := &hearkat.Message{
		Notification: &hearkat.Notification{
			deal.Title,
			deal.Description,
			&deal.DealURI,
			&deal.Image.URI,
		},
		ShopItem: &hearkat.ShopItem{
			deal.Title,
			deal.Price,
			deal.DealURI,
			&deal.Image.URI,
			&deal.Merchant.Name,
		},
		Metadata: &meta,
	}
	return msg
}
