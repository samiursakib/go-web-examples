package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	id int64
	username string
	password string
	createdAt time.Time
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "requested path: %s", r.URL.Path)
	})

	r.HandleFunc("/books/{title}/{page}", func (w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fmt.Println("title:", vars["title"], "page", vars["page"])
		fmt.Fprintf(w, "title: %s, page: %s", vars["title"], vars["page"])
	})
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading env file: %s", err)
	} else {
		fmt.Println("loading env variables")
	}

	uri := fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	fmt.Println(uri)
	db, err := sql.Open("mysql", uri)
	if (err != nil) {
		fmt.Println("error occured during opening database connection", err)
	}

	// uname, pass := "sakib", "sakib123"
	// lastIndex, err := createUser(db, uname, pass)
	// fmt.Println("last inserted id:", lastIndex)

	// users, err := getAllUsers(db)
	// fmt.Println(users)

	if lastUser, ok := getLastUser(db); ok {
		fmt.Println("last user:", lastUser)
	} else {
		fmt.Println("last user not found")
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("$", err)
	}
	http.ListenAndServe(":8080", r)
}





func getAllUsers(db *sql.DB) (users []User, err error) {
	selectQuery := `SELECT * FROM users`
	rows, err := db.Query(selectQuery)
	if (err != nil) {
		fmt.Println("error occured during query")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var u User
		err := rows.Scan(&u.id, &u.username, &u.password, &u.createdAt)
		if (err != nil) {
			fmt.Println("error occured in each element", err)
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func createUser(db *sql.DB, uname string, pass string) (id int64, err error) {
	insertQuery := `INSERT INTO users(username, password) VALUES(?, ?)`
	result, err := db.Exec(insertQuery, uname, pass)
	if (err != nil) {
		fmt.Println("error occured while creating user", err)
		return -1, err
	}
	lastIndex, err := result.LastInsertId()
	if (err != nil) {
		return -1, err
	}
	return lastIndex, nil
}

func getLastUser(db *sql.DB) (u User, ok bool) {
	findLastUserQuery := `SELECT * FROM users ORDER BY id DESC LIMIT 1`
	err := db.QueryRow(findLastUserQuery).Scan(&u.id, &u.username, &u.password, &u.createdAt)
	if (err != nil) {
		return u, false
	}
	return u, true
}
