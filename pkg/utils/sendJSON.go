// package utils contains utility functions
package utils

import (
	"encoding/json"
	"net/http"
)

// GenericJSONRes structures a json response
type GenericJSONRes struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// SendJSONResponse structures any response to JSON for sending over the wire
func SendJSONResponse(w http.ResponseWriter, status int, r *GenericJSONRes) {
	jsonRes, err := json.MarshalIndent(r, "", "  ")

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	_, err = w.Write(jsonRes)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
