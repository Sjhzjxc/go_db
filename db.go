package go_db

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func NewDb(config *Config, writer logger.Writer) (*DBModel, error) {
	dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=%s&parseTime=%t&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		"information_schema",
		config.Charset,
		config.ParseTime)
	if writer == nil {
		writer = log.New(os.Stdout, "\r\n", log.LstdFlags)
	}
	var logLevel logger.LogLevel
	if config.LogLevel == "silent" {
		logLevel = logger.Silent
	} else if config.LogLevel == "info" {
		logLevel = logger.Info
	} else if config.LogLevel == "warn" {
		logLevel = logger.Warn
	} else if config.LogLevel == "error" {
		logLevel = logger.Error
	} else {
		logLevel = logger.Info
	}
	newLogger := logger.New(
		writer,
		logger.Config{
			SlowThreshold:             time.Duration(config.SlowThreshold) * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logLevel,                                               // Log level
			IgnoreRecordNotFoundError: config.IgnoreRecordNotFoundError,                       // Ignore ErrRecordNotFound error for logger
			Colorful:                  config.Colorful,                                        // Disable color
		},
	)
	db, err := gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   newLogger,
			DryRun:                                   config.DryRun,
		})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(config.MaxLifeTime))
	// 查表，看数据库是否存在
	var schemas []string
	err = db.Table("SCHEMATA").Select("SCHEMA_NAME").Find(&schemas).Error
	if err != nil {
		return nil, err
	}
	if ArrayInBool(schemas, config.DbName) {
		err = db.Exec("use " + config.DbName + ";").Error
		if err != nil {
			return nil, err
		}
	} else {
		if config.CreateSchemaIfNotExist {
			err = recreateDatabase(db, config.DbName, config.Charset)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New(fmt.Sprintf("数据库%s不存在", config.DbName))
		}
	}
	return &DBModel{
		DB:      db,
		Configs: config,
	}, nil
}

func DefaultDb(username, password, host, dbName string, writer logger.Writer) (*DBModel, error) {
	config := &Config{
		Username:                  username,
		Password:                  password,
		Host:                      host,
		DbName:                    dbName,
		Charset:                   "utf8mb4",
		ParseTime:                 true,
		MaxIdleConns:              200,
		MaxOpenConns:              200,
		MaxLifeTime:               300,
		CreateSchemaIfNotExist:    true,
		SlowThreshold:             200,
		LogLevel:                  "info",
		Colorful:                  false,
		IgnoreRecordNotFoundError: true,
		DryRun:                    false,
	}
	return NewDb(config, writer)
}

func recreateDatabase(db *gorm.DB, dbName, charset string) error {
	err := db.Exec("use information_schema;").Error
	if err != nil {
		return err
	}
	err = db.Exec("drop database if exists " + dbName + ";").Error
	if err != nil {
		return err
	}
	err = db.Exec("create database " + dbName + " default character set " + charset + ";").Error
	if err != nil {
		return err
	}
	err = db.Exec("use " + dbName + ";").Error
	return err
}

func (m DBModel) RecreateDatabase() error {
	return recreateDatabase(m.DB, m.Configs.DbName, m.Configs.Charset)
}

func (m DBModel) AutoMigration(recreate bool, models ...interface{}) error {
	if recreate {
		err := m.RecreateDatabase()
		if err != nil {
			return err
		}
	}
	err := m.DB.Migrator().AutoMigrate(models...)
	return err
}

func Paginate(page, size int, tx *gorm.DB, dest *[]interface{}) (*Pagination, error) {
	if size < 1 {
		size = 10
	}
	if page < 1 {
		page = 1
	}
	var total int64
	offset := size * (page - 1)
	err := tx.Count(&total).Error
	if err != nil {
		return nil, err
	}

	err = tx.Limit(size).Offset(offset).Find(dest).Error
	if err != nil {
		return nil, err
	}

	return &Pagination{
		HasNext: int(total)-offset > size,
		Total:   total,
		Page:    page,
		Size:    size,
		Items:   dest,
	}, err
}
