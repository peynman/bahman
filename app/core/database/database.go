package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq/hstore"
	_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
	"avalanche/app/core/config"
	"time"
	"avalanche/app/core/interfaces"
)

var DefaultConnection *gorm.DB
var MigrationsManager *MigrationManager
var Connections = make(map[string]*gorm.DB)
var ConnectionPrefix = make(map[*gorm.DB]string)

func Initialize() {
	defaultConnectionName := config.GetString("database.default", "")

	gorm.LogFormatter = func(values ...interface{}) (messages []interface{}) {
		return values
	}

	connections := config.GetMap("database.connections", map[string]interface{} {})
	for connName, connParams := range connections {
		connMap := connParams.(map[string]interface{})
		if connName == defaultConnectionName {
			DefaultConnection = initDatabase("database.connections." + connName)
		} else if connMap["active"] == true {
			Connections[connName] = initDatabase("database.connections." + connName)
		}
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		if prefix, ok := ConnectionPrefix[db]; ok {
			return prefix + defaultTableName
		}
		return defaultTableName
	}

	MigrationsManager = new(MigrationManager)
	MigrationsManager.tableName = config.GetString("database.migrations", "migrations")

	MigrationsManager.Setup(DefaultConnection)
}

func DeployModule(module interfaces.Module) {
	MigrationsManager.Migrate(DefaultConnection, module.Migrations())
}

func Close() {
	DefaultConnection.Close()
}

func initDatabase(keyPrefix string) *gorm.DB {
	var connection *gorm.DB
	if 	connectionDriver := config.GetString(keyPrefix + ".driver", "");
		connectionDriver != "" {
		switch connectionDriver {
		case "mysql":
			connection = openMySQL(keyPrefix)
		case "sqlite3":
			connection = openSqlite3(keyPrefix)
		}
	} else {
		panic("Unknown database connection: " + keyPrefix)
	}

	if config.IsSet(keyPrefix + ".options") {
		maxIdleConnections := config.GetInt(keyPrefix + ".options.maxIdleConnections", 1)
		connection.DB().SetMaxIdleConns(maxIdleConnections)

		maxOpenConnections := config.GetInt(keyPrefix + ".options.maxOpenConnections", 1)
		connection.DB().SetMaxOpenConns(maxOpenConnections)

		maxConnectionLifetime := config.Get(keyPrefix + ".options.maxConnectionLifetime", time.Hour).(time.Duration)
		connection.DB().SetConnMaxLifetime(maxConnectionLifetime)

		if config.IsSet(keyPrefix + ".options.prefix") {
			ConnectionPrefix[connection] = config.GetString(keyPrefix + ".options.prefix", "")
		}
	} else {
		connection.DB().SetMaxIdleConns(1)
		connection.DB().SetMaxOpenConns(1)
		connection.DB().SetConnMaxLifetime(time.Hour)
	}

	connection.SetLogger(gorm.Logger{
		LogWriter: new(AvalancheDBLogWriter),
	})

	return connection
}

func openSqlite3(configPath string) *gorm.DB {
	db, err := gorm.Open("sqlite3", config.GetString(configPath + ".file", ""))
	if err != nil {
		panic(err)
	}
	return db
}

func openMySQL(configPath string) *gorm.DB {

	//engine := config.GetString(keyPrefix + ".options.engine", "InnoDB")
	//connection = DefaultConnection.Set("gorm:table_options", "ENGINE=" + engine)
	return nil
}


