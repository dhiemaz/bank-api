package handler

import (
	"errors"
	"github.com/dhiemaz/bank-api/domain/user/usecase"
	"github.com/dhiemaz/bank-api/entities"
	"github.com/dhiemaz/bank-api/middlewares"
	"github.com/dhiemaz/bank-api/utils"
	"github.com/dhiemaz/bank-api/utils/api_error"
	"github.com/dhiemaz/bank-api/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Usecase usecase.UserUseCase
}

func NewUserHandler(usecase usecase.UserUseCase) *Handler {
	return &Handler{
		Usecase: usecase,
	}
}

// Login godoc
//
//	@Summary		Login user and return session
//	@Description	Login user and return session
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		loginUserReq	true	"Login user"
//	@Success		200		{object}	response.JSON{data=loginUserRes}
//	@Failure		400  	{object}	response.JSON{}
//	@Router			/users/login [post]
func (user *Handler) LoginUser(ctx *gin.Context) {
	var req entities.LoginUserRequest
	if err := utils.ParseBody(ctx, &req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, entities.Err(err))
		return
	}

	response, err := user.Usecase.Login(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, entities.Err(err))
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}

// Register godoc
//
//	@Summary		Register user
//	@Description	Register user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		createUserReq	true	"Create user"
//	@Success		200		{object}	response.JSON{data=userResponse}
//	@Failure		409,500	{object}	response.JSON{}
//	@Router			/users/register [post]
func (user *Handler) Register(ctx *gin.Context) {
	var request entities.CreateUserRequest
	if err := utils.ParseBody(ctx, &request); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, entities.Err(err))
		return
	}

	userData, err := user.Usecase.UserRegistration(ctx, request)
	if err != nil {
		if err.Error() == "unique_violation" {
			ctx.JSON(http.StatusConflict, entities.Err(errors.New("username is already taken")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, entities.Err(err))
		return
	}

	ctx.JSON(http.StatusCreated, entities.Success(utils.MapUserToResponse(userData)))
}

// Get User godoc
//
//	@Summary		Get current user info
//	@Description	Get current user info
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	response.JSON{data=userResponse}
//	@Failure		400,500	{object}	response.JSON{}
//	@Security		bearerAuth
//	@Router			/users [get]
func (user *Handler) GetUser(ctx *gin.Context) {
	payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

	userData, err := user.Usecase.GetUser(ctx, payload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.Err(err))
	}

	ctx.JSON(http.StatusOK, entities.Success(utils.MapUserToResponse(userData)))
}

// UpdateUser godoc
//
//	@Summary		Update current user info
//	@Description	Update current user info
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		updateUserReq	true	"Update user"
//	@Success		200		{object}	response.JSON{data=userResponse}
//	@Failure		400,500	{object}	response.JSON{}
//	@Security		bearerAuth
//	@Router			/users [patch]
func (user *Handler) UpdateUser(ctx *gin.Context) {
	var request entities.UpdateUserRequest
	if err := utils.ParseBody(ctx, &request); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, entities.Err(err))
		return
	}

	payload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)
	dbUser, err := user.Usecase.UpdateUser(ctx, payload.Username, request)
	if err != nil {
		if errors.Is(err, api_error.ErrEmailSameAsOld) {
			ctx.JSON(http.StatusBadRequest, entities.Err(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, entities.Err(err))
		}
		return
	}

	ctx.JSON(http.StatusOK, entities.Success(utils.MapUserToResponse(dbUser)))
}
