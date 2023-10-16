package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Book struct {
	id     int
	name   string
	author string
}

func main() {
	db, err := sql.Open("sqlite3", "./books.db")
	if err != nil {
		log.Println(err)
	}

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS books (id INTEGER PRIMARY KEY, isbn INTEGER, author VARCHAR(64), name VARCHAR(64) NULL)")
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Successfully created table books!")
	}

	statement.Exec()
	dbOperations(db)
}

func dbOperations(db *sql.DB) {
	//Create
	// statement, _ := db.Prepare("insert into books (name, author, isbn) values(?,?,?)")
	// statement.Exec("HANDSON RESTful API WITH GO", "Yellavula", "4320985")
	// log.Println("Inserted the book into database!")

	//Read
	rows, _ := db.Query("select id, name, author from books")
	var tempBook Book
	for rows.Next() {
		rows.Scan(&tempBook.id, &tempBook.name, &tempBook.author)
		log.Printf("ID: %d Book: %s Author: %s", tempBook.id, tempBook.name, tempBook.author)
	}

	// //Update
	// statement, _ := db.Prepare("update books set name=? where id=?")
	// statement.Exec("Concurrency with Go", 5)
	//Delete
	statement, _ := db.Prepare("DELETE FROM books WHERE id=?")
	statement.Exec(5)

}
