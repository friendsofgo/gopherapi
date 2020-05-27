package redis

import (
	"context"
	"encoding/json"
	"errors"
	gopherapi "github.com/friendsofgo/gopherapi/pkg"
	"github.com/gomodule/redigo/redis"
	_ "github.com/lib/pq"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_GopherRepository_CreateGopher_RepositoryError(t *testing.T) {
	gopher := buildGopher("123ABC")

	conn := redigomock.NewConn()
	conn.Command("SET", gopher.ID, gopherToJSONString(gopher)).ExpectError(errors.New("something failed"))

	repo := NewRepository(wrapRedisConn(conn))
	err := repo.CreateGopher(context.Background(), &gopher)

	assert.Error(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func Test_GopherRepository_CreateGopher_Success(t *testing.T) {
	gopher := buildGopher("123ABC")

	conn := redigomock.NewConn()
	conn.Command("SET", gopher.ID, gopherToJSONString(gopher)).Expect("OK")

	repo := NewRepository(wrapRedisConn(conn))
	err := repo.CreateGopher(context.Background(), &gopher)

	assert.NoError(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func Test_GopherRepository_FetchGophers_RepositoryError(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("KEYS", "*").ExpectError(errors.New("something failed"))

	repo := NewRepository(wrapRedisConn(conn))
	_, err := repo.FetchGophers(context.Background())

	assert.Error(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func Test_GopherRepository_FetchGophers_NoRows(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("KEYS", "*").Expect([]interface{}{})

	repo := NewRepository(wrapRedisConn(conn))
	gophers, err := repo.FetchGophers(context.Background())

	assert.NoError(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
	assert.Len(t, gophers, 0)
}

func Test_GopherRepository_FetchGophers_RowWithInvalidData(t *testing.T) {
	conn := redigomock.NewConn()
	conn.Command("KEYS", "*").Expect([]interface{}{"123", "456"})
	conn.Command("MGET", "123", "456").Expect([]interface{}{"invalid-data"})

	repo := NewRepository(wrapRedisConn(conn))
	_, err := repo.FetchGophers(context.Background())

	assert.Error(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func Test_GopherRepository_FetchGophers_Succeeded(t *testing.T) {
	gopherA, gopherB := buildGopher("123ABC"), buildGopher("ABC123")
	expectedGophers := []gopherapi.Gopher{gopherA, gopherB}

	conn := redigomock.NewConn()
	conn.Command("KEYS", "*").Expect([]interface{}{gopherA.ID, gopherB.ID})
	conn.Command("MGET", gopherA.ID, gopherB.ID).Expect(
		[]interface{}{gopherToJSONString(gopherA), gopherToJSONString(gopherB)},
	)

	repo := NewRepository(wrapRedisConn(conn))
	gophers, err := repo.FetchGophers(context.Background())

	assert.NoError(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
	assert.Equal(t, expectedGophers, gophers)
}

func Test_GopherRepository_DeleteGopher_RepositoryError(t *testing.T) {
	gopherID := "123ABC"

	conn := redigomock.NewConn()
	conn.Command("DEL", gopherID).ExpectError(errors.New("something failed"))

	repo := NewRepository(wrapRedisConn(conn))
	err := repo.DeleteGopher(context.Background(), gopherID)

	assert.Error(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func Test_GopherRepository_DeleteGopher_Success(t *testing.T) {
	gopherID := "123ABC"

	conn := redigomock.NewConn()
	conn.Command("DEL", gopherID).Expect(1)

	repo := NewRepository(wrapRedisConn(conn))
	err := repo.DeleteGopher(context.Background(), gopherID)

	assert.NoError(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func Test_GopherRepository_UpdateGopher_RepositoryError(t *testing.T) {
	gopher := buildGopher("123ABC")

	conn := redigomock.NewConn()
	conn.Command("SET", gopher.ID, gopherToJSONString(gopher), "XX").ExpectError(errors.New("something failed"))

	repo := NewRepository(wrapRedisConn(conn))
	err := repo.UpdateGopher(context.Background(), gopher.ID, gopher)

	assert.Error(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func Test_GopherRepository_UpdateGopher_NotFound(t *testing.T) {
	gopher := buildGopher("123ABC")

	conn := redigomock.NewConn()
	conn.Command("SET", gopher.ID, gopherToJSONString(gopher), "XX").Expect(nil)

	repo := NewRepository(wrapRedisConn(conn))
	err := repo.UpdateGopher(context.Background(), gopher.ID, gopher)

	assert.Error(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func Test_GopherRepository_UpdateGopher_Success(t *testing.T) {
	gopher := buildGopher("123ABC")

	conn := redigomock.NewConn()
	conn.Command("SET", gopher.ID, gopherToJSONString(gopher), "XX").Expect("OK")

	repo := NewRepository(wrapRedisConn(conn))
	err := repo.UpdateGopher(context.Background(), gopher.ID, gopher)

	assert.NoError(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func Test_GopherRepository_FetchGopherByID_RepositoryError(t *testing.T) {
	gopherID := "123ABC"

	conn := redigomock.NewConn()
	conn.Command("GET", gopherID).ExpectError(errors.New("something failed"))

	repo := NewRepository(wrapRedisConn(conn))
	_, err := repo.FetchGopherByID(context.Background(), gopherID)

	assert.Error(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func Test_GopherRepository_FetchGopherByID_NoRows(t *testing.T) {
	gopherID := "123ABC"

	conn := redigomock.NewConn()
	conn.Command("GET", gopherID).Expect(nil)

	repo := NewRepository(wrapRedisConn(conn))
	_, err := repo.FetchGopherByID(context.Background(), gopherID)

	assert.Error(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func Test_GopherRepository_FetchGopherByID_RowWithInvalidData(t *testing.T) {
	gopherID := "123ABC"

	conn := redigomock.NewConn()
	conn.Command("GET", gopherID).Expect("invalid-data")

	repo := NewRepository(wrapRedisConn(conn))
	_, err := repo.FetchGopherByID(context.Background(), gopherID)

	assert.Error(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
}

func Test_GopherRepository_FetchGopherByID_Succeeded(t *testing.T) {
	gopherID := "123ABC"
	expectedGopher := buildGopher(gopherID)

	conn := redigomock.NewConn()
	conn.Command("GET", gopherID).Expect(gopherToJSONString(expectedGopher))

	repo := NewRepository(wrapRedisConn(conn))
	gopher, err := repo.FetchGopherByID(context.Background(), gopherID)

	assert.NoError(t, err)
	assert.NoError(t, conn.ExpectationsWereMet())
	assert.Equal(t, &expectedGopher, gopher)
}

func buildGopher(ID string) gopherapi.Gopher {
	return gopherapi.Gopher{
		ID:    ID,
		Name:  "The Saviour",
		Image: "https://via.placeholder.com/150.png",
		Age:   8,
	}
}

func gopherToJSONString(gopher gopherapi.Gopher) string {
	bytes, _ := json.Marshal(&gopher)
	return string(bytes)
}

func wrapRedisConn(conn redis.Conn) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return conn, nil },
	}
}
