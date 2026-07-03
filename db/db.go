package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // The _ shhh-ed the compilation error. Import library only for side effect, hidden init()
)

// Declared a package-level variable. The type is *sql.D aka a pointer to database/sql.DB. This variable will hold the database connection pool and connects to standard library's database/sql package.
var DB *sql.DB

func Init(path string) { // path is the server's file path, SQLite, of the database file.
var err error // panic fixed with: var err error declares err explicitly. Plain = then assigns into the existing package-level DB rather than creating a new shadow variable.

	DB, err = sql.Open("sqlite3", path+"?_foreign_keys=on") // WITH ?_foreign_keys=on, SQLite will enforce foreign key constraints, ensuring referential integrity between related tables. This is important for maintaining data consistency and preventing orphaned records in the database. Don't forget Niko by default SQLite does not enforce foreign key constraints, so enabling this option is crucial when working with relational data.
	if err != nil {
		log.Fatalf("could not open database: %v", err)
	}

	if err = DB.Ping(); err != nil { // Ping() method is used to verify that the database connection is still alive, establishing a connection if necessary.
		log.Fatalf("could not reach database: %v", err)
	}

	DB.Exec("PRAGMA journal_mode=WAL;") // Set the journal mode to Write-Ahead Logging (WAL) for better concurrency and performance. WAL, WebSocket and Go routines work hand to hand. WAL allows multiple readers and a single writer to access the database simultaneously, improving performance in concurrent scenarios.

	Migrate() // Call the Migrate function to create tables and perform any necessary migrations.

	log.Println("database ready")
}