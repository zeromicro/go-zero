package mongo

import (
	"testing"

	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/breaker"
)

func Test_rejectedQuery_All(t *testing.T) {
	assert.Equal(t, breaker.ErrServiceUnavailable, new(rejectedQuery).All(nil))
}

func Test_rejectedQuery_Apply(t *testing.T) {
	info, err := new(rejectedQuery).Apply(mgo.Change{}, nil)
	assert.Equal(t, breaker.ErrServiceUnavailable, err)
	assert.Nil(t, info)
}

func Test_rejectedQuery_Batch(t *testing.T) {
	var q rejectedQuery
	assert.Equal(t, q, q.Batch(1))
}

func Test_rejectedQuery_Collation(t *testing.T) {
	var q rejectedQuery
	assert.Equal(t, q, q.Collation(nil))
}

func Test_rejectedQuery_Comment(t *testing.T) {
	var q rejectedQuery
	assert.Equal(t, q, q.Comment(""))
}

func Test_rejectedQuery_Count(t *testing.T) {
	n, err := new(rejectedQuery).Count()
	assert.Equal(t, breaker.ErrServiceUnavailable, err)
	assert.Equal(t, 0, n)
}

func Test_rejectedQuery_Distinct(t *testing.T) {
	assert.Equal(t, breaker.ErrServiceUnavailable, new(rejectedQuery).Distinct("", nil))
}

func Test_rejectedQuery_Explain(t *testing.T) {
	assert.Equal(t, breaker.ErrServiceUnavailable, new(rejectedQuery).Explain(nil))
}

func Test_rejectedQuery_For(t *testing.T) {
	assert.Equal(t, breaker.ErrServiceUnavailable, new(rejectedQuery).For(nil, nil))
}

func Test_rejectedQuery_Hint(t *testing.T) {
	var q rejectedQuery
	assert.Equal(t, q, q.Hint())
}

func Test_rejectedQuery_Iter(t *testing.T) {
	assert.EqualValues(t, rejectedIter{}, new(rejectedQuery).Iter())
}

func Test_rejectedQuery_Limit(t *testing.T) {
	var q rejectedQuery
	assert.Equal(t, q, q.Limit(1))
}

func Test_rejectedQuery_LogReplay(t *testing.T) {
	var q rejectedQuery
	assert.Equal(t, q, q.LogReplay())
}

func Test_rejectedQuery_MapReduce(t *testing.T) {
	info, err := new(rejectedQuery).MapReduce(nil, nil)
	assert.Equal(t, breaker.ErrServiceUnavailable, err)
	assert.Nil(t, info)
}

func Test_rejectedQuery_One(t *testing.T) {
	assert.Equal(t, breaker.ErrServiceUnavailable, new(rejectedQuery).One(nil))
}

func Test_rejectedQuery_Prefetch(t *testing.T) {
	var q rejectedQuery
	assert.Equal(t, q, q.Prefetch(1))
}

func Test_rejectedQuery_Select(t *testing.T) {
	var q rejectedQuery
	assert.Equal(t, q, q.Select(nil))
}

func Test_rejectedQuery_SetMaxScan(t *testing.T) {
	var q rejectedQuery
	assert.Equal(t, q, q.SetMaxScan(0))
}

func Test_rejectedQuery_SetMaxTime(t *testing.T) {
	var q rejectedQuery
	assert.Equal(t, q, q.SetMaxTime(0))
}

func Test_rejectedQuery_Skip(t *testing.T) {
	var q rejectedQuery
	assert.Equal(t, q, q.Skip(0))
}

func Test_rejectedQuery_Snapshot(t *testing.T) {
	var q rejectedQuery
	assert.Equal(t, q, q.Snapshot())
}

func Test_rejectedQuery_Sort(t *testing.T) {
	var q rejectedQuery
	assert.Equal(t, q, q.Sort())
}

func Test_rejectedQuery_Tail(t *testing.T) {
	assert.EqualValues(t, rejectedIter{}, new(rejectedQuery).Tail(0))
}
