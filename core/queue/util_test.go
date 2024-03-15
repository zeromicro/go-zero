package queue

import (
	"errors"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mathx"
)

var (
	proba     = mathx.NewProba()
	failProba = 0.01
)

func init() {
	logx.Disable()
}

func TestGenerateName(t *testing.T) {
	pushers := []Pusher{
		&mockedPusher{name: "first"},
		&mockedPusher{name: "second"},
		&mockedPusher{name: "third"},
	}

	assert.Equal(t, "first,second,third", generateName(pushers))
}

func TestGenerateNameNil(t *testing.T) {
	var pushers []Pusher
	assert.Equal(t, "", generateName(pushers))
}

func calcMean(vals []int) float64 {
	if len(vals) == 0 {
		return 0
	}

	var result float64
	for _, val := range vals {
		result += float64(val)
	}
	return result / float64(len(vals))
}

func calcVariance(mean float64, vals []int) float64 {
	if len(vals) == 0 {
		return 0
	}

	var result float64
	for _, val := range vals {
		result += math.Pow(float64(val)-mean, 2)
	}
	return result / float64(len(vals))
}

type mockedPusher struct {
	name  string
	count int
}

func (p *mockedPusher) Name() string {
	return p.name
}

func (p *mockedPusher) Push(_ string) error {
	if proba.TrueOnProba(failProba) {
		return errors.New("dummy")
	}

	p.count++
	return nil
}
