package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func main() {
	files := []string{}

	// main.goファイルが存在する階層の絶対パスを取得する？
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// cwdを起点として、ツリー内の各ファイルまたはディレクトリに対してfnを呼び出します。
	err = filepath.WalkDir(cwd, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		files = append(files, path)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(files)	
}