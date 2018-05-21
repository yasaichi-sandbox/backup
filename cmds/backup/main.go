package main

import (
	"errors"
	"flag"
	"log"
)

func main() {
	var fatalErr error

	defer func() {
		if fatalErr != nil {
			flag.PrintDefaults()
			log.Fatal(fatalErr)
		}
	}()

	dbpath := flag.String("db", "./backupdata", "データベースのディレクトリへのパス")
	_ = dbpath

	flag.Parse()
	if len(flag.Args()) < 1 { // NOTE: `Args` returns the non-flag arguments.
		fatalErr = errors.New("エラー; コマンドを指定してください")
		return
	}
}
