package response

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{
		"error": message,
	})
}

func Message(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{
		"message": message,
	})
}
