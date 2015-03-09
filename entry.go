package main

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"time"
)

type Entry struct {
	Id    string
	Title string
	Date  time.Time
	Body  string
}

func LoadEntry(filename string) Entry {
	fp, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	entry := Entry{}
	var body bytes.Buffer
	isInBody := false
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "ID:") {
			entry.Id = strings.TrimPrefix(line, "ID: ")
			continue
		}

		if strings.HasPrefix(line, "TITLE:") {
			entry.Title = strings.TrimPrefix(line, "TITLE: ")
			continue
		}

		if strings.HasPrefix(line, "DATE:") {
			dateStr := strings.TrimPrefix(line, "DATE: ")
			format := "2006-01-02 15:04:05 -0700 MST"
			date, err := time.Parse(format, dateStr)
			if err != nil {
				panic(err)
			}
			entry.Date = date
			continue
		}

		if strings.HasPrefix(line, "BODY:") {
			isInBody = true
			continue
		}

		if isInBody {
			body.WriteString(strings.TrimSpace(line))
			body.WriteString("\n")
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	entry.Body = body.String()

	return entry
}
