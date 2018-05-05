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

	"AvConfig/lib"
)

var defaultConnection *gorm.DB

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


