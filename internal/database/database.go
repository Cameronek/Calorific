package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

type DB struct {
	*sql.DB
}

type Food struct {
	ID int
	Name string
	Calories int
}

func Initialize(dbPath string) (*DB, error) {
	// Create db if file doesnt exist
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			log.Println("Error: Could not create db")
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

		CREATE TABLE IF NOT EXISTS food (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL, 
			calories INTEGER NOT NULL
		);

		CREATE TABLE IF NOT EXISTS dailyConsumption (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL, 
			calories TEXT NOT NULL, 
			date DATE NOT NULL
		);

		CREATE TABLE IF NOT EXISTS dailyGoal (
	        id INTEGER PRIMARY KEY AUTOINCREMENT,
	        goalCalories INTEGER NOT NULL,
	        consumedCalories INTEGER NOT NULL,
	        date DATE NOT NULL
	    );
		`

	_, err := db.Exec(createTableSQL)
	return err
}


func GetFoods(db *DB) ([]Food, error) {
	rows, err := db.Query("SELECT * FROM food")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var foods []Food

	for rows.Next() {
		var food Food
		err := rows.Scan(&food.ID, &food.Name, &food.Calories)
		if err != nil {
			return nil, err
		}
		foods = append(foods, food)
	}

	return foods, nil
}

// Close DB connection (use pointer receiver such that actual DB is closed)
func (db *DB) Close() error {
	return db.DB.Close()
}
