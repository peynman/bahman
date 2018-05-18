package interfaces

type Migrator interface {
	Connection(string) Migrator

	AutoMigrate(...interface{}) error
	CreateTable(...interface{}) error
	DropTable(...interface{}) error
	DropTableIfExists(...interface{}) error
	DropColumn(string) error
	Migrate([]Migratable) error
	Rollback([]Migratable) error
}

type Entity interface {
	TableName() string
	PrimaryKey() string
}

type EntityConnection interface {
	ConnectionName() string
}

type QueryBuilder interface {
	Where(query interface{}, args...interface{}) QueryBuilder
	OrWhere(query interface{}, args...interface{}) QueryBuilder
	Limit(limit interface{}) QueryBuilder
	Offset(offset interface{}) QueryBuilder
	Select(query interface{}, args...interface{}) QueryBuilder

	Get(result interface{}) error
	Update(...interface{}) error
}

type Repository interface {
	Query(Entity) QueryBuilder
	Insert(Entity) error
	UpdateEntity(Entity) error
}

type Migratable interface {
	Up(db Migrator) bool
	Down(db Migrator) bool
}