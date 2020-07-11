// package utils contains utility functions
package utils

import (
	"encoding/json"
	"net/http"
)

// SendJSONResponse structures any response to JSON for sending over the wire
func SendJSONResponse(w http.ResponseWriter, status int, message string, payload interface{}) {
	type resStructure struct {
		Message string      `json:"message,omitempty"`
		Data    interface{} `json:"data,omitempty"`
	}

	// make res body
	var structuredRes = resStructure{message, payload}
	response, err := json.MarshalIndent(&structuredRes, "", "  ")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	_, err = w.Write(response)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
