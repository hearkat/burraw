package main

import (
	"./drivers"
	"github.com/hearkat/hearkat-go"
	"github.com/mmcdole/gofeed"
)

type Driver interface {
	Matches(source string) bool
	ToHearkat(item *gofeed.Item) *hearkat.Message
}

var driversList = []Driver{
	&drivers.Youtube{},
	&drivers.Crunchyroll{},
	&drivers.PresseCitron{},
	&drivers.Hitek{},
}

func findDriver(source string) Driver {
	for _, d := range driversList {
		if d.Matches(source) {
			return d
		}
	}
	return &drivers.Default{}
}
