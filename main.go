package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"

	_ "github.com/go-sql-driver/mysql"
)

const CONFIG_FILE = "config.json"

type Config struct {
	DBUser         string   `json:"db_user"`
	DBPasswd       string   `json:"db_passwd"`
	DBName         string   `json:"db_name"`
	ArticlePerPage int      `json:"article_per_page"`
	Port           int      `json:"port"`
	FileServer     []string `json:"file_server"`
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
	Id    int
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

func registerFileServer(paths []string) {
	for _, path := range paths {
		pattern := fmt.Sprintf("/%v/", path)
		handler := http.FileServer(http.Dir(path))

		http.Handle(pattern, http.StripPrefix(pattern, handler))
	}
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

func trim(s, prefix, suffix string) string {
	t := strings.TrimPrefix(s, prefix)
	u := strings.TrimSuffix(t, suffix)

	return u
}

func makeEntryHandler(conf Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("call handler: %+v", r.URL)
		id := trim(r.URL.Path, "/entry/", ".html")
		log.Println(id)
		db := ConnectDatabase(conf)
		entry, err := db.GetEntry(id)
		if err != nil {
			log.Printf("not found %+v.\n", r.URL)
			fmt.Fprintf(w, "not found %+v.\n", r.URL)
		} else {
			fmt.Fprintf(w, "%+v\n", entry)
		}
	}
}

func main() {
	log.Printf("run server.")
	conf := LoadConfig()

	http.HandleFunc("/", makeHandler(conf))
	registerFileServer(conf.FileServer)

	http.HandleFunc("/entry/", makeEntryHandler(conf))

	port := fmt.Sprintf(":%d", conf.Port)
	log.Fatal(http.ListenAndServe(port, nil))
}
