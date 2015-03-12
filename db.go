package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/go-sql-driver/mysql"
)

type Database interface {
	Size() int
	GetEntries(n, offset int) []Entry
	GetEntry(idString string) (Entry, error)
	Post(e Entry)
	PostFile(filename string)
	PostFiles(dir string)
}

type MySQL struct {
	db     *sql.DB
	config Config
}

func ConnectMySQL(c Config) (MySQL, error) {
	db, err := sql.Open("mysql", c.DBUser+":"+c.DBPasswd+"@/"+c.DBName)

	return MySQL{db: db, config: c}, err
}

func (this *MySQL) Size() int {
	query := "SELECT COUNT(*) FROM " + this.config.DBTable
	var size int
	err := this.db.QueryRow(query).Scan(&size)
	if err != nil {
		log.Fatal(err)
	}

	return size
}

func (this *MySQL) GetEntries(n, offset int) []Entry {
	query := "SELECT * FROM " + this.config.DBTable + " ORDER BY id DESC LIMIT ? OFFSET ?"
	rows, err := this.db.Query(query, n, offset)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var entries = make([]Entry, n, n)
	var i = 0
	for rows.Next() {
		var id string
		var title string
		var date mysql.NullTime
		var body string
		if err := rows.Scan(&id, &title, &date, &body); err != nil {
			log.Fatal(err)
		}

		if !date.Valid {
			log.Fatalf("invalid date %+v\n", date)
		}

		entries[i] = NewEntry(id, title, date.Time, body)
		i++
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("invoke query to get %d articles(offset=%d).", n, offset)

	return entries
}

func (this *MySQL) GetEntry(idString string) (Entry, error) {
	query := "SELECT * FROM " + this.config.DBTable + " WHERE id = ?"

	var id string
	var title string
	var date mysql.NullTime
	var body string

	row := this.db.QueryRow(query, idString)
	err := row.Scan(&id, &title, &date, &body)

	log.Printf("invoke query to get article: id=%s.", id)

	return NewEntry(id, title, date.Time, body), err
}

func (this *MySQL) Post(e Entry) error {
	query := "INSERT INTO " + this.config.DBTable + " (id, title, date, body) VALUES(?, ?, ?, ?)"

	body := fmt.Sprintf("%s", e.Body)
	_, err := this.db.Exec(query, e.Id, e.Title, e.Date, body)
	if err == nil {
		log.Printf("complete posting %+v.\n", e.Title)
	} else {
		log.Printf("posting %+v ends in failure.\n", e.Title)
	}

	return err
}

func (this *MySQL) PostFile(filename string) {
	entry := LoadEntry(filename)
	err := this.Post(entry)
	if err != nil {
		panic(err)
	}
}

func (this *MySQL) PostFiles(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		this.PostFile(dir + f.Name())
	}
}

func SetupMySQL(c Config) {
	log.Printf("Setup database.\n")
	db, err := sql.Open("mysql", c.DBUser+":"+c.DBPasswd+"@/")
	if err != nil {
		panic(err)
	}
	log.Printf("Connect database.\n")

	queries := []string{
		fmt.Sprintf("DROP DATABASE IF EXISTS %s", c.DBName),
		fmt.Sprintf("CREATE DATABASE %s", c.DBName),
		fmt.Sprintf("USE %s", c.DBName),
		fmt.Sprintf("CREATE TABLE %s (id CHAR(17) PRIMARY KEY, title VARCHAR(200) NOT NULL, date DATE NOT NULL, body TEXT NOT NULL);", c.DBTable),
	}

	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			panic(err)
		}
		log.Printf("Invoke query: %s\n", q)
	}
	log.Printf("Setup completed successfully.\n")
}
