package orm

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"reflect"
	"strings"
)

type MigrationManager struct {
	migrationsTableName string
	connection          *gorm.DB
	log                 services.Logger
	client				*RDBMSClients
}

func (m *MigrationManager) AutoMigrate(entity services.Entity) error {
	err := m.connection.AutoMigrate(entity).Error
	if err != nil {
		m.log.ErrorFields("Migratory:AutoMigrate", map[string]interface{}{
			"error": err,
			"table": entity.TableName(),
		})
	}
	return err
}

func (m *MigrationManager) CreateTable(entity services.Entity) error {
	err := m.connection.CreateTable(entity).Error
	if err != nil {
		m.log.ErrorFields("Migratory:CreateTable", map[string]interface{}{
			"error": err,
			"table": entity.TableName(),
		})
	}
	return err
}

func (m *MigrationManager) DropTable(table ...interface{}) error {
	err := m.connection.DropTable(table).Error
	if err != nil {
		m.log.ErrorFields("Migratory:DropTable", map[string]interface{}{
			"error": err,
			"table": table,
		})
	}
	return nil
}

func (m *MigrationManager) DropTableIfExists(table ...interface{}) error {
	err := m.connection.DropTableIfExists(table).Error
	if err != nil {
		m.log.ErrorFields("Migratory:DropTableIfExists", map[string]interface{}{
			"error": err,
			"table": table,
		})
	}
	return nil
}

func (m *MigrationManager) HasTable(table interface{}) bool {
	return m.connection.HasTable(table)
}

func (m *MigrationManager) DropColumn(column string) error {
	err := m.connection.DropColumn(column).Error
	if err != nil {
		m.log.ErrorFields("Migratory:DropColumn", map[string]interface{}{
			"error": err,
			"table": column,
		})
	}
	return nil
}

func (m *MigrationManager) ModifyColumn(column string, typ string) error {
	err := m.connection.ModifyColumn(column, typ).Error
	if err != nil {
		m.log.ErrorFields("Migratory:ModifyColumn", map[string]interface{}{
			"error":  err,
			"column": column,
			"type":   typ,
		})
	}
	return nil
}
func (m *MigrationManager) AddIndex(name string, columns ...string) error {
	err := m.connection.AddIndex(name, columns...).Error
	if err != nil {
		m.log.ErrorFields("Migratory:AddIndex", map[string]interface{}{
			"error":   err,
			"name":    name,
			"columns": columns,
		})
	}
	return nil
}
func (m *MigrationManager) AddUniqueIndex(name string, columns ...string) error {
	err := m.connection.AddUniqueIndex(name, columns...).Error
	if err != nil {
		m.log.ErrorFields("Migratory:AddUniqueIndex", map[string]interface{}{
			"error":   err,
			"name":    name,
			"columns": columns,
		})
	}
	return nil
}
func (m *MigrationManager) RemoveIndex(name string) error {
	err := m.connection.RemoveIndex(name).Error
	if err != nil {
		m.log.ErrorFields("Migratory:RemoveIndex", map[string]interface{}{
			"error": err,
			"name":  name,
		})
	}
	return nil
}
func (m *MigrationManager) RemoveForeignKey(name string, dest string) error {
	err := m.connection.RemoveForeignKey(name, dest).Error
	if err != nil {
		m.log.ErrorFields("Migratory:RemoveIndex", map[string]interface{}{
			"error": err,
			"name":  name,
		})
	}
	return nil
}
func (m *MigrationManager) AddForeignKey(name string, dest string, delete string, update string) error {
	err := m.connection.AddForeignKey(name, dest, delete, update).Error
	if err != nil {
		m.log.ErrorFields("Migratory:RemoveIndex", map[string]interface{}{
			"error": err,
			"name":  name,
		})
	}
	return nil
}

func (m *MigrationManager) Migrate(migrates []services.Migratable) error {
	return m.migrate(m.connection, migrates)
}

func (m *MigrationManager) Rollback(migrates []services.Migratable) error {
	return m.rollback(m.connection, migrates)
}

