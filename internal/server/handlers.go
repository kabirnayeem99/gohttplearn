package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func handleHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(map[string]string{
		"msg": "hello world",
	})

	if err != nil {
		slog.Error("failed to encode", "err", err)
	}
}

func handleGoodBye(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(map[string]string{
		"msg": "good bye",
	})

	if err != nil {
		slog.Error("failed to encode", "err", err)
	}
}

func handleHelloParameterized(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	hasUserQuery := r.URL.Query().Has("user")
	user := r.URL.Query().Get("user")

	if hasUserQuery && user == "" {
		http.Error(w, `{"error":"missing user"}`, http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	msg := "hello"
	if hasUserQuery {
		msg = fmt.Sprintf("hello, %s", user)
	}

	err := json.NewEncoder(w).Encode(map[string]string{
		"msg": msg,
	})

	if err != nil {
		slog.Error("failed to encode", "err", err)
	}
}
