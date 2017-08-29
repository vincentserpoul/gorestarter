package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/vincentserpoul/gorestarter/pkg/resourceone"
	"github.com/vincentserpoul/gorestarter/pkg/rest/mid"
)

// New instanciate the http server and return a channel
func New(httpPort int, db *sqlx.DB, logger *logrus.Logger) *http.Server {

	r := chi.NewRouter()
	r.Use(mid.RequestID())
	r.Use(mid.Header("Content-Type", "application/json"))
	r.Use(middleware.RealIP)
	r.Use(mid.Logger(logger))

	r.Mount("/v1", resourceone.Router(db))

	// Resourceone related things
	ddl := &resourceone.DDL{}
	err := ddl.MigrateUp(context.Background(), db)
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: r,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	return srv
}
