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

	err = createEntries(sqliteDB)
	if err != nil {
		log.Println("Could not initialize dates.")
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

// Ensure that entries for calories exist for the past 7 days 
// Uses default 2000, 0 in cases where the app was not used
func createEntries(db *sql.DB) error {
	createEntriesSQL := `
		WITH RECURSIVE dates(date) AS (
		  SELECT DATE('now', '-6 days')
		  UNION ALL
		  SELECT DATE(date, '+1 days')
		  FROM dates
		  WHERE date < DATE('now')
		),
		missing_dates AS (
		  SELECT dates.date
		  FROM dates
		  LEFT JOIN dailyGoal ON DATE(dailyGoal.date) = dates.date
		  WHERE dailyGoal.id IS NULL
		)
		INSERT INTO dailyGoal (goalCalories, consumedCalories, date)
		SELECT 2000, 0, date
		FROM missing_dates;
	`

	_, err := db.Exec(createEntriesSQL)
	return err
}


func GetFoods(db *DB) ([]Food, error) {
	rows, err := db.Query("SELECT * FROM food ORDER BY calories ASC")
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

func GetTarget(db *DB) (target int, err error ) {

	err = db.QueryRow("SELECT goalCalories FROM dailyGoal ORDER BY id DESC").Scan(&target)

	if err != nil {
		return 1000, err
	}
	// Default return 2000 in case where a target is not found
	return target, nil
}

// Close DB connection (use pointer receiver such that actual DB is closed)
func (db *DB) Close() error {
	return db.DB.Close()
}
