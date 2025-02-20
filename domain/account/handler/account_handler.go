package handler

import (
	"github.com/dhiemaz/bank-api/domain/account/usecase"
	"github.com/dhiemaz/bank-api/entities"
	"github.com/dhiemaz/bank-api/infrastructure/db/sqlc"
	"github.com/dhiemaz/bank-api/middlewares"
	"github.com/dhiemaz/bank-api/utils"
	"github.com/dhiemaz/bank-api/utils/api_error"
	"github.com/dhiemaz/bank-api/utils/token"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Usecase usecase.AccountUseCase
}

func NewAccountHandler(usecase usecase.AccountUseCase) *Handler {
	return &Handler{
		Usecase: usecase,
	}
}

// CreateAccount godoc
//
//	@Summary		creates a new account for the currently logged-in user
//	@Description	creates a new account for the currently logged-in user
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			body	body		createAccountReq	true	"Account to create"
//	@Success		200		{object}	response.JSON{data=accountResponse}
//	@Failure		400,500	{object}	response.JSON{}
//	@Security		bearerAuth
//	@Router			/accounts [post]
func (account *Handler) CreateAccount(ctx *gin.Context) {
	var request entities.CreateAccountRequest
	if err := utils.ParseBody(ctx, &request); err != nil {
		return
	}

	payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)
	accountData, err := account.Usecase.AccountRegistration(ctx, payload.Username, request)
	if err != nil {
		if err.Error() == "foreign_key_violation" && err.Error() == "unique_violation" {
			ctx.JSON(http.StatusForbidden, entities.Err(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, entities.Err(err))
		return
	}

	ctx.JSON(http.StatusCreated, entities.Success(utils.MapAccountToResponse(accountData)))
}

// GetAccount godoc
//
//	@Summary		gets an account by id
//	@Description	gets an account by id
//	@Tags			accounts
//	@Produce		json
//	@Param			id		path		int64	true	"Account ID"
//	@Success		200		{object}	response.JSON{data=accountResponse}
//	@Failure		400,500	{object}	response.JSON{}
//	@Security		bearerAuth
//	@Router			/accounts/{id} [get]
func (account *Handler) GetAccount(ctx *gin.Context) {
	var request entities.GetAccountRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, entities.Err(err))
		return
	}

	accountData, err := account.Usecase.IsValidAccount(ctx, request.ID)
	if err != nil {
		if err.Error() == "account not found" {
			ctx.JSON(http.StatusNotFound, entities.Err(err))
		}

		ctx.JSON(http.StatusInternalServerError, entities.Err(err))
	}

	if !isUserAccountOwner(ctx, accountData) {
		ctx.JSON(http.StatusUnauthorized, api_error.ErrNotAccountOwner)
		return
	}

	ctx.JSON(http.StatusOK, entities.Success(utils.MapAccountToResponse(accountData)))
}

// GetDeletedAccounts godoc
//
//	@Summary
//	@Description	gets a list of accounts for the currently logged-in user
//	@Tags			accounts
//	@Produce		json
//	@Success		200		{object}	response.JSON{data=[]accountResponse}
//	@Failure		400,500	{object}	response.JSON{}
//	@Security		bearerAuth
//	@Router			/accounts/del [get]
func (account *Handler) GetDeletedAccounts(ctx *gin.Context) {
	payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)
	accounts, err := account.Usecase.GetDeletedAccounts(ctx, payload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.Err(err))
		return
	}

	var resp []entities.AccountResponse
	for _, account := range accounts {
		resp = append(resp, utils.MapAccountToResponse(&account))
	}

	ctx.JSON(http.StatusOK, entities.Success(resp))
}

// GetAccounts godoc
//
//	@Summary		gets a list of accounts for the currently logged-in user
//	@Description	gets a list of accounts for the currently logged-in user
//	@Tags			accounts
//	@Produce		json
//	@Success		200		{object}	response.JSON{data=[]accountResponse}
//	@Failure		400,500	{object}	response.JSON{}
//	@Security		bearerAuth
//	@Router			/accounts [get]
func (account *Handler) GetAccounts(ctx *gin.Context) {
	payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

	accounts, err := account.Usecase.GetAccounts(ctx, payload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.Err(err))
		return
	}

	var accountResponses []entities.AccountResponse
	for _, account := range accounts {
		accountResponses = append(accountResponses, utils.MapAccountToResponse(&account))
	}

	ctx.JSON(http.StatusOK, entities.Success(accountResponses))
}

// DeleteAccount godoc
//
//	@Summary		deletes an account by id for the currently logged-in user
//	@Description	deletes an account by id for the currently logged-in user
//	@Tags			accounts
//	@Produce		json
//	@Param			id		path		int64	true	"Account ID"
//	@Success		200		{object}	response.JSON{data=int64}
//	@Failure		400,500	{object}	response.JSON{}
//	@Security		bearerAuth
//	@Router			/accounts/{id} [delete]
func (account *Handler) DeleteAccount(ctx *gin.Context) {
	var request entities.DeleteAccountRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, entities.Err(err))
		return
	}

	accountData, err := account.Usecase.IsValidAccount(ctx, request.ID)
	if err != nil {
		if err.Error() == "account not found" {
			ctx.JSON(http.StatusNotFound, entities.Err(err))
		}
		ctx.JSON(http.StatusInternalServerError, entities.Err(err))
	}

	if !isUserAccountOwner(ctx, accountData) {
		ctx.JSON(http.StatusUnauthorized, entities.Err(api_error.ErrNotAccountOwner))
		return
	}

	err = account.Usecase.DeleteAccount(ctx, request.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.Err(err))
		return
	}

	ctx.JSON(http.StatusOK, entities.Success(request.ID))
}

// RestoreAccount godoc
//
//	@Summary		deletes an account by id for the currently logged-in user
//	@Description	deletes an account by id for the currently logged-in user
//	@Tags			accounts
//	@Produce		json
//	@Param			id		path		int64	true	"Account ID"
//	@Success		200		{object}	response.JSON{data=int64}
//	@Failure		400,500	{object}	response.JSON{}
//	@Security		bearerAuth
//	@Router			/accounts/res/{id} [patch]
func (account *Handler) RestoreAccount(ctx *gin.Context) {
	var request entities.RestoreAccountRequest
	if err := utils.ParseURI(ctx, &request); err != nil {
		return
	}

	accountData, err := account.Usecase.IsValidAccount(ctx, request.ID)
	if err != nil {
		if err.Error() == "account not found" {
			ctx.JSON(http.StatusNotFound, entities.Err(err))
		}
		ctx.JSON(http.StatusInternalServerError, entities.Err(err))
	}

	if !isUserAccountOwner(ctx, accountData) {
		ctx.JSON(http.StatusUnauthorized, entities.Err(api_error.ErrNotAccountOwner))
		return
	}

	err = account.Usecase.RestoreAccount(ctx, request.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.Err(err))
		return
	}

	ctx.JSON(http.StatusOK, entities.Success(request.ID))
}

func isUserAccountOwner(ctx *gin.Context, account *db.Account) bool {
	payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)
	log.Println(payload.Username, account.Owner)
	return payload.Username == account.Owner
}
