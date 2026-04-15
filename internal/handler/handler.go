package handler

import (
	"net/http"
	"time"
	"vul-parser/internal/domain/dto"
	"vul-parser/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Services *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		Services: service,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api")
	{
		api.GET("/health", h.Health)
		analyze := api.Group("/analyze")
		{
			analyze.POST("", h.Analyze)
			analyze.POST("/files", h.AnalyzeFile)
		}
	}

	return router
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, dto.HealthResponse{
		Status:  "ok",
		Version: "1.0.0",
	})
}
