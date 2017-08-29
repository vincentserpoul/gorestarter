package mid

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestLoggerCompleteness(t *testing.T) {
	expectedNonEmptyFields := []string{"http_scheme",
		"http_proto", "http_method", "remote_addr", "user_agent", "uri",
		"process_time", "http_status", "resp_length"}
	errTest := errors.New("test error")

	tc := []struct {
		name                   string
		handler                http.HandlerFunc
		expectedNonEmptyFields []string
		expectedLogLevel       logrus.Level
		expectedHTTPStatus     int
		withHTTPS              bool
	}{
		{
			name: "classic request log",
			handler: http.HandlerFunc(
				func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			expectedNonEmptyFields: expectedNonEmptyFields,
			expectedLogLevel:       logrus.InfoLevel,
			expectedHTTPStatus:     http.StatusOK,
		},
		{
			name: "request log with error",
			handler: http.HandlerFunc(
				func(w http.ResponseWriter, req *http.Request) {
					*req = *req.WithContext(context.WithValue(req.Context(), ErrRequestContextKey, errTest))
					w.WriteHeader(http.StatusInternalServerError)
				}),
			expectedNonEmptyFields: expectedNonEmptyFields,
			expectedLogLevel:       logrus.ErrorLevel,
			expectedHTTPStatus:     http.StatusInternalServerError,
		},
		{
			name: "404 request log",
			handler: http.HandlerFunc(
				func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}),
			expectedNonEmptyFields: expectedNonEmptyFields,
			expectedLogLevel:       logrus.WarnLevel,
			expectedHTTPStatus:     http.StatusNotFound,
		},
		{
			name: "400 request log",
			handler: http.HandlerFunc(
				func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
				}),
			expectedNonEmptyFields: expectedNonEmptyFields,
			expectedLogLevel:       logrus.WarnLevel,
			expectedHTTPStatus:     http.StatusBadRequest,
		},
		{
			name: "https request log",
			handler: http.HandlerFunc(
				func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			expectedLogLevel:   logrus.InfoLevel,
			withHTTPS:          true,
			expectedHTTPStatus: http.StatusOK,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			logger, hook := test.NewNullLogger()
			midWared := Logger(logger)(tt.handler)
			rr := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", ``, nil)
			r.RemoteAddr = "127.0.0.1"
			r.Header.Set("User-Agent", "test")
			if tt.withHTTPS {
				r.TLS = &tls.ConnectionState{}
			}
			midWared.ServeHTTP(rr, r)

			if len(hook.Entries) > 1 {
				t.Errorf("got %d logs instead of 1", len(hook.Entries))
				return
			}

			if tt.expectedLogLevel != hook.LastEntry().Level {
				t.Errorf("wrong log level, expected %v, got %v", tt.expectedLogLevel, hook.LastEntry().Level)
				return
			}

			if tt.withHTTPS && hook.LastEntry().Data["http_scheme"] != "https" {
				t.Errorf("wrong http scheme detected, expected https, got %s", hook.LastEntry().Data["http_scheme"])
				return
			}

			if hook.LastEntry().Data["http_status"] != tt.expectedHTTPStatus {
				t.Errorf("wrong http status, expected %d, got %d", tt.expectedHTTPStatus, hook.LastEntry().Data["http_status"])
				return
			}

			for _, field := range tt.expectedNonEmptyFields {
				if hook.LastEntry().Data[field] == nil {
					t.Errorf("missing field %s in log", field)
					return
				}
			}
		})
	}
}

func TestLoggerData(t *testing.T) {
	reqIDTest := ksuid.New()
	tc := []struct {
		name                string
		handler             http.HandlerFunc
		requestID           ksuid.KSUID
		RequestErrorMessage string
		expectedLen         int
	}{
		{
			name: "classic request log",
			handler: http.HandlerFunc(
				func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			RequestErrorMessage: "",
		},
		{
			name:      "classic request log with request ID",
			requestID: reqIDTest,
			handler: http.HandlerFunc(
				func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			RequestErrorMessage: "",
		},
		{
			name: "classic request log with some bytes written",
			handler: http.HandlerFunc(
				func(w http.ResponseWriter, req *http.Request) {
					_, _ = w.Write([]byte{0x00})
					w.WriteHeader(http.StatusOK)
				}),
			RequestErrorMessage: "",
			expectedLen:         1,
		},
		{
			name: "request log with error",
			handler: http.HandlerFunc(
				func(w http.ResponseWriter, req *http.Request) {
					*req = *req.WithContext(context.WithValue(req.Context(), ErrRequestContextKey, errors.New("test error big")))
					w.WriteHeader(http.StatusInternalServerError)
				}),
			RequestErrorMessage: "test error big",
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			logger, hook := test.NewNullLogger()
			midWared := Logger(logger)(tt.handler)
			rr := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", ``, nil)
			r.RemoteAddr = "127.0.0.1"
			r.Header.Set("User-Agent", "test")
			if tt.requestID != ksuid.Nil {
				*r = *r.WithContext(context.WithValue(r.Context(), contextKeyRequestID, reqIDTest))
			}
			midWared.ServeHTTP(rr, r)

			if len(hook.Entries) > 1 {
				t.Errorf("got %d logs instead of 1", len(hook.Entries))
				return
			}

			if tt.requestID != ksuid.Nil && hook.LastEntry().Data["request_id"] != tt.requestID {
				t.Errorf("missing requestid, expected %s, got %s", tt.requestID, hook.LastEntry().Data["request_id"])
				return
			}

			if string(hook.LastEntry().Message) != tt.RequestErrorMessage {
				t.Errorf("missing error message `%s` in log, got `%s`", tt.RequestErrorMessage, hook.LastEntry().Message)
				return
			}

			if tt.expectedLen != hook.LastEntry().Data["resp_length"] {
				t.Errorf("expected length %d got %d", tt.expectedLen, hook.LastEntry().Data["resp_length"])
				return
			}

		})
	}
}
