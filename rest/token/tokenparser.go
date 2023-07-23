package token

import (
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
	"github.com/zeromicro/go-zero/core/timex"
)

const claimHistoryResetDuration = time.Hour * 24

type (
	// ParseOption defines the method to customize a TokenParser.
	ParseOption func(parser *TokenParser)

	// A TokenParser is used to parse tokens.
	TokenParser struct {
		resetTime     time.Duration
		resetDuration time.Duration
		history       sync.Map
	}
)

// NewTokenParser returns a TokenParser.
func NewTokenParser(opts ...ParseOption) *TokenParser {
	parser := &TokenParser{
		resetTime:     timex.Now(),
		resetDuration: claimHistoryResetDuration,
	}

	for _, opt := range opts {
		opt(parser)
	}

	return parser
}

// ParseToken parses token from given r, with passed in secret and prevSecret.
func (tp *TokenParser) ParseToken(r *http.Request, secret, prevSecret string) (*jwt.Token, error) {
	var token *jwt.Token
	var err error

	if len(prevSecret) > 0 {
		count := tp.loadCount(secret)
		prevCount := tp.loadCount(prevSecret)

		var first, second string
		if count > prevCount {
			first = secret
			second = prevSecret
		} else {
			first = prevSecret
			second = secret
		}

		token, err = tp.doParseToken(r, first)
		if err != nil {
			token, err = tp.doParseToken(r, second)
			if err != nil {
				return nil, err
			}

			tp.incrementCount(second)
		} else {
			tp.incrementCount(first)
		}
	} else {
		token, err = tp.doParseToken(r, secret)
		if err != nil {
			return nil, err
		}
	}

	return token, nil
}

func (tp *TokenParser) doParseToken(r *http.Request, secret string) (*jwt.Token, error) {
	return request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		}, request.WithParser(newParser()))
}

func (tp *TokenParser) incrementCount(secret string) {
	now := timex.Now()
	if tp.resetTime+tp.resetDuration < now {
		tp.history.Range(func(key, value any) bool {
			tp.history.Delete(key)
			return true
		})
	}

	value, ok := tp.history.Load(secret)
	if ok {
		atomic.AddUint64(value.(*uint64), 1)
	} else {
		var count uint64 = 1
		tp.history.Store(secret, &count)
	}
}

func (tp *TokenParser) loadCount(secret string) uint64 {
	value, ok := tp.history.Load(secret)
	if ok {
		return *value.(*uint64)
	}

	return 0
}

// WithResetDuration returns a func to customize a TokenParser with reset duration.
func WithResetDuration(duration time.Duration) ParseOption {
	return func(parser *TokenParser) {
		parser.resetDuration = duration
	}
}

func newParser() *jwt.Parser {
	return jwt.NewParser(jwt.WithJSONNumber())
}
