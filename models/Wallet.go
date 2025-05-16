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
	Currency  string          `json:"currency"`
	Type      string          `json:"currency"`
	Balance   decimal.Decimal `json:"balance"`
	CreatedAt time.Time       `json:"created_at`
}

func GetDefaultWalletByUserID(db *sql.DB, userID int64, walletType string) (*Wallet, error) {

	if walletType == "" {
		walletType = "saving"
	}

	query := `
		SELECT id, user_id, balance, currency, type, created_at
		FROM wallets
		WHERE user_id = $1 AND type = $2
		ORDER BY created_at DESC
	`
	rows, err := db.Query(query, userID, walletType)
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

	if len(wallets) > 1 {
		return nil, fmt.Errorf("expected 1 wallet but found %d", len(wallets))
	}
	if len(wallets) == 0 {
		return nil, nil
	}
	return &wallets[0], nil

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
		SELECT id, user_id, balance, currency, type, created_at
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
