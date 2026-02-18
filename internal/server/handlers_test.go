package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
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
