package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rudithu/CRYPTO-WalletApp/db"
	"github.com/rudithu/CRYPTO-WalletApp/models"
)

// HandleWithdrawMoney handles withdrawal requests from a specific wallet.
// It validates the wallet ID, parses the withdrawal amount, checks the wallet balance,
// and updates the wallet balance accordingly.
func (h *HandlerDB) HandleWithdrawMoney(w http.ResponseWriter, r *http.Request) {
	// Extract wallet ID from URL path variables
	vars := mux.Vars(r)
	walletIdStr := vars["id"]

	walletId, err := strconv.ParseInt(walletIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid Wallet Id", http.StatusBadRequest)
		return
	}

	// Decode the JSON request body into TransactionRequest struct
	var msg models.TransactionRequest
	err = json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Validate the transaction request for withdrawal type
	if err = msg.ValidateRequest(models.TxnTypeWithdraw); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve wallet details from database by wallet ID
	wallet, err := db.GetWalletById(h.DB, walletId)
	if err != nil {
		http.Error(w, "failed to get wallet info", http.StatusInternalServerError)
		return
	}

	if wallet == nil {
		http.Error(w, "wallet not found", http.StatusNotFound)
		return
	}

	// Check if wallet balance is sufficient for the withdrawal amount
	if wallet.Balance.LessThan(msg.Amount) {
		http.Error(w, "withdrawal is not allowed", http.StatusBadRequest)
		return
	}

	t := models.Transaction{
		WalletId: walletId,
		Amount:   msg.Amount,
		Type:     models.TxnTypeWithdraw,
	}

	// Perform the withdrawal update on the database
	err = db.WithdrawUpdate(h.DB, &t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Send HTTP status 204 No Content indicating success with no response body
	w.WriteHeader(http.StatusNoContent)

}
