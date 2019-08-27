package tracer

import (
	"fmt"

	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/reporter/http"
)

// NewTracer creates a new tracer with the necessary dependencies
func NewTracer(serviceName string, reporterURL string) (*zipkin.Tracer, error) {

	reporter := http.NewReporter(fmt.Sprintf("%s/api/v2/spans", reporterURL))

	endpoint := &model.Endpoint{
		ServiceName: serviceName,
	}

	// sampler indicate the range of how many traces are going to be sampled
	sampler, err := zipkin.NewCountingSampler(1)
	if err != nil {
		return nil, err
	}

	t, err := zipkin.NewTracer(
		reporter,
		zipkin.WithSampler(sampler),
		zipkin.WithLocalEndpoint(endpoint),
	)
	if err != nil {
		return nil, err
	}
	return t, nil
}
