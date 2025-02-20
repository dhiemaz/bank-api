package usecase

import (
	"database/sql"
	"errors"
	"github.com/dhiemaz/bank-api/entities"
	"github.com/dhiemaz/bank-api/infrastructure/db/sqlc"
	"github.com/dhiemaz/bank-api/utils/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// AccountUseCase :
type AccountUseCase interface {
	AccountRegistration(ctx *gin.Context, username string, request entities.CreateAccountRequest) (*db.Account, error)
	IsValidAccount(ctx *gin.Context, accountID int64) (*db.Account, error)
	GetDeletedAccounts(ctx *gin.Context, username string) ([]db.Account, error)
	GetAccounts(ctx *gin.Context, username string) ([]db.Account, error)
	DeleteAccount(ctx *gin.Context, accountID int64) error
	RestoreAccount(ctx *gin.Context, accountID int64) error
}

type UseCase struct {
	db  db.Querier
	jwt token.JWTMaker
}

func NewAccountUseCase(db db.Querier) *UseCase {
	return &UseCase{db: db}
}

func (account *UseCase) AccountRegistration(ctx *gin.Context, username string, request entities.CreateAccountRequest) (*db.Account, error) {
	accountData, err := account.db.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    username,
		Balance:  1000,
		Currency: request.Currency,
	})

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				//ctx.JSON(http.StatusForbidden, entities.Err(err))
				return nil, errors.New(pqErr.Code.Name())
			}
		}
		//ctx.JSON(http.StatusInternalServerError, entities.Err(err))
		return nil, err
	}

	return &accountData, nil
}

func (account *UseCase) IsValidAccount(ctx *gin.Context, accountID int64) (*db.Account, error) {
	accountData, err := account.db.GetAccount(ctx, accountID)
	if err != nil {
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				//ctx.JSON(http.StatusNotFound, entities.Err(err))
				return nil, errors.New("not found account")
			}

			//ctx.JSON(http.StatusInternalServerError, entities.Err(err))
			return nil, err
		}
	}

	return &accountData, nil
}

func (account *UseCase) GetDeletedAccounts(ctx *gin.Context, username string) ([]db.Account, error) {
	accountData, err := account.db.GetDeletedAccounts(ctx, username)
	if err != nil {
		return nil, err
	}

	return accountData, nil
}

func (account *UseCase) GetAccounts(ctx *gin.Context, username string) ([]db.Account, error) {
	accountData, err := account.db.GetAccounts(ctx, username)
	if err != nil {
		return nil, err
	}
	return accountData, nil
}

func (account *UseCase) DeleteAccount(ctx *gin.Context, accountID int64) error {
	err := account.db.DeleteAccount(ctx, accountID)
	return err
}

func (account *UseCase) RestoreAccount(ctx *gin.Context, accountID int64) error {
	err := account.db.RestoreAccount(ctx, accountID)
	return err
}
