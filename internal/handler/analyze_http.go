package handler

import (
	"net/http"
	"vul-parser/internal/domain/dto"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Analyze(c *gin.Context) {
	var req dto.AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request body",
			Code:    http.StatusBadRequest,
			Details: err.Error(),
		})
		return
	}

	if req.Config == nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Config is required",
			Code:  http.StatusBadRequest,
		})
		return
	}

	resp, err := h.Services.AnalyzisHTTP.Analyze(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Analysis failed",
			Code:    http.StatusInternalServerError,
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) AnalyzeFile(c *gin.Context) {
	filePath := c.Query("path")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "File path is required",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var req dto.AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req = dto.AnalyzeRequest{}
	}

	resp, err := h.Services.AnalyzisHTTP.AnalyzeWithFile(filePath, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Analysis failed",
			Code:    http.StatusInternalServerError,
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}