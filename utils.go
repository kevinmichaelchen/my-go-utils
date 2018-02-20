package utils

import (
	"encoding/json"
	"net/http"
)

// RespondWithError writes a JSON error message to the client.
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

// RespondWithJSON writes a JSON struct to the client.
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// GetInt64 returns the given route parameter as an int64.
func GetInt64(w http.ResponseWriter, routeVars map[string]string, varKey string) (int64, bool) {
	routeVar := routeVars[varKey]
	if !IsParseableAsInt64(routeVar) {
		RespondWithError(w, http.StatusBadRequest, "Invalid route var: "+varKey)
		return 0, false
	}
	return StringToInt64(routeVar), true
}
