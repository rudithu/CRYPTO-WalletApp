package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/rudithu/CRYPTO-WalletApp/models"
	"github.com/shopspring/decimal"
)

func GetDefaultWalletOrCurrencyByUserID(db *sql.DB, userID int64, currency string) ([]models.Wallet, error) {

	query := `
		SELECT id, user_id, balance, currency, type, is_default, created_at
		FROM wallets
		WHERE user_id = $1 
		AND (is_default = TRUE %s)
		ORDER BY created_at DESC
	`

	var rows *sql.Rows
	var err error

	if currency != "" {
		rows, err = db.Query(fmt.Sprintf(query, "OR currency = $2 "), userID, currency)
	} else {
		rows, err = db.Query(fmt.Sprintf(query, ""), userID)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	defer rows.Close()

	var wallets []models.Wallet
	for rows.Next() {
		var w models.Wallet
		err := rows.Scan(
			&w.ID,
			&w.UserId,
			&w.Balance,
			&w.Currency,
			&w.Type,
			&w.IsDefault,
			&w.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, w)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(wallets) == 0 {
		return nil, nil
	}
	return wallets, nil

}

func GetWalletById(db *sql.DB, walletId int64) (*models.Wallet, error) {
	query := `
		SELECT id, user_id, balance, currency, type, is_default, created_at
		FROM wallets
		WHERE id = $1
	`
	var wallet models.Wallet
	err := db.QueryRow(query, walletId).Scan(
		&wallet.ID,
		&wallet.UserId,
		&wallet.Balance,
		&wallet.Currency,
		&wallet.Type,
		&wallet.IsDefault,
		&wallet.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &wallet, nil
}

func GetWalletByUserIDs(db *sql.DB, userIDs []int64) ([]models.Wallet, error) {

	if userIDs == nil {
		return nil, nil
	}

	placeholders := make([]string, len(userIDs))
	args := make([]interface{}, len(userIDs))

	for i, id := range userIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, balance, currency, type, is_default, created_at
		FROM wallets
		WHERE user_id IN (%s)
		ORDER BY created_at DESC
	`, strings.Join(placeholders, ", "))

	rows, err := db.Query(query, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	var wallets []models.Wallet

	for rows.Next() {
		var w models.Wallet
		err := rows.Scan(
			&w.ID,
			&w.UserId,
			&w.Balance,
			&w.Currency,
			&w.Type,
			&w.IsDefault,
			&w.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, w)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return wallets, nil
}

func getWalletBalance(tx *sql.Tx, walletId int64) (*decimal.Decimal, error) {
	query := `SELECT balance FROM wallets WHERE id = $1`

	var balance decimal.Decimal
	err := tx.QueryRow(query, walletId).Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &balance, nil
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
