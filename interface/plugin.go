package _interface

type Plugin interface {
	Name() string
	Start(Burraw)
	Stop()
}