package database

import (
	"github.com/peyman-abdi/avalanche/app/interfaces"
	"github.com/jinzhu/gorm"
)

type RepositoryManager struct {
	entity interfaces.Entity
	query *gorm.DB
	db *gorm.DB
}
var _ interfaces.Repository = (*RepositoryManager)(nil)

func (r *RepositoryManager) Query(entity interfaces.Entity) interfaces.QueryBuilder  {
	var db *gorm.DB
	if conn, ok := r.entity.(interfaces.EntityConnection); ok {
		db = connections[conn.ConnectionName()]
		if db != nil {
			db = db.Table(entity.TableName())
		}
	} else {
		db = connections["app"].Table(entity.TableName())
	}

	c := &RepositoryManager{
		entity: entity,
		query: db,
		db: db,
	}
	return c
}

func (r *RepositoryManager) Insert(entity interfaces.Entity) error  {
	var db *gorm.DB
	if conn, ok := r.entity.(interfaces.EntityConnection); ok {
		db = connections[conn.ConnectionName()]
		if db != nil {
			db = db.Table(entity.TableName())
		}
	} else {
		db = connections["app"].Table(entity.TableName())
	}

	db.Create(entity)
	return nil
}

func (r *RepositoryManager) UpdateEntity(entity interfaces.Entity) error  {
	var db *gorm.DB
	if conn, ok := r.entity.(interfaces.EntityConnection); ok {
		db = connections[conn.ConnectionName()]
		if db != nil {
			db = db.Table(entity.TableName())
		}
	} else {
		db = connections["app"].Table(entity.TableName())
	}

	db.Update(entity)
	return nil
}

func (r *RepositoryManager) Update(values...interface{}) error  {
	r.db.Updates(values, true)
	return nil
}

func (r *RepositoryManager) Where(query interface{}, args ...interface{}) interfaces.QueryBuilder {
	r.db = r.db.Where(query, args)
	return r
}

func (r *RepositoryManager) Select(query interface{}, args ...interface{}) interfaces.QueryBuilder {
	r.db = r.db.Select(query, args)
	return r
}

func (r *RepositoryManager) Limit(limit interface{}) interfaces.QueryBuilder {
	r.db = r.db.Limit(limit)
	return r
}

func (r *RepositoryManager) Offset(offset interface{}) interfaces.QueryBuilder {
	r.db = r.db.Offset(offset)
	return r
}

func (r *RepositoryManager) Get(result interface{}) error {
	rows, err := r.db.Rows()
	r.db = r.query

	if err != nil {
		return err
	}

	err = rows.Scan(result)
	if err != nil {
		return err
	}

	return err
}

func (r *RepositoryManager) OrWhere(query interface{}, args ...interface{}) interfaces.QueryBuilder {
	r.db = r.db.Or(query, args)
	return r
}










