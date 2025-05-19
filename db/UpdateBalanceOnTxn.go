package db

import (
	"database/sql"
	"fmt"

	"github.com/rudithu/CRYPTO-WalletApp/models"
)

func DepositUpdate(db *sql.DB, txn *models.Transaction) error {
	return withTx(db, func(tx *sql.Tx) error {
		return depositInternal(tx, txn)
	})
}

func WithdrawUpdate(db *sql.DB, txn *models.Transaction) error {
	return withTx(db, func(tx *sql.Tx) error {
		return withdrawInternal(tx, txn)
	})
}

func TransferUpdate(db *sql.DB, srcTxn *models.Transaction, targetTxn *models.Transaction) error {
	return withTx(db, func(tx *sql.Tx) error {
		err := withdrawInternal(tx, srcTxn)
		if err != nil {
			return err
		}
		return depositInternal(tx, targetTxn)
	})
}

func depositInternal(tx *sql.Tx, txn *models.Transaction) error {
	err := createTransaction(tx, txn)
	if err != nil {
		return fmt.Errorf("failed to create incoming-transaction: %w", err)
	}

	err = incrementBalanceByWalletID(tx, txn.WalletId, txn.Amount)
	if err != nil {
		return fmt.Errorf("failed to update incoming-balance: %w", err)
	}

	return nil
}

func withdrawInternal(tx *sql.Tx, txn *models.Transaction) error {
	balance, err := getWalletBalance(tx, txn.WalletId)
	if err != nil || balance == nil {
		return fmt.Errorf("failed to get balance")
	}

	if balance.LessThan(txn.Amount) {
		return fmt.Errorf("not enough balance")
	}

	err = createTransaction(tx, txn)
	if err != nil {
		return fmt.Errorf("failed to create outgoing-transaction: %w", err)
	}

	err = updateBalanceByWalletID(tx, txn.WalletId, balance.Sub(txn.Amount))
	if err != nil {
		return fmt.Errorf("failed to update outgoing-balance: %w", err)
	}
	return nil
}

func withTx(db *sql.DB, fn func(tx *sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after rollback
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return
}
