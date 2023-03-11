package main

import (
	"log"
	"os"
	"time"
	"github.com/sofuetakuma112/mattn-go/functional-options-pattern/server"
)

func main() {
	f, err := os.Create("server.log")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags) // ログはファイルに出力される
	svr := server.New("localhost", 8888, server.WithTimeout(time.Minute), server.WithLogger(logger))
	if err := svr.Start(); err != nil {
		log.Fatal(err)
	}
}