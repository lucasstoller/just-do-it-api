package database

import (
	"gorm.io/gorm"
)

type Database interface {
	Find(dest interface{}, conds ...interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	AutoMigrate(dst ...interface{}) error
}

type GormDB struct {
	db *gorm.DB
}

func (g *GormDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	return g.db.Find(dest, conds...)
}

func (g *GormDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	return g.db.First(dest, conds...)
}

func (g *GormDB) Create(value interface{}) *gorm.DB {
	return g.db.Create(value)
}

func (g *GormDB) Save(value interface{}) *gorm.DB {
	return g.db.Save(value)
}

func (g *GormDB) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	return g.db.Delete(value, conds...)
}

func (g *GormDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return g.db.Where(query, args...)
}

func (g *GormDB) AutoMigrate(dst ...interface{}) error {
	return g.db.AutoMigrate(dst...)
}

var db Database

func CreateConnection() Database {
	if db != nil {
		return db
	}

	gormDB := initDB()
	db = &GormDB{db: gormDB}
	return db
}

// For testing purposes
func SetTestDB(testDB Database) {
	db = testDB
}
