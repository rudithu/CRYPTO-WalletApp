package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/rudithu/CRYPTO-WalletApp/adapters"
	"github.com/rudithu/CRYPTO-WalletApp/db"
	"github.com/rudithu/CRYPTO-WalletApp/models"
)

// HandleBalance handles the GET request to retrieve a user's wallet balance.
// It supports optional filtering by a specific wallet ID via query parameters.
func (h *HandlerDB) HandleBalance(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from the URL path
	vars := mux.Vars(r)
	userIdStr := vars["id"]

	// Extract optional wallet ID from query parameters
	walletIdStr := r.URL.Query().Get("wallet_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	var walletId int64
	if walletIdStr != "" {
		walletId, err = strconv.ParseInt(walletIdStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid Wallet ID", http.StatusBadRequest)
			return
		}
	}

	// Retrieve user information from the database
	userInfo, err := db.GetUserById(h.DB, userId)
	if err != nil {
		http.Error(w, "Error Getting Wallet Info", http.StatusInternalServerError)
		return
	}

	userIds := make([]int64, 1)
	userIds[0] = userId

	var selectedWallets []models.Wallet
	if walletIdStr == "" {
		// If no wallet ID specified, fetch all wallets for the user
		selectedWallets, err = db.GetWalletByUserIDs(h.DB, userIds)
		if err != nil || selectedWallets == nil {
			http.Error(w, "Error Getting Wallet Info", http.StatusInternalServerError)
			return
		}
	} else {
		// If wallet ID specified, fetch only that wallet
		wallet, err := db.GetWalletById(h.DB, walletId)
		if err != nil || wallet == nil {
			http.Error(w, "Error Getting Wallet Info", http.StatusInternalServerError)
			return
		}
		selectedWallets = []models.Wallet{*wallet}
	}

	// Collect non-base currencies to get conversion rates
	var ccys []string
	for _, sw := range selectedWallets {
		if sw.Currency != models.BaseCcy {
			ccys = append(ccys, sw.Currency)
		}
	}

	// Retrieve currency conversion rates if needed
	ccyMap := make(map[string]models.CcyRateToBaseCcy)
	if len(ccys) > 0 {
		rates, err := db.GetCcyRateToBaseCcy(h.DB, ccys)
		if err != nil {
			http.Error(w, "error to get currency rate", http.StatusInternalServerError)
			return
		}

		for _, rate := range rates {
			ccyMap[rate.Ccy] = rate
		}
	}

	txns := make([]models.Transaction, 0)
	// Convert and format the wallet response
	resp := adapters.ToWalletDetailsResp(userInfo, selectedWallets, txns, ccyMap)

	// Send the response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
