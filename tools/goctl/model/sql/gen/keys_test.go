package gen

import (
	"log"
	"testing"

	"github.com/tal-tech/go-zero/core/logx"
)

func TestKeys(t *testing.T) {
	var table = OuterTable{
		Table:          "user_info",
		CreateNotFound: true,
		Fields: []*OuterFiled{
			{
				IsPrimaryKey: true,
				Name:         "user_id",
				DataBaseType: "bigint",
				Comment:      "主键id",
			},
			{
				Name:         "campus_id",
				DataBaseType: "bigint",
				Comment:      "整校id",
				QueryType:    QueryAll,
				Cache:        false,
			},
			{
				Name:         "name",
				DataBaseType: "varchar",
				Comment:      "用户姓名",
				QueryType:    QueryOne,
			},
			{
				Name:         "id_number",
				DataBaseType: "varchar",
				Comment:      "身份证",
				Cache:        false,
				QueryType:    QueryNone,
				WithFields: []OuterWithField{
					{
						Name:         "name",
						DataBaseType: "varchar",
					},
				},
			},
			{
				Name:         "age",
				DataBaseType: "int",
				Comment:      "年龄",
				Cache:        false,
				QueryType:    QueryNone,
			},
			{
				Name:         "gender",
				DataBaseType: "tinyint",
				Comment:      "性别，0-男，1-女，2-不限",
				QueryType:    QueryLimit,
				WithFields: []OuterWithField{
					{
						Name:         "campus_id",
						DataBaseType: "bigint",
					},
				},
				OuterSort: []OuterSort{
					{
						Field: "create_time",
						Asc:   false,
					},
				},
			},
			{
				Name:         "mobile",
				DataBaseType: "varchar",
				Comment:      "手机号",
				QueryType:    QueryOne,
				Cache:        true,
			},
			{
				Name:         "create_time",
				DataBaseType: "timestamp",
				Comment:      "创建时间",
			},
			{
				Name:         "update_time",
				DataBaseType: "timestamp",
				Comment:      "更新时间",
			},
		},
	}
	innerTable, err := TableConvert(table)
	if err != nil {
		log.Fatalln(err)
	}
	tp, err := GenModel(innerTable)
	if err != nil {
		log.Fatalln(err)
	}
	logx.Info(tp)
}
