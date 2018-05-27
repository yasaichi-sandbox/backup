package main

import (
	"flag"
	"github.com/matryer/filedb"
	"github.com/yasaichi-sandbox/backup"
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
			log.Fatalln(fatalErr)
		}
	}()

	interval := flag.Int("interval", 10, "チェックの間隔（秒単位）")
	_ = interval
	archive := flag.String("archive", "archive", "アーカイブの保存先")
	dbpath := flag.String("db", "./db", "filedbデータベースへのパス")
	flag.Parse()

	m := &backup.Monitor{
		Destination: *archive,
		Archiver:    backup.ZIP,
		Paths:       make(map[string]string),
	}
	_ = m

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
