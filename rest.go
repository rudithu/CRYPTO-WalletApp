package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rudithu/CRYPTO-WalletApp/db"
	"github.com/rudithu/CRYPTO-WalletApp/handler"
)

func main() {

	database := db.Connnect()
	testHandler := handler.TestHandler{DB: database}

	dbHandler := handler.HandlerDB{DB: database}

	http.HandleFunc("/hello", testHandler.HandleHello)
	http.HandleFunc("/echo", testHandler.HandleEcho)
	http.HandleFunc("/user/balance", dbHandler.HandleBalance)

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
