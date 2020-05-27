package redis

import (
	"context"
	"encoding/json"
	"errors"
	gopherapi "github.com/friendsofgo/gopherapi/pkg"
	"github.com/gomodule/redigo/redis"
	_ "github.com/lib/pq"
)

const (
	onlyIfExists = "XX"
)

type gopherRepository struct {
	pool *redis.Pool
}

// NewRepository instances a Redis implementation of the gopherapi.Repository
func NewRepository(pool *redis.Pool) gopherapi.Repository {
	return gopherRepository{
		pool: pool,
	}
}

// CreateGopher satisfies the gopherapi.Repository interface
func (r gopherRepository) CreateGopher(ctx context.Context, gopher *gopherapi.Gopher) error {
	bytes, err := json.Marshal(gopher)
	if err != nil {
		return err
	}

	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", gopher.ID, string(bytes))
	return err
}

func (r gopherRepository) FetchGophers(ctx context.Context) ([]gopherapi.Gopher, error) {
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}

	keys, err := redis.Strings(conn.Do("KEYS", "*"))
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return []gopherapi.Gopher{}, nil
	}

	args := make([]interface{}, 0, len(keys))
	for _, key := range keys {
		args = append(args, key)
	}

	results, err := redis.Strings(conn.Do("MGET", args...))
	if err != nil {
		return nil, err
	}

	gophers := make([]gopherapi.Gopher, 0, len(results))
	for _, result := range results {
		gopher := gopherapi.Gopher{}

		err := json.Unmarshal([]byte(result), &gopher)
		if err != nil {
			return nil, err
		}

		gophers = append(gophers, gopher)
	}
	return gophers, nil
}

func (r gopherRepository) DeleteGopher(ctx context.Context, ID string) error {
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return err
	}

	_, err = conn.Do("DEL", ID)
	return err
}

func (r gopherRepository) UpdateGopher(ctx context.Context, ID string, gopher gopherapi.Gopher) error {
	bytes, err := json.Marshal(gopher)
	if err != nil {
		return err
	}

	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return err
	}

	result, err := conn.Do("SET", ID, string(bytes), onlyIfExists)
	if result == nil {
		return errors.New("not found")
	}
	return err
}

func (r gopherRepository) FetchGopherByID(ctx context.Context, ID string) (*gopherapi.Gopher, error) {
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}

	result, err := redis.String(conn.Do("GET", ID))
	if err != nil {
		return nil, err
	}

	if result == "" {
		return nil, errors.New("not found")
	}

	gopher := &gopherapi.Gopher{}
	err = json.Unmarshal([]byte(result), gopher)

	return gopher, err
}
