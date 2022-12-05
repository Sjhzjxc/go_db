package go_db

import (
	"gorm.io/gorm"
)

type DBModel struct {
	DB      *gorm.DB
	Configs *Config
}

// SlowThreshold SQL慢查询 毫秒
// LogLevel info warn error silent
// CreateSchemaIfNotExist 账号权限不足会导致无法创建,报错
// DryRun true 所有sql都不会执行
type Config struct {
	Username                  string
	Password                  string
	Host                      string
	DbName                    string
	Charset                   string
	ParseTime                 bool
	MaxIdleConns              int
	MaxOpenConns              int
	MaxLifeTime               int
	CreateSchemaIfNotExist    bool
	SlowThreshold             int // 毫秒
	LogLevel                  string
	Colorful                  bool
	IgnoreRecordNotFoundError bool
	DryRun                    bool
}

type Pagination struct {
	HasNext bool           `json:"has_next"`
	Total   int64          `json:"total"`
	Page    int            `json:"page"`
	Size    int            `json:"size"`
	Items   *[]interface{} `json:"items"`
}
