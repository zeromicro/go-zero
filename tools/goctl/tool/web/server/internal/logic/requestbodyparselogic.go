package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/svc"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/types"
	"github.com/zeromicro/go-zero/tools/goctl/util"
)

const (
	typeInterface  = "interface{}"
	typeInt64      = "int64"
	typeFloat64    = "float64"
	typeString     = "string"
	prefixSlice    = "[]"
	typeJSONNumber = "json.Number"
	typeBool       = "bool"
)

type RequestBodyParseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRequestBodyParseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RequestBodyParseLogic {
	return &RequestBodyParseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RequestBodyParseLogic) RequestBodyParse(req *types.ParseJsonRequest) (resp *types.ParseJsonResponse, err error) {
	resp = new(types.ParseJsonResponse)
	data, err := parseJSON(req.JSON)
	if err != nil {
		return nil, err
	}

	resp.Form = data
	return resp, nil
}

func parseJSON(s string) ([]*types.FormItem, error) {
	if util.IsEmptyStringOrWhiteSpace(s) {
		return []*types.FormItem{}, nil
	}

	var v any
	if err := jsonDecode(s, &v); err != nil {
		return nil, err
	}

	m, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected map[string]any, got %T", v)
	}

	var resp []*types.FormItem
	for fieldName, fieldType := range m {
		tp, err := parseType(fieldName, fieldType)
		if err != nil {
			return nil, err
		}

		resp = append(resp, &types.FormItem{
			Name:     fieldName,
			Type:     tp,
			Optional: false,
		})
	}

	sort.SliceStable(resp, func(i, j int) bool {
		return resp[i].Name < resp[j].Name
	})
	return resp, nil
}

func parseType(childName string, v any) (string, error) {
	tp := reflect.TypeOf(v)
	if tp == nil {
		return typeInterface, nil
	}

	switch {
	case tp.String() == typeJSONNumber:
		number := v.(json.Number)
		_, err := number.Int64()
		if err == nil {
			return typeInt64, nil
		}

		_, err = number.Float64()
		if err == nil {
			return typeFloat64, nil
		}

		return typeString, nil
	case tp.Kind() >= reflect.Bool && tp.Kind() <= reflect.Float64,
		tp.Kind() == reflect.String:
		return tp.String(), nil
	case tp.Kind() == reflect.Slice:
		slice := v.([]any)
		if len(slice) == 0 {
			return prefixSlice + typeInterface, nil
		}

		elemType, err := parseType(childName, slice[0])
		if err != nil {
			return "", err
		}

		if strings.HasPrefix(elemType, prefixSlice) {
			return "", fmt.Errorf("child %q, slice item must be basic type, but got %s", childName, elemType)
		}
		return prefixSlice + elemType, nil
	default:
		return "", fmt.Errorf("child %q, expected basic golang type, got go golang type %s", childName, tp.String())
	}
}
