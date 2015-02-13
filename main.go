package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func registerFileServer(paths []string) {
	for _, path := range paths {
		pattern := fmt.Sprintf("/%v/", path)
		handler := http.FileServer(http.Dir(path))

		http.Handle(pattern, http.StripPrefix(pattern, handler))
	}
}

type HomePage struct {
	Entries []Entry
	Offset  int
}

func makeHomeHandler(conf Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := ConnectDatabase(conf)
		if err != nil {
			log.Printf("internal server error: %s", err.Error())

			status := http.StatusInternalServerError
			http.Error(w, err.Error(), status)
			return
		}

		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if offset <= 0 || err != nil {
			offset = 0
		}
		entries := db.GetEntries(conf.ArticlePerPage, offset)
		page := HomePage{
			Entries: entries,
			Offset:  offset + conf.ArticlePerPage,
		}

		err = tmpl.ExecuteTemplate(w, "index", page)
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

func run() {
	log.Println("Load templates.")
	tmpl = template.Must(template.ParseGlob("templates/*.tmpl"))

	log.Printf("Run server.")
	conf := LoadConfig()

	http.HandleFunc("/", makeHomeHandler(conf))
	registerFileServer(conf.FileServer)

	http.HandleFunc("/entry/", makeEntryHandler(conf))
	http.HandleFunc("/archive.html", makeArchiveHandler(conf))

	port := fmt.Sprintf(":%d", conf.Port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func main() {
	run()
}
