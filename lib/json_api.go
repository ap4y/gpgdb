package lib

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, body interface{}) {
	data, _ := json.Marshal(body)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func ErrorJSON(w http.ResponseWriter, error string, code int) {
	errors := map[string]string{"error": error}
	data, _ := json.Marshal(errors)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
