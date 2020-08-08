package limit

import (
	"errors"
	"strconv"
	"time"

	"github.com/tal-tech/go-zero/core/stores/redis"
)

const (
	// to be compatible with aliyun redis, we cannot use `local key = KEYS[1]` to reuse the key
	periodScript = `local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local current = redis.call("INCRBY", KEYS[1], 1)
if current == 1 then
    redis.call("expire", KEYS[1], window)
    return 1
elseif current < limit then
    return 1
elseif current == limit then
    return 2
else
    return 0
end`
	zoneDiff = 3600 * 8 // GMT+8 for our services
)

const (
	Unknown = iota
	Allowed
	HitQuota
	OverQuota

	internalOverQuota = 0
	internalAllowed   = 1
	internalHitQuota  = 2
)

var ErrUnknownCode = errors.New("unknown status code")

type (
	LimitOption func(l *PeriodLimit)

	PeriodLimit struct {
		period     int
		quota      int
		limitStore *redis.Redis
		keyPrefix  string
		align      bool
	}
)

func NewPeriodLimit(period, quota int, limitStore *redis.Redis, keyPrefix string,
	opts ...LimitOption) *PeriodLimit {
	limiter := &PeriodLimit{
		period:     period,
		quota:      quota,
		limitStore: limitStore,
		keyPrefix:  keyPrefix,
	}

	for _, opt := range opts {
		opt(limiter)
	}

	return limiter
}

func (h *PeriodLimit) Take(key string) (int, error) {
	resp, err := h.limitStore.Eval(periodScript, []string{h.keyPrefix + key}, []string{
		strconv.Itoa(h.quota),
		strconv.Itoa(h.calcExpireSeconds()),
	})
	if err != nil {
		return Unknown, err
	}

	code, ok := resp.(int64)
	if !ok {
		return Unknown, ErrUnknownCode
	}

	switch code {
	case internalOverQuota:
		return OverQuota, nil
	case internalAllowed:
		return Allowed, nil
	case internalHitQuota:
		return HitQuota, nil
	default:
		return Unknown, ErrUnknownCode
	}
}

func (h *PeriodLimit) calcExpireSeconds() int {
	if h.align {
		unix := time.Now().Unix() + zoneDiff
		return h.period - int(unix%int64(h.period))
	} else {
		return h.period
	}
}

func Align() LimitOption {
	return func(l *PeriodLimit) {
		l.align = true
	}
}
