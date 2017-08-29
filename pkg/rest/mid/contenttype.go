package mid

import "net/http"

// Header will return an error if the required params are not there
func Header(key string, value string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(key, value)
			h.ServeHTTP(w, r) // all params present, proceed
		})
	}
}
