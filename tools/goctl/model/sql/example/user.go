package main

import (
	"flag"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/example/internal/config"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/example/model"
)

var configFile = flag.String("f", "etc/config.json", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)

	logx.MustSetup(c.LogConf)
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	userModel := model.NewUserModel(conn, c.CacheRedis, c.Mysql.Table.User)

	var user model.User
	user.Name = "test"
	user.Password = "123456"
	user.Mobile = "136****0001"
	user.Gender = "男"
	user.Nickname = "Keson"
	insert, err := userModel.Insert(user)
	logx.Must(err)
	insertId, err := insert.LastInsertId()
	logx.Must(err)
	logx.Infof("insert success,the last insert id:%v", insertId)
}
