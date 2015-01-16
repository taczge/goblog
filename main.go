package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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
