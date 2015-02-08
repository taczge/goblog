package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

const CONFIG_FILE = "config.json"

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
			log.Printf("internal server error: %s", err.Error())

			status := http.StatusInternalServerError
			http.Error(w, err.Error(), status)
			return
		}
		entries := db.GetEntries(conf.ArticlePerPage, 0)

		err = tmpl.ExecuteTemplate(w, "index", entries)
		if err != nil {
			log.Printf("internal server error: %s", err.Error())

			status := http.StatusInternalServerError
			http.Error(w, err.Error(), status)
			return
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
		id := trim(r.URL.Path, "/entry/", ".html")
		db, err := ConnectDatabase(conf)
		if err != nil {
			log.Printf("internal server error: %s", err.Error())

			status := http.StatusInternalServerError
			http.Error(w, err.Error(), status)
			return
		}

		entry, err := db.GetEntry(id)
		if err != nil {
			log.Printf("internal server error: %s", err.Error())

			status := http.StatusInternalServerError
			http.Error(w, err.Error(), status)
			return
		}

		err = tmpl.ExecuteTemplate(w, "entry", entry)
		if err != nil {
			log.Printf("internal server error: %s", err.Error())

			status := http.StatusInternalServerError
			http.Error(w, err.Error(), status)
			return
		}
	}
}

func makeArchiveHandler(conf Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := ConnectDatabase(conf)
		if err != nil {
			log.Printf("internal server error: %s", err.Error())

			status := http.StatusInternalServerError
			http.Error(w, err.Error(), status)
			return
		}

		entries := db.GetEntries(conf.ArchiveListSize, 0)
		if err != nil {
			log.Printf("internal server error: %s", err.Error())

			status := http.StatusInternalServerError
			http.Error(w, err.Error(), status)
			return
		}

		err = tmpl.ExecuteTemplate(w, "archive", entries)
		if err != nil {
			log.Printf("internal server error: %s", err.Error())

			status := http.StatusInternalServerError
			http.Error(w, err.Error(), status)
			return
		}
	}
}

var tmpl *template.Template

func init() {
	log.Println("Load templates.")
	tmpl = template.Must(template.ParseGlob("templates/*.tmpl"))
}

func run() {
	log.Printf("Run server.")
	conf := LoadConfig()

	http.HandleFunc("/", makeHandler(conf))
	registerFileServer(conf.FileServer)

	http.HandleFunc("/entry/", makeEntryHandler(conf))
	http.HandleFunc("/archive.html", makeArchiveHandler(conf))

	port := fmt.Sprintf(":%d", conf.Port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func main() {
	run()
}
