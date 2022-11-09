package ent

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"ariga.io/entcache"
	"entgo.io/ent/dialect"
	"entgo.io/ent"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

const DefaultMaxOpenCon = 100

type DatabaseConf struct {
	DbName       string `json:",optional"`
	SSLMode      bool   `json:",optional"`
	Host         string `json:",optional"`
	Port         int    `json:",optional"`
	User         string `json:",optional"`
	Password     string `json:",optional"`
	DbPath       string `json:",optional"`
	Type         string `json:",optional"` // "postgres" or "sqlite3" or "mysql"
	MaxOpenConns *int   `json:",optional,default=100"`
	Debug        bool   `json:",optional,default=false"`
	AutoMigrate  bool   `json:",optional,default=false"`
}

func (config DatabaseConf) NewDatabase(redisConfig RedisConf) (*ent., error) {
	var entOpts []ent.Option

	if config.Debug {
		logx.Info("Enabling Ent Client Request Debug")
		entOpts = append(entOpts, ent.Log(logx.Info))
		entOpts = append(entOpts, ent.Debug())
	}

	switch config.Type {
	case "sqlite":
		/*if it's the first startup, we want to touch and chmod file*/
		if _, err := os.Stat(config.DbPath); os.IsNotExist(err) {
			f, err := os.OpenFile(config.DbPath, os.O_CREATE|os.O_RDWR, 0600)
			if err != nil {
				return nil, fmt.Errorf("failed to create SQLite database file %q", config.DbPath)
			}
			if err := f.Close(); err != nil {
				return nil, fmt.Errorf("failed to create SQLite database file %q", config.DbPath)
			}
		} else {
			/*ensure file perms*/
			if err := os.Chmod(config.DbPath, 0660); err != nil {
				return nil, fmt.Errorf("unable to set perms on %s: %v", config.DbPath, err)
			}
		}
		drv, err := config.getEntDriver("sqlite3", dialect.SQLite, fmt.Sprintf("file:%s?_busy_timeout=100000&_fk=1", config.DbPath), redisConfig)
		if err != nil {
			return nil, errors.Wrapf(err, "failed opening connection to sqlite: %v", config.DbPath)
		}
		entOpts = append(entOpts, ent.Driver(drv))
	case "mysql":
		drv, err := config.getEntDriver("mysql", dialect.MySQL, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True", config.User, config.Password, config.Host, config.Port, config.DbName), redisConfig)
		if err != nil {
			return nil, fmt.Errorf("failed opening connection to mysql: %v", err)
		}
		entOpts = append(entOpts, ent.Driver(drv))
	case "postgres", "postgresql":
		drv, err := config.getEntDriver("postgres", dialect.Postgres, fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s", config.Host, config.Port, config.User, config.DbName, config.Password, config.SSLMode), redisConfig)
		if err != nil {
			return nil, fmt.Errorf("failed opening connection to postgresql: %v", err)
		}
		entOpts = append(entOpts, ent.Driver(drv))
	default:
		return nil, fmt.Errorf("unknown database type '%s'", config.Type)
	}

	return ent.NewClient(entOpts...), nil
}

func (config DatabaseConf) getEntDriver(dbtype string, dbdialect string, dsn string, redisConfig RedisConf) (*entcache.Driver, error) {
	db, err := sql.Open(dbtype, dsn)

	if err != nil {
		logx.Infof("failed opening connection to %s: %v", dbtype, err)
		return nil, err
	}

	if config.MaxOpenConns == nil {
		logx.Infof("MaxOpenConns is 0, defaulting to %d", DefaultMaxOpenCon)
	}

	db.SetMaxOpenConns(*config.MaxOpenConns)
	drv := entsql.OpenDB(dbdialect, db)

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprint(redisConfig.Host, ":", redisConfig.Port),
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logx.Error("cache: redis ping error")
		return nil, fmt.Errorf("cache: redis ping error: %v", err)
	}

	cacheDrv := entcache.NewDriver(
		drv,
		entcache.TTL(time.Minute),
		entcache.Levels(
			entcache.NewLRU(256),
			entcache.NewRedis(rdb),
		),
	)

	return cacheDrv, nil
}
