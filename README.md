## my-go-utils
### Usage
```
go get github.com/kevinmichaelchen/my-go-utils
```

### Features
This library provides several helper functions.

#### Environment variables
- reading strings
- reading int64
- reading booleans

#### Type conversions
- parsing strings to int64

#### Request and Response
- writing errors and structs to a http.ResponseWriter
- parsing int64 vars from a map of route variables

#### DB
- initializing a DB connection with retries and intermittent sleeping

#### Snowflake
- initializing new Snowflake nodes
- generating new Snowflake IDs
