package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	DBUser         string `json:"db_user"`
	DBPasswd       string `json:"db_passwd"`
	DBName         string `json:"db_name"`
	ArticlePerPage int    `json:"article_per_page"`
	Port           int    `json:"port"`
}

func LoadConfig() Config {
	filename := "config.json"
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	var c Config
	dec.Decode(&c)

	log.Printf("load %s.", filename)

	return c
}

type Entry struct {
	Title string
	Date  time.Time
	Tags  []string
	Body  string
}

type Entries interface {
	Size() int
	GetLatesed(int) []Entry
}

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
	query := "select id, title, body from entry order by id limit ?"
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
		var body string
		if err := rows.Scan(&id, &title, &body); err != nil {
			log.Fatal(err)
		}

		e := Entry{Title: title, Body: body}
		entries[i] = e
		i++
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("invoke query to get latest %d articles.", n)

	return entries
}

func handler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/index.html"))

	db := ConnectDatabase(conf)
	entries := db.GetLatesed(conf.ArticlePerPage)
	err := t.Execute(w, entries)
	if err != nil {
		panic(err)
	}
}

var conf Config

func main() {
	conf = LoadConfig()

	http.HandleFunc("/", handler)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
}
