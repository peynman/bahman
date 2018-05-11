package interfaces

type Model interface {

}

type Module interface {
	Migrations() []Migratable
}