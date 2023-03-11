package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type Todo struct {
	bun.BaseModel `bun:"table:todos,alias:t"`

	ID        int64     `bun:"id,pk,autoincrement"`
	Content   string    `bun:"content,notnull"`
	Done      bool      `bun:"done"`
	Until     time.Time `bun:"until,nullzero"`
	CreatedAt time.Time
	UpdatedAt time.Time `bun:",nullzero"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"` // ,soft_deleteをつけているので、deleted_atがNULLのレコードのみ取得される
}

func main() {
	// PostgreSQLに接続するための情報
	const (
		host     = "localhost"
		port     = 5432
		user     = "admin"
		password = "admin"
		dbname   = "todos"
	)

	// PostgreSQLに接続する
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	sqldb, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer sqldb.Close()

	db := bun.NewDB(sqldb, pgdialect.New())
	defer db.Close()

	// untilがNULLではなく、doneがfalseのtodoをすべて取得する
	var todos []Todo
	ctx := context.Background()
	err = db.NewSelect().Model(&todos).Order("created_at").Where("until is not null").Where("done is false").Scan(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if len(todos) == 0 {
		return
	}

	// メール送信
	from := mail.Address{Name: "TODO Reminder", Address: os.Getenv("MAIL_FROM")}
	var buf bytes.Buffer
	buf.WriteString("From: " + from.String() + "\r\n")
	buf.WriteString("To: " + os.Getenv("MAIL_TO") + "\r\n")
	buf.WriteString("Subject: TODO Reminder\r\n")
	buf.WriteString("\r\n")
	buf.WriteString("This is your todo list\n\n")
	for _, todo := range todos {
		fmt.Fprintf(&buf, "%s %s\n", todo.Until, todo.Content)
	}

	smtpAuth := smtp.PlainAuth(
		os.Getenv("MAIL_DOMAIN"),
		os.Getenv("MAIL_USER"),
		os.Getenv("MAIL_PASSWORD"),
		os.Getenv("MAIL_AUTHSERVER"),
	)
	err = smtp.SendMail(
		os.Getenv("MAIL_SERVER"),
		smtpAuth,
		from.Address,
		[]string{os.Getenv("MAIL_TO")},
		buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
}
