package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rudithu/CRYPTO-WalletApp/db"
	"github.com/rudithu/CRYPTO-WalletApp/models"
)

func (h *HandlerDB) HandleWithdrawMoney(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	walletIdStr := vars["id"]

	walletId, err := strconv.ParseInt(walletIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid Wallet Id", http.StatusBadRequest)
		return
	}

	var msg models.TransactionRequest
	err = json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err = msg.ValidateRequest(models.TxnTypeWithdraw); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	wallet, err := db.GetWalletById(h.DB, walletId)
	if err != nil {
		http.Error(w, "failed to get wallet info", http.StatusInternalServerError)
		return
	}

	if wallet == nil {
		http.Error(w, "wallet not found", http.StatusNotFound)
		return
	}

	if wallet.Balance.LessThan(msg.Amount) {
		http.Error(w, "withdrawal is not allowed", http.StatusBadRequest)
		return
	}

	t := models.Transaction{
		WalletId: walletId,
		Amount:   msg.Amount,
		Type:     models.TxnTypeWithdraw,
	}

	err = db.WithdrawUpdate(h.DB, &t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)

}
