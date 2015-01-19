package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

const CONFIG_FILE = "config.json"

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
		db, err := ConnectDatabase(conf)
		if err != nil {
			panic(err) // いまだけ，じき直す
		}
		entries := db.GetEntries(conf.ArticlePerPage, 0)

		err = tmpl.ExecuteTemplate(w, "index", entries)
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
		db, err := ConnectDatabase(conf)
		if err != nil {
			panic(err) // いまだけ，じき直す
		}

		entry, err := db.GetEntry(id)
		if err != nil {
			log.Printf("not found %+v.\n", r.URL)
		}

		err = tmpl.ExecuteTemplate(w, "entry", entry)
		if err != nil {
			panic(err) // TODO: internal server error?
		}
	}
}

func makeArchiveHandler(conf Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := ConnectDatabase(conf)
		if err != nil {
			panic(err) // TODO: internal server error?
		}

		entries := db.GetEntries(conf.ArchiveListSize, 0)
		if err != nil {
			log.Printf("not found %+v.\n", r.URL)
			fmt.Fprintf(w, "not found %+v.\n", r.URL)
		}

		err = tmpl.ExecuteTemplate(w, "archive", entries)
		if err != nil {
			panic(err) // TODO: internal server error?
		}
	}
}

var tmpl *template.Template
func init() {
	log.Println("Load templates.")
	tmpl = template.Must(template.ParseGlob("templates/*.tmpl"))
}

func main() {
	log.Printf("Run server.")
	conf := LoadConfig()

	http.HandleFunc("/", makeHandler(conf))
	registerFileServer(conf.FileServer)

	http.HandleFunc("/entry/", makeEntryHandler(conf))
	http.HandleFunc("/archive.html", makeArchiveHandler(conf))

	port := fmt.Sprintf(":%d", conf.Port)
	log.Fatal(http.ListenAndServe(port, nil))
}
