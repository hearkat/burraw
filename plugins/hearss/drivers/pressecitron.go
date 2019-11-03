package drivers

import (
	"regexp"

	"github.com/hearkat/hearkat-go"
	"github.com/mmcdole/gofeed"
)

type PresseCitron struct{}

func (def *PresseCitron) Matches(source string) bool {
	match, _ := regexp.MatchString("https://www.presse-citron.net.*/feed", source)
	return match
}

func (def *PresseCitron) ToHearkat(item *gofeed.Item) *hearkat.Message {
	return &hearkat.Message{
		Notification: &hearkat.Notification{
			Title:   "New PresseCitron news",
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
		},
	}
}
