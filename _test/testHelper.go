package _test

import (
	"math"
	"os"

	"github.com/gomodule/redigo/redis"
	"github.com/mgironi/operation-fire-quasar/model"
	"github.com/mgironi/operation-fire-quasar/store"
	"github.com/rafaeljusto/redigomock"
)

func AreFloats32Equals(a float32, b float32) bool {
	return AreFloatsEquals(float64(a), float64(b))
}

func AreFloatsEquals(a float64, b float64) bool {
	diff := math.Abs(a - b)
	return diff < model.FLOAT_COMPARISION_TOLERANCE
}

// Cleans the satelite info environment variables
func CleanSatelitesInfoEnvs() {
	//clean envs
	envs := []string{store.SATELITE_KENOBI_ENV, store.SATELITE_SKYWALKER_ENV, store.SATELITE_SATO_ENV}
	for _, key := range envs {
		os.Unsetenv(key)
	}
}

func InitRedisMockConnection() *redigomock.Conn {
	conn := redigomock.NewConn()
	store.GetRedisConnection = func() redis.Conn {
		return conn
	}
	return conn
}
