package drivers

import (
	"github.com/hearkat/hearkat-go"
	"github.com/mmcdole/gofeed"
)

type Hitek struct{}

func (def *Hitek) Matches(source string) bool {
	return source == "https://hitek.fr/rss"
}

func (def *Hitek) ToHearkat(item *gofeed.Item) *hearkat.Message {
	return &hearkat.Message{
		Notification: &hearkat.Notification{
			Title:   "New Hitek news",
			Message: item.Title,
			Link:    &item.Link,
		},
		Metadata: &hearkat.Metadata{
			"rss": item,
		},
		News: &hearkat.News{
			Title:   item.Title,
			Summary: item.Description,
			Link:    item.Link,
		},
		Tags: &hearkat.Tags{
			"rss",
			"news",
			"anime",
		},
	}
}
