package transport

import (
	"encoding/json"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteResponse(w, status, map[string]string{"error": message})
}

func DecodeJSON(r *http.Request, dest any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dest)
}
