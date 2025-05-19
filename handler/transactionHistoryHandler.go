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

func (h *HandlerDB) HandleTxHistory(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userIdStr := vars["id"]
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	walletIdStr := r.URL.Query().Get("wallet_id")
	var walletId int64
	if walletIdStr != "" {
		walletId, err = strconv.ParseInt(walletIdStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid wallet id", http.StatusBadRequest)
			return
		}
	}

	userInfo, err := db.GetUserById(h.DB, userId)
	if err != nil {
		http.Error(w, "error getting user info", http.StatusInternalServerError)
		return
	}

	var selectedWallets []models.Wallet
	if walletIdStr == "" {
		selectedWallets, err = db.GetWalletByUserIDs(h.DB, []int64{userId})
		if err != nil || selectedWallets == nil {
			http.Error(w, "error getting wallet info", http.StatusInternalServerError)
			return
		}
	} else {
		wallet, err := db.GetWalletById(h.DB, walletId)
		if err != nil || wallet == nil {
			http.Error(w, "error getting wallet info", http.StatusInternalServerError)
			return
		}
		selectedWallets = []models.Wallet{*wallet}
	}

	var walletIds []int64
	for _, w := range selectedWallets {
		walletIds = append(walletIds, w.ID)
	}

	var txns []models.Transaction
	if walletIds != nil {
		transactions, err := db.GetTransactionsByWalletIDs(h.DB, walletIds)
		if err != nil {
			http.Error(w, "Error Getting Transaction Details", http.StatusInternalServerError)
			return
		}
		txns = transactions
	}

	var ccys []string
	for _, sw := range selectedWallets {
		if sw.Currency != models.BaseCcy {
			ccys = append(ccys, sw.Currency)
		}
	}

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

	resp := adapters.ToWalletDetailsResp(userInfo, selectedWallets, txns, ccyMap)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
