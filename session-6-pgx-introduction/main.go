package main

import (
	"context"
	"fmt"
	"log"

	"solution1/session-6-db-pgx/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// "postgres://YourUserName:YourPassword@YourHostname:5432/YourDatabaseName"
	dsn := "postgresql://postgres:postgres@localhost:5432/postgres"
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Println(err)
	}

	// ping untuk cek koneksi ke database
	err = pool.Ping(context.Background())
	if err != nil {
		log.Print(err)
	}
	fmt.Println("successfully connect to db")

	// query untuk mengambil row
	var u entity.User
	ctx := context.Background()
	err = pool.QueryRow(ctx, "select id,name from users order by id desc limit 1").Scan(&u.ID, &u.Name)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("user retrieved", u)

	// exec untuk menjalankan perintah terkait insert/update/delete
	_, err = pool.Exec(ctx, "insert into users (name,email,password,created_at,updated_at) values "+
		"('test','test@test.com','pass',NOW(),NOW())")
	if err != nil {
		log.Println(err)
	}

	// query untuk mengambil row
	rows, err := pool.Query(ctx, "select id,name from users order by id desc")
	var users []entity.User
	for rows.Next() {
		var u2 entity.User
		rows.Scan(&u2.ID, &u2.Name)
		if err != nil {
			log.Println(err)
		}
		users = append(users, u2)
	}
	fmt.Println("all user retrieved", users)

}
