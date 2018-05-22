package main

import (
	"errors"
	"flag"
	"github.com/matryer/filedb"
	"log"
)

type path struct {
	Path string
	Hash string
}

func main() {
	var fatalErr error
	defer func() {
		if fatalErr != nil {
			flag.PrintDefaults()
			log.Fatal(fatalErr)
		}
	}()

	dbpath := flag.String("db", "./backupdata", "データベースのディレクトリへのパス")

	flag.Parse()
	if len(flag.Args()) < 1 { // NOTE: `Args` returns the non-flag arguments.
		fatalErr = errors.New("エラー; コマンドを指定してください")
		return
	}

	db, err := filedb.Dial(*dbpath)
	if err != nil {
		fatalErr = err
		return
	}
	defer db.Close()

	col, err := db.C("paths")
	_ = col
	if err != nil {
		fatalErr = err
		return
	}
}
