/**
 * @Author: chentong
 * @Date: 2025/02/20 01:12
 */

package token

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4/request"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

type CookieExtractor []string

func (e CookieExtractor) ExtractToken(req *http.Request) (string, error) {
	for _, name := range e {
		if ah, err := req.Cookie(name); err == nil {
			return ah.Value, nil
		}
	}
	return "", request.ErrNoTokenInRequest
}

type ParamExtractor []string

func (e ParamExtractor) ExtractToken(req *http.Request) (string, error) {
	vars := pathvar.Vars(req)
	for _, name := range e {
		if ah, ok := vars[name]; ok {
			return ah, nil
		}

	}
	return "", request.ErrNoTokenInRequest
}
