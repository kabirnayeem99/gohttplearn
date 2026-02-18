package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func  handleHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(map[string]string{
		"msg": "hello world",
	})

	if err != nil {
		slog.Error("failed to encode", "err", err)
	}
}
