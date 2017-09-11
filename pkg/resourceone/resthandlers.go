package resourceone

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/vincentserpoul/gorestarter/pkg/rest/renderer"

	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
)

// Router is returning the handler for resourceone rest handler
func Router(db *sqlx.DB) http.Handler {
	r := chi.NewRouter()
	// RESTy routes for resourceone resource
	r.Route("/resourceone", func(r chi.Router) {
		r.Post("/", POSTHandler(db))
		r.Get("/", GETListHandler(db))

		// Subrouters:
		r.Route("/{resourceoneID}", func(r chi.Router) {
			r.Get("/", GETHandler(db))
			r.Put("/", PUTHandler(db))
			r.Delete("/", DELETEHandler(db))
		})
	})
	return r
}

// POSTHandler will handle data from request and returns bytes to be written to response
func POSTHandler(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handler errors after rendering
		var errRender error
		defer func() {
			if errRender != nil {
				log.Fatalf("PUTHandler: render error %v", errRender)
			}
		}()

		e := &Resourceone{}

		errJSON := json.NewDecoder(r.Body).Decode(e)
		if errJSON != nil {
			errRender = render.Render(w, r, renderer.ErrInvalidRequest(errJSON))
			return
		}

		err := e.Create(
			r.Context(),
			db,
		)

		if err != nil {
			errRender = render.Render(w, r, renderer.ErrRender(err))
			return
		}

		w.WriteHeader(http.StatusCreated)
		renderer.ResponseJSONRender(w, r, e)
	}
}

// GETListHandler will handle data from request and returns bytes to be written to response
func GETListHandler(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handler errors after rendering
		var errRender error
		defer func() {
			if errRender != nil {
				log.Fatalf("PUTHandler: render error %v", errRender)
			}
		}()

		es, errS := SelectByTimeUpdated(r.Context(), db, time.Now().Add(-6*time.Hour*24))
		if errS == ErrSQLNotFound {
			errRender = render.Render(w, r, renderer.ErrNotFound)
			return
		}
		if errS != nil {
			errRender = render.Render(w, r, renderer.ErrRender(errS))
			return
		}

		w.WriteHeader(http.StatusOK)
		renderer.ResponseJSONRender(w, r, es)
	}
}

// GETHandler will handle data from request and returns bytes to be written to response
func GETHandler(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handler errors after rendering
		var errRender error
		defer func() {
			if errRender != nil {
				log.Fatalf("PUTHandler: render error %v", errRender)
			}
		}()

		resourceoneIDstr := chi.URLParam(r, "resourceoneID")
		resourceoneID, errConv := strconv.ParseInt(resourceoneIDstr, 10, 64)
		if resourceoneIDstr == "" || errConv != nil {
			errRender = render.Render(w, r, renderer.ErrInvalidRequest(errConv))
			return
		}

		e, errS := SelectByID(r.Context(), db, resourceoneID)
		if errS != nil && errS != ErrSQLNotFound {
			errRender = render.Render(w, r, renderer.ErrRender(errS))
			return
		}
		if e == nil && errS == ErrSQLNotFound {
			errRender = render.Render(w, r, renderer.ErrNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		renderer.ResponseJSONRender(w, r, e)
	}
}

// PUTHandler will handle data from request and update the specified resourceone
func PUTHandler(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handler errors after rendering
		var errRender error
		defer func() {
			if errRender != nil {
				log.Fatalf("PUTHandler: render error %v", errRender)
			}
		}()

		resourceoneIDstr := chi.URLParam(r, "resourceoneID")
		resourceoneID, errConv := strconv.ParseInt(resourceoneIDstr, 10, 64)
		if resourceoneIDstr == "" || errConv != nil {
			errRender = render.Render(w, r, renderer.ErrInvalidRequest(errConv))
			return
		}

		e := &Resourceone{
			ID: resourceoneID,
		}

		errJSON := json.NewDecoder(r.Body).Decode(e)
		if errJSON != nil {
			errRender = render.Render(w, r, renderer.ErrInvalidRequest(errJSON))
			return
		}

		errU := Update(r.Context(), db, resourceoneID, e)
		if errU != nil && errU != ErrSQLNotFound {
			errRender = render.Render(w, r, renderer.ErrRender(errU))
			return
		}
		if errU == ErrSQLNotFound {
			errRender = render.Render(w, r, renderer.ErrNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		renderer.ResponseJSONRender(w, r, e)
	}
}

// DELETEHandler will handle data from request and delete the specified resourceone
func DELETEHandler(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handler errors after rendering
		var errRender error
		defer func() {
			if errRender != nil {
				log.Fatalf("PUTHandler: render error %v", errRender)
			}
		}()

		resourceoneIDstr := chi.URLParam(r, "resourceoneID")
		resourceoneID, errConv := strconv.ParseInt(resourceoneIDstr, 10, 64)
		if resourceoneIDstr == "" || errConv != nil {
			errRender = render.Render(w, r, renderer.ErrInvalidRequest(errConv))
			return
		}

		errD := Delete(
			r.Context(),
			db,
			resourceoneID,
		)
		if errD != nil && errD != ErrSQLNotFound {
			errRender = render.Render(w, r, renderer.ErrRender(errD))
			return
		}
		if errD == ErrSQLNotFound {
			errRender = render.Render(w, r, renderer.ErrNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.WriteHeader(http.StatusNoContent)
	}
}
