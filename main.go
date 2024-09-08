package main

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/microsoft/go-mssqldb"
)

var files []string

func main() {
	errENV := godotenv.Load(".env")
	if errENV != nil {
		log.Fatalf("Env dosyası bulunamadı.")
		os.Exit(1)
	}

	tick := time.Tick(300000 * time.Millisecond) // 300000 = 5 dakika

	for range tick {
		root := "pdf/"
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) == ".pdf" {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			panic(err)
		}

		for _, file := range files {
			f, _ := os.Open(file)
			reader := bufio.NewReader(f)
			content, _ := ioutil.ReadAll(reader)
			encoded := base64.StdEncoding.EncodeToString(content)
			server := GetEnvWithKey("server")
			port := GetEnvWithKey("port")
			user := GetEnvWithKey("user")
			password := GetEnvWithKey("password")
			database := GetEnvWithKey("database")
			connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;encrypt=disable;connection+timeout=30", server, user, password, port, database)
			db, err := sql.Open("sqlserver", connString)
			if err != nil {
				log.Fatalf("SQL Server bağlantı sorunu : %s", err)
			}
			defer db.Close()
			query := "insert into belgeler (belge, belge_turu, kullanici_id, notes, status) VALUES ('" + encoded + "', '9', '0', 'WhatsApp Grubu -" + file + "', '0')"
			_, err2 := db.Exec(query)
			if err2 != nil {
				log.Fatalf("belgeler tablosuna erişim yok : %s", err)
			}
			//fmt.Println(file)
			errPDF := removeFile(filepath.Join("", file))
			if errPDF != nil {
				fmt.Println(errPDF)
			}
		}
		files = nil
	}
}
func removeFile(dir string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return os.Remove(path)
	})
	if err != nil {
		return err
	}
	return nil
}
func GetEnvWithKey(key string) string {
	return os.Getenv(key)
}
