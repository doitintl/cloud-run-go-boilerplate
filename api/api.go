//nolint:contextcheck
// Package api implements the interface of the RESTful interface of the
// application. The API methods also contain the core application logic.
package api

import (
	"context"
	"net/http"

	"github.com/doitintl/cloud-run-go-boilerplate/service"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

// API is the central element of the GetCRE server. It is instanciated
// once and kept for the entire server lifecycle. All service clients
// are stored here.
type API struct {
	// Mux is the *http.ServerMux instance to pass to
	// http.ListenAndServe when running this API.
	Mux *http.ServeMux
	log *logrus.Entry
}

// New is the recommended way to construct a new API instances. It
// relies on configuration values in viper.
func New(
	ctx context.Context,
	log *logrus.Entry,
) (*API, error) {
	api := API{
		Mux: nil,
		log: log,
	}

	api.setupMux()

	return &api, nil
}

func (a *API) setupMux() {
	a.Mux = http.NewServeMux()
	a.Mux.HandleFunc("/", a.traceMiddleware(a.indexHandler))
}

func (a *API) indexHandler(writer http.ResponseWriter, req *http.Request) {
	log := a.loggerFromContext(req.Context())

	log.Info("index handler running ...")

	result, err := service.HeavyComputation()
	if err != nil {
		log.WithError(err).Error("could not complete heavy calculation")
	}

	log.WithField("heavy computatio result", result)

	_, _ = writer.Write([]byte(result))
}

func (a *API) loggerFromContext(ctx context.Context) *logrus.Entry {
	return a.log.WithField("span_context", trace.SpanContextFromContext(ctx))
}
