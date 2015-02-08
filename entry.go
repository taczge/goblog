package main

import (
	"bufio"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var intToMonth = map[int]time.Month{
	1:  time.January,
	2:  time.February,
	3:  time.March,
	4:  time.April,
	5:  time.May,
	6:  time.June,
	7:  time.July,
	8:  time.August,
	9:  time.September,
	10: time.October,
	11: time.November,
	12: time.December,
}

func NewYMD(y, m, d int) time.Time {
	return time.Date(y, intToMonth[m], d, 0, 0, 0, 0, time.UTC)
}

type Entry struct {
	Id    int
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
	body := ""
	for scanner.Scan() {
		body += scanner.Text() + "\n"
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return Entry{Title: ToTitle(body), Date: NewTime(filename), Body: body}
}

func NewTime(filepath string) time.Time {
	filename := path.Base(filepath)

	y, _ := strconv.Atoi(filename[:4])
	m, _ := strconv.Atoi(filename[4:6])
	d, _ := strconv.Atoi(filename[6:8])

	return NewYMD(y, m, d)
}

func ToTitle(body string) string {
	tag := "<h1>"
	begin := strings.Index(body, tag) + len(tag)
	end := strings.Index(body, "</h1>")

	return body[begin:end]
}
