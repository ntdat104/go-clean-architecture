package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ntdat104/go-clean-architecture/internal/form/response"
	"github.com/ntdat104/go-clean-architecture/internal/model"
	"github.com/ntdat104/go-clean-architecture/internal/service"
)

type UserHandler interface {
	Create(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type userHandler struct {
	router      *gin.Engine
	userService service.UserService
}

func NewUserHandler(router *gin.Engine, userService service.UserService) {
	h := &userHandler{
		router:      router,
		userService: userService,
	}
	h.initRoutes()
}

func (h *userHandler) initRoutes() {
	v1 := h.router.Group("/api/v1/users")
	{
		v1.POST("", h.Create)
		v1.GET("/:id", h.GetByID)
		v1.PUT("/:id", h.Update)
		v1.DELETE("/:id", h.Delete)
	}
}

// Create user
func (h *userHandler) Create(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		response.JSON(ctx, http.StatusBadRequest, "invalid request payload")
		return
	}

	created, err := h.userService.CreateUser(ctx, &user)
	if err != nil {
		response.JSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(ctx, created)
}

// Get user by ID
func (h *userHandler) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.JSON(ctx, http.StatusBadRequest, "invalid user ID")
		return
	}

	user, err := h.userService.GetUserByID(ctx, id)
	if err != nil {
		response.JSON(ctx, http.StatusNotFound, err.Error())
		return
	}

	response.Success(ctx, user)
}

// Update user
func (h *userHandler) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.JSON(ctx, http.StatusBadRequest, "invalid user ID")
		return
	}

	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		response.JSON(ctx, http.StatusBadRequest, "invalid request payload")
		return
	}
	user.ID = id

	if err := h.userService.UpdateUser(ctx, &user); err != nil {
		response.JSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(ctx, "updated successfully")
}

// Delete user
func (h *userHandler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.JSON(ctx, http.StatusBadRequest, "invalid user ID")
		return
	}

	if err := h.userService.DeleteUser(ctx, id); err != nil {
		response.JSON(ctx, http.StatusNotFound, err.Error())
		return
	}

	response.Success(ctx, "deleted successfully")
}
