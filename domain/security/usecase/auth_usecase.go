package usecase

import (
	"database/sql"
	"errors"
	"github.com/dhiemaz/bank-api/entities"
	db "github.com/dhiemaz/bank-api/infrastructure/db/sqlc"
	common "github.com/dhiemaz/bank-api/utils/api_error"
	"github.com/dhiemaz/bank-api/utils/token"
	"github.com/gin-gonic/gin"
	"time"
)

// AuthUseCase :
type AuthUseCase interface {
	RenewToken(ctx *gin.Context, request entities.RenewAccessTokenRequest) (string, *token.Payload, error)
}

type UseCase struct {
	db  db.Store
	jwt token.JWTMaker
}

func NewAuthUseCase(db db.Store) *UseCase {
	return &UseCase{db: db}
}

func (auth *UseCase) RenewToken(ctx *gin.Context, request entities.RenewAccessTokenRequest) (string, *token.Payload, error) {

	refreshPayload, err := auth.jwt.VerifyToken(request.RefreshToken)
	if err != nil {
		return "", nil, err
	}

	session, err := auth.db.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil, errors.New("invalid token")
		}

		return "", nil, err
	}

	if session.IsBlocked {
		return "", nil, common.ErrBlockedRefreshToken
	}

	if time.Now().After(session.ExpiresAt) {
		return "", nil, common.ErrExpiredRefreshToken
	}

	if session.Username != refreshPayload.Username {
		return "", nil, common.ErrMismatchedRefreshTokens
	}

	if session.RefreshToken != request.RefreshToken {
		return "", nil, common.ErrMismatchedRefreshTokens
	}

	// Generate New Access Token for User
	accessToken, accessPayload, err := auth.jwt.CreateToken(session.Username)
	if err != nil {
		return "", nil, err
	}

	return accessToken, accessPayload, nil
}
