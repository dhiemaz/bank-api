package utils

import (
	"github.com/dhiemaz/bank-api/entities"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseBody(ctx *gin.Context, obj interface{}) error {
	err := ctx.ShouldBindJSON(obj)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, entities.Err(err))
		return err
	}
	return nil
}

func ParseURI(ctx *gin.Context, obj interface{}) error {
	if err := ctx.ShouldBindUri(obj); err != nil {
		ctx.JSON(http.StatusBadRequest, entities.Err(err))
		return err
	}
	return nil
}

type PaginationQuery struct {
	Offset int32 `form:"offset" binding:"required"`
	Limit  int32 `form:"limit" binding:"required"`
}

func ParsePagination(ctx *gin.Context) (*PaginationQuery, error) {
	pgQuery := &PaginationQuery{}

	offset := ctx.DefaultQuery("offset", "1")
	limit := ctx.DefaultQuery("limit", "1")

	limit32, err := strconv.Atoi(limit)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, entities.Err(err))
		return nil, err
	}
	offset32, err := strconv.Atoi(offset)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, entities.Err(err))
		return nil, err
	}

	pgQuery.Limit = int32(limit32)
	pgQuery.Offset = int32(offset32)

	return pgQuery, nil
}
