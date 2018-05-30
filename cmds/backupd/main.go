package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/matryer/filedb"
	"github.com/yasaichi-sandbox/backup"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	archive := flag.String("archive", "archive", "アーカイブの保存先")
	dbpath := flag.String("db", "./db", "filedbデータベースへのパス")
	flag.Parse()

	m := &backup.Monitor{
		Destination: *archive,
		Archiver:    backup.ZIP,
		Paths:       make(map[string]string),
	}

	db, err := filedb.Dial(*dbpath)
	if err != nil {
		fatalErr = err
		return
	}
	defer db.Close()

	col, err := db.C("paths")
	if err != nil {
		fatalErr = err
		return
	}

	var path path
	col.ForEach(func(_ int, data []byte) (stop bool) {
		if err := json.Unmarshal(data, &path); err != nil {
			fatalErr = err
			return true
		}

		m.Paths[path.Path] = path.Hash
		return
	})
	if fatalErr != nil {
		return
	}
	if len(m.Paths) == 0 {
		fatalErr = errors.New("パスがありません。backupツールを使って追加してください")
		return
	}

	check(m, col)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

Loop:
	for {
		select {
		case <-time.After(time.Duration(*interval) * time.Second):
			check(m, col)
		case <-signalChan:
			fmt.Println()
			log.Println("終了します...")
			break Loop // NOTE: Go out of the for loop
		}
	}
}

func check(m *backup.Monitor, col *filedb.C) {
	log.Println("チェックします...")

	counter, err := m.Now()
	if err != nil {
		log.Panicln("バックアップに失敗しました:", err)
	}
	if counter == 0 {
		log.Println(" 変更はありません")
		return
	}

	log.Printf(" %d個のディレクトリがアーカイブされました\n", counter)

	var path path
	col.SelectEach(func(_ int, data []byte) (include bool, newData []byte, stop bool) {
		if err := json.Unmarshal(data, &path); err != nil {
			log.Println(
				"JSONデータの読み込みに失敗しました。"+
					"次の項目に進みます:",
				err,
			)

			return true, data, false
		}

		path.Hash, _ = m.Paths[path.Path]
		newData, err = json.Marshal(&path)
		if err != nil {
			log.Println(
				"JSONデータの書き出しに失敗しました。"+
					"次の項目に進みます:",
				err,
			)

			return true, data, false
		}

		return true, newData, false
	})
}
