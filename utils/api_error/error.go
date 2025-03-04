package api_error

import (
	"errors"
	"fmt"
)

var (
	ErrNotAccountOwner         = errors.New("account doesn't belong to authenticated user")
	ErrEmailSameAsOld          = errors.New("new email is the same as old email")
	ErrBlockedRefreshToken     = errors.New("refresh token is blocked")
	ErrMismatchedRefreshTokens = errors.New("refresh token doesn't match with stored refresh token")
	ErrExpiredRefreshToken     = errors.New("refresh token has expired")
	ErrPasswordWrong           = errors.New("old password is different from the one stored in the database")

	ErrSameAccountTransfer = func(from, to int64) error {
		return fmt.Errorf(fmt.Sprintf("can't transfer to the same account, req.FromAccountId=%d, req.ToAccount=%d", from, to))
	}
	ErrCurrencyMismatch = func(from, to string) error {
		return fmt.Errorf(fmt.Sprintf("currency mismatch account1.currency=%s, account2.currency=%s", from, to))
	}

	ErrAccountDeleted = func(id int64) error {
		return fmt.Errorf("account %d is deleted", id)
	}
)
