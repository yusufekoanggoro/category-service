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
	Close() error
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
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
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
			return fmt.Errorf("failed to migrate: %w", err)
		}
	}
	return nil
}

func (g *GormDatabase) GetDB() *gorm.DB {
	return g.db
}

func (g *GormDatabase) Close() error {
	if g.db != nil {
		sqlDB, err := g.db.DB()
		if err != nil {
			return fmt.Errorf("failed to get database instance: %w", err)
		}
		return sqlDB.Close()
	}
	return nil
}
