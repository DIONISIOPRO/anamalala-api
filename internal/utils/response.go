package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response estrutura padrão para respostas da API
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta contém metadados para paginação
type Meta struct {
	Total       int `json:"total"`
	Page        int `json:"page"`
	PerPage     int `json:"per_page"`
	TotalPages  int `json:"total_pages"`
	HasNextPage bool `json:"has_next_page"`
	HasPrevPage bool `json:"has_prev_page"`
}

// NewMeta cria uma estrutura de metadados para paginação
func NewMeta(total, page, perPage int) Meta {
	totalPages := total / perPage
	if total%perPage > 0 {
		totalPages++
	}

	return Meta{
		Total:       total,
		Page:        page,
		PerPage:     perPage,
		TotalPages:  totalPages,
		HasNextPage: page < totalPages,
		HasPrevPage: page > 1,
	}
}

// SuccessResponse envia uma resposta de sucesso
func SuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// CreatedResponse envia uma resposta de recurso criado
func CreatedResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// BadRequestResponse envia uma resposta de requisição inválida
func BadRequestResponse(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: message,
	})
}

// UnauthorizedResponse envia uma resposta de não autorizado
func UnauthorizedResponse(c *gin.Context, message string) {
	if message == "" {
		message = "não autorizado"
	}
	c.JSON(http.StatusUnauthorized, Response{
		Success: false,
		Message: message,
	})
}

// ForbiddenResponse envia uma resposta de acesso proibido
func ForbiddenResponse(c *gin.Context, message string) {
	if message == "" {
		message = "acesso proibido"
	}
	c.JSON(http.StatusForbidden, Response{
		Success: false,
		Message: message,
	})
}

// NotFoundResponse envia uma resposta de recurso não encontrado
func NotFoundResponse(c *gin.Context, message string) {
	if message == "" {
		message = "recurso não encontrado"
	}
	c.JSON(http.StatusNotFound, Response{
		Success: false,
		Message: message,
	})
}

// InternalServerErrorResponse envia uma resposta de erro interno
func InternalServerErrorResponse(c *gin.Context, message string) {
	if message == "" {
		message = "erro interno do servidor"
	}
	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Message: message,
	})
}
