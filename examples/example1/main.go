package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/nikhilbhatia08/EuphoriaDB/driver"
)

func main() {
	dbDirectory := "./testdb1"
	defer func() {
		if err := os.RemoveAll(dbDirectory); err != nil {
			log.Fatalf("Failed to clean up database directory: %v\n", err)
		}
	}()

	db, err := sql.Open("euphoriadb", dbDirectory)
	if err != nil {
		log.Fatalf("Failed to open euphoriaDB: %v\n", err)
	}
	defer db.Close()

	createQuery := "CREATE TABLE employees(id int, name VARCHAR(32), age int)"
	if _, err := db.Exec(createQuery); err != nil {
		log.Fatalf("error creating employees table: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("error beginning transaction: %v", err)
	}

	if _, err := tx.Exec("INSERT INTO employees(id, name, age) VALUES (1, 'new name 1', 22)"); err != nil {
		_ = tx.Rollback()
		log.Fatalf("error inserting to table: %v", err)
	}

	if _, err := tx.Exec("INSERT INTO employees(id, name, age) VALUES (2, 'new name 2', 25)"); err != nil {
		_ = tx.Rollback()
		log.Fatalf("error inserting to table: %v", err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("error committing transaction: %v", err)
	}

	fmt.Println("querying rows: ")
	selectQuery := "SELECT id, name, age FROM employees"
	rows, err := db.Query(selectQuery)
	if err != nil {
		log.Fatalf("error querying rows: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var age int

		if err := rows.Scan(&id, &name, &age); err != nil {
			log.Fatalf("error scanning rows: %v", err)
		}

		fmt.Printf("Employee details: id: %d, name: %s, age: %d\n", id, name, age)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Rows iteration error: %v\n", err)
	}
}
