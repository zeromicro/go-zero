package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/rest"
	"github.com/tal-tech/go-zero/rest/httpx"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

const jwtUserField = "user"

type (
	Config struct {
		rest.RestConf
		AccessSecret  string
		AccessExpire  int64 `json:",default=1209600"` // 2 weeks
		RefreshSecret string
		RefreshExpire int64 `json:",default=2419200"` // 4 weeks
		RefreshAfter  int64 `json:",default=604800"`  // 1 week
	}

	TokenOptions struct {
		AccessSecret  string
		AccessExpire  int64
		RefreshSecret string
		RefreshExpire int64
		RefreshAfter  int64
		Fields        map[string]interface{}
	}

	Tokens struct {
		// Access token to access the apis
		AccessToken string `json:"access_token"`
		// Access token expire time, generated like: time.Now().Add(time.Day*14).Unix()
		AccessExpire int64 `json:"access_expire"`
		// Refresh token, use this to refresh the token
		RefreshToken string `json:"refresh_token"`
		// Refresh token expire time, generated like: time.Now().Add(time.Month).Unix()
		RefreshExpire int64 `json:"refresh_expire"`
		// Recommended time to refresh the access token
		RefreshAfter int64 `json:"refresh_after"`
	}

	UserCredentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	User struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	Response struct {
		Data string `json:"data"`
	}

	Token struct {
		Token string `json:"token"`
	}

	AuthRequest struct {
		User string `json:"u"`
	}
)

func main() {
	var c Config
	conf.MustLoad("user.json", &c)

	engine, err := rest.NewServer(c.RestConf)
	if err != nil {
		log.Fatal(err)
	}
	defer engine.Stop()

	engine.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/login",
		Handler: LoginHandler(c),
	})
	engine.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/resource",
		Handler: ProtectedHandler,
	}, rest.WithJwt(c.AccessSecret))
	engine.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/refresh",
		Handler: RefreshHandler(c),
	}, rest.WithJwt(c.RefreshSecret))

	fmt.Println("Now listening...")
	engine.Start()
}

func RefreshHandler(c Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var authReq AuthRequest

		if err := httpx.Parse(r, &authReq); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println(err)
			return
		}

		token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(c.RefreshSecret), nil
			})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("Unauthorized access to this resource")
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("Token is not valid")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("not a valid jwt.MapClaims")
			return
		}

		user, ok := claims[jwtUserField]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("no user info in fresh token")
			return
		}

		userStr, ok := user.(string)
		if !ok || authReq.User != userStr {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("user info not match in query and fresh token")
			return
		}

		respond(w, c, userStr)
	}
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{"Gained access to protected resource"}
	JsonResponse(response, w)
}

func LoginHandler(c Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user UserCredentials

		if err := httpx.Parse(r, &user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Error in request")
			return
		}

		if strings.ToLower(user.Username) != "someone" {
			if user.Password != "p@ssword" {
				w.WriteHeader(http.StatusForbidden)
				fmt.Println("Error logging in")
				fmt.Fprint(w, "Invalid credentials")
				return
			}
		}

		respond(w, c, user.Username)
	}
}

func JsonResponse(response interface{}, w http.ResponseWriter) {
	content, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(content)
}

type ()

func buildTokens(opt TokenOptions) (Tokens, error) {
	var tokens Tokens

	accessToken, err := genToken(opt.AccessSecret, opt.Fields, opt.AccessExpire)
	if err != nil {
		return tokens, err
	}

	refreshToken, err := genToken(opt.RefreshSecret, opt.Fields, opt.RefreshExpire)
	if err != nil {
		return tokens, err
	}

	now := time.Now().Unix()
	tokens.AccessToken = accessToken
	tokens.AccessExpire = now + opt.AccessExpire
	tokens.RefreshAfter = now + opt.RefreshAfter
	tokens.RefreshToken = refreshToken
	tokens.RefreshExpire = now + opt.RefreshExpire

	return tokens, nil
}

func genToken(secretKey string, payloads map[string]interface{}, seconds int64) (string, error) {
	now := time.Now().Unix()
	claims := make(jwt.MapClaims)
	claims["exp"] = now + seconds
	claims["iat"] = now
	for k, v := range payloads {
		claims[k] = v
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secretKey))
}

func respond(w http.ResponseWriter, c Config, user string) {
	tokens, err := buildTokens(TokenOptions{
		AccessSecret:  c.AccessSecret,
		AccessExpire:  c.AccessExpire,
		RefreshSecret: c.RefreshSecret,
		RefreshExpire: c.RefreshExpire,
		RefreshAfter:  c.RefreshAfter,
		Fields: map[string]interface{}{
			jwtUserField: user,
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Println(err)
		return
	}

	httpx.OkJson(w, tokens)
}
