package dao

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"

	"github.com/mariusmagureanu/burlan/src/pkg/entities"
)

type DAO struct {
	dbSession *gorm.DB

	userDao UserDao
}

func (dao *DAO) Init(dbPath string) error {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	dao.dbSession = db

	dao.userDao = userDao{base: base{db: db}}

	return sqlDB.Ping()
}

func (dao *DAO) CreateTables() error {
	return dao.dbSession.AutoMigrate(&entities.User{}, &entities.Group{})
}

func (dao *DAO) DropTables() error {
	return dao.dbSession.Migrator().DropTable(&entities.User{}, &entities.Group{})
}

func (dao *DAO) Close() error {
	sqlDB, err := dao.dbSession.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

func (dao *DAO) Users() UserDao {
	return dao.userDao
}
