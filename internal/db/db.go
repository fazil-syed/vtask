package db

// we will have a function that initializes the database connection using gorm and returns the gorm DB instance
import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Path  string
	Debug bool
}

func Open(cfg Config) (*gorm.DB, error) {
	gormCfg := &gorm.Config{}
	if cfg.Debug {
		gormCfg.Logger = logger.Default.LogMode(logger.Info)
	}
	db, err := gorm.Open(sqlite.Open(cfg.Path), gormCfg)
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	sqlDB.SetMaxOpenConns(1)
	return db, nil

}
