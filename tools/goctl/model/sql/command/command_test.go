package command

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/config"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/gen"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

var (
	sql = "-- 用户表 --\nCREATE TABLE `user` (\n  `id` bigint(10) NOT NULL AUTO_INCREMENT,\n  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名称',\n  `password` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户密码',\n  `mobile` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机号',\n  `gender` char(5) COLLATE utf8mb4_general_ci NOT NULL COMMENT '男｜女｜未公开',\n  `nickname` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '用户昵称',\n  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,\n  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,\n  PRIMARY KEY (`id`),\n  UNIQUE KEY `name_index` (`name`),\n  UNIQUE KEY `mobile_index` (`mobile`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;\n\n"
	cfg = &config.Config{
		NamingFormat: "gozero",
	}
)

func TestFromDDl(t *testing.T) {
	err := gen.Clean()
	assert.Nil(t, err)

	err = fromDDL("./user.sql", util.MustTempDir(), cfg, true, false, "go_zero")
	assert.Equal(t, errNotMatched, err)

	// case dir is not exists
	unknownDir := filepath.Join(util.MustTempDir(), "test", "user.sql")
	err = fromDDL(unknownDir, util.MustTempDir(), cfg, true, false, "go_zero")
	assert.True(t, func() bool {
		switch err.(type) {
		case *os.PathError:
			return true
		default:
			return false
		}
	}())

	// case empty src
	err = fromDDL("", util.MustTempDir(), cfg, true, false, "go_zero")
	if err != nil {
		assert.Equal(t, "expected path or path globbing patterns, but nothing found", err.Error())
	}

	tempDir := filepath.Join(util.MustTempDir(), "test")
	err = util.MkdirIfNotExist(tempDir)
	if err != nil {
		return
	}

	user1Sql := filepath.Join(tempDir, "user1.sql")
	user2Sql := filepath.Join(tempDir, "user2.sql")

	err = ioutil.WriteFile(user1Sql, []byte(sql), os.ModePerm)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(user2Sql, []byte(sql), os.ModePerm)
	if err != nil {
		return
	}

	_, err = os.Stat(user1Sql)
	assert.Nil(t, err)

	_, err = os.Stat(user2Sql)
	assert.Nil(t, err)

	filename := filepath.Join(tempDir, "usermodel.go")
	fromDDL := func(db string) {
		err = fromDDL(filepath.Join(tempDir, "user*.sql"), tempDir, cfg, true, false, db)
		assert.Nil(t, err)

		_, err = os.Stat(filename)
		assert.Nil(t, err)
	}

	fromDDL("go_zero")
	_ = os.Remove(filename)
	fromDDL("go-zero")
	_ = os.Remove(filename)
	fromDDL("1gozero")
}
