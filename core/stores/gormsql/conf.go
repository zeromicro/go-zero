package gormsql

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GORMConf struct {
	Type        string `json:"Type"`        // type of database: mysql, postgres
	Host        string `json:"Host"`        // address
	Port        int    `json:"Port"`        // port
	Config      string `json:"Config"`      // extra config such as charset=utf8mb4&parseTime=True
	DbName      string `json:"DBName"`      // database name
	Username    string `json:"Username"`    // username
	Password    string `json:"Password"`    // password
	MaxIdleConn int    `json:"MaxIdleConn"` // the maximum number of connections in the idle connection pool
	MaxOpenConn int    `json:"MaxOpenConn"` // the maximum number of open connections to the database
	LogMode     string `json:"LogMode"`     // open gorm's global logger
}

func (g GORMConf) MysqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", g.Username, g.Password, g.Host, g.Port, g.DbName, g.Config)
}

func (g GORMConf) PostgresDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d %s", g.Host, g.Username, g.Password,
		g.DbName, g.Port, g.Config)
}

func (g GORMConf) NewGORM() (*gorm.DB, error) {
	switch g.Type {
	case "mysql":
		return GormMysql(g)
	case "pgsql":
		return GormPgSql(g)
	default:
		return GormMysql(g)
	}
}

func GormMysql(c GORMConf) (*gorm.DB, error) {
	if c.DbName == "" {
		return nil, errors.New("database name cannot be nil")
	}
	mysqlConfig := mysql.Config{
		DSN:                       c.MysqlDSN(),
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // autoconfiguration based on currently MySQL version
	}

	if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		Logger: logger.New(gormWriter{}, logger.Config{
			SlowThreshold:             1 * time.Second,
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  getLevel(c.LogMode),
		}),
	}); err != nil {
		return nil, err
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(c.MaxIdleConn)
		sqlDB.SetMaxOpenConns(c.MaxOpenConn)
		return db, nil
	}
}

func GormPgSql(c GORMConf) (*gorm.DB, error) {
	if c.DbName == "" {
		return nil, errors.New("database name cannot be nil")
	}
	pgsqlConfig := postgres.Config{
		DSN:                  c.PostgresDSN(),
		PreferSimpleProtocol: false, // disables implicit prepared statement usage
	}

	if db, err := gorm.Open(postgres.New(pgsqlConfig), &gorm.Config{
		Logger: logger.New(gormWriter{}, logger.Config{
			SlowThreshold:             1 * time.Second,
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  getLevel(c.LogMode),
		}),
	}); err != nil {
		return nil, err
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(c.MaxIdleConn)
		sqlDB.SetMaxOpenConns(c.MaxOpenConn)
		return db, nil
	}
}

func getLevel(logMode string) logger.LogLevel {
	var level logger.LogLevel
	switch logMode {
	case "info":
		level = logger.Info
	case "warn":
		level = logger.Warn
	case "error":
		level = logger.Error
	case "silent":
		level = logger.Silent
	default:
		level = logger.Error
	}
	return level
}
