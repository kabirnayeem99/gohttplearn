package server

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kabirnayeem99/gohttplearn/internal/model"
)

func TestHandleHello(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handleHello(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status mismatch: got %v want %v",
			res.StatusCode, http.StatusOK)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	want := `{"msg":"hello world"}` + "\n"

	if string(body) != want {
		t.Fatalf("body mismatch: got %q want %q",
			string(body), want)
	}
}

func TestHandleGoodByte(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/goodbye", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handleGoodBye(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status mismatch: got %v, want %v", res.StatusCode, http.StatusOK)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		t.Fatal(err)
	}

	want := `{"msg":"good bye"}` + "\n"

	if string(body) != want {
		t.Fatalf("body mismatch: got %v, want %v", string(body), want)
	}
}

func TestHelloParameterized(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "with user returns greeting",
			url:        "/hello?user=nabil",
			wantStatus: http.StatusOK,
			wantMsg:    "hello, nabil",
		},
		{
			name:       "no user returns default greeting",
			url:        "/hello",
			wantStatus: http.StatusOK,
			wantMsg:    "hello",
		},
		{
			name:       "empty user param returns 400",
			url:        "/hello?user=",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rr := httptest.NewRecorder()

			handleHelloParameterized(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("status: got %d want %d", res.StatusCode, tt.wantStatus)
			}

			if ct := rr.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
				t.Fatalf("content-type: got %q want %q", ct, "application/json")
			}

			if tt.wantStatus != http.StatusOK {
				return
			}

			var body struct {
				Msg string `json:"msg"`
			}
			if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}

			if body.Msg != tt.wantMsg {
				t.Fatalf("msg: got %q want %q", body.Msg, tt.wantMsg)
			}
		})
	}
}

func TestHandlerGreetingsUserHello(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		username   string
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "with user returns greeting",
			url:        "/greetings/nabil/hello",
			username:   "nabil",
			wantStatus: http.StatusOK,
			wantMsg:    "hello, nabil",
		},
		{
			name:       "no user returns not found",
			url:        "/greetings/hello",
			username:   "",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid user returns bad request",
			url:        "/greetings/nabil-bhul/hello",
			username:   "nabil-bhul",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name,
			func(t *testing.T) {
				req, err := http.NewRequest(http.MethodGet, tt.url, nil)
				if err != nil {
					t.Fatal(err)
				}

				req.SetPathValue("user", tt.username)

				rr := httptest.NewRecorder()

				handlerGreetingsUserHello(rr, req)

				res := rr.Result()

				defer res.Body.Close()

				if res.StatusCode != tt.wantStatus {
					t.Fatalf("status: got %d want %d", res.StatusCode, tt.wantStatus)
				}

				if ct := rr.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
					t.Fatalf("content-type: got %q want %q", ct, "application/json")
				}

				if tt.wantStatus != http.StatusOK {
					return
				}

				var body struct {
					Msg string `json:"msg"`
				}

				if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}

				if body.Msg != tt.wantMsg {
					t.Fatalf("msg: got %q want %q", body.Msg, tt.wantMsg)
				}

			},
		)
	}
}

func TestHandlerGreetingsHello(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		username   string
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "with user returns greeting",
			url:        "/greetings/hello",
			username:   "nabil",
			wantStatus: http.StatusOK,
			wantMsg:    "hello, nabil",
		},
		{
			name:       "no user returns not found",
			url:        "/greetings/hello",
			username:   "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid user returns bad request",
			url:        "/greetings/hello",
			username:   "nabil-bhul",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name,
			func(t *testing.T) {
				req, err := http.NewRequest(http.MethodGet, tt.url, nil)
				if err != nil {
					t.Fatal(err)
				}

				req.Header.Set("user", tt.username)

				rr := httptest.NewRecorder()

				handleGreetingsHello(rr, req)

				res := rr.Result()

				defer res.Body.Close()

				if res.StatusCode != tt.wantStatus {
					t.Fatalf("status: got %d want %d", res.StatusCode, tt.wantStatus)
				}

				if ct := rr.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
					t.Fatalf("content-type: got %q want %q", ct, "application/json")
				}

				if tt.wantStatus != http.StatusOK {
					return
				}

				var body struct {
					Msg string `json:"msg"`
				}

				if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}

				if body.Msg != tt.wantMsg {
					t.Fatalf("msg: got %q want %q", body.Msg, tt.wantMsg)
				}

			},
		)
	}
}

func TestHandleGreetingsNewHello(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		url        string
		body       any
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "with user returns greeting",
			method:     http.MethodPost,
			url:        "/greetings/hello",
			body:       model.UserDataDto{Name: "nabil"},
			wantStatus: http.StatusOK,
			wantMsg:    "hello, nabil",
		},
		{
			name:       "missing user returns bad request",
			method:     http.MethodPost,
			url:        "/greetings/hello",
			body:       model.UserDataDto{Name: ""},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid user returns bad request",
			method:     http.MethodPost,
			url:        "/greetings/hello",
			body:       model.UserDataDto{Name: "nabil-bhul"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json returns bad request",
			method:     http.MethodPost,
			url:        "/greetings/hello",
			body:       `{"name":`, // broken JSON
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request

			switch v := tt.body.(type) {
			case string:
				req = httptest.NewRequest(tt.method, tt.url, strings.NewReader(v))
			default:
				b, err := json.Marshal(v)
				if err != nil {
					t.Fatal(err)
				}
				req = httptest.NewRequest(tt.method, tt.url, bytes.NewReader(b))
			}

			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handleGreetingsNewHello(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("status: got %d want %d", res.StatusCode, tt.wantStatus)
			}

			ct := res.Header.Get("Content-Type")
			mediaType, _, err := mime.ParseMediaType(ct)
			if err != nil {
				t.Fatalf("bad content-type %q: %v", ct, err)
			}
			if mediaType != "application/json" {
				t.Fatalf("content-type: got %q want %q", mediaType, "application/json")
			}

			if tt.wantStatus != http.StatusOK {
				return
			}

			var body struct {
				Msg string `json:"msg"`
			}
			if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if body.Msg != tt.wantMsg {
				t.Fatalf("msg: got %q want %q", body.Msg, tt.wantMsg)
			}
		})
	}
}
