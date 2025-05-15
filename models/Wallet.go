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
		rows.Scan(
			&w.ID,
			&w.UserId,
			&w.Balance,
			&w.Currency,
			&w.Type,
			&w.CreatedAt,
		)

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
