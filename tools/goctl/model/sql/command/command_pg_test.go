package command

import (
	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/config"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/gen"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"path/filepath"
	"testing"
)

// generate test table sql
/*
CREATE TABLE "public"."users" (
  "id" serial NOT NULL,
  "account" varchar(20) COLLATE "pg_catalog"."default" NOT NULL,
  "avatar" text COLLATE "pg_catalog"."default",
  "nick_name" varchar(60) COLLATE "pg_catalog"."default",
  "register_time" timestamp(6) NOT NULL,
  "update_time" timestamp(6),
  "password" varchar(255) COLLATE "pg_catalog"."default",
  "email" varchar(100) COLLATE "pg_catalog"."default",
  "reset_key" varchar(10) COLLATE "pg_catalog"."default",
  "active" bool NOT NULL DEFAULT true,
  CONSTRAINT "users_pk" PRIMARY KEY ("id")
)
;

ALTER TABLE "public"."users"
  OWNER TO "postgres";
*/

func TestFromDatasource(t *testing.T) {
	err := gen.Clean()
	assert.Nil(t, err)

	url := "postgres://postgres:postgres@127.0.0.1:5432/demo?sslmode=disable"

	pattern := "users" // table name

	cfg, err := config.NewConfig("")
	tempDir := filepath.Join(util.MustTempDir(), "test")
	err = fromPostgreSqlDataSource(url, pattern, tempDir, "public", cfg, false, false)
	assert.Nil(t, err)
}
