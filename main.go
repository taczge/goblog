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

	"github.com/go-sql-driver/mysql"

	_ "github.com/go-sql-driver/mysql"
)

const CONFIG_FILE = "config.json"

type Config struct {
	DBUser         string `json:"db_user"`
	DBPasswd       string `json:"db_passwd"`
	DBName         string `json:"db_name"`
	ArticlePerPage int    `json:"article_per_page"`
	Port           int    `json:"port"`
}

func LoadConfig() Config {
	file, err := os.Open(CONFIG_FILE)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	var conf Config
	dec.Decode(&conf)

	log.Printf("load %s.", CONFIG_FILE)

	return conf
}

type Entry struct {
	Title string
	Date  time.Time
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
	query := "select title, date, body from entry order by id limit ?"
	rows, err := this.db.Query(query, n)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var entries = make([]Entry, n, n)
	var i = 0
	for rows.Next() {
		var title string
		var date mysql.NullTime
		var body string
		if err := rows.Scan(&title, &date, &body); err != nil {
			log.Fatal(err)
		}

		if !date.Valid {
			log.Fatalf("invalid date %+v\n", date)
		}

		e := Entry{Title: title, Date: date.Time, Body: body}
		entries[i] = e
		i++
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("invoke query to get latest %d articles.", n)

	return entries
}

func makeHandler(conf Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.ParseFiles("templates/index.html"))
		db := ConnectDatabase(conf)
		entries := db.GetLatesed(conf.ArticlePerPage)

		err := t.Execute(w, entries)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	log.Printf("run server.")
	conf := LoadConfig()

	http.HandleFunc("/", makeHandler(conf))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
}
