package api

import (
	"bytes"
	"encoding/json"

	db "github.com/anilbolat/simple-bank/db/sqlc"
	"github.com/anilbolat/simple-bank/util"
)

func getAccountFromResponseBody(resBody *bytes.Buffer, account *db.Account) error {
	return json.NewDecoder(resBody).Decode(account)
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
