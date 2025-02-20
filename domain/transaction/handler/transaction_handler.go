package handler

import (
	"github.com/dhiemaz/bank-api/domain/transaction/usecase"
	"github.com/dhiemaz/bank-api/entities"
	"github.com/dhiemaz/bank-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Usecase usecase.TransferUseCase
}

func NewTransactionHandler(usecase usecase.TransferUseCase) *Handler {
	return &Handler{
		Usecase: usecase,
	}
}

// CreateTransfer godoc
//
//	@Summary		creates a new transfer between two accounts
//	@Description	creates a new transfer between two accounts
//	@Tags			transfers
//	@Accept			json
//	@Produce		json
//	@Param			body	body		createTransferReq	true	"Transfer to create"
//	@Success		200		{object}	response.JSON{data=transferResponse}
//	@Failure		400,500	{object}	response.JSON{}
//	@Security		bearerAuth
//	@Router			/transfers [post]
func (transaction *Handler) CreateTransfer(ctx *gin.Context) {
	var request entities.CreateTransferRequest
	if err := utils.ParseBody(ctx, &request); err != nil {
		return
	}

	fromAccount, toAccount, err := transaction.Usecase.ValidateTransfer(ctx, request.FromAccountID, request.ToAccountID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.Err(err))
		return
	}

	request.ToAccountID = toAccount.ID
	request.FromAccountID = fromAccount.ID

	result, err := transaction.Usecase.CreateTransfer(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.Err(err))
		return
	}

	res := utils.FromTransferTxToTransferResponse(result)
	ctx.JSON(http.StatusOK, res)
}

// GetTransfers godoc
//
//	@Summary		gets all transfers for an account
//	@Description	gets all transfers for an account
//	@Tags			transfers
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int64	true	"Account ID"
//	@Param			page_id		query		int32	true	"Page ID"
//	@Param			page_size	query		int32	true	"Page Size"
//	@Success		200			{object}	response.JSON{data=transferResponse}
//	@Failure		400,500		{object}	response.JSON{}
//	@Security		bearerAuth
//	@Router			/transfers/{id} [get]
func (transaction *Handler) GetTransfersList(ctx *gin.Context) {
	var request entities.GetTransferRequest
	var pgQuery *utils.PaginationQuery
	var err error

	if err := utils.ParseURI(ctx, &request); err != nil {
		return
	}

	if pgQuery, err = utils.ParsePagination(ctx); err != nil {
		ctx.JSON(http.StatusBadRequest, entities.Err(err))
		return
	}

	transfers, err := transaction.Usecase.GetListTransfer(ctx, request, pgQuery)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, entities.Err(err))
	}

	var responsesTransfer []*entities.TransferResponse
	for _, transfer := range transfers {
		responsesTransfer = append(responsesTransfer, utils.MapTransferToResponse(transfer))
	}

	ctx.JSON(http.StatusOK, entities.Success(responsesTransfer))
}
