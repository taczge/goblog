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
	query := "select count(*) from ?"
	var nEntry int
	err := this.db.QueryRow(query, ENTRY_TABLE_NAME).Scan(&nEntry)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Database has %d articles.\n", nEntry)

	return nEntry
}

func (this *Database) GetLatesed(n int) []Entry {
	query := "SELECT * FROM " + ENTRY_TABLE_NAME + " ORDER BY id DESC LIMIT ?"
	rows, err := this.db.Query(query, n)
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

	log.Printf("invoke query to get latest %d articles.", n)

	return entries
}

func (this *Database) GetEntry(idString string) (Entry, error) {
	query := "SELECT * FROM ? WHERE id = ?"

	var id int
	var title string
	var date mysql.NullTime
	var body string

	row := this.db.QueryRow(query, ENTRY_TABLE_NAME, idString)
	err := row.Scan(&id, &title, &date, &body)

	return Entry{Id: id, Title: title, Date: date.Time, Body: body}, err
}

func (this *Database) Post(e Entry) error {
	query := "INSERT INTO ? (title, date, body) VALUES(?, ?, ?)"

	_, err := this.db.Exec(query, ENTRY_TABLE_NAME, e.Title, e.Date, e.Body)
	if err != nil {
		log.Printf("posting %+v ends in failure.\n", e.Title)
	} else {
		log.Printf("complete posting %+v.\n", e.Title)
	}

	return err
}
