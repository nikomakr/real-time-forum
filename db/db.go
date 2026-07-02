package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // The _ shhh-ed the compilation error. Import library only for side effect, init()
)

// Declared a package-level variable. The type is *sql.D aka a pointer to database/sql.DB. This variable will hold the database connection pool and connects to standard library's database/sql package.
var DB *sql.DB

func Init(path string) { // path is the path to the database file.

	DB, err := sql.Open("sqlite3", path+"?_foreign_keys=on")
	if err != nil {
		log.Fatalf("could not open database: %v", err) // 
	}

	if err = DB.Ping(); err != nil { // Ping() method is used to verify that the database connection is still alive, establishing a connection if necessary.
		log.Fatalf("could not reach database: %v", err)
	}

	DB.Exec("PRAGMA journal_mode=WAL;") // Set the journal mode to Write-Ahead Logging (WAL) for better concurrency and performance.

	Migrate() // Call the Migrate function to create tables and perform any necessary migrations.

	log.Println("database ready")
}