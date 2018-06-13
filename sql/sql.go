package sql

import (
	"database/sql"
	"net/http"
	"time"
	"log"
	"strings"
	"encoding/hex"
	"github.com/google/uuid"
	requestUtils "github.com/kevinmichaelchen/my-go-utils/request"
)

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
	} else {
		log.Printf("Successfully connected to DB: %s\n", dataSourceName)
	}
	return db
}

// StartTransaction starts a transaction within the context of a request.
// If the transaction cannot be started, the ResponseWriter returns an error.
func StartTransaction(w http.ResponseWriter, DB *sql.DB) (*sql.Tx, error) {
	var err error
	tx, err := DB.Begin()
	if err != nil {
		tx.Rollback()
		requestUtils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return nil, err
	}
	return tx, err
}

// UnhexUuid converts a (google) UUID to unhexed bytes.
// This method is useful for MySQL apps, since MySQL has no native support for UUIDs (unlike Postgres).
// This method is functionally equivalent to running
// SELECT UNHEX(REPLACE(UUID(), "-", ""));
// in MySQL.
func UnhexUuid(uuid uuid.UUID) ([]byte, error) {
	s := strings.Replace(uuid.String(), "-", "", -1)
	return hex.DecodeString(s)
}