package database

import (
	"database/sql"
	"log"
	"os"
	//	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
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
		SELECT 
		    (SELECT goalCalories 
		    FROM dailyGoal 
		    WHERE dailyGoal.date <= missing_dates.date 
		    ORDER BY dailyGoal.date DESC 
		    LIMIT 1),
		    0,
		    date
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

func GetTarget(db *DB, date time.Time) (target int, err error ) {

	fmtDate := date.Format("2006-01-02")
	err = db.QueryRow("SELECT goalCalories FROM dailyGoal WHERE date = ?", fmtDate).Scan(&target)

	if err != nil {
		log.Println(err)
		return 1000, err
	}
	// Default return 2000 in case where a target is not found
	return target, nil
}

/*
func GetTarget(db *DB, day string) (target int, err error ) {

	err = db.QueryRow("SELECT goalCalories FROM dailyGoal WHERE strftime('%d', date) = ? ORDER BY id DESC LIMIT 1", day).Scan(&target)

	if err != nil {
		log.Println(err)
		return 1000, err
	}
	// Default return 2000 in case where a target is not found
	return target, nil
}
*/

// Should insert this into dailyGoal table, however had trouble with these;
// Skipping that step and doing direct aggregation
func GetDailyConsumption(db *DB, date time.Time)(sum int, err error) {
	//rows, err := db.Query("SELECT calories FROM dailyConsumption WHERE date = DATE('now')")
	//fmt.Println(date)
	//rows, err := db.Query("SELECT calories FROM dailyConsumption WHERE strftime('%d', date) = ?", strconv.Itoa(date.Day()))
	
	fmtDate := date.Format("2006-01-02")
	rows, err := db.Query("SELECT calories FROM dailyConsumption WHERE date = ?", fmtDate)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows == nil {
		return 0, nil
	}

	for rows.Next() {
		var cals int
		err := rows.Scan(&cals)
		if err != nil {
			return 0, err
		}
		sum += cals
	}

	return sum, nil
}
/*
func GetDailyConsumption(db *DB, day string)(sum int, err error) {
	//rows, err := db.Query("SELECT calories FROM dailyConsumption WHERE date = DATE('now')")
	rows, err := db.Query("SELECT calories FROM dailyConsumption WHERE strftime('%d', date) = ?", day)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows == nil {
		return 0, nil
	}

	for rows.Next() {
		var cals int
		err := rows.Scan(&cals)
		if err != nil {
			return 0, err
		}
		sum += cals
	}

	return sum, nil
}
*/

func GetDailyFoods(db *DB)([]Food, error) {
	rows, err := db.Query("SELECT id, name, calories FROM dailyConsumption WHERE date = DATE('now')")
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

func GetStreak(db *DB)(streak int, err error) {
	for {
		//target, err := db.Query("SELECT goalCalories FROM dailyGoal WHERE date = DATE('now', '-' || ? || ' days')", streak)
		target, err := GetTarget(db, time.Now().AddDate(0, 0, -streak))
		if err != nil {
			return 0, err
		}

		consumption, err := GetDailyConsumption(db, time.Now().AddDate(0, 0, -streak))
		if err != nil {
			return 0, err
		}

		if consumption >= target {
			streak++
		} else {
			break
		}
	}
	return streak, nil
}

// Close DB connection (use pointer receiver such that actual DB is closed)
func (db *DB) Close() error {
	return db.DB.Close()
}
