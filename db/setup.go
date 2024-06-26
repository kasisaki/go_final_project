package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const DefaultDbName = "scheduler.db"

var DB *sql.DB
var err error

func initializeDB(dbPath string, initScript string) (*sql.DB, error) {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Printf("Database file '%s' does not exist. Creating a new one.", dbPath)
		_, err := os.Create(dbPath)
		if err != nil {
			return nil, err
		}
	}

	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	requests := strings.Split(initScript, ";\n")

	for _, request := range requests {
		_, err := DB.Exec(request)
		if err != nil {
			log.Fatalf("DATABASE setup completed with an error: %s\n", err)
		}
	}

	return DB, nil
}

func SetupDb() *sql.DB {
	initScriptPath := filepath.Join("schema.sql")
	dbName, exists := os.LookupEnv("DBFILE")
	if !exists {
		log.Println("DB name is not provided... Setting to default")
		dbName = DefaultDbName
	}

	initScript, errLoad := os.ReadFile(initScriptPath)
	if errLoad != nil {
		log.Fatal("DB setup script was not loaded properly")
	}

	db, _ := initializeDB(dbName, string(initScript))
	return db
}
