package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rudithu/CRYPTO-WalletApp/db"
	"github.com/rudithu/CRYPTO-WalletApp/routes"
)

func main() {

	database, err := db.Connnect()
	if err != nil {
		log.Fatal("failed to connect db")
		return
	}

	r := mux.NewRouter()
	routes.Route(database, r)

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
