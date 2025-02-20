package utils

import (
	"github.com/dhiemaz/bank-api/entities"
	"github.com/dhiemaz/bank-api/infrastructure/db/sqlc"
)

func MapUserToResponse(user *db.User) entities.UserResponse {
	return entities.UserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		CreatedAt:         user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
	}
}

func MapAccountToResponse(account *db.Account) entities.AccountResponse {
	return entities.AccountResponse{
		ID:        account.ID,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: account.CreatedAt,
	}
}

func MapTransferToResponse(transfer db.Transfer) *entities.TransferResponse {
	return &entities.TransferResponse{
		ID: transfer.ID,
		// FromAccount: transfer.FromAccountID,
		ToAccountID: transfer.ToAccountID,
		// FromEntry:   transfer.FromEntryID,
		Amount:    transfer.Amount,
		CreatedAt: transfer.CreatedAt,
	}
}

func FromTransferTxToTransferResponse(result *db.TransferTxResult) entities.TransferResponse {
	return entities.TransferResponse{
		ID:          result.Transfer.ID,
		FromAccount: result.FromAccount,
		ToAccountID: result.ToAccount.ID,
		FromEntry:   result.FromEntry,
		Amount:      result.Transfer.Amount,
	}
}
