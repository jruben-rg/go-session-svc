package adapters

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

var redisServer *miniredis.Miniredis
var cache redisCache
var ctx = context.Background()

func TestShouldStoreDataInRedis(t *testing.T) {
	setup()
	defer teardown()

	sessionKey := "someSessionTestKey"
	sessionValue := `{"someSessionTest":"Value"}`
	err := cache.Set(ctx, sessionKey, sessionValue)
	if err != nil {
		t.Errorf("got error when storing session value in redis %s\n", err)
	}

	val, err := cache.Get(ctx, sessionKey)
	if err != nil {
		t.Errorf("got error when retrieving session value in redis %s\n", err)
	}

	assert.True(t, val == sessionValue, "Expect session to be stored in redis")
}

func TestShouldNotRetrieveDataIfDataDoesntExist(t *testing.T) {
	setup()
	defer teardown()

	sessionKey := "thisSessionKeyShouldNotExist"
	val, err := cache.Get(ctx, sessionKey)

	assert.NotNil(t, err, "Expect redis to return nil for non existing session")
	assert.True(t, val == "", "Expect session data not to be in redis")
}

func TestShouldDeleteData(t *testing.T) {
	setup()
	defer teardown()

	sessionKey := "someDeleteTestKey"
	sessionValue := `{"someDeleteTest":"Value"}`
	err := cache.Set(ctx, sessionKey, sessionValue)
	assert.Nil(t, err, "Expect err is nil when inserting session key")

	_, err = cache.Delete(ctx, sessionKey)
	assert.Nil(t, err, "Expect err is nil when deleting session key")

	val, err := cache.Get(ctx, sessionKey)
	assert.True(t, val == "", "Expect session data to have been deleted from redis")
	assert.NotNil(t, err, "Expect redis to return nil for non existing session")
}

func mockRedis() *miniredis.Miniredis {
	s, err := miniredis.Run()

	if err != nil {
		panic(err)
	}

	return s
}

func setup() {

	redisServer = mockRedis()
	cache = redisCache{
		client: redis.NewClient(&redis.Options{
			Addr: redisServer.Addr(),
		}),
	}

}

func teardown() {
	redisServer.Close()
}
