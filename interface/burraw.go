package _interface

import "github.com/hearkat/hearkat-go"

type Burraw interface {
	Stream() chan *hearkat.MessageContainer
	OnMessage(func(*hearkat.MessageContainer))
	Push(channel string, message *hearkat.Message) error
	GetConfig(interface{}) error
	GetSave() *Save
}
