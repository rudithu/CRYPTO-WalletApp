package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID        int64           `json:"id"`
	UserId    int64           `json:"user_id"`
	Type      string          `json:"currency"`
	IsDefault bool            `json:"is_default"`
	Currency  string          `json:"currency"`
	Balance   decimal.Decimal `json:"balance"`
	CreatedAt time.Time       `json:"created_at`
}

func GetDefaultWalletOrCurrencyByUserID(db *sql.DB, userID int64, currency string) ([]Wallet, error) {

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

	var wallets []Wallet
	for rows.Next() {
		var w Wallet
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

func GetWalleBalanceById(db *sql.DB, walletId int64) (decimal.Decimal, error) {
	query := `
		SELECT balance
		FROM wallets
		WHERE id = $1
	`
	var balance decimal.Decimal
	err := db.QueryRow(query, walletId).Scan(&balance)
	if err != nil {
		return decimal.Zero, err
	}
	return balance, nil
}

func GetWalletById(db *sql.DB, walletId int64) (*Wallet, error) {
	query := `
		SELECT id, user_id, balance, currency, type, is_default, created_at
		FROM wallets
		WHERE id = $1
	`
	var wallet Wallet
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
		return nil, err
	}
	return &wallet, nil
}

func GetWalletByUserIDs(db *sql.DB, userIDs []int64) ([]Wallet, error) {

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

	var wallets []Wallet

	for rows.Next() {
		var w Wallet
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

func IncrementBalanceByWalletID(db *sql.DB, walletID int64, delta decimal.Decimal) error {
	query := `UPDATE wallets SET balance = balance + $1 WHERE id = $2`
	_, err := db.Exec(query, delta, walletID)
	return err
}

func UpdateBalanceByWalletID(db *sql.DB, walletID int64, newBalance decimal.Decimal) error {
	query := `UPDATE wallets SET balance = $1 WHERE id = $2`
	_, err := db.Exec(query, newBalance, walletID)
	return err
}
