package renderer

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/render"
)

// ResponseJSONRender will simply JSON the response and add the specific headers
func ResponseJSONRender(w http.ResponseWriter, r *http.Request, e interface{}) {
	// Render JSON
	eJSON, errJSON := json.Marshal(e)
	if errJSON != nil {
		errRend := render.Render(w, r, ErrRender(errJSON))
		if errRend != nil {
			log.Errorf("ResponseJSONRender %s", errRend)
		}
		return
	}

	_, errW := fmt.Fprintf(w, string(eJSON))
	if errW != nil {
		errRend := render.Render(w, r, ErrRender(errW))
		if errRend != nil {
			log.Errorf("ResponseJSONRender %s", errRend)
		}
		return
	}
}
