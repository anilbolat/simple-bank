package api

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	db "github.com/anilbolat/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	err := ctx.ShouldBindJSON(&req) // request is in ctx (gin).
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err)) // writes the status and the data into response
		return
	}

	account, err := server.store.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errNotFound := fmt.Errorf("account ID %d does not exist", req.ID)
			log.Printf("%v", errNotFound.Error())
			ctx.JSON(http.StatusNotFound, errorResponse(errNotFound))
			return
		}

		errServer := fmt.Errorf("error occurred for account ID %d: %w", req.ID, err)
		log.Printf("%v", errServer.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(errServer))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type ListAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccount(ctx *gin.Context) {
	var req ListAccountRequest
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	accounts, err := server.store.ListAccounts(ctx, db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})
	if err != nil {
		errServer := fmt.Errorf("error occurred while listing accounts: %w", err)
		log.Printf("%v", errServer.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(errServer))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
