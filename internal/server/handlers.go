package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/kabirnayeem99/gohttplearn/internal/model"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("failed to encode json", "err", err)
	}
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"msg": "hello world"})
}

func handleGoodBye(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"msg": "good bye"})
}

func handleHelloParameterized(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	user, hasUser := q.Get("user"), q.Has("user")

	if hasUser && user == "" {
		writeJSONError(w, http.StatusBadRequest, "missing user")
		return
	}

	msg := "hello"
	if hasUser {
		msg = fmt.Sprintf("hello, %s", user)
	}

	writeJSON(w, http.StatusOK, map[string]string{"msg": msg})
}

func handlerGreetingsUserHello(w http.ResponseWriter, r *http.Request) {
	user := r.PathValue("user")
	if user == "" {
		writeJSONError(w, http.StatusNotFound, "missing user")
		return
	}
	if !isValidUsername(user) {
		writeJSONError(w, http.StatusBadRequest, "invalid username")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"msg": fmt.Sprintf("hello, %s", user)})
}

func handleGreetingsHello(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get("user")
	if user == "" || !isValidUsername(user) {
		writeJSONError(w, http.StatusBadRequest, "missing or invalid user")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"msg": fmt.Sprintf("hello, %s", user)})
}

func handleGreetingsNewHello(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var dto model.UserDataDto
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		slog.Error("failed to decode body", "err", err)
		writeJSONError(w, http.StatusBadRequest, "missing or invalid user")
		return
	}

	user := dto.Name
	if user == "" || !isValidUsername(user) {
		writeJSONError(w, http.StatusBadRequest, "missing or invalid user")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"msg": fmt.Sprintf("hello, %s", user)})
}

var usernameRe = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{2,19}$`)

func isValidUsername(u string) bool {
	return usernameRe.MatchString(u)
}
