package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Thitiphong/finalexam/customer"

	_ "github.com/lib/pq"
)

var db *sql.DB

// export DATABASE_URL=postgres://zrkcqukc:kVMmAKKN7appxvVWa-g2GKRq15Ndg4Q_@baasu.db.elephantsql.com:5432/zrkcqukc
//
// docker run --name some-postgres -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 -d postgres
// export DATABASE_URL="postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable"
// docker run -it --rm --link some-postgres:postgres postgres psql -h postgres -U postgres
func InsertCustomer(name, email, status string) *sql.Row {
	return db.QueryRow("insert into customers (name, email, status) values ($1, $2, $3) returning id", name, email, status)
}

func SelectByKeyCustomer(id int) (*sql.Row, error) {

	stmt, err := db.Prepare("SELECT id, name, email, status FROM customers WHERE id=$1")
	if err != nil {
		return nil, err
	}

	return stmt.QueryRow(id), nil
}

// func UpdateCustomer(name, email, status string) *sql.Row {
// 	stmt, err := db.Prepare("UPDATE customers SET  name=$1, email=$2 ,status=$3  WHERE id=$1;")

// }

func DeleteCustomer(id int) (sql.Result, error) {
	stmt, err := db.Prepare("delete from customers WHERE id=$1;")
	if err != nil {
		log.Fatal("can't prepare statment update", err)
	}
	return stmt.Exec(id)
}

func UpdateCustomer(id int, t customer.Customer) (sql.Result, error) {
	t.ID = id

	stmt, err := db.Prepare("UPDATE customers SET  name=$2, email=$3 ,status=$4  WHERE id=$1;")

	if err != nil {

		return nil, err
	}
	return stmt.Exec(id, t.Name, t.Email, t.Status)

}

func CreateTable() {
	createTb := `CREATE TABLE IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY, 
		name TEXT, 
		email TEXT, 
		status TEXT
	)`
	if _, err := db.Exec(createTb); err != nil {
		log.Fatal("Cann't create table", err)
	}
	fmt.Println("create table success")
}

func Conn() *sql.DB {
	if db != nil {
		return db
	}
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("connot connect to database", err)
	}
	return db
}

func DisconnectDB() {
	db.Close()
}
