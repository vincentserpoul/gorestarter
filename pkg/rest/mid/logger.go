package mid

import (
	"fmt"
	"net/http"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
)

// ErrRequestContextKey will allow the error to be passed down
const ErrRequestContextKey = ContextKey("error request")

type augmentedResponseWriter struct {
	http.ResponseWriter
	length     int
	httpStatus int
}

// WriteHeader will not only write b to w but also save the http status in the struct
func (w *augmentedResponseWriter) WriteHeader(httpStatus int) {
	w.ResponseWriter.WriteHeader(httpStatus)
	w.httpStatus = httpStatus
}

// Write will not only write b to w but also save the byte length in the struct
func (w *augmentedResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.length = n

	return n, err
}

func newAugmentedResponseWriter(w http.ResponseWriter) *augmentedResponseWriter {
	return &augmentedResponseWriter{ResponseWriter: w}
}

// Logger will return an error if the required params are not there
func Logger(l *logrus.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logFields := logrus.Fields{}

			if reqID := GetRequestID(r.Context()); reqID != ksuid.Nil {
				logFields["request_id"] = reqID
			}

			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}
			logFields["http_scheme"] = scheme
			logFields["http_proto"] = r.Proto
			logFields["http_method"] = r.Method

			logFields["remote_addr"] = r.RemoteAddr
			logFields["user_agent"] = r.UserAgent()

			logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

			startTime := time.Now()

			naw := newAugmentedResponseWriter(w)

			// Write log after the request is finished
			defer func() {
				logFields["process_time"] = time.Since(startTime)
				logFields["http_status"] = naw.httpStatus
				logFields["resp_length"] = naw.length

				// Get response status and size
				if naw.httpStatus == http.StatusNotFound ||
					naw.httpStatus == http.StatusBadRequest {
					l.WithFields(logFields).Warnln()
					return
				}

				// Get request error if there is one
				if err, ok := r.Context().Value(ErrRequestContextKey).(error); ok && err != nil {
					l.WithFields(logFields).Errorln(err)
					return
				}

				l.WithFields(logFields).Infoln()
			}()

			h.ServeHTTP(naw, r)
		})
	}
}
