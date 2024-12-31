package util

import (
	"encoding/json"
	"io"
	"net/http"
)

// ReadJSON reads and decodes JSON from an HTTP response body.
func ReadJSON(resp *http.Response, v interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.Unmarshal(body, v)
}

// WriteJSON writes JSON to an HTTP response writer.
func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// IsJSON checks if a string is valid JSON.
func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}
