package database

import (
	"github.com/peyman-abdi/avalanche/app/interfaces"
	"github.com/jinzhu/gorm"
	"reflect"
)

type RepositoryManager struct {
	entity interfaces.Entity
	includeDeletes bool
	query *gorm.DB
	db *gorm.DB
	log interfaces.Logger
}
var _ interfaces.Repository = (*RepositoryManager)(nil)

func (r *RepositoryManager) Query(entity interfaces.Entity) interfaces.QueryBuilder  {
	var db *gorm.DB
	conn, ok := entity.(interfaces.EntityConnection)
	if ok {
		db = connections[conn.ConnectionName()]
		if db != nil {
			db = db.Table(entity.TableName())
		}
	}

	if db == nil {
		db = connections["app"].Table(entity.TableName())
	}

	c := &RepositoryManager{
		entity: entity,
		query: db,
		db: db,
		log: r.log,
		includeDeletes: false,
	}
	return c
}

func (r *RepositoryManager) Insert(entity interfaces.Entity) error  {
	var db *gorm.DB
	if conn, ok := entity.(interfaces.EntityConnection); ok {
		db = connections[conn.ConnectionName()]
		if db != nil {
			db = db.Table(entity.TableName())
		}
	} else {
		db = connections["app"].Table(entity.TableName())
	}

	if err := db.Create(entity).Error; err != nil {
		r.log.ErrorFields("Repository Error in Insert", map[string]interface{} {
			"error": err,
		})
		return err
	}
	return nil
}

func (r *RepositoryManager) UpdateEntity(entity interfaces.Entity) error  {
	var db *gorm.DB
	if conn, ok := entity.(interfaces.EntityConnection); ok {
		db = connections[conn.ConnectionName()]
		if db != nil {
			db = db.Table(entity.TableName())
		}
	} else {
		db = connections["app"].Table(entity.TableName())
	}

	if err := db.Save(entity).Error; err != nil {
		r.log.ErrorFields("Repository Error in UpdateEntity", map[string]interface{} {
			"error": err,
		})
		return err
	}
	return nil
}

func (r *RepositoryManager) DeleteEntity(entity interfaces.Entity) error  {
	var db *gorm.DB
	if conn, ok := entity.(interfaces.EntityConnection); ok {
		db = connections[conn.ConnectionName()]
		if db != nil {
			db = db.Table(entity.TableName())
		}
	} else {
		db = connections["app"].Table(entity.TableName())
	}

	if err := db.Delete(entity).Error; err != nil {
		r.log.ErrorFields("Repository Error in DeleteEntity", map[string]interface{} {
			"error": err,
		})
		return err
	}
	return nil
}

func (r *RepositoryManager) EraseEntity(entity interfaces.Entity) error  {
	var db *gorm.DB
	if conn, ok := entity.(interfaces.EntityConnection); ok {
		db = connections[conn.ConnectionName()]
		if db != nil {
			db = db.Table(entity.TableName())
		}
	} else {
		db = connections["app"].Table(entity.TableName())
	}

	if err := db.Unscoped().Delete(entity).Error; err != nil {
		r.log.ErrorFields("Repository Error in EraseEntity", map[string]interface{} {
			"error": err,
		})
		return err
	}
	return nil
}

func (r *RepositoryManager) GetValue(result interface{}) error {
	if !r.includeDeletes {
		ty := reflect.TypeOf(r.entity)
		if ty.Kind() == reflect.Ptr {
			ty = ty.Elem()
		}
		if _, has := ty.FieldByName("DeletedAt"); has {
			r.db = r.db.Where("deleted_at is null")
		}
	}

	err := r.db.Row().Scan(result)
	r.Reset()

	if err != nil {
		r.log.ErrorFields("Repository Error in GetValue", map[string]interface{} {
			"error": err,
		})
		return err
	}

	return nil
}

