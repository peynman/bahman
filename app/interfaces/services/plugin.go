package services

type AvalanchePlugin interface {
	Initialize(services Services) bool
	Interface() interface{}
}
