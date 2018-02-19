package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
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

// StringToInt64 converts a string to int64, since strconv doesn't provide this straight up.
func StringToInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

// IsParseableAsInt64 checks whether a string is parseable as int64.
func IsParseableAsInt64(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return false
	}
	return true
}
