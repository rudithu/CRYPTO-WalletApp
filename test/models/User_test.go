package models_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rudithu/CRYPTO-WalletApp/models"
	"github.com/rudithu/CRYPTO-WalletApp/test/testutils"
	"github.com/stretchr/testify/assert"
	// Import your actual package here
)

func TestGetUserById_Success(t *testing.T) {
	testutils.WithDBMock(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		expectedID := int64(1)
		expectedName := "Alice"
		expectedTime := time.Now()

		rows := sqlmock.NewRows([]string{"id", "name", "created_at"}).
			AddRow(expectedID, expectedName, expectedTime)

		mock.ExpectQuery("SELECT id, name, created_at FROM users where id=\\$1").
			WithArgs(expectedID).
			WillReturnRows(rows)

		user, err := models.GetUserById(db, expectedID)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedID, user.ID)
		assert.Equal(t, expectedName, user.Name)
		assert.WithinDuration(t, expectedTime, user.CreatedAt, time.Second)
	})
}

func TestGetUserById_NotFound(t *testing.T) {
	testutils.WithDBMock(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectQuery("SELECT id, name, created_at FROM users where id=\\$1").
			WithArgs(int64(2)).
			WillReturnError(sql.ErrNoRows)
		user, err := models.GetUserById(db, 2)
		assert.NoError(t, err)
		assert.Nil(t, user)
	})

}

func TestGetUserById_DBError(t *testing.T) {
	testutils.WithDBMock(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectQuery("SELECT id, name, created_at FROM users where id=\\$1").
			WithArgs(int64(3)).
			WillReturnError(errors.New("db failed"))

		user, err := models.GetUserById(db, 3)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}
