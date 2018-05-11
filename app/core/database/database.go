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
	"github.com/spf13/viper"
	"time"
	"avalanche/app/core/interfaces"
)

var defaultConnection *gorm.DB
var migrationsManager *MigrationManager

func Initialize() {
	defaultConnectionName := config.GetString("database.default", "")
	if defaultConnectionName == "" {
		panic("Database connection is set to nothing, Please provide one!")
	} else {
		if 	connectionDriver := config.GetString("database.connections." +defaultConnectionName+ ".driver", "");
			connectionDriver != "" {
			switch connectionDriver {
			case "mysql":
				defaultConnection = openMySQL("database.connections." + defaultConnectionName)
			case "sqlite3":
				defaultConnection = openSqlite3("database.connections." + defaultConnectionName)
			}
		} else {
			panic("Unknown database connection: " + defaultConnectionName)
		}
	}

	if viper.IsSet("database.connections." +defaultConnectionName+ ".options") {
		engine := config.GetString("database.connections." +defaultConnectionName+ ".options.engine", "InnoDB")
		maxIdleConnections := config.GetInt("database.connections." +defaultConnectionName+ ".options.maxIdleConnections", 1)
		maxOpenConnections := config.GetInt("database.connections." +defaultConnectionName+ ".options.maxOpenConnections", 1)
		maxConnectionLifetime := config.GetInt64("database.connections." +defaultConnectionName+ ".options.maxConnectionLifetime", int64(time.Hour))

		defaultConnection = defaultConnection.Set("gorm:table_options", "ENGINE=" + engine)
		defaultConnection.DB().SetMaxIdleConns(maxIdleConnections)
		defaultConnection.DB().SetMaxOpenConns(maxOpenConnections)
		defaultConnection.DB().SetConnMaxLifetime(time.Duration(maxConnectionLifetime))
	}

	defaultConnection.SetLogger(gorm.Logger{
		LogWriter: new(AvalancheDBLogWriter),
	})

	migrationsManager = new(MigrationManager)
	migrationsManager.Setup(defaultConnection)
}

func DeployModule(module interfaces.Module) {
	migrationsManager.Migrate(defaultConnection, module.Migrations())
}

func Close() {
	defaultConnection.Close()
}

func openSqlite3(configPath string) *gorm.DB {
	db, err := gorm.Open("sqlite3", config.GetString(configPath + ".file", ""))
	if err != nil {
		panic(err)
	}
	return db
}

func openMySQL(configPath string) *gorm.DB {

	return nil
}


