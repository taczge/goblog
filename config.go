package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	DBUser         string   `json:"db_user"`
	DBPasswd       string   `json:"db_passwd"`
	DBName         string   `json:"db_name"`
	ArticlePerPage int      `json:"article_per_page"`
	Port           int      `json:"port"`
	FileServer     []string `json:"file_server"`
}

func LoadConfig() Config {
	file, err := os.Open(CONFIG_FILE)
	if err != nil {
		path, _ := filepath.Abs(CONFIG_FILE)
		log.Fatalf("Cannot find %s.", path)
	}
	defer file.Close()

	var conf Config
	dec := json.NewDecoder(file)
	if err := dec.Decode(&conf); err != nil {
		log.Fatal(err)
	}

	log.Printf("Load %s.", CONFIG_FILE)

	return conf
}
