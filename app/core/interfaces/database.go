package interfaces

type Migrator interface {
}

type Migratable interface {
	Up(db Migrator) bool
	Down(db Migrator) bool
}