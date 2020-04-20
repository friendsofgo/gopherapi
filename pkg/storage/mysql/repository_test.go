package mysql

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	gopherapi "github.com/friendsofgo/gopherapi/pkg"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_GopherRepository_CreateGopher_RepositoryError(t *testing.T) {
	gopher := buildGopher()

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectExec(
		"INSERT INTO gophers (id, name, image, age, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)").
		WithArgs(gopher.ID, gopher.Name, gopher.Image, gopher.Age, gopher.CreatedAt, gopher.UpdatedAt).
		WillReturnError(errors.New("database failed"))

	repo := NewRepository("gophers", db)
	err = repo.CreateGopher(context.Background(), &gopher)

	assert.Error(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func Test_GopherRepository_CreateGopher_Success(t *testing.T) {
	gopher := buildGopher()

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectExec(
		"INSERT INTO gophers (id, name, image, age, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)").
		WithArgs(gopher.ID, gopher.Name, gopher.Image, gopher.Age, gopher.CreatedAt, gopher.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewRepository("gophers", db)
	err = repo.CreateGopher(context.Background(), &gopher)

	assert.NoError(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func Test_GopherRepository_FetchGophers_RepositoryError(t *testing.T) {
	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectQuery(
		"SELECT gophers.id, gophers.name, gophers.image, gophers.age, gophers.created_at, gophers.updated_at FROM gophers").
		WillReturnError(errors.New("something-failed"))

	repo := NewRepository("gophers", db)
	_, err = repo.FetchGophers(context.Background())

	assert.Error(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func Test_GopherRepository_FetchGophers_NoRows(t *testing.T) {
	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectQuery(
		"SELECT gophers.id, gophers.name, gophers.image, gophers.age, gophers.created_at, gophers.updated_at FROM gophers").
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "name", "image", "age", "created_at", "updated_at"}),
		)

	repo := NewRepository("gophers", db)
	gophers, err := repo.FetchGophers(context.Background())

	assert.NoError(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())

	assert.Len(t, gophers, 0)
}

func Test_GopherRepository_FetchGophers_RowWithInvalidData(t *testing.T) {
	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectQuery(
		"SELECT gophers.id, gophers.name, gophers.image, gophers.age, gophers.created_at, gophers.updated_at FROM gophers").
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "name", "image", "age", "created_at", "updated_at"}).
			AddRow(nil, nil, nil, nil, nil, nil), // This is a row failure as the data type is wrong
		)

	repo := NewRepository("gophers", db)
	_, err = repo.FetchGophers(context.Background())

	assert.Error(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func Test_GopherRepository_FetchGophers_Succeeded(t *testing.T) {
	expectedGophers := []gopherapi.Gopher{
		buildGopher(),
		buildGopher(),
	}

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectQuery(
		"SELECT gophers.id, gophers.name, gophers.image, gophers.age, gophers.created_at, gophers.updated_at FROM gophers").
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "name", "image", "age", "created_at", "updated_at"}).
			AddRow(expectedGophers[0].ID, expectedGophers[0].Name, expectedGophers[0].Image, expectedGophers[0].Age, expectedGophers[0].CreatedAt, expectedGophers[0].UpdatedAt).
			AddRow(expectedGophers[1].ID, expectedGophers[1].Name, expectedGophers[1].Image, expectedGophers[1].Age, expectedGophers[1].CreatedAt, expectedGophers[1].UpdatedAt),
		)

	repo := NewRepository("gophers", db)
	gophers, err := repo.FetchGophers(context.Background())

	assert.NoError(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.Equal(t, expectedGophers, gophers)
}

func Test_GopherRepository_DeleteGopher_RepositoryError(t *testing.T) {
	gopherID := "123ABC"

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectExec(
		"DELETE FROM gophers WHERE id = ?").
		WithArgs(gopherID).
		WillReturnError(errors.New("database failed"))

	repo := NewRepository("gophers", db)
	err = repo.DeleteGopher(context.Background(), gopherID)

	assert.Error(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func Test_GopherRepository_DeleteGopher_Success(t *testing.T) {
	gopherID := "123ABC"

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectExec(
		"DELETE FROM gophers WHERE id = ?").
		WithArgs(gopherID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewRepository("gophers", db)
	err = repo.DeleteGopher(context.Background(), gopherID)

	assert.NoError(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func Test_GopherRepository_UpdateGopher_RepositoryError(t *testing.T) {
	gopher := buildGopher()

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectExec(
		"UPDATE gophers SET id = ?, name = ?, image = ?, age = ?, created_at = ?, updated_at = ? WHERE id = ?").
		WithArgs(gopher.ID, gopher.Name, gopher.Image, gopher.Age, gopher.CreatedAt, gopher.UpdatedAt, gopher.ID).
		WillReturnError(errors.New("database failed"))

	repo := NewRepository("gophers", db)
	err = repo.UpdateGopher(context.Background(), gopher.ID, gopher)

	assert.Error(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func Test_GopherRepository_UpdateGopher_NotFound(t *testing.T) {
	gopher := buildGopher()

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectExec(
		"UPDATE gophers SET id = ?, name = ?, image = ?, age = ?, created_at = ?, updated_at = ? WHERE id = ?").
		WithArgs(gopher.ID, gopher.Name, gopher.Image, gopher.Age, gopher.CreatedAt, gopher.UpdatedAt, gopher.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	repo := NewRepository("gophers", db)
	err = repo.UpdateGopher(context.Background(), gopher.ID, gopher)

	assert.Error(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func Test_GopherRepository_UpdateGopher_Success(t *testing.T) {
	gopher := buildGopher()

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectExec(
		"UPDATE gophers SET id = ?, name = ?, image = ?, age = ?, created_at = ?, updated_at = ? WHERE id = ?").
		WithArgs(gopher.ID, gopher.Name, gopher.Image, gopher.Age, gopher.CreatedAt, gopher.UpdatedAt, gopher.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewRepository("gophers", db)
	err = repo.UpdateGopher(context.Background(), gopher.ID, gopher)

	assert.NoError(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func Test_GopherRepository_FetchGopherByID_RepositoryError(t *testing.T) {
	gopherID := "123ABC"

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectQuery(
		"SELECT gophers.id, gophers.name, gophers.image, gophers.age, gophers.created_at, gophers.updated_at FROM gophers WHERE id = ?").
		WithArgs(gopherID).
		WillReturnError(errors.New("something-failed"))

	repo := NewRepository("gophers", db)
	_, err = repo.FetchGopherByID(context.Background(), gopherID)

	assert.Error(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func Test_GopherRepository_FetchGopherByID_NoRows(t *testing.T) {
	gopherID := "123ABC"

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectQuery(
		"SELECT gophers.id, gophers.name, gophers.image, gophers.age, gophers.created_at, gophers.updated_at FROM gophers WHERE id = ?").
		WithArgs(gopherID).
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "name", "image", "age", "created_at", "updated_at"}),
		)

	repo := NewRepository("gophers", db)
	_, err = repo.FetchGopherByID(context.Background(), gopherID)

	assert.Error(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func Test_GopherRepository_FetchGopherByID_RowWithInvalidData(t *testing.T) {
	gopherID := "123ABC"

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectQuery(
		"SELECT gophers.id, gophers.name, gophers.image, gophers.age, gophers.created_at, gophers.updated_at FROM gophers WHERE id = ?").
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "name", "image", "age", "created_at", "updated_at"}).
			AddRow(nil, nil, nil, nil, nil, nil), // This is a row failure as the data type is wrong
		)

	repo := NewRepository("gophers", db)
	_, err = repo.FetchGopherByID(context.Background(), gopherID)

	assert.Error(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func Test_GopherRepository_FetchGopherByID_Succeeded(t *testing.T) {
	expectedGopher := buildGopher()

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		assert.NoError(t, err)
	}

	sqlMock.ExpectQuery(
		"SELECT gophers.id, gophers.name, gophers.image, gophers.age, gophers.created_at, gophers.updated_at FROM gophers WHERE id = ?",
	).
		WithArgs(expectedGopher.ID).
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "name", "image", "age", "created_at", "updated_at"}).
			AddRow(expectedGopher.ID, expectedGopher.Name, expectedGopher.Image, expectedGopher.Age, expectedGopher.CreatedAt, expectedGopher.UpdatedAt),
		)

	repo := NewRepository("gophers", db)
	gopher, err := repo.FetchGopherByID(context.Background(), expectedGopher.ID)

	assert.NoError(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.Equal(t, &expectedGopher, gopher)
}

func buildGopher() gopherapi.Gopher {
	now := time.Now()
	return gopherapi.Gopher{
		ID:        "123ABC",
		Name:      "The Saviour",
		Image:     "https://via.placeholder.com/150.png",
		Age:       8,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}
