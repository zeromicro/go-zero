package gormsql

import (
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GORMConf struct {
	Type        string `json:"Type"`        // type of database: mysql, postgres
	Path        string `json:"Path"`        // address
	Port        int    `json:"Port"`        // port
	Config      string `json:"Config"`      // extra config such as charset=utf8mb4&parseTime=True
	Dbname      string `json:"DBName"`      // database name
	Username    string `json:"Username"`    // username
	Password    string `json:"Password"`    // password
	MaxIdleConn int    `json:"MaxIdleConn"` // the maximum number of connections in the idle connection pool
	MaxOpenConn int    `json:"MaxOpenConn"` // the maximum number of open connections to the database
	LogMode     string `json:"LogMode"`     // open gorm's global logger
	LogZap      bool   `json:"LogZap"`
}

func (g GORMConf) MysqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", g.Username, g.Password, g.Path, g.Port, g.Dbname, g.Config)
}

func (g GORMConf) PostgresDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d %s", g.Path, g.Username, g.Password,
		g.Dbname, g.Port, g.Config)
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
	if c.Dbname == "" {
		return nil, errors.New("database name cannot be nil")
	}
	mysqlConfig := mysql.Config{
		DSN:                       c.MysqlDSN(),
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{}); err != nil {
		return nil, err
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(c.MaxIdleConn)
		sqlDB.SetMaxOpenConns(c.MaxOpenConn)
		return db, nil
	}
}

func GormPgSql(c GORMConf) (*gorm.DB, error) {
	if c.Dbname == "" {
		return nil, errors.New("database name cannot be nil")
	}
	pgsqlConfig := postgres.Config{
		DSN:                  c.PostgresDSN(),
		PreferSimpleProtocol: false, // disables implicit prepared statement usage
	}
	if db, err := gorm.Open(postgres.New(pgsqlConfig), &gorm.Config{}); err != nil {
		return nil, err
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(c.MaxIdleConn)
		sqlDB.SetMaxOpenConns(c.MaxOpenConn)
		return db, nil
	}
}
