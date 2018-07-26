package services

type Plugin interface {
	Initialize(services Services) bool
	Interface() interface{}
}
