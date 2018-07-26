package orm

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"time"
)

type RDBMSClients struct {
	appConnection     *gorm.DB
	runtimeConnection *gorm.DB
	connections       map[string]*gorm.DB
	connectionPrefix  map[*gorm.DB]string
	migrationsManager *MigrationManager
	repoManager *RepositoryManager
	loggerRef   services.Logger
	migrator    services.Migratory
}

func New(config services.Config, logger services.Logger) (services.Repository, services.Migratory) {
	client := new(RDBMSClients)
	client.connections = make(map[string]*gorm.DB)
	client.connectionPrefix = make(map[*gorm.DB]string)
	client.loggerRef = logger
	appConnectionName := config.GetString("rdbms.app", "sqlite3")
	runtimeConnectionName := config.GetString("rdbms.runtime.connection", "sqlite3")

	gorm.LogFormatter = func(values ...interface{}) (messages []interface{}) {
		return values
	}

	connectionDefs := config.GetMap("rdbms.connections", map[string]interface{}{
		"sqlite3": map[string]interface{} {
			"active": true,
			"driver": "sqlite3",

		},
	})
	for connName, connParams := range connectionDefs {
		connMap := connParams.(map[string]interface{})
		var connection *gorm.DB = nil
		if connName == appConnectionName || connName == runtimeConnectionName {
			connection = initDatabase(client, config, "rdbms.connections."+connName)
			if connName == appConnectionName {
				client.appConnection = connection
				client.connections["app"] = connection
			}
			if connName == runtimeConnectionName {
				client.runtimeConnection = connection
				client.connections["runtime"] = connection
			}
		}

		if connMap["active"] == true {
			if connection == nil {
				connection = initDatabase(client, config, "rdbms.connections."+connName)
			}

			client.connections[connName] = connection
		}
	}

	if config.GetBoolean("app.debug", false) {
		for _, conn := range client.connections {
			conn.LogMode(true)
		}
		client.appConnection.LogMode(true)
		client.runtimeConnection.LogMode(true)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		if prefix, ok := client.connectionPrefix[db]; ok {
			return prefix + defaultTableName
		}
		return defaultTableName
	}

	client.repoManager = NewRepository(
		logger,
		client,
	)
	client.migrationsManager = NewMigrator(
		logger,
		client,
		config.GetString("rdbms.runtime.migrations", "migrations"),
	)
	client.migrator = client.migrationsManager.Connection("runtime")
	client.migrationsManager.setup(client.runtimeConnection)

	return client.repoManager, client.migrator
}

func initDatabase(client *RDBMSClients, config services.Config, keyPrefix string) *gorm.DB {
	var connection *gorm.DB
	if connectionDriver := config.GetString(keyPrefix+".driver", "sqlite3"); connectionDriver != "" {
		switch connectionDriver {
		case "mysql":
			connection = openMySQL(config, keyPrefix)
		case "sqlite3":
			connection = openSqlite3(config, keyPrefix)
		}
	} else {
		panic("Unknown orm connection: " + keyPrefix)
	}

	if config.IsSet(keyPrefix + ".options") {
		maxIdleConnections := config.GetInt(keyPrefix+".options.maxIdleConnections", 1)
		connection.DB().SetMaxIdleConns(maxIdleConnections)

		maxOpenConnections := config.GetInt(keyPrefix+".options.maxOpenConnections", 1)
		connection.DB().SetMaxOpenConns(maxOpenConnections)

		maxConnectionLifetime := config.Get(keyPrefix+".options.maxConnectionLifetime", time.Hour).(time.Duration)
		connection.DB().SetConnMaxLifetime(maxConnectionLifetime)

		if config.IsSet(keyPrefix + ".options.prefix") {
			client.connectionPrefix[connection] = config.GetString(keyPrefix+".options.prefix", "")
		}
	} else {
		connection.DB().SetMaxIdleConns(1)
		connection.DB().SetMaxOpenConns(1)
		connection.DB().SetConnMaxLifetime(time.Hour)
	}

	connection.SetLogger(gorm.Logger{
		LogWriter: new(bahmanDBLogWriter),
	})

	return connection
}

func openSqlite3(config services.Config, configPath string) *gorm.DB {
	filename := config.GetString(configPath+".file", "")
	db, err := gorm.Open("sqlite3", filename)
	if err != nil {
		panic(err)
	}
	return db
}

func openMySQL(config services.Config, configPath string) *gorm.DB {
	//engine := config.GetString(keyPrefix + ".options.engine", "InnoDB")
	//connection = appConnection.Set("gorm:table_options", "ENGINE=" + engine)
	return nil
}
