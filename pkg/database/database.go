package database

import (
	"category-service/config"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database interface {
	Connect(cfg config.ConfigProvider) error
	AutoMigrate(models ...interface{}) error
	GetDB() *gorm.DB
}

type GormDatabase struct {
	db *gorm.DB
}

func (g *GormDatabase) Connect(cfg config.ConfigProvider) error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.GetDBHost(), cfg.GetDBPort(), cfg.GetDBUser(), cfg.GetDBPassword(), cfg.GetDBName(), cfg.GetSSLMode())

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("gagal menghubungkan ke database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("gagal mendapatkan instance database: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	g.db = db
	return nil
}

func (g *GormDatabase) AutoMigrate(models ...interface{}) error {
	if len(models) > 0 {
		if err := g.db.AutoMigrate(models...); err != nil {
			return fmt.Errorf("gagal melakukan migrasi: %w", err)
		}
	}
	return nil
}

func (g *GormDatabase) GetDB() *gorm.DB {
	return g.db
}
