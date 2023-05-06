package main

import (
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	c "github.com/mertcaliskanlnx/golang-test/Clients"
)

func main() {
	// Router oluşturma
	key := c.DecryptFile()
	fmt.Println(key, "key")
	r := mux.NewRouter()

	// Tüm müşterileri listeleme
	r.HandleFunc("/clients", c.GetClientList).Methods("GET")
	// Yeni müşteri ekleme işlemi
	r.HandleFunc("/client", c.AddClient).Methods("POST")
	// Müşteri güncelleme işlemi
	r.HandleFunc("/client/{id}", c.UpdateClient).Methods("PUT")
	// Müşteri silme işlemi
	r.HandleFunc("/client/{firstname}/{tpmkey}", c.DeleteClient).Methods("DELETE")

	// Sunucuyu başlatma
	http.ListenAndServe(":8080", r)
}
