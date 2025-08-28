package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ntdat104/go-clean-architecture/application/response"
	"github.com/ntdat104/go-clean-architecture/application/service"
)

type SystemHandler interface {
	GetTime(ctx *gin.Context)
}

type systemHandler struct {
	router        *gin.Engine
	systemService service.SystemService
}

func NewSystemHandler(router *gin.Engine, systemService service.SystemService) {
	h := &systemHandler{
		router:        router,
		systemService: systemService,
	}
	h.initRoutes()
}

func (h *systemHandler) initRoutes() {
	v1 := h.router.Group("/api/v1/system")
	{
		v1.GET("/time", h.GetTime)
	}
}

func (h *systemHandler) GetTime(ctx *gin.Context) {
	response.Success(ctx, h.systemService.GetTime())
}
