package usecase

import (
	"database/sql"
	"errors"
	"github.com/dhiemaz/bank-api/entities"
	"github.com/dhiemaz/bank-api/infrastructure/db/sqlc"
	"github.com/dhiemaz/bank-api/infrastructure/logger"
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
				logger.WithFields(logger.Fields{"component": "usecase", "action": "create account", "data": accountData}).
					Errorf("failed to create account due to foreign key violation")

				return nil, errors.New(pqErr.Code.Name())
			}
		}

		logger.WithFields(logger.Fields{"component": "usecase", "action": "create account", "data": accountData}).
			Errorf("failed to create account, error : %v", err)

		return nil, err
	}

	return &accountData, nil
}

func (account *UseCase) IsValidAccount(ctx *gin.Context, accountID int64) (*db.Account, error) {
	accountData, err := account.db.GetAccount(ctx, accountID)
	if err != nil {
		if err != nil {

			logger.WithFields(logger.Fields{"component": "usecase", "action": "check if is valid account with", "account_id": accountID}).
				Errorf("failed get account with id %v, error : %v", accountID, err)

			if errors.Is(err, sql.ErrNoRows) {
				return nil, errors.New("not found account")
			}
			return nil, err
		}
	}

	return &accountData, nil
}

func (account *UseCase) GetDeletedAccounts(ctx *gin.Context, username string) ([]db.Account, error) {
	accountData, err := account.db.GetDeletedAccounts(ctx, username)
	if err != nil {

		logger.WithFields(logger.Fields{"component": "usecase", "action": "get deleted accounts", "username": username}).
			Errorf("failed get deleted account, error : %v", err)

		return nil, err
	}

	return accountData, nil
}

func (account *UseCase) GetAccounts(ctx *gin.Context, username string) ([]db.Account, error) {
	accountData, err := account.db.GetAccounts(ctx, username)
	if err != nil {
		logger.WithFields(logger.Fields{"component": "usecase", "action": "get accounts", "username": username}).
			Errorf("failed get accounts, error : %v", err)

		return nil, err
	}
	return accountData, nil
}

func (account *UseCase) DeleteAccount(ctx *gin.Context, accountID int64) error {
	err := account.db.DeleteAccount(ctx, accountID)
	if err != nil {
		logger.WithFields(logger.Fields{"component": "usecase", "action": "delete account", "account_id": accountID}).
			Errorf("failed delete account, error : %v", err)
	}
	return err
}

func (account *UseCase) RestoreAccount(ctx *gin.Context, accountID int64) error {
	err := account.db.RestoreAccount(ctx, accountID)
	if err != nil {
		logger.WithFields(logger.Fields{"component": "usecase", "action": "restore an account", "account_id": accountID}).
			Errorf("failed restoring an account, error : %v", err)
	}
	return err
}
