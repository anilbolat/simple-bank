package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/anilbolat/simple-bank/db/mock"
	db "github.com/anilbolat/simple-bank/db/sqlc"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
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
	assertResponse(t, recorder.Body, account)
}

func assertResponse(t *testing.T, resBody *bytes.Buffer, expectedAccount db.Account) {
	require.NotEmpty(t, resBody)
	var actualAccount db.Account
	err := getAccountFromResponseBody(resBody, &actualAccount)
	require.NoError(t, err)
	require.Equal(t, expectedAccount, actualAccount)
}
