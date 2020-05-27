package redis

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	gopher "github.com/friendsofgo/gopherapi/pkg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GopherRepository_Example(t *testing.T) {
	// GIVEN a miniredis instance and a Redis implementation of result.Repository
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	repo := NewRepository(NewConn(s.Addr()))

	// WHEN two gophers are created
	gopherA, gopherB := buildGopher("123ABC"), buildGopher("ABC123")

	err = repo.CreateGopher(context.Background(), &gopherA)
	assert.NoError(t, err)

	err = repo.CreateGopher(context.Background(), &gopherB)
	assert.NoError(t, err)

	// THEN they can be fetched by ID
	result, err := repo.FetchGopherByID(context.Background(), gopherA.ID)
	assert.NoError(t, err)
	assert.Equal(t, gopherA, *result)

	result, err = repo.FetchGopherByID(context.Background(), gopherB.ID)
	assert.NoError(t, err)
	assert.Equal(t, gopherB, *result)

	// AND they can be fetched in batch
	results, err := repo.FetchGophers(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, []gopher.Gopher{gopherA, gopherB}, results)
}
