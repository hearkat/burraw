package main

import (
	"github.com/hearkat/hearkat-go"
	"github.com/skyforce77/jobtracker/providers"
)

func ToHearkat(job *providers.Job) *hearkat.Message {
	meta := hearkat.Metadata{}
	meta["job"] = *job

	msg := &hearkat.Message{
		Notification: &hearkat.Notification{
			job.Title,
			job.Desc,
			&job.Link,
			nil,
		},
		Metadata: &meta,
	}
	return msg
}
