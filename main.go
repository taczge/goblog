package main

import (
	"html/template"
	"net/http"
	"time"
)

type Entry struct {
	Title string
	Date  time.Time
	Tags  []string
	Body  []byte
}

func createEntry() []Entry {
	t1 := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	t2 := time.Date(2001, time.November, 10, 23, 0, 0, 0, time.UTC)
	t3 := time.Date(2000, time.November, 10, 23, 0, 0, 0, time.UTC)

	return []Entry{
		Entry{Title: "なんちゃら", Date: t1, Tags: []string{"a"}, Body: []byte("本文")},
		Entry{Title: "タイトル", Date: t2, Tags: []string{"t", "u"}, Body: []byte("本文だぃ")},
		Entry{Title: "ぱんぱん", Date: t3, Tags: []string{"o", "o"}, Body: []byte("あららら")}}
}

func handler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/index.html"))

	entries := createEntry()
	err := t.Execute(w, entries)
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.ListenAndServe(":8080", nil)
}
