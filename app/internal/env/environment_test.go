package env

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type EnvSuite struct {
	suite.Suite
	e *Env
}

func (es *EnvSuite) SetupSubTest() {
	es.e = NewEnv(nil)
}

func TestEnvSuite(t *testing.T) {
	suite.Run(t, new(EnvSuite))
}

func (es *EnvSuite) TestPutAndGet() {
	es.Run("test get after put", func() {
		es.e.Put("a", 1.0)
		v, ok := es.e.Get("a")
		es.True(ok)
		es.Equal(1.0, v)
	})
	es.Run("test put overwrites existing value", func() {
		es.e.Put("a", 1.0)
		es.e.Put("a", "hello")
		v, ok := es.e.Get("a")
		es.True(ok)
		es.Equal("hello", v)
	})
	es.Run("test get missing key", func() {
		v, ok := es.e.Get("missing")
		es.False(ok)
		es.Nil(v)
	})
	es.Run("test put nil value", func() {
		es.e.Put("a", nil)
		v, ok := es.e.Get("a")
		es.True(ok)
		es.Nil(v)
		es.True(es.e.Contains("a"))
	})
}

func (es *EnvSuite) TestContains() {
	es.Run("test contains after put", func() {
		es.e.Put("a", true)
		es.True(es.e.Contains("a"))
	})
	es.Run("test does not contain missing key", func() {
		es.False(es.e.Contains("a"))
	})
}

func (es *EnvSuite) TestCheckAndGet() {
	es.Run("test existing key", func() {
		es.e.Put("a", 42.0)
		v, ok := es.e.CheckAndGet("a")
		es.True(ok)
		es.Equal(42.0, v)
	})
	es.Run("test existing key with nil value", func() {
		es.e.Put("a", nil)
		v, ok := es.e.CheckAndGet("a")
		es.True(ok)
		es.Nil(v)
	})
	es.Run("test missing key", func() {
		v, ok := es.e.CheckAndGet("a")
		es.False(ok)
		es.Nil(v)
	})
}

func (es *EnvSuite) TestAssign() {
	es.Run("test assign existing key", func() {
		es.e.Put("a", 1.0)
		err := es.e.Assign("a", 2.0)
		es.NoError(err)
		v, _ := es.e.Get("a")
		es.Equal(2.0, v)
	})
	es.Run("test assign missing key returns error", func() {
		err := es.e.Assign("a", 2.0)
		es.EqualError(err, "key does not exist: a")
		es.False(es.e.Contains("a"))
	})
}

func (es *EnvSuite) TestWrappedEnv() {
	es.Run("test get falls back to enclosing env", func() {
		es.e.Put("a", 1.0)
		inner := NewWrappedEnv(es.e)
		v, ok := inner.Get("a")
		es.True(ok)
		es.Equal(1.0, v)
	})
	es.Run("test inner put shadows enclosing env", func() {
		es.e.Put("a", 1.0)
		inner := NewWrappedEnv(es.e)
		inner.Put("a", 2.0)
		v, _ := inner.Get("a")
		es.Equal(2.0, v)
		v, _ = es.e.Get("a")
		es.Equal(1.0, v)
	})
	es.Run("test inner put not visible in enclosing env", func() {
		inner := NewWrappedEnv(es.e)
		inner.Put("a", 1.0)
		_, ok := es.e.Get("a")
		es.False(ok)
	})
	es.Run("test assign updates enclosing env", func() {
		es.e.Put("a", 1.0)
		inner := NewWrappedEnv(es.e)
		err := inner.Assign("a", 2.0)
		es.NoError(err)
		es.False(inner.Contains("a"))
		v, _ := es.e.Get("a")
		es.Equal(2.0, v)
	})
	es.Run("test assign missing key in any env returns error", func() {
		inner := NewWrappedEnv(es.e)
		err := inner.Assign("a", 2.0)
		es.EqualError(err, "key does not exist: a")
	})
	es.Run("test get missing key in any env", func() {
		inner := NewWrappedEnv(NewWrappedEnv(es.e))
		v, ok := inner.Get("a")
		es.False(ok)
		es.Nil(v)
	})
}
