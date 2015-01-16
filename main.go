package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
		path, _ := filepath.Abs(CONFIG_FILE)
		log.Fatalf("Cannot find %s.", path)
	}
	defer file.Close()

	var conf Config
	dec := json.NewDecoder(file)
	if err := dec.Decode(&conf); err != nil {
		log.Fatal(err)
	}

	log.Printf("Load %s.", CONFIG_FILE)

	return conf
}

type Entry struct {
	Id    int
	Title string
	Date  time.Time
	Body  string
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
		db, err := ConnectDatabase(conf)
		if err != nil {
			panic(err) // いまだけ，じき直す
		}
		entries := db.GetLatesed(conf.ArticlePerPage)

		err = t.Execute(w, entries)
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
		db, err := ConnectDatabase(conf)
		if err != nil {
			panic(err) // いまだけ，じき直す
		}

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
	log.Printf("Run server.")
	conf := LoadConfig()

	http.HandleFunc("/", makeHandler(conf))
	registerFileServer(conf.FileServer)

	http.HandleFunc("/entry/", makeEntryHandler(conf))

	port := fmt.Sprintf(":%d", conf.Port)
	log.Fatal(http.ListenAndServe(port, nil))
}
