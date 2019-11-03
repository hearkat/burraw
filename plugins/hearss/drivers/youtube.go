package drivers

import (
	"github.com/hearkat/hearkat-go"
	"github.com/mmcdole/gofeed"
	"strings"
)

type Youtube struct{}

func (def *Youtube) Matches(source string) bool {
	return strings.HasPrefix(source, "https://www.youtube.com/feeds/videos.xml")
}

func (def *Youtube) ToHearkat(item *gofeed.Item) *hearkat.Message {
	image := item.Extensions["media"]["group"][0].Children["thumbnail"][0].Attrs["url"]
	return &hearkat.Message{
		Notification: &hearkat.Notification{
			Title:   "New youtube video",
			Message: item.Author.Name + ": " + item.Title,
			Link:    &item.Link,
			Image:   &image,
		},
		Metadata: &hearkat.Metadata{
			"rss": item,
		},
		Tags: &hearkat.Tags{
			"rss",
		},
	}
}
