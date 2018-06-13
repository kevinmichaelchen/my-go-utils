package request

import (
	"encoding/json"
	"net/http"
	stringUtils "github.com/kevinmichaelchen/my-go-utils/string"
)

// RespondWithError writes a JSON error message to the client.
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

// RespondWithMessage writes a simple JSON message to the client.
func RespondWithMessage(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"msg": message})
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
	if !stringUtils.IsParseableAsInt64(routeVar) {
		RespondWithError(w, http.StatusBadRequest, "Invalid route var: "+varKey)
		return 0, false
	}
	return stringUtils.StringToInt64(routeVar), true
}

// GetInt32 returns the given route parameter as an int32.
func GetInt32(w http.ResponseWriter, routeVars map[string]string, varKey string) (int32, bool) {
	routeVar := routeVars[varKey]
	if !stringUtils.IsParseableAsInt32(routeVar) {
		RespondWithError(w, http.StatusBadRequest, "Invalid route var: "+varKey)
		return 0, false
	}
	return stringUtils.StringToInt32(routeVar), true
}