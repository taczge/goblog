package main

import "time"

type Entry struct {
	Id    int
	Title string
	Date  time.Time
	Body  string
}
