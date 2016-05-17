package main

import (
	_ "github.com/jinzhu/gorm"
)

type PackageBase struct {
	ID          int
	PkgBase     string
	Packages    []Package
	Version     string
	Release     int
	Depenencies []Package `gorm:"many2many:package_dependencies;"`
	Directory   string
}

type Package struct {
	ID            int
	PackageBaseID int    `sql:"index"`
	Name          string `sql:"index"`
}
