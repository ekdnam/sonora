package psql

import (
	"database/sql"
	"fmt"
	"leo/src/utils"
	"strconv"

	_ "github.com/lib/pq"
)

// DB encapsulates a PostgreSQL database connection pool and provides high-level database operations.
//
// The DB struct serves as the primary interface for database interactions in the Leo application,
// managing a connection pool through sql.DB. It implements connection lifecycle management and
// provides thread-safe access to the underlying PostgreSQL database.
//
// The embedded sql.DB connection pool automatically handles:
// - Connection pooling and reuse
// - Connection health checks
// - Concurrent access through connection multiplexing
//
// Usage:
//
//	db, err := NewDB(".env")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer db.Close()
//
// Thread Safety:
// All methods on DB are safe for concurrent use by multiple goroutines.
// The underlying sql.DB connection pool handles connection management and thread safety.
type DB struct {
	conn *sql.DB
}

// initDB initializes a connection to the PostgreSQL database using credentials from the .env file.
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
//	db, err := initDB(".env")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer db.Close()
func initDB(envPath string) (*sql.DB, error) {
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

// NewDB creates a new database connection wrapper using configuration from the specified env file.
//
// The envPath parameter should point to a valid .env file containing PostgreSQL connection details
// including host, port, username, password (optional), and database name.
//
// Returns:
//   - A pointer to an initialized DB struct containing the database connection
//   - An error if the connection could not be established or configuration is invalid
//
// Example usage:
//
//	db, err := NewDB(".env")
//	if err != nil {
//	    log.Fatalf("Failed to connect to database: %v", err)
//	}
//	defer db.Close()
func NewDB(envPath string) (*DB, error) {
	conn, err := initDB(envPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize DB: %w", err)
	}
	return &DB{conn: conn}, nil
}

// GetConnection returns the underlying *sql.DB connection object.
//
// This method provides direct access to the database connection, which should be used
// judiciously and only when necessary. In most cases, you should prefer using the DB
// struct's higher-level methods rather than accessing the raw connection.
//
// Warning: The returned connection should not be closed directly. Use the DB.Close()
// method instead to ensure proper cleanup.
//
// Returns:
//   - *sql.DB: The underlying database connection object
//
// Example usage:
//
//	db, _ := NewDB(".env")
//	conn := db.GetConnection()
//	// Use conn for raw SQL operations when absolutely necessary
//	// defer db.Close() // Always close using DB.Close(), not conn.Close()
func (db *DB) GetConnection() *sql.DB {
	return db.conn
}

// Close releases all database resources and closes the connection pool.
//
// This method should be called when the database connection is no longer needed,
// typically using defer immediately after creating a new DB instance. It ensures
// proper cleanup of connection pool resources and prevents connection leaks.
//
// The Close method is idempotent - calling it multiple times on the same DB instance
// is safe, though only the first call will have an effect. After Close is called,
// any other methods on the DB instance will return errors.
//
// Returns:
//   - error: nil if the connection was closed successfully, or an error if the
//     close operation failed. Even if an error is returned, the connection pool
//     will be closed.
//
// Example usage:
//
//	db, err := NewDB(".env")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer db.Close() // Always defer Close after successful DB creation
//
// Note: This is a blocking operation that waits for all in-flight queries to
// complete before closing the connection pool. In high-traffic scenarios, consider
// implementing a graceful shutdown mechanism with appropriate timeouts.
func (db *DB) Close() error {
	return db.conn.Close()
}
