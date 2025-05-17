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
	CreatedAt            time.Time       `json:"created_at"`
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

	query := fmt.Sprintf(`
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
