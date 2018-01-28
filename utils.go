package utils

import (
	"database/sql"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"log"
	"net/http"
	"strconv"
	"time"
)

// InitSnowflakeNode initializes a Snowflake node.
//
// Example usage:
//
//   a.snowflakeNode = utils.InitSnowflakeNode(1)
func InitSnowflakeNode(nodeNumber int64) *snowflake.Node {
	node, err := snowflake.NewNode(nodeNumber)
	if err != nil {
		log.Fatalf("Could not generate Snowflake node for node number: %d", nodeNumber)
	}
	return node
}

// InitDatabase initializes the DB connection.
//
// Example usage:
//
//   connectionString := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", user, password, host, dbname)
//   a.DB = utils.InitDatabase("postgres", connectionString, 3, time.Second*5)
func InitDatabase(driverName, dataSourceName string, numRetries int, sleepDuration time.Duration) *sql.DB {
	var err error
	var db *sql.DB
	for i := 0; i < numRetries; i++ {
		db, err = sql.Open(driverName, dataSourceName)
		if err == nil {
			log.Println("DB connection initialized...")
			break
		}
		log.Println("DB connection failed to initialize... Sleeping...")
		time.Sleep(sleepDuration)
	}
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// NewPrimaryKey generates a Snowflake ID and returns it as an int64.
func NewPrimaryKey(snowflakeNode *snowflake.Node) int64 {
	return snowflakeNode.Generate().Int64()
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
