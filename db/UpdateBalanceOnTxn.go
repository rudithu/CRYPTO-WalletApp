package db

import (
	"database/sql"
	"fmt"

	"github.com/rudithu/CRYPTO-WalletApp/models"
	"github.com/shopspring/decimal"
)

func DepositUpdate(db *sql.DB, txn *models.Transaction) error {

	return withTx(db, func(tx *sql.Tx) error {
		err := createTransaction(tx, txn)
		if err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		err = incrementBalanceByWalletID(tx, txn.WalletId, txn.Amount)
		if err != nil {
			return fmt.Errorf("failed to update balance: %w", err)
		}

		return nil
	})
}

func WithdrawUpdate(db *sql.DB, initialBalance decimal.Decimal, txn *models.Transaction) error {

	return withTx(db, func(tx *sql.Tx) error {
		err := createTransaction(tx, txn)
		if err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		err = updateBalanceByWalletID(tx, txn.WalletId, initialBalance.Sub(txn.Amount))
		if err != nil {
			return fmt.Errorf("failed to update balance: %w", err)
		}
		return nil
	})

}

func TransferUpdate(db *sql.DB, initialSrcBalance decimal.Decimal, srcTxn *models.Transaction, targetTxn *models.Transaction) error {
	return withTx(db, func(tx *sql.Tx) error {

		err := createTransaction(tx, srcTxn)
		if err != nil {
			return fmt.Errorf("failed to create outgoing transaction: %w", err)
		}
		err = updateBalanceByWalletID(tx, srcTxn.WalletId, initialSrcBalance.Sub(srcTxn.Amount))
		if err != nil {
			return fmt.Errorf("failed to update outgoing balance: %w", err)
		}
		err = createTransaction(tx, targetTxn)
		if err != nil {
			return fmt.Errorf("failed to create incoming transaction: %w", err)
		}
		err = incrementBalanceByWalletID(tx, targetTxn.WalletId, targetTxn.Amount)
		if err != nil {
			return fmt.Errorf("failed to update incoming balance: %w", err)
		}
		return nil

	})

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

func createTransaction(tx *sql.Tx, t *models.Transaction) error {
	query := `
		INSERT INTO transactions (wallet_id, type, amount, counterparty_wallet_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	err := tx.QueryRow(
		query,
		t.WalletId,
		t.Type,
		t.Amount,
		t.CounterpartyWalletId,
	).Scan(&t.ID, &t.CreatedAt)

	return err
}

func incrementBalanceByWalletID(tx *sql.Tx, walletID int64, delta decimal.Decimal) error {
	query := `UPDATE wallets SET balance = balance + $1 WHERE id = $2`
	_, err := tx.Exec(query, delta, walletID)
	return err
}

func updateBalanceByWalletID(tx *sql.Tx, walletID int64, newBalance decimal.Decimal) error {
	query := `UPDATE wallets SET balance = $1 WHERE id = $2`
	_, err := tx.Exec(query, newBalance, walletID)
	return err
}
