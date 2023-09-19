package api

import (
	"encoding/json"
	"fmt"
	mockdb "github.com/anilbolat/simple-bank/db/mock"
	db "github.com/anilbolat/simple-bank/db/sqlc"
	"github.com/anilbolat/simple-bank/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccountAPI(t *testing.T) {
	// given
	account := randomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := mockdb.NewMockStore(ctrl)

	// stubs
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	// test
	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)

	// assert
	require.Equal(t, http.StatusOK, recorder.Code)
	require.NotEmpty(t, recorder.Body)
	var actualAccount db.Account
	err = json.NewDecoder(recorder.Body).Decode(&actualAccount)
	require.NoError(t, err)
	require.Equal(t, account, actualAccount)
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
