package Clients

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Client struct {
	ID           int    `json:"id"`
	FirstName    string `json:"firstname"`
	LastName     string `json:"lastname"`
	Password     string `json:"password"`
	Descriptions string `json:"descriptions"`
}

// Veritabanı bağlantısı
func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "mertlinux:123@tcp/api_database")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Veritabanı bağlantısı başarılı")

	return db, nil
}

func GetClientList(w http.ResponseWriter, r *http.Request) {
	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM client")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	clients := []Client{}

	for rows.Next() {
		var c Client
		if err := rows.Scan(&c.ID, &c.FirstName, &c.LastName, &c.Password, &c.Descriptions); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		clients = append(clients, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clients)
}

func AddClient(w http.ResponseWriter, r *http.Request) {

	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var c Client

	json.NewDecoder(r.Body).Decode(&c)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO client(firstname, lastname, password, descriptions) VALUES(?, ?, ?, ?)", c.FirstName, c.LastName, c.Password, c.Descriptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.ID = int(lastInsertID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func UpdateClient(w http.ResponseWriter, r *http.Request) {

	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var c Client

	json.NewDecoder(r.Body).Decode(&c)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE client SET firstname=?, lastname=?, password=?, descriptions=? WHERE id=?", c.FirstName, c.LastName, c.Password, c.Descriptions, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.ID = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func DeleteClient(w http.ResponseWriter, r *http.Request) {
	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM client WHERE id=?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
