package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rudithu/CRYPTO-WalletApp/db"
	"github.com/rudithu/CRYPTO-WalletApp/models"
	"github.com/shopspring/decimal"
)

// HandleTransferMoney handles transferring money from one wallet to another.
// It supports transfers either by specifying a destination wallet ID or a destination user ID.
func (h *HandlerDB) HandleTransferMoney(w http.ResponseWriter, r *http.Request) {
	// Extract wallet ID from the URL path variables
	vars := mux.Vars(r)
	walletIdStr := vars["id"]

	// Validate and parse wallet ID to int64
	walletId, err := strconv.ParseInt(walletIdStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid wallet id", http.StatusBadRequest)
		return
	}

	// Decode the JSON request body into TransactionRequest struct
	var msg models.TransactionRequest
	err = json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Validate the request payload according to transfer out transaction type
	if err = msg.ValidateRequest(models.TxnTypeTransferOut); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the source wallet from database
	sourceWallet, err := db.GetWalletById(h.DB, walletId)
	if err != nil {
		http.Error(w, "failed to get source wallet info", http.StatusBadRequest)
		return
	}

	// Check if source wallet exists
	if sourceWallet == nil {
		http.Error(w, "source wallet not found", http.StatusNotFound)
		return
	}

	// Validate that source wallet has enough balance for the transfer amount
	if sourceWallet.Balance.LessThan(msg.Amount) {
		http.Error(w, "transferred is not allowed", http.StatusBadRequest)
		return
	}

	var targetWallet models.Wallet

	// Determine the target wallet based on provided request data
	if msg.DestinationUserID != nil {
		// Transfer to another user's wallet

		// Prevent transferring to own wallet using user ID (must use wallet ID instead)
		if sourceWallet.UserId == *msg.DestinationUserID {
			http.Error(w, "please use destination_wallet_id to transfer for the same user", http.StatusBadRequest)
			return
		}

		// Get default wallet(s) or wallets with matching currency for the target user
		targetWallets, err := db.GetDefaultWalletOrCurrencyByUserID(h.DB, *msg.DestinationUserID, sourceWallet.Currency)
		if err != nil {
			http.Error(w, "failed to get target wallet", http.StatusBadRequest)
			return
		}

		// Pick the wallet with matching currency or the default wallet if none match
		for _, w := range targetWallets {
			if sourceWallet.Currency == w.Currency {
				targetWallet = w
				break
			} else if w.IsDefault {
				targetWallet = w
			}
		}

	} else if msg.DestinationWalletID != nil {
		// Transfer to a specific wallet by wallet ID
		tWallet, err := db.GetWalletById(h.DB, *msg.DestinationWalletID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to find target wallet %d", msg.DestinationWalletID), http.StatusBadRequest)
			return
		}
		targetWallet = *tWallet
	}

	// Create transaction record for transfer out from source wallet
	txnOut := models.Transaction{
		WalletId:             walletId,
		Type:                 models.TxnTypeTransferOut,
		Amount:               msg.Amount,
		CounterpartyWalletId: sql.NullInt64{Int64: targetWallet.ID, Valid: true},
	}

	var targetAmount decimal.Decimal

	// Calculate target amount considering currency conversion if necessary
	if sourceWallet.Currency == targetWallet.Currency {
		targetAmount = msg.Amount
	} else {
		rate, err := db.GetCcyRate(h.DB, sourceWallet.Currency, targetWallet.Currency)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		targetAmount = msg.Amount.Mul(rate)
	}

	// Create transaction record for transfer in to target wallet
	txnIn := models.Transaction{
		WalletId:             targetWallet.ID,
		Type:                 models.TxnTypeTransferIn,
		Amount:               targetAmount,
		CounterpartyWalletId: sql.NullInt64{Int64: walletId, Valid: true},
	}

	// Perform the transfer update atomically in the database
	err = db.TransferUpdate(h.DB, &txnOut, &txnIn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Respond with no content status indicating success
	w.WriteHeader(http.StatusNoContent)

}
