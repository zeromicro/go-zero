package sqlmodel

import (
	"go/format"

	"zero/core/logx"
)

type (
	Column struct {
		DataType string `json:"dataType"`
		Name     string `json:"name"`
		Comment  string `json:"comment"`
	}

	Table struct {
		Name       string    `json:"table"`
		PrimaryKey string    `json:"primaryKey"`
		Columns    []*Column `json:"columns"`
	}
)

func GenMysqlGoModel(req *Table, conditions []string) (string, error) {
	resp, err := generateTypeModel(req.Name, req.Columns)
	if err != nil {
		return "", err
	}
	resp.PrimaryKey = req.PrimaryKey
	s, err := NewStructExp(*resp)
	if len(conditions) == 0 {
		s.conditions = []string{resp.PrimaryKey}
	} else {
		s.conditions = conditions
	}
	if err != nil {
		return "", err
	}
	result, err := s.genMysqlCRUD()
	// code format
	bts, err := format.Source([]byte(result))
	if err != nil {
		logx.Errorf("%+v", err)
		return "", err
	}
	return string(bts), err
}