func (m *MigrationManager) Connection(connection string) services.Migratory {
	conn := m.client.connections[connection]
	if conn == nil {
		conn = m.client.runtimeConnection
	}

	n := &MigrationManager{
		migrationsTableName: m.migrationsTableName,
		connection:          m.client.connections[connection],
		log:                 m.log,
		client: 			 m.client,
	}
	return n
}

func (m *MigrationManager) setup(connection *gorm.DB) {
	if !connection.HasTable(&MigrationModel{}) {
		connection.AutoMigrate(&MigrationModel{})
	}
}

func (m *MigrationManager) migrate(connection *gorm.DB, migrates []services.Migratable) error {
	var migrations []*MigrationModel
	var migratableInterfaces []string

	for _, migratable := range migrates {
		migratableInterfaces = append(migratableInterfaces, reflect.TypeOf(migratable).String())
	}

	repo := m.client.repoManager.Query(&MigrationModel{})
	err := repo.Where("interface IN (?)", strings.Join(migratableInterfaces, ",")).GetAll(&migrations)
	if err != nil {
		m.log.ErrorFields("Migratory:Migrate", map[string]interface{}{
			"error":      err,
			"connection": reflect.TypeOf(connection).String(),
			"migrates":   reflect.TypeOf(migrates).String(),
		})
		return err
	}

	var maxStepV sql.NullInt64
	err = repo.Select("max(step)").GetValue(&maxStepV)
	if err != nil {
		m.log.ErrorFields("Migratory:Migrate", map[string]interface{}{
			"error":      err,
			"connection": reflect.TypeOf(connection).String(),
			"migrates":   reflect.TypeOf(migrates).String(),
		})
		return err
	}
	var maxStep int64 = 0
	if maxStepV.Valid {
		maxStep = maxStepV.Int64
	}

	for index, migratable := range migrates {
		if migration := getMigration(migrations, migratable); migration == nil {
			if err = migratable.Up(m); err == nil {
				migration = &MigrationModel{
					Step:      maxStep + 1,
					Interface: migratableInterfaces[index],
				}
				err = m.client.repoManager.Insert(migration)
				if err != nil {
					m.log.ErrorFields("Migratory:Migrate:Insert", map[string]interface{}{
						"error":      err,
						"connection": reflect.TypeOf(connection).String(),
						"migrates":   reflect.TypeOf(migrates).String(),
					})
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}

func (m *MigrationManager) rollback(connection *gorm.DB, migrates []services.Migratable) error {
	var migrations []*MigrationModel
	var migratableInterfaces []string

	for _, migratable := range migrates {
		migratableInterfaces = append(migratableInterfaces, reflect.TypeOf(migratable).String())
	}

	repo := m.client.repoManager.Query(&MigrationModel{})

	err := repo.Where("interface IN (?)", strings.Join(migratableInterfaces, ",")).GetAll(&migrations)
	if err != nil {
		m.log.ErrorFields("Migratory:Rollback", map[string]interface{}{
			"error":      err,
			"connection": reflect.TypeOf(connection).String(),
			"migrates":   reflect.TypeOf(migrates).String(),
		})
		return err
	}

	for _, migratable := range migrates {
		if migration := getMigration(migrations, migratable); migration != nil {
			if err = migratable.Down(m); err == nil {
				err = m.client.repoManager.DeleteEntity(migration)
				if err != nil {
					m.log.ErrorFields("Migratory:Rollback", map[string]interface{}{
						"error":      err,
						"connection": reflect.TypeOf(connection).String(),
						"migrates":   reflect.TypeOf(migrates).String(),
					})
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}

func getMigration(migrations []*MigrationModel, migratable services.Migratable) *MigrationModel {
	interfaceName := reflect.TypeOf(migratable).String()
	for _, migration := range migrations {
		if migration.Interface == interfaceName {
			return migration
		}
	}
	return nil
}

func NewMigrator(logger services.Logger, client *RDBMSClients, migrationsTable string) *MigrationManager {
	m := new(MigrationManager)
	m.log = logger
	m.client = client
	m.migrationsTableName = migrationsTable
	return m
}
