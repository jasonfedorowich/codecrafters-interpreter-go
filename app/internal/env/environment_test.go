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
	es.e = NewEnv()
}

func TestEnvSuite(t *testing.T) {
	suite.Run(t, new(EnvSuite))
}

func (es *EnvSuite) TestPutAndGet() {
	es.Run("test get after put", func() {
		es.e.Put("a", 1.0)
		es.Equal(1.0, es.e.Get("a"))
	})
	es.Run("test put overwrites existing value", func() {
		es.e.Put("a", 1.0)
		es.e.Put("a", "hello")
		es.Equal("hello", es.e.Get("a"))
	})
	es.Run("test get missing key returns nil", func() {
		es.Nil(es.e.Get("missing"))
	})
	es.Run("test put nil value", func() {
		es.e.Put("a", nil)
		es.Nil(es.e.Get("a"))
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
		es.Equal(2.0, es.e.Get("a"))
	})
	es.Run("test assign missing key returns error", func() {
		err := es.e.Assign("a", 2.0)
		es.EqualError(err, "key does not exist: a")
		es.False(es.e.Contains("a"))
	})
}
