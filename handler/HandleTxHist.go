package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/rudithu/CRYPTO-WalletApp/adapters"
	"github.com/rudithu/CRYPTO-WalletApp/models"
)

func (h *HandlerDB) HandleTxHistory(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userIdStr := vars["id"]
	walletIdStr := r.URL.Query().Get("wallet_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
	}

	userInfo, err := models.GetUserById(h.DB, userId)
	if err != nil {
		http.Error(w, "Error Getting Wallet Info", http.StatusInternalServerError)
		return
	}

	var userIds []int64 = make([]int64, 1)
	userIds[0] = userId
	wallets, err := models.GetWalletByUserIDs(h.DB, userIds)
	if err != nil {
		http.Error(w, "Error Getting Wallet Info", http.StatusInternalServerError)
		return
	}

	var selectedWallets []models.Wallet
	if walletIdStr != "" {
		walletId, err := strconv.ParseInt(walletIdStr, 10, 64)
		if err == nil {
			for _, w := range wallets {
				if walletId == w.ID {
					selectedWallets = []models.Wallet{w}
					break
				}
			}
		}
	} else {
		selectedWallets = wallets
	}

	var walletIds []int64
	for _, w := range selectedWallets {
		walletIds = append(walletIds, w.ID)
	}

	var txns []models.Transaction
	if walletIds != nil {
		transactions, err := models.GetTransactionsByWalletIDs(h.DB, walletIds)
		if err != nil {
			http.Error(w, "Error Getting Transaction Details", http.StatusInternalServerError)
			return
		}
		txns = transactions
	}

	resp := adapters.ToWalletDetailsResp(userInfo, selectedWallets, txns)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
