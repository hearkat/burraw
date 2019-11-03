package _interface

import "github.com/hearkat/hearkat-go"

type Burraw interface {
	Stream() chan *hearkat.MessageContainer
	Push(channel string, message *hearkat.Message) error
}
