package infrastructure

import (
	"fmt"
	"github.com/dhiemaz/bank-api/config"
	accountHandler "github.com/dhiemaz/bank-api/domain/account/handler"
	accountUsecase "github.com/dhiemaz/bank-api/domain/account/usecase"
	securityHandler "github.com/dhiemaz/bank-api/domain/security/handler"
	securityUsecase "github.com/dhiemaz/bank-api/domain/security/usecase"
	transactionHandler "github.com/dhiemaz/bank-api/domain/transaction/handler"
	transactionUsecase "github.com/dhiemaz/bank-api/domain/transaction/usecase"
	userHandler "github.com/dhiemaz/bank-api/domain/user/handler"
	userUsecase "github.com/dhiemaz/bank-api/domain/user/usecase"
	"github.com/dhiemaz/bank-api/infrastructure/db/sqlc"
	"github.com/dhiemaz/bank-api/middlewares"
	"github.com/dhiemaz/bank-api/swagger/docs"
	"github.com/dhiemaz/bank-api/utils"
	"github.com/dhiemaz/bank-api/utils/token"

	_ "github.com/dhiemaz/bank-api/swagger/docs"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type GinServer struct {
	config             *config.Config
	dbStore            db.Store
	dbQueries          db.Querier
	tm                 token.Maker
	authHandler        *securityHandler.Handler
	userHandler        *userHandler.Handler
	accountHandler     *accountHandler.Handler
	transactionHandler *transactionHandler.Handler
	router             *gin.Engine
}

func NewServer(config *config.Config, dbStore db.Store, dbQueries db.Querier) (*GinServer, error) {
	maker, err := token.NewPasetoMaker(config.SymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create tokenMaker, %w", err)
	}

	// authentication
	authUC := securityUsecase.NewAuthUseCase(dbStore)
	authHandler := securityHandler.NewAuthHandler(authUC)

	// user
	userUC := userUsecase.NewUserUseCase(dbQueries)
	userHandler := userHandler.NewUserHandler(userUC)

	// account
	accountUC := accountUsecase.NewAccountUseCase(dbQueries)
	accountHandler := accountHandler.NewAccountHandler(accountUC)

	// transaction
	transactionUC := transactionUsecase.NewTransferUseCase(dbStore, accountUC)
	transactionHandler := transactionHandler.NewTransactionHandler(transactionUC)

	s := &GinServer{
		config:             config,
		tm:                 maker,
		dbStore:            dbStore,
		dbQueries:          dbQueries,
		authHandler:        authHandler,
		transactionHandler: transactionHandler,
		userHandler:        userHandler,
		accountHandler:     accountHandler,
	}

	gin.SetMode(gin.ReleaseMode)
	s.setupValidator()
	s.setupRouter()
	s.setupSwagger(config)
	return s, nil
}

func (s *GinServer) Start(address string) error {
	if err := s.router.Run(address); err != nil {
		return err
	}
	return nil

}

func (s *GinServer) setupValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", utils.ValidCurrency)
	}
}

func (s *GinServer) setupRouter() {
	router := gin.Default()

	auth := router.Group("/").Use(middlewares.AuthMiddleware(s.tm))

	// Account Routes
	auth.POST("/api/accounts", s.accountHandler.CreateAccount)
	auth.GET("/api/accounts/:id", s.accountHandler.GetAccount)
	auth.GET("/api/accounts", s.accountHandler.GetAccounts)
	auth.GET("/api/accounts/del", s.accountHandler.GetDeletedAccounts)
	auth.PATCH("/api/accounts/res/:id", s.accountHandler.RestoreAccount)
	auth.DELETE("/api/accounts/:id", s.accountHandler.DeleteAccount)

	// Transfer Routes
	auth.GET("/api/transfers/:id", s.transactionHandler.GetTransfersList)
	auth.POST("/api/transfers", s.transactionHandler.CreateTransfer)

	// User Routes
	auth.GET("api/users", s.userHandler.GetUser)
	auth.PATCH("api/users", s.userHandler.UpdateUser)

	// Unauthenticated Routes
	router.POST("api/users/register", s.userHandler.Register)
	router.POST("api/users/login", s.userHandler.LoginUser)
	router.POST("api/users/renew", s.authHandler.RenewAccessToken)

	s.router = router
}

func (s *GinServer) setupSwagger(config *config.Config) {
	if config.Env == "development" {
		docs.SwaggerInfo.BasePath = "/api"
	} else {
		docs.SwaggerInfo.BasePath = "/bank-api/api"
	}
	s.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
