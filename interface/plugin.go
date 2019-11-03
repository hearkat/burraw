package _interface

type Plugin interface {
	Name() string
	Start(Burraw)
	Stop()
}

type Save struct {
	Store    map[string]interface{}
	Tags     []interface{}
	Metadata map[interface{}]interface{}
}
