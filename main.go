package main

import (
	"flag"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

var database = flag.String("database", "", "Path to sqlite3 database")
var logSql = flag.Bool("log-sql", false, "Log SQL")

func main() {
	flag.Parse()
	if *database == "" {
		log.Fatal("--database must be set")
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
