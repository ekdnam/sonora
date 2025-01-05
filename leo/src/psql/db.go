package psql

import (
	"database/sql"
	"fmt"
	"leo/src/utils"
	"strconv"

	_ "github.com/lib/pq"
)

// InitDB initializes a connection to the PostgreSQL database using credentials from the .env file.
//
// This function reads database configuration from the specified .env file and establishes
// a connection to PostgreSQL. It implements robust error handling and follows security
// best practices for database connectivity.
//
// The following environment variables are required:
//   - PSQL_UNAME: Database username
//   - PSQL_PORT: Database port number
//   - PSQL_HOST: Database host address
//   - PSQL_DBNAME: Database name
//   - PSQL_PWD: Database password (optional, defaults to empty string)
//
// Parameters:
//   - envPath: Path to the .env file containing database credentials
//
// Returns:
//   - *sql.DB: Database connection pool if successful
//   - error: Error if connection fails or configuration is invalid
//
// Example usage:
//
//	db, err := InitDB(".env")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer db.Close()
func InitDB(envPath string) (*sql.DB, error) {
	// Read database configuration from .env
	username, err := utils.LoadConfig(envPath, "PSQL_UNAME")
	if err != nil {
		return nil, fmt.Errorf("failed to load PSQL_UNAME from %s: %w", envPath, err)
	}

	portStr, err := utils.LoadConfig(envPath, "PSQL_PORT")
	if err != nil {
		return nil, fmt.Errorf("failed to load PSQL_PORT from %s: %w", envPath, err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PSQL_PORT value from %s: %w", envPath, err)
	}

	host, err := utils.LoadConfig(envPath, "PSQL_HOST")
	if err != nil {
		return nil, fmt.Errorf("failed to load PSQL_HOST from %s: %w", envPath, err)
	}

	dbname, err := utils.LoadConfig(envPath, "PSQL_DBNAME")
	if err != nil {
		return nil, fmt.Errorf("failed to load PSQL_DBNAME from %s: %w", envPath, err)
	}

	// Password is optional, defaults to empty string
	password, err := utils.LoadConfig(envPath, "PSQL_PWD")
	if err != nil {
		password = ""
	}

	// Construct connection string
	var connStr string
	if len(password) > 0 {
		connStr = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			host, port, username, password, dbname,
		)
	} else {
		connStr = fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s sslmode=disable",
			host, port, username, dbname,
		)
	}

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Verify connection is working
	if err = db.Ping(); err != nil {
		db.Close() // Clean up before returning error
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
