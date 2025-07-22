package helpers

import (
	"encoding/json"
)

// GetErrorResponseBody returns the json identated error response body
func GetErrorResponseBody(res []byte) string {
	m := map[string]any{}
	json.Unmarshal(res, &m) //nolint:errcheck
	mb, _ := json.MarshalIndent(m, "", " ")
	return string(mb)
}
