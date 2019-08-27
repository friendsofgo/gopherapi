package inmem

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/openzipkin/zipkin-go"

	gopher "github.com/friendsofgo/gopherapi/pkg"
)

type gopherRepository struct {
	mtx    sync.RWMutex
	tracer *zipkin.Tracer

	gophers map[string]gopher.Gopher
}

// NewRepository creates a inmem repository with the necessary dependencies
func NewRepository(gophers map[string]gopher.Gopher, tracer *zipkin.Tracer) gopher.Repository {
	if gophers == nil {
		gophers = make(map[string]gopher.Gopher)
	}

	return &gopherRepository{
		gophers: gophers,
		tracer:  tracer,
	}
}

func (r *gopherRepository) CreateGopher(ctx context.Context, g *gopher.Gopher) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if err := r.checkIfExists(ctx, g.ID); err != nil {
		return err
	}
	r.gophers[g.ID] = *g
	return nil
}

func (r *gopherRepository) FetchGophers(ctx context.Context) ([]gopher.Gopher, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	values := make([]gopher.Gopher, 0, len(r.gophers))
	for _, value := range r.gophers {
		values = append(values, value)
	}
	return values, nil
}

func (r *gopherRepository) DeleteGopher(ctx context.Context, ID string) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	delete(r.gophers, ID)

	return nil
}

func (r *gopherRepository) UpdateGopher(ctx context.Context, ID string, g gopher.Gopher) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.gophers[ID] = g
	return nil
}

func (r *gopherRepository) FetchGopherByID(ctx context.Context, ID string) (*gopher.Gopher, error) {
	span, _ := r.tracer.StartSpanFromContext(ctx, "FetchGopherByID")
	span.Tag("Repository", "in memory")
	span.Annotate(time.Now(), "Transaction Start")

	defer func() {
		span.Annotate(time.Now(), "Transaction End")
		span.Finish()
	}()

	r.mtx.Lock()
	defer r.mtx.Unlock()

	for _, v := range r.gophers {
		if v.ID == ID {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("Error has ocurred while finding gopher %s", ID)
}

func (r *gopherRepository) checkIfExists(ctx context.Context, ID string) error {
	for _, v := range r.gophers {
		if v.ID == ID {
			return fmt.Errorf("The gopher %s is already exist", ID)
		}
	}

	return nil
}
