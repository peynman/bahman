package interfaces

type Migrator interface {
	Connection(string) Migrator

	AutoMigrate(entity Entity) error
	CreateTable(entity Entity) error
	DropTable(...interface{}) error
	DropTableIfExists(...interface{}) error
	DropColumn(string) error
	HasTable(interface{}) bool
	Migrate([]Migratable) error
	Rollback([]Migratable) error
}

type Entity interface {
	TableName() string
}

type EntityConnection interface {
	ConnectionName() string
}

const (
	DESC = 1
	ASC = 2
)

type QueryBuilder interface {
	IncludeDeleted() QueryBuilder
	Where(query interface{}, args...interface{}) QueryBuilder
	OrWhere(query interface{}, args...interface{}) QueryBuilder
	Limit(limit interface{}) QueryBuilder
	Offset(offset interface{}) QueryBuilder
	Select(query interface{}, args...interface{}) QueryBuilder
	Order(column interface{}) QueryBuilder

	GetValues(valuesArrPtr interface{}) error
	GetValue(valuePtr interface{}) error
	GetFirst(entityPtr interface{}) error
	GetLast(entityPtr interface{}) error
	Get(entitiesArrRef interface{}) error
	Update(entityRef interface{}, params map[string]interface{}) error
	Updates(params map[string]interface{}) error
	DeleteAll() error
	EraseAll() error

	Reset() QueryBuilder
}

type Repository interface {
	Query(entity Entity) QueryBuilder
	Insert(entity Entity) error
	UpdateEntity(entity Entity) error
	DeleteEntity(entity Entity) error
	EraseEntity(entity Entity) error
}

type Migratable interface {
	Up(db Migrator) bool
	Down(db Migrator) bool
}