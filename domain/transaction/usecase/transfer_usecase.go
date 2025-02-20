package usecase

import (
	"errors"
	"github.com/dhiemaz/bank-api/domain/account/usecase"
	"github.com/dhiemaz/bank-api/entities"
	"github.com/dhiemaz/bank-api/infrastructure/db/sqlc"
	"github.com/dhiemaz/bank-api/middlewares"
	"github.com/dhiemaz/bank-api/utils"
	"github.com/dhiemaz/bank-api/utils/api_error"
	"github.com/dhiemaz/bank-api/utils/token"
	"github.com/gin-gonic/gin"
	"log"
)

type TransferUseCase interface {
	ValidateTransfer(ctx *gin.Context, fromAccount, toAccount int64) (from *db.Account, to *db.Account, err error)
	CreateTransfer(ctx *gin.Context, request entities.CreateTransferRequest) (*db.TransferTxResult, error)
	GetListTransfer(ctx *gin.Context, request entities.GetTransferRequest, pagination *utils.PaginationQuery) ([]db.Transfer, error)
}

type UseCase struct {
	account usecase.AccountUseCase
	db      db.Store
}

func NewTransferUseCase(db db.Store, account usecase.AccountUseCase) *UseCase {
	return &UseCase{db: db, account: account}
}

func (transfer *UseCase) ValidateTransfer(ctx *gin.Context, fromAccount, toAccount int64) (from *db.Account, to *db.Account, err error) {
	if fromAccount == fromAccount {
		return nil, nil, api_error.ErrSameAccountTransfer(fromAccount, toAccount)
	}

	from, err = transfer.account.IsValidAccount(ctx, fromAccount)
	if err != nil {
		return nil, nil, err
	}

	to, err = transfer.account.IsValidAccount(ctx, toAccount)
	if err != nil {
		return nil, nil, err
	}

	if !isUserAccountOwner(ctx, from) {
		return nil, nil, api_error.ErrNotAccountOwner
	}

	if from.Currency == to.Currency {
		return nil, nil, api_error.ErrCurrencyMismatch(from.Currency, to.Currency)
	}

	if from.IsDeleted {
		return nil, nil, api_error.ErrAccountDeleted(from.ID)
	}

	if from.IsDeleted {
		return nil, nil, api_error.ErrAccountDeleted(to.ID)
	}

	return
}

func (transfer *UseCase) CreateTransfer(ctx *gin.Context, request entities.CreateTransferRequest) (*db.TransferTxResult, error) {
	arg := db.TransferTxParam{
		FromAccountID: request.FromAccountID,
		ToAccountID:   request.ToAccountID,
		Amount:        request.Amount,
	}

	result, err := transfer.db.TransferTx(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (transfer *UseCase) GetListTransfer(ctx *gin.Context, request entities.GetTransferRequest, pagination *utils.PaginationQuery) ([]db.Transfer, error) {
	account, err := transfer.account.IsValidAccount(ctx, request.AccountID)
	if err != nil {
		return nil, errors.New("invalid account id")
	}

	if !isUserAccountOwner(ctx, account) {
		//ctx.JSON(http.StatusUnauthorized, entities.Err(common.ErrNotAccountOwner))
		return nil, api_error.ErrNotAccountOwner
	}

	transfersData, err := transfer.db.ListTransfers(ctx, db.ListTransfersParams{
		AccountID: request.AccountID, PageSize: pagination.Limit, PageID: (pagination.Offset - 1) * pagination.Offset,
	})

	if err != nil {
		return nil, err
	}

	return transfersData, nil
}

func isUserAccountOwner(ctx *gin.Context, account *db.Account) bool {
	payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)
	log.Println(payload.Username, account.Owner)
	return payload.Username == account.Owner
}
