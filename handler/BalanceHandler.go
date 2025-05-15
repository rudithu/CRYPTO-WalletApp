package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rudithu/CRYPTO-WalletApp/adapters"
	"github.com/rudithu/CRYPTO-WalletApp/models"
)

type HandlerDB struct {
	DB *sql.DB
}

func (h *HandlerDB) HandleBalance(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userIdStr := r.URL.Query().Get("id")
	if userIdStr == "" {
		http.Error(w, "User ID not found", http.StatusBadRequest)
		return
	}

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
	}

	userInfo, err := models.GetUserById(h.DB, userId)
	if err != nil {
		http.Error(w, "Error Getting Wallet Info", http.StatusInternalServerError)
		return
	}

	userIds := make([]int64, 1)
	userIds[0] = userId

	wallets, err := models.GetWalletByUserIDs(h.DB, userIds)
	if err != nil {
		http.Error(w, "Error Getting Wallet Info", http.StatusInternalServerError)
		return
	}

	resp := adapters.ToWalletBalanceResponse(userInfo, wallets)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
