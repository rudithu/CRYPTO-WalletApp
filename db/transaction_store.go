package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/rudithu/CRYPTO-WalletApp/models"
)

func GetTransactionsByWalletIDs(db *sql.DB, walletIDs []int64) ([]models.Transaction, error) {

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

	var transactions []models.Transaction

	for rows.Next() {
		var t models.Transaction
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
