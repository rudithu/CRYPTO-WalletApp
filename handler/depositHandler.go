package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rudithu/CRYPTO-WalletApp/db"
	"github.com/rudithu/CRYPTO-WalletApp/models"
)

// HandleDepositMoney processes a deposit request to add money to a specific wallet.
// It validates the wallet ID, parses the request body, validates the transaction request,
// and updates the wallet balance accordingly.
func (h *HandlerDB) HandleDepositMoney(w http.ResponseWriter, r *http.Request) {
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

	if err = msg.ValidateRequest(models.TxnTypeDeposit); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t := models.Transaction{
		WalletId: walletId,
		Amount:   msg.Amount,
		Type:     models.TxnTypeDeposit,
	}

	// Perform the deposit update in the database
	err = db.DepositUpdate(h.DB, &t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Return HTTP 204 No Content to indicate success without body content
	w.WriteHeader(http.StatusNoContent)
}
