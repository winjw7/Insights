package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// The database file path
const dbFile string = "./insights.db"

// Connect returns a connection to the database for regular use
// (assumes the database has already been set up)
func Connect() (*sql.DB, error) {
    // Connect to the existing database
    db, err := sql.Open("sqlite3", dbFile)
    if err != nil {
        return nil, err
    }
    
    // Verify connection works
    if err := db.Ping(); err != nil {
        db.Close()
        return nil, err
    }
    
    return db, nil
}

// Setup initializes the database and creates the schema
func Setup() error {
    // Create database connection
    db, err := sql.Open("sqlite3", dbFile)
    if err != nil {
        return err
    }

	//Close connection when done
    defer db.Close()
    
    // Create tables if they don't exist
    if err := setupTables(db); err != nil {
        return err
    }
    
    log.Println("Database setup completed successfully")
    return nil
}

// Setup the tables that we need
func setupTables(db *sql.DB) error {
    loginEvents := `CREATE TABLE IF NOT EXISTS login_events (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        tenant TEXT NOT NULL,
        user TEXT NOT NULL,
        origin TEXT NOT NULL,
        status TEXT NOT NULL,
        timestamp DATETIME NOT NULL
    )`
    
    _, err := db.Exec(loginEvents)
    if err != nil {
        return err
    }
    
    return nil
}