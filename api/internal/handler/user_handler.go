package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"

	"github.com/khalifaalhasan/go-gin-clean-arch/internal/database"
	"github.com/khalifaalhasan/go-gin-clean-arch/internal/dto"
	"github.com/khalifaalhasan/go-gin-clean-arch/util"

)

type UserHandler struct {
	store *database.Queries
}

func NewUserHandler(store *database.Queries) *UserHandler {
	return &UserHandler{store : store}
}

func (h * UserHandler) CreateUser(ctx *gin.Context) {
	var req dto.CreateUserRequest

	if err := ctx.ShouldBindJSON(&req);err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	arg := database.CreateUserParams{
		Username: req.Username,
		Email: req.Email,
		Password: hashedPassword,
		FullName: req.FullName,
	}

	user, err := h.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, gin.H{"error": "username or email already exists"})
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rsp := dto.NewUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}