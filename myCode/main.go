package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func dbConnection() (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbDialect := os.Getenv("DB_DIALECT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	log.Println(dsn)
	db, err := sql.Open(dbDialect, dsn)
	if err != nil {
		log.Printf("error while sql connection, err:%v", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Printf("error while pinging, err:%v", err)
		return nil, err
	}

	return db, nil
}

type user struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
}

type handler struct {
	db *sql.DB
}

func (h *handler) insertData(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`error while reading body`))
		return
	}

	var userData user

	if err = json.Unmarshal(body, &userData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`Invalid json body`))
		return
	}

	query := "INSERT INTO users (user_name, email) VALUES (?, ?)"

	_, err = h.db.Exec(query, userData.UserName, userData.Email)
	if err != nil {
		log.Printf("error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`error from db`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`created successfully`))
}

func main() {
	db, err := dbConnection()
	if err != nil {
		return
	}

	defer db.Close()

	h := &handler{db: db}

	http.HandleFunc("/user", h.insertData)

	log.Println("started listening")

	if err = http.ListenAndServe(":8088", nil); err != nil {
		return
	}
}
