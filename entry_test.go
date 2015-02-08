package main

import (
	"testing"
	"time"
)

func TestNewTime(t *testing.T) {
	cases := []struct {
		in   string
		want time.Time
	}{
		{"20120312001", NewYMD(2012, 3, 12)},
		{"dir/20120312001", NewYMD(2012, 3, 12)},
	}

	for _, c := range cases {
		got := NewTime(c.in)
		check(t, "NewTime", c.in, got, c.want)
	}
}

func TestLoadEntry(t *testing.T) {
	in := "testdata/20150120001"
	got := LoadEntry(in)
	want := Entry{
		Title: "エントリのタイトル",
		Date: NewYMD(2015, 1, 20),
		Body: `<h1>エントリのタイトル</h1>
<p>本文の内容</p>
<h2>見出し</h2>
<p>あいうえお</p>
<p>かきくけこ</p>
`}

	check(t, "LoadEntry", in, got, want)
}

func check(t *testing.T, testname string, in, got, want interface{}) {
	if got != want {
		t.Errorf("%s: in: %q, got: %q, want %q", testname, in, got, want)
	}
}
