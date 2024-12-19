package database

import (
	"database/sql"
	"log"
	"os"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func Initialize(dbPath string) (*DB, error) {
	// Create db if file doesnt exist
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			return nil, err
		}		
		file.Close()
	}

	// Open db connection
	sqliteDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error: Could not open db connection.")
		return nil, err
	}

	// Test the connection 
	err = sqliteDB.Ping()
	if err != nil {
		log.Println("Error: Could not ping the db.")
		return nil, err
	}

	// Create tables if necessary
	err = createTables(sqliteDB)
	if err != nil {
		log.Println("Could not create db tables.")
		return nil, err
	}

	// Return pointer to DB struct and error (nil = no error)
	return &DB{sqliteDB}, nil
}

// Function to initialize tables in DB if not already created
func createTables(db *sql.DB) error {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS foods (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			calories INTEGER NOT NULL,
			creationDate DATETIME DEFAULT CURRENT_TIMESTAMP,
			updatedDate DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS users (
	        id INTEGER PRIMARY KEY AUTOINCREMENT,
	        username TEXT NOT NULL UNIQUE,
	        email TEXT NOT NULL UNIQUE,
	        password TEXT NOT NULL,
	        creationDate DATETIME DEFAULT CURRENT_TIMESTAMP,
	        updatedDate DATETIME DEFAULT CURRENT_TIMESTAMP
    	);

	    CREATE TABLE IF NOT EXISTS dailyConsumption (
	        id INTEGER PRIMARY KEY AUTOINCREMENT,
	        userID INTEGER NOT NULL,
	        foodID INTEGER NOT NULL,
	        calories INTEGER NOT NULL,
	        date DATE NOT NULL,
	        creationDate DATETIME DEFAULT CURRENT_TIMESTAMP,
	        updatedDate DATETIME DEFAULT CURRENT_TIMESTAMP,
	        FOREIGN KEY (userID) REFERENCES users(id),
	        FOREIGN KEY (foodID) REFERENCES foods(id)
	    );

	    CREATE TABLE IF NOT EXISTS dailyGoals (
	        id INTEGER PRIMARY KEY AUTOINCREMENT,
	        userID INTEGER NOT NULL,
	        goalCalories INTEGER NOT NULL,
	        consumedCalories INTEGER NOT NULL,
	        date DATE NOT NULL,
	        creationDate DATETIME DEFAULT CURRENT_TIMESTAMP,
	        updatedDate DATETIME DEFAULT CURRENT_TIMESTAMP,
	        FOREIGN KEY (userID) REFERENCES users(id)
	    );`

	_, err := db.Exec(createTableSQL)
	return err
}

// Close DB connection (use pointer receiver such that actual DB is closed)
func (db *DB) Close() error {
	return db.DB.Close()
} 