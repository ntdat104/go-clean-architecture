package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ntdat104/go-clean-architecture/api/dto"
	"github.com/ntdat104/go-clean-architecture/api/error_code"
	"github.com/ntdat104/go-clean-architecture/api/http/handle"
	"github.com/ntdat104/go-clean-architecture/api/http/validator"
	"github.com/ntdat104/go-clean-architecture/application/service"
	"github.com/ntdat104/go-clean-architecture/pkg/logger"
)

type ExampleHandler interface {
	Create(ctx *gin.Context)
	// Delete(ctx *gin.Context)
	// Update(ctx *gin.Context)
	// Get(ctx *gin.Context)
	// FindByName(ctx *gin.Context)
}

type exampleHandler struct {
	router         *gin.Engine
	exampleService service.IExampleService
}

func NewExampleHandler(router *gin.Engine, exampleService service.IExampleService) {
	h := &exampleHandler{
		router:         router,
		exampleService: exampleService,
	}
	h.initRoutes()
}

func (h *exampleHandler) initRoutes() {
	v1 := h.router.Group("/api/v1/examples")
	{
		v1.POST("", h.Create)
		// v1.GET("/:id", h.Get)
		// v1.PUT("/:id", h.Update)
		// v1.DELETE("/:id", h.Delete)
		// v1.GET("/name/:name", h.FindByName)
	}
}

func (h *exampleHandler) Create(ctx *gin.Context) {
	response := handle.NewResponse(ctx)
	body := dto.CreateExampleReq{}

	valid, errs := validator.BindAndValid(ctx, &body, ctx.ShouldBindJSON)
	if !valid {
		logger.SugaredLogger.Errorf("CreateExample.BindAndValid errs: %v", errs)
		err := error_code.InvalidParams.WithDetails(errs.Errors()...)
		response.ToErrorResponse(err)
		return
	}

	// Execute use case
	var result any
	ctx.JSON(http.StatusCreated, result)
}
