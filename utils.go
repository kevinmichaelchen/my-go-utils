package utils

import (
	"database/sql"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"google.golang.org/grpc"
)

func StartTransaction(w http.ResponseWriter, DB *sql.DB) (*sql.Tx, error) {
	var err error
	tx, err := DB.Begin()
	if err != nil {
		tx.Rollback()
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return nil, err
	}
	return tx, err
}

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

// NewPrimaryKey generates a Snowflake ID and returns it as an int64.
func NewPrimaryKey(snowflakeNode *snowflake.Node) int64 {
	return snowflakeNode.Generate().Int64()
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

func InitGrpcConn(address string, numRetries int, sleepDuration time.Duration) *grpc.ClientConn {
	var err error
	var conn *grpc.ClientConn
	for i := 0; i < numRetries; i++ {
		conn, err = grpc.Dial(address, grpc.WithInsecure())
		if err == nil {
			log.Println("Grpc connection initialized...")
			break
		}
		log.Println("Grpc connection failed to initialize... Sleeping...")
		time.Sleep(sleepDuration)
	}
	if err != nil {
		log.Fatal(err)
	}
	return conn
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

// EnvOrInt returns the environment variable for the given key as an int,
// or the default value if no environment variable is found.
func EnvOrInt(key string, defaultVal int) int {
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("No env var found for %s. Using default value: %d\n", key, defaultVal)
		return defaultVal
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("Env var %s needs to be an int...\n", key)
	}
	return i
}

// EnvOrString returns the environment variable for the given key,
// or the default value if no environment variable is found.
func EnvOrString(key, defaultVal string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("No env var found for %s. Using default value: %s\n", key, defaultVal)
		return defaultVal
	}
	return v
}

// MustEnv returns the environment variable for the given key, or exits if no such variable is found.
func MustEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("No env var for: %s", key)
	}
	return v
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
