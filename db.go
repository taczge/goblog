package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(user, passwd, name string) Database {
	db, err := sql.Open("mysql", user+":"+passwd+"@/"+name)
	if err != nil {
		panic(err)
	}

	return Database{db}
}

func ConnectDatabase(c Config) Database {
	return NewDatabase(c.DBUser, c.DBPasswd, c.DBName)
}

func (this *Database) Size() int {
	query := "select count(*) from entry"
	var nEntry int
	err := this.db.QueryRow(query).Scan(&nEntry)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Database has %d articles.\n", nEntry)

	return nEntry
}

func (this *Database) GetLatesed(n int) []Entry {
	query := "SELECT * FROM entry ORDER BY id DESC LIMIT ?"
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
	query := "SELECT * FROM entry WHERE id = ?"

	var id int
	var title string
	var date mysql.NullTime
	var body string

	err := this.db.QueryRow(query, idString).Scan(&id, &title, &date, &body)

	return Entry{Id: id, Title: title, Date: date.Time, Body: body}, err
}
