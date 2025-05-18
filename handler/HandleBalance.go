package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/rudithu/CRYPTO-WalletApp/adapters"
	"github.com/rudithu/CRYPTO-WalletApp/models"
)

func (h *HandlerDB) HandleBalance(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userIdStr := vars["id"]
	walletIdStr := r.URL.Query().Get("wallet_id")

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
	}

	userInfo, err := models.GetUserById(h.DB, userId)
	if err != nil {

		fmt.Println("Error Here 1", err.Error())

		http.Error(w, "Error Getting Wallet Info", http.StatusInternalServerError)
		return
	}

	userIds := make([]int64, 1)
	userIds[0] = userId

	wallets, err := models.GetWalletByUserIDs(h.DB, userIds)
	if err != nil {

		fmt.Println("Error Here 2", err.Error())

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

	var ccys []string
	for _, sw := range selectedWallets {
		if sw.Currency != models.BaseCcy {
			ccys = append(ccys, sw.Currency)
		}
	}

	ccyMap := make(map[string]models.CcyRateToBaseCcy)
	if len(ccys) > 0 {
		rates, err := models.GetCcyRateToBaseCcy(h.DB, ccys)
		if err != nil {

			fmt.Println("Error Here 3", err.Error())

			http.Error(w, "error to get currency rate", http.StatusInternalServerError)
			return
		}

		for _, rate := range rates {
			ccyMap[rate.Ccy] = rate
		}

	}

	txns := make([]models.Transaction, 0)
	resp := adapters.ToWalletDetailsResp(userInfo, selectedWallets, txns, ccyMap)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
