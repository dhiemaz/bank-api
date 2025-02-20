package usecase

import (
	"database/sql"
	"errors"
	"github.com/dhiemaz/bank-api/entities"
	db "github.com/dhiemaz/bank-api/infrastructure/db/sqlc"
	"github.com/dhiemaz/bank-api/utils"
	"github.com/dhiemaz/bank-api/utils/api_error"
	"github.com/dhiemaz/bank-api/utils/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// UserUseCase :
type UserUseCase interface {
	Login(ctx *gin.Context, request entities.LoginUserRequest) (*entities.LoginUserResponse, error)
	UserRegistration(ctx *gin.Context, request entities.CreateUserRequest) (*db.User, error)
	GetUser(ctx *gin.Context, username string) (*db.User, error)
	CheckUserExist(ctx *gin.Context, username string) (*db.User, error)
	UpdateUser(ctx *gin.Context, username string, request entities.UpdateUserRequest) (*db.User, error)
}

type UseCase struct {
	db  db.Querier
	jwt token.JWTMaker
}

func NewUserUseCase(db db.Querier) *UseCase {
	return &UseCase{db: db}
}

// Login : user login
func (user *UseCase) Login(ctx *gin.Context, request entities.LoginUserRequest) (*entities.LoginUserResponse, error) {
	userData, err := user.CheckUserExist(ctx, request.Username)
	if err != nil {
		return nil, err
	}

	// Check User's Password
	err = utils.CheckHashedPassword(userData.HashedPassword, request.Password)
	if err != nil {
		return nil, err
	}

	// Generate New Access Token for User
	accessToken, accessPayload, err := user.jwt.CreateToken(request.Username)
	if err != nil {
		return nil, err
	}

	// Generate New Refresh Token for User
	refreshToken, refreshPayload, err := user.jwt.CreateRefreshToken(request.Username)
	if err != nil {
		return nil, err
	}

	session, err := user.db.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     request.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		ExpiresAt:    refreshPayload.ExpireAt,
	})

	response := entities.LoginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpireAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpireAt,
		User:                  utils.MapUserToResponse(userData),
	}

	return &response, nil
}

// UserRegistration : register a new user
func (user *UseCase) UserRegistration(ctx *gin.Context, request entities.CreateUserRequest) (*db.User, error) {
	hashPassword, err := utils.GenerateHashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	userData, err := user.db.CreateUser(ctx, db.CreateUserParams{
		Username:       request.Username,
		HashedPassword: hashPassword,
		FullName:       request.FullName,
		Email:          request.Email,
	})

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, errors.New("unique violation")
			}
		}
		return nil, err
	}

	return &userData, nil
}

// GetUser : get single user
func (user *UseCase) GetUser(ctx *gin.Context, username string) (*db.User, error) {
	userData, err := user.CheckUserExist(ctx, username)
	if err != nil {
		return nil, err
	}

	return userData, nil
}

// UpdateUser : update existing user data
func (user *UseCase) UpdateUser(ctx *gin.Context, username string, request entities.UpdateUserRequest) (*db.User, error) {

	var params db.UpdateUserParams

	userData, err := user.CheckUserExist(ctx, username)
	if err != nil {
		return nil, err
	}

	if request.Email == userData.Email {
		return nil, api_error.ErrEmailSameAsOld
	}

	if err := utils.CheckHashedPassword(request.OldPassword, userData.HashedPassword); err != nil {
		return nil, api_error.ErrPasswordWrong
	}

	hashedPassword, err := utils.GenerateHashPassword(request.NewPassword)
	if err != nil {
		return nil, err
	}

	params.HashedPassword = sql.NullString{String: hashedPassword, Valid: true}

	dbUser, err := user.db.UpdateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return &dbUser, err
}

// CheckUserExist : check if user is exist in database
func (user *UseCase) CheckUserExist(ctx *gin.Context, username string) (*db.User, error) {
	userData, err := user.db.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return &userData, nil
}
