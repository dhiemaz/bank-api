package usecase

import (
	"database/sql"
	"errors"
	"github.com/dhiemaz/bank-api/entities"
	"github.com/dhiemaz/bank-api/infrastructure/db/sqlc"
	"github.com/dhiemaz/bank-api/infrastructure/logger"
	"github.com/dhiemaz/bank-api/utils/api_error"
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

		logger.WithFields(logger.Fields{"component": "usecase", "action": "renew token", "payload": request}).
			Errorf("failed renew token, error : %v", err)

		if errors.Is(err, sql.ErrNoRows) {
			return "", nil, errors.New("invalid token")
		}

		return "", nil, err
	}

	if session.IsBlocked {

		logger.WithFields(logger.Fields{"component": "usecase", "action": "renew token", "payload": request, "session": session}).
			Errorf("failed renew token, error : session is blocked")

		return "", nil, api_error.ErrBlockedRefreshToken
	}

	if time.Now().After(session.ExpiresAt) {

		logger.WithFields(logger.Fields{"component": "usecase", "action": "renew token", "payload": request, "session": session}).
			Errorf("failed renew token, error : session is expired refresh token")

		return "", nil, api_error.ErrExpiredRefreshToken
	}

	if session.Username != refreshPayload.Username {

		logger.WithFields(logger.Fields{"component": "usecase", "action": "renew token", "payload": request, "session": session}).
			Errorf("failed renew token, error : username is not match")

		return "", nil, api_error.ErrMismatchedRefreshTokens
	}

	if session.RefreshToken != request.RefreshToken {

		logger.WithFields(logger.Fields{"component": "usecase", "action": "renew token", "payload": request, "session": session}).
			Errorf("failed renew token, error : refresh token is not match")

		return "", nil, api_error.ErrMismatchedRefreshTokens
	}

	// Generate New Access Token for User
	accessToken, accessPayload, err := auth.jwt.CreateToken(session.Username)
	if err != nil {

		logger.WithFields(logger.Fields{"component": "usecase", "action": "renew token", "payload": request, "session": session}).
			Errorf("failed renew token, error : %v", err)

		return "", nil, err
	}

	return accessToken, accessPayload, nil
}
