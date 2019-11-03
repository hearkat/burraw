package drivers

import (
	"strings"

	"github.com/hearkat/hearkat-go"
	"github.com/mmcdole/gofeed"
)

type Crunchyroll struct{}

func (def *Crunchyroll) Matches(source string) bool {
	return strings.HasPrefix(source, "https://www.crunchyroll.com/rss")
}

func (def *Crunchyroll) ToHearkat(item *gofeed.Item) *hearkat.Message {
	image := item.Extensions["media"]["thumbnail"][0].Attrs["url"]
	return &hearkat.Message{
		Notification: &hearkat.Notification{
			Title:   "New Crunchyroll video",
			Message: item.Title,
			Link:    &item.Link,
			Image:   &image,
		},
		Metadata: &hearkat.Metadata{
			"rss": item,
		},
		Tags: &hearkat.Tags{
			"rss",
			"crunchyroll",
		},
	}
}
