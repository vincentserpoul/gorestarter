package resourceone

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi"
)

func TestPOSTHandler(t *testing.T) {

	tests := []struct {
		name         string
		requestBody  string
		wantedStatus int
		wantedLabel  string
	}{
		{
			name:         "Working POST",
			requestBody:  `{"label": "test"}`,
			wantedStatus: http.StatusCreated,
			wantedLabel:  `test`,
		},
		{
			name:         "Non Working POST",
			requestBody:  `"label": "test"}`,
			wantedStatus: http.StatusBadRequest,
			wantedLabel:  ``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			request, errR := http.NewRequest("POST", `http://dummy/resourceone`,
				bytes.NewBufferString(tt.requestBody))
			if errR != nil {
				t.Fatalf("request creation failed %v", errR)
			}
			POSTHandler(pool)(rr, request)
			resp := rr.Result()
			defer resp.Body.Close()

			if status := resp.StatusCode; status != tt.wantedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantedStatus)
			}

			if resp.StatusCode == http.StatusOK {
				e := &Resourceone{}
				_ = json.NewDecoder(resp.Body).Decode(e)
				if e.Label != tt.wantedLabel || e.ID == 0 {
					t.Errorf("POSTHandler rendered %d as ID and %s as Label instead of %s",
						e.ID, e.Label, tt.wantedLabel)
					return
				}
			}
		})
	}
}

var testResourceoneIDsHandler []int64

func BenchmarkPOSTHandler(b *testing.B) {
	jsonRequestOK, errR := http.NewRequest("POST", `http://dummy/resourceone`,
		bytes.NewBufferString(`{"label": "test"}`))
	if errR != nil {
		b.Fatal(errR)
	}

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		POSTHandler(pool)(rr, jsonRequestOK)
		res := rr.Result()
		defer res.Body.Close()

		e := &Resourceone{}
		_ = json.NewDecoder(rr.Result().Body).Decode(e)
		testResourceoneIDsHandler = append(testResourceoneIDsHandler, e.ID)
	}
}

func TestGETListHandler(t *testing.T) {
	// Preinsert a  list of resourceone
	for i := 0; i < 3; i++ {
		ec := &Resourceone{Label: `test`}
		_ = ec.Create(context.Background(), pool)
	}

	tests := []struct {
		name         string
		wantedStatus int
	}{
		{
			name:         "Working GETList",
			wantedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, _ := http.NewRequest("GET", ``, nil)
			rr := httptest.NewRecorder()

			GETListHandler(pool)(rr, request)
			res := rr.Result()
			defer res.Body.Close()

			if status := rr.Code; status != tt.wantedStatus {
				t.Errorf("GETListHandler returned wrong status code: got %v want %v",
					status, tt.wantedStatus)
			}

			var e []*Resourceone
			_ = json.NewDecoder(rr.Body).Decode(&e)
			if len(e) < 3 {
				t.Errorf("GETListHandler only rendered %d entities", len(e))
				return
			}
		})
	}
}

func BenchmarkGETListHandler(b *testing.B) {
	jsonRequestOK, _ := http.NewRequest("GET", `http://dummy/resourceone`, nil)

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		GETListHandler(pool)(rr, jsonRequestOK)
		res := rr.Result()
		defer res.Body.Close()
	}
}

func TestGETHandler(t *testing.T) {
	// Pre-insert a resourceone
	ec := &Resourceone{Label: `test`}
	_ = ec.Create(context.Background(), pool)

	tests := []struct {
		name         string
		resourceID   string
		wantedStatus int
		wantedLabel  string
	}{
		{
			name:         "Working GET specific ID",
			resourceID:   strconv.FormatInt(ec.ID, 10),
			wantedStatus: http.StatusOK,
			wantedLabel:  `test`,
		},
		{
			name:         "Non Working GET not found",
			resourceID:   `99999`,
			wantedStatus: http.StatusNotFound,
			wantedLabel:  ``,
		},
		{
			name:         "Non Working GET bad request",
			resourceID:   `sdr`,
			wantedStatus: http.StatusBadRequest,
			wantedLabel:  ``,
		},
	}

	request, _ := http.NewRequest("GET", ``, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			ctx := getTestContextWithResourceID(tt.resourceID)
			GETHandler(pool)(rr, request.WithContext(ctx))
			res := rr.Result()
			defer res.Body.Close()

			if status := rr.Code; status != tt.wantedStatus {
				t.Errorf("GETHandler returned wrong status code: got %v want %v",
					status, tt.wantedStatus)
			}

			if tt.wantedStatus == http.StatusOK {
				e := Resourceone{}
				errJSON := json.NewDecoder(rr.Body).Decode(&e)
				if errJSON != nil {
					t.Fatal(errJSON)
					return
				}
				if e.Label != tt.wantedLabel || e.ID == 0 {
					t.Errorf("GETHandler rendered %d as ID and %s as Label instead of %s",
						e.ID, e.Label, tt.wantedLabel)
					return
				}
			}
		})
	}
}

