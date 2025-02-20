package handler

import (
	"errors"
	"github.com/dhiemaz/bank-api/domain/security/usecase"
	"github.com/dhiemaz/bank-api/entities"
	"github.com/dhiemaz/bank-api/utils"
	"github.com/dhiemaz/bank-api/utils/api_error"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	Usecase usecase.AuthUseCase
}

func NewAuthHandler(usecase usecase.AuthUseCase) *Handler {
	return &Handler{
		Usecase: usecase,
	}
}

// RenewAccessToken godoc
//
//	@Summary		renews an access token
//	@Description	renews an access token
//	@Tags			users
//	@Produce		json
//	@Param			body	body		renewAccessTokenReq	true	"Refresh token"
//	@Success		200		{object}	response.JSON{data=renewAccessTokenRes}
//	@Failure		400,500	{object}	response.JSON{}
//	@Router			/users/renew [post]
func (auth *Handler) RenewAccessToken(ctx *gin.Context) {
	var request entities.RenewAccessTokenRequest
	if err := utils.ParseBody(ctx, &request); err != nil {
		return
	}

	accessToken, accessPayload, err := auth.Usecase.RenewToken(ctx, request)
	if err != nil {
		if err.Error() == "invalid token" {
			ctx.JSON(http.StatusNotFound, entities.Err(err))
			return
		}

		if errors.Is(err, api_error.ErrBlockedRefreshToken) && errors.Is(err, api_error.ErrExpiredRefreshToken) && errors.Is(err, api_error.ErrMismatchedRefreshTokens) {
			ctx.JSON(http.StatusUnauthorized, entities.Err(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, entities.Err(err))
		return
	}

	response := entities.RenewAccessTokenResponse{AccessToken: accessToken, AccessTokenExpiresAt: accessPayload.ExpireAt}
	ctx.JSON(http.StatusAccepted, entities.JSON{Data: response})
}
