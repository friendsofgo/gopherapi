package cockroach

import (
	"context"
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"

	gopher "github.com/friendsofgo/gopherapi/pkg"
	"github.com/openzipkin/zipkin-go"
)

type gopherRepository struct {
	db     *sql.DB
	tracer *zipkin.Tracer
}

// NewRepository creates a crockoach repository with the necessary dependencies
func NewRepository(db *sql.DB, tracer *zipkin.Tracer) gopher.Repository {
	return gopherRepository{db: db, tracer: tracer}
}

func (r gopherRepository) CreateGopher(_ context.Context, g *gopher.Gopher) error {
	sqlStm := `INSERT INTO gophers (id, name, age, image, created_at) 
	VALUES ($1, $2, $3, $4, NOW())`
	_, err := r.db.Exec(sqlStm, g.ID, g.Name, g.Age, g.Image)
	if err != nil {
		return err
	}
	return nil
}

func (r gopherRepository) FetchGophers(ctx context.Context) ([]gopher.Gopher, error) {
	sqlStm := `SELECT id, name, age, image, created_at, updated_at FROM gophers`
	rows, err := r.db.Query(sqlStm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var gophers []gopher.Gopher

	for rows.Next() {
		var g gopher.Gopher
		if err := rows.Scan(&g.ID, &g.Name, &g.Age, &g.Image, &g.CreatedAt, &g.UpdatedAt); err != nil {
			log.Println(err)
			continue
		}
		gophers = append(gophers, g)
	}
	return gophers, nil
}

func (r gopherRepository) DeleteGopher(ctx context.Context, ID string) error {
	return errors.New("method not implemented")
}

func (r gopherRepository) UpdateGopher(ctx context.Context, ID string, g gopher.Gopher) error {
	return errors.New("method not implemented")
}

func (r gopherRepository) FetchGopherByID(ctx context.Context, ID string) (*gopher.Gopher, error) {
	return nil, errors.New("method not implemented")
}
