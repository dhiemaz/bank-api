package entities

import (
	db "github.com/dhiemaz/bank-api/infrastructure/db/sqlc"
	"github.com/google/uuid"
	"time"
)

type AccountResponse struct {
	ID        int64     `json:"id"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type TransferResponse struct {
	ID          int64      `json:"id"`
	FromAccount db.Account `json:"from_account"`
	FromEntry   db.Entry   `json:"from_entry"`
	ToAccountID int64      `json:"to_account_id"`
	Amount      int64      `json:"amount"`
	CreatedAt   time.Time  `json:"created_at"`
}

type UserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	CreatedAt         time.Time `json:"created_at"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}

type LoginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	RefreshToken          string       `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_expires_at"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_expires_at"`
	User                  UserResponse `json:"user"`
}

type RenewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_expires_at"`
}

type JSON struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty" swaggerignore:"true"`
	Error   *Error      `json:"error,omitempty" swaggerignore:"true"`
}

type Error struct {
	Error string `json:"error,omitempty"`
}

func Success(data interface{}) JSON {
	return JSON{Success: true, Data: data}
}

func Err(err error) JSON {
	return JSON{Success: false, Error: &Error{Error: err.Error()}}
}
