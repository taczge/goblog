package main

import (
	"testing"
	"time"
)

func TestLoadEntry(t *testing.T) {
	in := "testdata/2012-01-18-183032"
	got := LoadEntry(in)
	want := Entry{
		Id: "2012-01-18-183032",
		Title: "タイトルタイトル",
		Date:  time.Date(2012, time.January, 18, 18, 30, 32, 0, time.UTC),
		Body: `<p>内容</p>

<p>あいうえお</p>
`}

	check(t, "LoadEntry", in, got, want)
}

func check(t *testing.T, testname string, in, got, want interface{}) {
	if got != want {
		t.Errorf("%s: in: %q, got: %q, want %q", testname, in, got, want)
	}
}
