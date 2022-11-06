package api

import (
	"net/http"
	"search-service/db"

	"github.com/gin-gonic/gin"
)

type getUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type getUserListRequest struct {
	Offset int32 `form:"offset"`
	Limit  int32 `form:"limit" binding:"required,min=1,max=20"`
}

type createUserRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) GetListingByID(ctx *gin.Context) {

	// Check if request has ID field in URI.
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		ctx.Abort()
		return
	}

	// Execute query.
	result, err := server.store.GetListingByID(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) GetAllListings(ctx *gin.Context) {

	// Check if request has parameters offset and limit for pagination.
	var req getUserListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		ctx.Abort()
		return
	}

	arg := db.ListListingParam{
		Offset: req.Offset,
		Limit:  req.Limit,
	}

	// Execute query.
	result, err := server.store.GetAllListings(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) CreateListing(ctx *gin.Context) {

	// Check if request has all required fields in json body.
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		ctx.Abort()
		return
	}

	arg := db.CreateListingParam{
		Name: req.Name,
	}

	// Execute query.
	result, err := server.store.CreateListing(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusCreated, result)
}
