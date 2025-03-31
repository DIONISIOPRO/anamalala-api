package handlers

import (
	"net/http"
	"strconv"

	"github.com/anamalala/internal/models"
	"github.com/anamalala/internal/services"
	"github.com/gin-gonic/gin"
)

type SuggestionHandler struct {
	suggestionService services.SuggestionService
}

func NewSuggestionHandler(suggestionService services.SuggestionService) SuggestionHandler {
	return SuggestionHandler{
		suggestionService: suggestionService,
	}
}

func (h *SuggestionHandler) CreateSuggestion(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, "Usuário não autenticado")
		return
	}

	var suggestion models.Suggestion
	if err := c.ShouldBindJSON(&suggestion); err != nil {
		c.JSON(http.StatusBadRequest, "Dados inválidos")
		return
	}

	// Validação de campos obrigatórios
	if suggestion.Description == "" {
		c.JSON(http.StatusBadRequest, "Conteúdo é obrigatório")
		return
	}

	suggestion.UserID, _ = userID.(string)
	suggestion.Status = "pending" // Status inicial: pendente de avaliação

	createdSuggestion, err := h.suggestionService.CreateSuggestion(c, suggestion, suggestion.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao criar sugestão")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Sugestão enviada com sucesso",
		"data":    createdSuggestion,
	})
}

func (h *SuggestionHandler) GetUserSuggestions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, "Usuário não autenticado")
		return
	}

	// Parâmetros para paginação
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))

	if err != nil {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if err != nil {
		limit = 10
	}

	suggestions, total, err := h.suggestionService.GetSuggestionsByUserID(c, userID.(string), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao buscar sugestões")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Sugestões obtidas com sucesso",
		"data": gin.H{
			"suggestions": suggestions,
			"total":       total,
			"page":        page,
			"limit":       limit,
			"totalPages":  (total + limit - 1) / limit,
		},
	})
}

func (h *SuggestionHandler) GetAllSuggestions(c *gin.Context) {
	// Verificando se é um usuário administrador
	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, "Acesso restrito a administradores")
		return
	}

	// Parâmetros para paginação
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))

	if err != nil {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if err != nil {
		limit = 10
	}

	// Filtros opcionais

	suggestions, total, err := h.suggestionService.GetAllSuggestions(c, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao buscar sugestões")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Sugestões obtidas com sucesso",
		"data": gin.H{
			"suggestions": suggestions,
			"total":       total,
			"page":        page,
			"limit":       limit,
			"totalPages":  (total + limit - 1) / limit,
		},
	})
}

func (h *SuggestionHandler) GetSuggestionByID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, "Usuário não autenticado")
		return
	}

	suggestionID := c.Param("id")
	if suggestionID == "" {
		c.JSON(http.StatusBadRequest, "ID da sugestão não fornecido")
		return
	}

	suggestion, err := h.suggestionService.GetSuggestionByID(c, suggestionID)
	if err != nil {
		c.JSON(http.StatusNotFound, "Sugestão não encontrada")
		return
	}

	// Verificar se o usuário é o autor da sugestão ou um administrador
	isAdmin, adminExists := c.Get("isAdmin")
	if suggestion.UserID != userID.(string) && (!adminExists || !isAdmin.(bool)) {
		c.JSON(http.StatusForbidden, "Sem permissão para visualizar esta sugestão")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Sugestão obtida com sucesso",
		"data":    suggestion,
	})
}

func (h *SuggestionHandler) UpdateSuggestionStatus(c *gin.Context) {
	// Verificando se é um usuário administrador
	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, "Acesso restrito a administradores")
		return
	}

	suggestionID := c.Param("id")
	if suggestionID == "" {
		c.JSON(http.StatusBadRequest, "ID da sugestão não fornecido")
		return
	}

	var updateData struct {
		Status     string `json:"status" binding:"required"`
		AdminNotes string `json:"admin_notes"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, "Dados inválidos")
		return
	}

	// Validar status
	validStatus := map[string]bool{
		"pending":   true,
		"approved":  true,
		"rejected":  true,
		"completed": true,
	}

	if !validStatus[updateData.Status] {
		c.JSON(http.StatusBadRequest, "Status inválido")
		return
	}

	_, err := h.suggestionService.UpdateSuggestionStatus(c, suggestionID, updateData.Status, updateData.AdminNotes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao atualizar status da sugestão")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Status da sugestão atualizado com sucesso",
		"data": gin.H{
			"id":          suggestionID,
			"status":      updateData.Status,
			"admin_notes": updateData.AdminNotes,
		},
	})
}
