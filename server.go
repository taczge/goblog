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

type EntriesPage struct {
	Entries        []Entry
	ExistsPrevPage bool
	ExistsNextPage bool
	PrevOffset     int
	NextOffset     int
}

func NewEntriesPageByDB(conf Config, db Database, offset int) *EntriesPage {
	entries := db.GetEntries(conf.ArticlePerPage, offset)

	prevOffset := offset - conf.ArticlePerPage
	nextOffset := offset + conf.ArticlePerPage
	existsPrevPage := offset > 0
	existsNextPage := nextOffset < db.Size()

	return &EntriesPage{
		Entries:        entries,
		ExistsPrevPage: existsPrevPage,
		ExistsNextPage: existsNextPage,
		PrevOffset:     prevOffset,
		NextOffset:     nextOffset,
	}
}

func getOffset(r *http.Request) int {
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 || err != nil {
		offset = 0
	}

	return offset
}

func makeHomeHandler(conf Config, db Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		offset := getOffset(r)
		page := NewEntriesPageByDB(conf, db, offset)
		err := tmpl.ExecuteTemplate(w, "index", page)
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

func makeEntryHandler(conf Config, db Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := trim(r.URL.Path, "/entry/", ".html")

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

func makeArchiveHandler(conf Config, db Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entries := db.GetEntries(conf.ArchiveListSize, 0)
		err := tmpl.ExecuteTemplate(w, "archive", entries)
		if err != nil {
			log.Printf("internal server error: %s", err.Error())

			status := http.StatusInternalServerError
			http.Error(w, err.Error(), status)
			return
		}
	}
}

var tmpl *template.Template

func Run() {
	log.Println("Load templates.")
	tmpl = template.Must(template.ParseGlob("templates/*.tmpl"))

	conf := LoadConfig()
	db, err := ConnectMySQL(conf)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", makeHomeHandler(conf, db))
	registerFileServer(conf.FileServer)

	http.HandleFunc("/entry/", makeEntryHandler(conf, db))
	http.HandleFunc("/archive.html", makeArchiveHandler(conf, db))

	log.Printf("Run server.")
	port := fmt.Sprintf(":%d", conf.Port)
	log.Fatal(http.ListenAndServe(port, nil))
}
