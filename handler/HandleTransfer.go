package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rudithu/CRYPTO-WalletApp/models"
	"github.com/shopspring/decimal"
)

func (h *HandlerDB) HandleTransferMoney(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	walletIdStr := vars["id"]

	//validate wallet id
	walletId, err := strconv.ParseInt(walletIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid Wallet Id", http.StatusBadRequest)
	}

	//parse request body
	var msg models.TransactionRequest
	err = json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	//validate request body
	if err = msg.ValidateRequest(models.TxnTypeTransferOut); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sourceWallet, err := models.GetWalletById(h.DB, walletId)
	if err != nil {
		http.Error(w, "failed to get source wallet info", http.StatusBadRequest)
		return
	}

	//get current wallet balance and validate
	if sourceWallet.Balance.LessThan(msg.Amount) {
		http.Error(w, "transferred is not allowed", http.StatusBadRequest)
		return
	}

	var targetWallet models.Wallet
	if msg.DestinationUserID != nil {
		//transfer money to another user id

		if sourceWallet.UserId == *msg.DestinationUserID {
			http.Error(w, "please use destination_wallet_id to transfer for the same user", http.StatusBadRequest)
			return
		}

		targetWallets, err := models.GetDefaultWalletOrCurrencyByUserID(h.DB, *msg.DestinationUserID, sourceWallet.Currency)
		if err != nil {
			http.Error(w, "failed to get target wallet", http.StatusBadRequest)
			return
		}

		for _, w := range targetWallets {
			if sourceWallet.Currency == w.Currency {
				targetWallet = w
				break
			} else if w.IsDefault {
				targetWallet = w
			}
		}

	} else if msg.DestinationWalletID != nil {
		//transfer money to another wallet id
		tWallet, err := models.GetWalletById(h.DB, *msg.DestinationWalletID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to find target wallet %d", msg.DestinationWalletID), http.StatusBadRequest)
			return
		}
		targetWallet = *tWallet
	}

	txnOut := models.Transaction{
		WalletId:             walletId,
		Type:                 models.TxnTypeTransferOut,
		Amount:               msg.Amount,
		CounterpartyWalletId: sql.NullInt64{Int64: targetWallet.ID, Valid: true},
	}

	var targetAmount decimal.Decimal

	if sourceWallet.Currency == targetWallet.Currency {
		targetAmount = msg.Amount
	} else {
		rate, err := models.GetCcyRate(h.DB, sourceWallet.Currency, targetWallet.Currency)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		targetAmount = msg.Amount.Mul(rate)
	}

	txnIn := models.Transaction{
		WalletId:             targetWallet.ID,
		Type:                 models.TxnTypeTransferIn,
		Amount:               targetAmount,
		CounterpartyWalletId: sql.NullInt64{Int64: walletId, Valid: true},
	}

	models.CreateTransaction(h.DB, &txnOut)
	models.UpdateBalanceByWalletID(h.DB, walletId, sourceWallet.Balance.Sub(msg.Amount))

	models.CreateTransaction(h.DB, &txnIn)
	models.IncrementBalanceByWalletID(h.DB, targetWallet.ID, targetAmount)

	w.WriteHeader(http.StatusNoContent)

}
