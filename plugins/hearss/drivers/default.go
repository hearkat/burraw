package drivers

import (
	"github.com/hearkat/hearkat-go"
	"github.com/mmcdole/gofeed"
)

type Default struct{}

func (def *Default) Matches(source string) bool {
	return true
}

func (def *Default) ToHearkat(item *gofeed.Item) *hearkat.Message {
	return &hearkat.Message{
		Notification: &hearkat.Notification{
			Title:   "RSS update",
			Message: item.Title,
			Link:    &item.Link,
			Image: func() *string {
				if item.Image != nil {
					return &item.Image.URL
				}
				return nil
			}(),
		},
		Metadata: &hearkat.Metadata{
			"rss": item,
		},
		Tags: &hearkat.Tags{
			"rss",
		},
	}
}