func BenchmarkGETHandler(b *testing.B) {
	request, _ := http.NewRequest("GET", ``, nil)
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		ctx := getTestContextWithResourceID(strconv.FormatInt(testResourceoneIDsHandler[b.N%len(testResourceoneIDsHandler)], 10))
		GETHandler(pool)(rr, request.WithContext(ctx))
		res := rr.Result()
		defer res.Body.Close()
	}
}

func TestPUTHandler(t *testing.T) {
	// Pre-insert a resourceone
	ec := &Resourceone{Label: `test`}
	_ = ec.Create(context.Background(), pool)

	tests := []struct {
		name           string
		resourceID     string
		requestURLBody string
		wantedStatus   int
		wantedLabel    string
	}{
		{
			name:           "Working PUT specific ID",
			resourceID:     strconv.FormatInt(ec.ID, 10),
			requestURLBody: `{"label": "testUpdate"}`,
			wantedStatus:   http.StatusOK,
			wantedLabel:    `testUpdate`,
		},
		{
			name:           "Non Working PUT not found",
			resourceID:     `99999`,
			requestURLBody: `{"label": "testUpdate"}`,
			wantedStatus:   http.StatusNotFound,
			wantedLabel:    ``,
		},
		{
			name:           "Non Working PUT bad request",
			resourceID:     `sdr`,
			requestURLBody: ``,
			wantedStatus:   http.StatusBadRequest,
			wantedLabel:    ``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, _ := http.NewRequest("PUT", ``, bytes.NewBufferString(tt.requestURLBody))
			rr := httptest.NewRecorder()
			ctx := getTestContextWithResourceID(tt.resourceID)
			PUTHandler(pool)(rr, request.WithContext(ctx))
			res := rr.Result()
			defer res.Body.Close()

			if status := rr.Code; status != tt.wantedStatus {
				t.Errorf("PUTHandler returned wrong status code: got %v want %v",
					status, tt.wantedStatus)
			}

			if tt.wantedStatus == http.StatusOK {
				ep := &Resourceone{}
				errJSON := json.NewDecoder(rr.Body).Decode(ep)
				if errJSON != nil {
					fmt.Println(errJSON)
					return
				}
				if ep.Label != tt.wantedLabel || ep.ID == 0 {
					t.Errorf("PUTHandler rendered %d as ID and %s as Label instead of %s",
						ep.ID, ep.Label, tt.wantedLabel)
					return
				}
			}
		})
	}
}

func BenchmarkPUTHandler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		request, _ := http.NewRequest("PUT", ``, bytes.NewBufferString(`{"label": "testUpdate"}`))
		rr := httptest.NewRecorder()
		ctx := getTestContextWithResourceID(strconv.FormatInt(testResourceoneIDsHandler[b.N%len(testResourceoneIDsHandler)], 10))
		PUTHandler(pool)(rr, request.WithContext(ctx))
		res := rr.Result()
		defer res.Body.Close()
	}
}

func TestDELETEHandler(t *testing.T) {
	// Preinsert a resourceone
	ec := &Resourceone{Label: `test`}
	_ = ec.Create(context.Background(), pool)

	tests := []struct {
		name         string
		resourceID   string
		wantedStatus int
		wantedLabel  string
	}{
		{
			name:         "Working DELETE specific ID",
			resourceID:   strconv.FormatInt(ec.ID, 10),
			wantedStatus: http.StatusNoContent,
			wantedLabel:  `testUpdate`,
		},
		{
			name:         "Non Working DELETE not found",
			resourceID:   `999999`,
			wantedStatus: http.StatusNotFound,
			wantedLabel:  ``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			request, _ := http.NewRequest("DELETE", ``, nil)
			ctx := getTestContextWithResourceID(tt.resourceID)

			DELETEHandler(pool)(rr, request.WithContext(ctx))
			res := rr.Result()
			defer res.Body.Close()

			if status := rr.Code; status != tt.wantedStatus {
				t.Errorf("DELETEHandler returned wrong status code: got %v want %v",
					status, tt.wantedStatus)
			}
		})
	}
}

func BenchmarkDELETEHandler(b *testing.B) {
	for i := 0; i < b.N && i < len(testResourceoneIDsHandler); i++ {
		request, _ := http.NewRequest("DELETE", ``, nil)
		rr := httptest.NewRecorder()
		ctx := getTestContextWithResourceID(strconv.FormatInt(testResourceoneIDsHandler[i], 10))

		DELETEHandler(pool)(rr, request.WithContext(ctx))
		res := rr.Result()
		defer res.Body.Close()
	}
}

func getTestContextWithResourceID(resourceID string) context.Context {
	// Set the URL param
	ctxR := chi.NewRouteContext()
	ctxR.URLParams.Add("resourceoneID", resourceID)
	ctx := context.WithValue(context.Background(), chi.RouteCtxKey, ctxR)

	return ctx
}

// func TestRouter(t *testing.T) {
// 	db := &sqlx.DB{}
// 	router := Router(db)
// 	tests := []struct {
// 		name   string
// 		method string
// 		URL    string
// 		want   http.HandlerFunc
// 	}{
// 		{
// 			name:   `GET`,
// 			method: `GET`,
// 			URL:    `/resourceone`,
// 			want:   GETHandler(db),
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := Router(tt.args.db); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Router() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
