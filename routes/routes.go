package routes

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/rudithu/CRYPTO-WalletApp/handler"
)

func Route(database *sql.DB, r *mux.Router) {
	dbHandler := handler.HandlerDB{DB: database}

	r.HandleFunc("/users/{id}/wallets/balance", dbHandler.HandleBalance).Methods("GET")
	r.HandleFunc("/users/{id}/wallets/transactions", dbHandler.HandleTxHistory).Methods("GET")
	r.HandleFunc("/wallets/{id}/deposit", dbHandler.HandleDepositMoney).Methods("POST")
	r.HandleFunc("/wallets/{id}/withdraw", dbHandler.HandleWithdrawMoney).Methods("POST")
	r.HandleFunc("/wallets/{id}/transfer", dbHandler.HandleTransferMoney).Methods("POST")
}
