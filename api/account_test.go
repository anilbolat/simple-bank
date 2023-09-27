package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
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

	testCases := []struct {
		name            string
		accountID       int64
		stubFn          func(store *mockdb.MockStore)
		checkResponseFn func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			stubFn: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponseFn: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				assertAccountInResponse(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			stubFn: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponseFn: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				assertErrorInResponse(t, recorder.Body, "does not exist")
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			stubFn: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponseFn: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				assertErrorInResponse(t, recorder.Body, "error occurred for account ID")
			},
		},
		{
			name:      "InvalidID",
			accountID: 0,
			stubFn: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponseFn: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				assertErrorInResponse(t, recorder.Body, "Field validation for 'ID' failed")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// stub
			tc.stubFn(store)

			// test
			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%d", tc.accountID), nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)

			// assert
			tc.checkResponseFn(t, recorder)
		})
	}
}

func assertAccountInResponse(t *testing.T, resBody *bytes.Buffer, expectedAccount db.Account) {
	require.NotEmpty(t, resBody)
	var actualAccount db.Account
	err := getAccountFromResponseBody(resBody, &actualAccount)
	require.NoError(t, err)
	require.Equal(t, expectedAccount, actualAccount)
}

func assertErrorInResponse(t *testing.T, resBody *bytes.Buffer, expectedError string) {
	var actualError struct {
		Error string `json:"error"`
	}
	err := json.NewDecoder(resBody).Decode(&actualError)
	require.NoError(t, err)
	require.Contains(t, actualError.Error, expectedError)
}
