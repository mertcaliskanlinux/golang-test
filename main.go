package main

import (
	"net/http"

	c "github.com/mertcaliskanlnx/Clients"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	// Router oluşturma

	r := mux.NewRouter()

	// Tüm müşterileri listeleme
	r.HandleFunc("/clients", c.GetClientList).Methods("GET")
	// Yeni müşteri ekleme işlemi
	r.HandleFunc("/client", c.AddClient).Methods("POST")
	// Müşteri güncelleme işlemi
	r.HandleFunc("/client/{id}", c.UpdateClient).Methods("PUT")
	// Müşteri silme işlemi
	r.HandleFunc("/client/{id}", c.DeleteClient).Methods("DELETE")

	// Sunucuyu başlatma
	http.ListenAndServe(":8080", r)
}
