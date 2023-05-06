package Clients

import (
	"crypto/rand"
	"crypto/sha256"
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
	TpmKey       string `json:"tpmkey"`
}

func KeyGenerateSH256() (string, error) {

	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("error:", err)
		return err.Error(), err
	}

	h := sha256.New() //machine language genarete
	fmt.Println(h)    //machine language result
	h.Write(b)        //machine language write b variable added 32 byte random number and string

	return fmt.Sprintf("%x", h.Sum(nil)), err
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

	rows, err := db.Query("SELECT * FROM Clients")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	clients := []Client{}

	for rows.Next() {
		var c Client
		if err := rows.Scan(&c.ID, &c.FirstName, &c.LastName, &c.Password, &c.Descriptions, &c.TpmKey); err != nil {
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

	tpm_key, _ := KeyGenerateSH256()

	c := Client{
		FirstName:    r.FormValue("firstname"),
		LastName:     r.FormValue("lastname"),
		Password:     r.FormValue("password"),
		Descriptions: r.FormValue("descriptions"),
		TpmKey:       r.FormValue("tpmkey"),
	}

	json.NewDecoder(r.Body).Decode(&c)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO Clients (firstname, lastname, password, descriptions,tpmkey) VALUES(?, ?, ?, ?, ?)", c.FirstName, c.LastName, c.Password, c.Descriptions, tpm_key)
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

	_, err = db.Exec("UPDATE Clients SET firstname=?, lastname=?, password=?, descriptions=? WHERE id=?", c.FirstName, c.LastName, c.Password, c.Descriptions, id)
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

	firstname := mux.Vars(r)["firstname"]
	tpm_key := mux.Vars(r)["tpmkey"]

	fmt.Println(firstname, "asd", tpm_key)
	result, err := db.Exec("DELETE FROM Clients WHERE firstname=? AND tpmkey=?", firstname, tpm_key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "No rows found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
