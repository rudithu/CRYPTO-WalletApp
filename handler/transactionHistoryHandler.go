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

// HandleTxHistory handles the request to fetch a user's wallet transaction history.
// It supports optional filtering by wallet ID via query parameter.
func (h *HandlerDB) HandleTxHistory(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL path variables
	vars := mux.Vars(r)
	userIdStr := vars["id"]

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	// Get optional wallet_id query parameter
	walletIdStr := r.URL.Query().Get("wallet_id")
	var walletId int64
	if walletIdStr != "" {
		walletId, err = strconv.ParseInt(walletIdStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid wallet id", http.StatusBadRequest)
			return
		}
	}

	// Fetch user information from the database
	userInfo, err := db.GetUserById(h.DB, userId)
	if err != nil {
		http.Error(w, "error getting user info", http.StatusInternalServerError)
		return
	}

	var selectedWallets []models.Wallet
	if walletIdStr == "" {
		// If wallet ID not specified, get all wallets for the user
		selectedWallets, err = db.GetWalletByUserIDs(h.DB, []int64{userId})
		if err != nil || selectedWallets == nil {
			http.Error(w, "error getting wallet info", http.StatusInternalServerError)
			return
		}
	} else {
		// If wallet ID specified, fetch that specific wallet
		wallet, err := db.GetWalletById(h.DB, walletId)
		if err != nil || wallet == nil {
			http.Error(w, "error getting wallet info", http.StatusInternalServerError)
			return
		}
		selectedWallets = []models.Wallet{*wallet}
	}

	// Collect all wallet IDs to query transactions
	var walletIds []int64
	for _, w := range selectedWallets {
		walletIds = append(walletIds, w.ID)
	}

	// Retrieve transactions for the selected wallets
	var txns []models.Transaction
	if walletIds != nil {
		transactions, err := db.GetTransactionsByWalletIDs(h.DB, walletIds)
		if err != nil {
			http.Error(w, "Error Getting Transaction Details", http.StatusInternalServerError)
			return
		}
		txns = transactions
	}

	// Collect currencies from selected wallets that are not the base currency
	var ccys []string
	for _, sw := range selectedWallets {
		if sw.Currency != models.BaseCcy {
			ccys = append(ccys, sw.Currency)
		}
	}

	ccyMap := make(map[string]models.CcyRateToBaseCcy)
	if len(ccys) > 0 {
		// Fetch currency rates from the database
		rates, err := db.GetCcyRateToBaseCcy(h.DB, ccys)
		if err != nil {
			http.Error(w, "error to get currency rate", http.StatusInternalServerError)
			return
		}

		for _, rate := range rates {
			ccyMap[rate.Ccy] = rate
		}

	}

	// Prepare the response using adapter to convert DB models into response format
	resp := adapters.ToWalletDetailsResp(userInfo, selectedWallets, txns, ccyMap)

	// Set response content type to JSON and write the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