func (r *RepositoryManager) GetValues(resultsArray interface{}) error {
	if !r.includeDeletes {
		ty := reflect.TypeOf(r.entity)
		if ty.Kind() == reflect.Ptr {
			ty = ty.Elem()
		}
		if _, has := ty.FieldByName("DeletedAt"); has {
			r.db = r.db.Where("deleted_at is null")
		}
	}
	rows, err := r.db.Rows()
	defer rows.Close()

	r.db = r.query

	if err != nil {
		r.log.ErrorFields("Repository Error in GetValues", map[string]interface{} {
			"error": err,
		})
		return err
	}

	elemArrType := reflect.TypeOf(resultsArray).Elem()
	elemPlaceholder := reflect.New(elemArrType.Elem()).Interface()
	var results = reflect.MakeSlice(elemArrType, 0, 0)

	for rows.Next() {
		rows.Scan(elemPlaceholder)
		results = reflect.Append(results, reflect.ValueOf(elemPlaceholder).Elem())
	}

	v := reflect.ValueOf(resultsArray).Elem()
	v.Set(reflect.ValueOf(results.Interface()))

	return nil
}


func (r *RepositoryManager) GetFirst(entityPtr interface{}) error {
	ty := reflect.ValueOf(entityPtr).Type()
	if ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
	}
	ref := reflect.New(reflect.SliceOf(ty))

	err := r.db.First(ref.Interface()).Error
	r.Reset()

	if ref.Elem().Len() == 1 {
		result := ref.Elem().Index(0).Interface()
		v := reflect.ValueOf(entityPtr)
		v.Elem().Set(reflect.ValueOf(result))
	}

	if err != nil {
		r.log.ErrorFields("Repository Error in GetFirst", map[string]interface{} {
			"error": err,
		})
		return err
	}

	return nil
}

func (r *RepositoryManager) GetLast(entity interface{}) error {
	ty := reflect.ValueOf(entity).Type()
	if ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
	}
	ref := reflect.New(reflect.SliceOf(ty))

	err := r.db.Last(ref.Interface()).Error
	r.Reset()

	if ref.Elem().Len() == 1 {
		result := ref.Elem().Index(0).Interface()
		v := reflect.ValueOf(entity)
		v.Elem().Set(reflect.ValueOf(result))
	}

	if err != nil {
		r.log.ErrorFields("Repository Error in GetLast", map[string]interface{} {
			"error": err,
		})
		return err
	}

	return nil
}

func (r *RepositoryManager) Get(entitiesArray interface{}) error {
	err := r.db.Find(entitiesArray).Error
	r.Reset()

	if err != nil {
		r.log.ErrorFields("Repository Error in Get", map[string]interface{} {
			"error": err,
		})
		return err
	}

	return err
}

func (r *RepositoryManager) Update(entityRef interface{}, values map[string]interface{}) error  {
	err := r.db.Model(entityRef).Updates(values).Error
	r.Reset()
	return err
}

func (r *RepositoryManager) Updates(values map[string]interface{}) error  {
	err := r.db.Updates(values).Error
	r.Reset()
	return err
}

func (r *RepositoryManager) DeleteAll() error  {
	err := r.db.Delete(r.entity).Error
	r.Reset()
	return err
}

func (r *RepositoryManager) EraseAll() error  {
	err := r.db.Unscoped().Delete(r.entity).Error
	r.Reset()
	return err
}

func (r *RepositoryManager) Reset() interfaces.QueryBuilder {
	r.db = r.query
	r.includeDeletes = false
	return r
}

func (r *RepositoryManager) IncludeDeleted() interfaces.QueryBuilder {
	r.db = r.db.Unscoped()
	r.includeDeletes = true
	return r
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


func (r *RepositoryManager) Order(value interface{}) interfaces.QueryBuilder {
	r.db = r.db.Order(value)
	return r
}

func (r *RepositoryManager) OrWhere(query interface{}, args ...interface{}) interfaces.QueryBuilder {
	r.db = r.db.Or(query, args)
	return r
}










