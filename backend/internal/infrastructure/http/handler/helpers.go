package handler

import (
	"encoding/json"
	"net/http"
)

// decodeJSON reads and decodes a JSON request body into dst.
// Returns an error if the body is malformed or empty.
func decodeJSON(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}
