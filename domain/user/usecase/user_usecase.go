package usecase

import (
	"database/sql"
	"errors"
	"github.com/dhiemaz/bank-api/entities"
	db "github.com/dhiemaz/bank-api/infrastructure/db/sqlc"
	"github.com/dhiemaz/bank-api/infrastructure/logger"
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

		logger.WithFields(logger.Fields{"component": "usecase", "action": "user login", "payload": request}).
			Errorf("failed user login, err : %v", err)

		return nil, err
	}

	// Check User's Password
	err = utils.CheckHashedPassword(userData.HashedPassword, request.Password)
	if err != nil {
		logger.WithFields(logger.Fields{"component": "usecase", "action": "user login", "payload": request}).
			Errorf("failed check hashed password login, err : %v", err)

		return nil, err
	}

	// Generate New Access Token for User
	accessToken, accessPayload, err := user.jwt.CreateToken(request.Username)
	if err != nil {

		logger.WithFields(logger.Fields{"component": "usecase", "action": "user login", "payload": request}).
			Errorf("failed create token, err : %v", err)

		return nil, err
	}

	// Generate New Refresh Token for User
	refreshToken, refreshPayload, err := user.jwt.CreateRefreshToken(request.Username)
	if err != nil {

		logger.WithFields(logger.Fields{"component": "usecase", "action": "user login", "payload": request}).
			Errorf("failed create refresh token, err : %v", err)

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
		logger.WithFields(logger.Fields{"component": "usecase", "action": "user registration", "payload": request}).
			Errorf("failed register new user, err : %v", err)

		return nil, err
	}

	userData, err := user.db.CreateUser(ctx, db.CreateUserParams{
		Username:       request.Username,
		HashedPassword: hashPassword,
		FullName:       request.FullName,
		Email:          request.Email,
	})

	if err != nil {

		logger.WithFields(logger.Fields{"component": "usecase", "action": "user registration", "payload": request}).
			Errorf("failed register new user, err : %v", err)

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
		logger.WithFields(logger.Fields{"component": "usecase", "action": "get user", "username": username}).
			Errorf("failed get user, err : %v", err)

		return nil, err
	}

	return userData, nil
}

// UpdateUser : update existing user data
func (user *UseCase) UpdateUser(ctx *gin.Context, username string, request entities.UpdateUserRequest) (*db.User, error) {
	var params db.UpdateUserParams

	userData, err := user.CheckUserExist(ctx, username)
	if err != nil {
		logger.WithFields(logger.Fields{"component": "usecase", "action": "update user", "username": username, "payload": request}).
			Errorf("failed check if user exist, err : %v", err)

		return nil, err
	}

	if request.Email == userData.Email {
		logger.WithFields(logger.Fields{"component": "usecase", "action": "update user", "username": username, "payload": request}).
			Errorf("failed update user due to email is same")

		return nil, api_error.ErrEmailSameAsOld
	}

	if err := utils.CheckHashedPassword(request.OldPassword, userData.HashedPassword); err != nil {
		logger.WithFields(logger.Fields{"component": "usecase", "action": "update user", "username": username, "payload": request}).
			Errorf("failed update user due to password is not valid")

		return nil, api_error.ErrPasswordWrong
	}

	hashedPassword, err := utils.GenerateHashPassword(request.NewPassword)
	if err != nil {
		logger.WithFields(logger.Fields{"component": "usecase", "action": "update user", "username": username, "payload": request}).
			Errorf("failed generate hash password, err : %v", err)

		return nil, err
	}

	params.HashedPassword = sql.NullString{String: hashedPassword, Valid: true}

	dbUser, err := user.db.UpdateUser(ctx, params)
	if err != nil {
		logger.WithFields(logger.Fields{"component": "usecase", "action": "update user", "username": username, "payload": request}).
			Errorf("failed update user, err : %v", err)

		return nil, err
	}

	return &dbUser, nil
}

// CheckUserExist : check if user is exist in database
func (user *UseCase) CheckUserExist(ctx *gin.Context, username string) (*db.User, error) {
	userData, err := user.db.GetUser(ctx, username)
	if err != nil {
		logger.WithFields(logger.Fields{"component": "usecase", "action": "check user exist", "username": username}).
			Errorf("failed check if user exist, err : %v", err)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return &userData, nil
}
