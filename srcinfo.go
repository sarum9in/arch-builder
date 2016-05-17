package main

import (
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/jinzhu/gorm"
	"github.com/sarum9in/archutil/srcinfo"
)

var pkgNameReg = regexp.MustCompile("[0-9a-zA-Z_-]+")

type SrcInfoProcessor func(path string) error

func WalkSrcInfo(root string, processor SrcInfoProcessor) error {
	return filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				if err != nil {
					return filepath.SkipDir
				}
			} else {
				if err != nil {
					return nil
				}
				if filepath.Base(path) == ".SRCINFO" {
					return processor(path)
				}
			}
			return nil
		})
}

func StripDependencyPkgName(dependency string) string {
	return pkgNameReg.FindString(dependency)
}

func FillSrcInfo(path string, db *gorm.DB) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	info, err := srcinfo.ParseSrcInfo(f)
	if err != nil {
		return err
	}
	var oldPkgBase PackageBase
	db.Select("id").First(&oldPkgBase, "pkg_base = ?", info.Global.PkgBase)
	pkgBase := PackageBase{
		ID:          oldPkgBase.ID,
		PkgBase:     info.Global.PkgBase,
		Packages:    []Package{},
		Version:     info.Global.PkgVer,
		Release:     info.Global.PkgRel,
		Depenencies: []Package{},
		Directory:   filepath.Dir(path),
	}
	log.Printf("  {%s}", pkgBase.PkgBase)
	depend := func(dependency string) {
		pkgName := StripDependencyPkgName(dependency)
		log.Printf("    %s => %s", dependency, pkgName)
		var dep Package
		db.FirstOrCreate(&dep, Package{Name: pkgName})
		log.Printf("    %+v", dep)
		pkgBase.Depenencies = append(pkgBase.Depenencies, dep)
	}
	for _, dependency := range info.Global.Depends {
		depend(dependency)
	}
	for _, dependency := range info.Global.MakeDepends {
		depend(dependency)
	}
	for _, pkginfo := range info.Packages {
		for _, dependency := range pkginfo.Depends {
			depend(dependency)
		}
		var pkg Package
		db.FirstOrCreate(&pkg, Package{
			PackageBaseID: pkgBase.ID,
			Name:          pkginfo.PkgName,
		})
		log.Printf("  %s", pkginfo.PkgName)
		pkgBase.Packages = append(pkgBase.Packages, pkg)
	}
	db.Save(&pkgBase)
	return nil
}
