package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rudithu/CRYPTO-WalletApp/db"
	"github.com/rudithu/CRYPTO-WalletApp/handler"
)

func main() {

	database := db.Connnect()
	testHandler := handler.TestHandler{DB: database}

	dbHandler := handler.HandlerDB{DB: database}

	http.HandleFunc("/hello", testHandler.HandleHello)
	http.HandleFunc("/echo", testHandler.HandleEcho)
	// http.HandleFunc("/user/balance", dbHandler.HandleBalance)

	r := mux.NewRouter()
	r.HandleFunc("/users/{id}/wallets/balance", dbHandler.HandleBalance).Methods("GET")
	r.HandleFunc("/users/{id}/wallets/transactions", dbHandler.HandleTxHistory).Methods("GET")
	r.HandleFunc("/wallets/{id}/deposit", dbHandler.HandleDepositMoney).Methods("POST")
	r.HandleFunc("/wallets/{id}/withdraw", dbHandler.HandleWithdrawMoney).Methods("POST")

	/*
		GET /users/{id}/transactions?type=deposit
		GET /users/{id}/transactions?wallet_id=10

		GET /wallets/{wallet_id}/transactions

		POST /users/{fromUserId}/transfer - "to_user_id": 456
		POST /wallets/{walletId}/transfer - "to_wallet_id": 987

		POST /users/{userId}/deposit
		POST /wallets/{walletId}/deposit

		POST /users/{userId}/withdraw
		POST /wallets/{walletId}/withdraw
	*/

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
