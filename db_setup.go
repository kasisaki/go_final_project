package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func setupDb() {
	path := filepath.Join("schema.sql")
	dbName, exists := os.LookupEnv("TODO_DBFILE")
	if !exists {
		log.Println("DB name is not provided... Setting to default")
		dbName = "scheduler.db"
	}
	dbFile := filepath.Join(dbName)
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	if install {
		file, errLoad := os.ReadFile(path)
		if errLoad != nil {
			log.Fatal("DB setup script was not loaded properly")
		}

		db, err := sql.Open("sqlite3", dbFile)
		if err != nil {
			log.Println("Error connecting to DB")
		}
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				log.Println("Error closing DB session")
			}
		}(db)

		requests := strings.Split(string(file), ";\n")

		for _, request := range requests {
			_, err := db.Exec(request)
			if err != nil {
				log.Fatalf("DATABASE setup completed with an error: %s\n", err)
			}
		}
	}
}
