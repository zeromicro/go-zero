package docker

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var hooks = make(map[string]HookFunc)

type HookFunc func(*container) error

func Register(name string, hookFunc HookFunc) {
	if hooks[name] != nil {
		panic(fmt.Sprintf("%s already registed", name))
	}
	hooks[name] = hookFunc
}

func init() {
	Register("refresh_mysql", mysqlHook)
}

func mysqlHook(c *container) error {
	dsn := getDSN(c)
	if err := cleanMysql(dsn); err != nil {
		return err
	}
	return initMysql(c, dsn)
}

func getDSN(c *container) string {
	var (
		user = "root"
		host = "127.0.0.1"
		port = "13306"
		pwd  string
	)
	envs := make(map[string]string)
	for _, env := range c.env {
		a := strings.Split(env, "=")
		if len(a) == 2 {
			envs[a[0]] = a[1]
		}
	}
	if len(c.imageCfg.Ports) > 0 {
		for _, k := range c.imageCfg.Ports {
			a := strings.Split(k, ":")
			if len(a) != 2 {
				continue
			}
			if a[1] == "3306" {
				port = a[0]
			}
		}
	}
	if envs["MYSQL_ROOT_PASSWORD"] != "" {
		pwd = envs["MYSQL_ROOT_PASSWORD"]
	}
	if envs["MYSQL_USER"] != "" {
		user = envs["MYSQL_USER"]
	}
	if envs["MYSQL_PASSWORD"] != "" {
		pwd = envs["MYSQL_PASSWORD"]
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/?multiStatements=true", user, pwd, host, port)
}

func sqlPath(ct *container) (res string) {
	for _, m := range ct.mounts {
		if m.Target == "/docker-entrypoint-initdb.d" {
			return m.Source
		}
	}
	return
}

func cleanMysql(dsn string) (err error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return
	}
	c := context.Background()
	defer db.Close()
	rows, err := db.QueryContext(c, "show databases")
	if err != nil {
		return
	}
	var dbs []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		dbs = append(dbs, name)
	}
	dbs = businessDbs(dbs)
	for _, name := range dbs {
		_, err = db.ExecContext(c, fmt.Sprintf("drop database %s", name))
		if err != nil {
			return
		}
	}
	return
}

func initMysql(ct *container, dsn string) (err error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return
	}
	c := context.Background()
	defer db.Close()
	pathDir := sqlPath(ct)
	files, err := ioutil.ReadDir(pathDir)
	if err != nil {
		return
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if !strings.HasSuffix(f.Name(), ".sql") {
			continue
		}
		content, err := ioutil.ReadFile(filepath.Join(pathDir, f.Name()))
		if err != nil {
			return err
		}
		_, err = db.ExecContext(c, string(content))
		if err != nil {
			return err
		}
	}
	return
}

func businessDbs(dbs []string) (res []string) {
	for _, db := range dbs {
		if db == "information_schema" || db == "mysql" || db == "performance_schema" {
			continue
		}
		res = append(res, db)
	}
	return
}
