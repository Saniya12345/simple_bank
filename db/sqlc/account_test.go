package db

import (
	"context"
	"database/sql"
	"simplebank/db/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) Account {
	owner := util.RandomOwner()
	balance := util.RandomMoney()
	currency := util.RandomCurrency()
	account, err := testQueries.CreateAccount(context.Background(), owner, int64(balance), currency)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, owner, account.Owner)
	require.EqualValues(t, int64(balance), account.Balance)
	require.Equal(t, currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Owner, account2.Owner)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
	require.Equal(t, account1.Currency, account2.Currency)

}

func TestUpdateAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	bal := util.RandomMoney()
	account2, err := testQueries.UpdateAccount(context.Background(), account1.ID, bal)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, bal, account2.Balance)
	require.Equal(t, account1.Owner, account2.Owner)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
	require.Equal(t, account1.Currency, account2.Currency)
}

func TestDeleteAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)

	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomAccount(t)
	}

	limit := int32(5)
	offset := 5

	accounts, err := testQueries.ListAccounts(context.Background(), limit, int32(offset))

	require.NoError(t, err)

	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
