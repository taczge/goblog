package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"

	_ "github.com/go-sql-driver/mysql"
)

const ENTRY_TABLE_NAME = "entry"

type Database struct {
	db *sql.DB
}

func ConnectDatabase(c Config) (Database, error) {
	db, err := sql.Open("mysql", c.DBUser+":"+c.DBPasswd+"@/"+c.DBName)

	return Database{db: db}, err
}

func (this *Database) Size() int {
	query := "select count(*) from " + ENTRY_TABLE_NAME
	var size int
	err := this.db.QueryRow(query).Scan(&size)
	if err != nil {
		log.Fatal(err)
	}

	return size
}

func (this *Database) GetEntries(n, offset int) []Entry {
	query := "SELECT * FROM " + ENTRY_TABLE_NAME + " ORDER BY id DESC LIMIT ? OFFSET ?"
	rows, err := this.db.Query(query, n, offset)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var entries = make([]Entry, n, n)
	var i = 0
	for rows.Next() {
		var id int
		var title string
		var date mysql.NullTime
		var body string
		if err := rows.Scan(&id, &title, &date, &body); err != nil {
			log.Fatal(err)
		}

		if !date.Valid {
			log.Fatalf("invalid date %+v\n", date)
		}

		e := Entry{Id: id, Title: title, Date: date.Time, Body: body}
		entries[i] = e
		i++
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("invoke query to get %d articles(offset=%d).", n, offset)

	return entries
}

func (this *Database) GetEntry(idString string) (Entry, error) {
	query := "SELECT * FROM " + ENTRY_TABLE_NAME + " WHERE id = ?"

	var id int
	var title string
	var date mysql.NullTime
	var body string

	row := this.db.QueryRow(query, idString)
	err := row.Scan(&id, &title, &date, &body)

	return Entry{Id: id, Title: title, Date: date.Time, Body: body}, err
}

func (this *Database) Post(e Entry) error {
	query := "INSERT INTO " + ENTRY_TABLE_NAME + " (title, date, body) VALUES(?, ?, ?)"

	_, err := this.db.Exec(query, e.Title, e.Date, e.Body)
	if err == nil {
		log.Printf("complete posting %+v.\n", e.Title)
	} else {
		log.Printf("posting %+v ends in failure.\n", e.Title)
	}

	return err
}

func (this *Database) PostFile(filename string) error {
	entry := LoadEntry(filename)

	return this.Post(entry)
}
