package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID                   int64           `json:"id"`
	WalletId             int64           `json:"wallet_id"`
	Type                 string          `json:"type"`
	Amount               decimal.Decimal `json:"amount"`
	CounterpartyWalletId sql.NullInt64   `json:"counterparty_wallet_id"`
	CreatedAt            time.Time       `json:"created_at`
}

func CreateTransaction(db *sql.DB, t *Transaction) error {
	query := `
		INSERT INTO transactions (wallet_id, type, amount, counterparty_wallet_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	err := db.QueryRow(
		query,
		t.WalletId,
		t.Type,
		t.Amount,
		t.CounterpartyWalletId,
		time.Now(),
	).Scan(&t.ID, &t.CreatedAt)

	return err
}

func GetTransactionsByWalletIDs(db *sql.DB, walletIDs []int64) ([]Transaction, error) {

	if walletIDs == nil {
		return nil, nil
	}

	placeholders := make([]string, len(walletIDs))
	args := make([]interface{}, len(walletIDs))

	for i, id := range walletIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprint(`
        SELECT id, wallet_id, type, amount, counterparty_wallet_id, created_at
        FROM transactions
        WHERE wallet_id in (%s)
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

	var transactions []Transaction

	for rows.Next() {
		var t Transaction
		err := rows.Scan(
			&t.ID,
			&t.WalletId,
			&t.Type,
			&t.Amount,
			&t.CounterpartyWalletId,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
