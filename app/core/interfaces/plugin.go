package interfaces

type AvalanchePlugin interface {
	Version() string
	VersionCode() int
	AvalancheVersionCode() int
	Title() string
	Description() string
	Initialize() bool
	Interface() interface{}
}