package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

var root = flag.String("root", "", "Path to source repository root")
var database = flag.String("database", "", "Path to sqlite3 database")
var logSql = flag.Bool("log-sql", false, "Log SQL")

func main() {
	flag.Parse()
	if *root == "" {
		log.Fatal("--root must be set")
	}
	if *database == "" {
		*database = filepath.Join(*root, ".database.sql")
	}

	db, err := gorm.Open("sqlite3", *database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.LogMode(*logSql)
	db.CreateTable(&PackageBase{})
	db.CreateTable(&Package{})
	db.AutoMigrate(&PackageBase{}, &Package{})
	if db.Error != nil {
		log.Fatal(db.Error)
	}
	log.Println("It's OK!")
}
