package handlers

import (
	"net/http"
	"strconv"

	"github.com/anamalala/internal/models"
	"github.com/anamalala/internal/services"
	"github.com/gin-gonic/gin"
)

type InformationHandler struct {
	informationService services.InformationService
}

func NewInformationHandler(informationService services.InformationService) InformationHandler {
	return InformationHandler{
		informationService: informationService,
	}
}

func (h *InformationHandler) Create(c *gin.Context) {
	// Verificando se é um usuário administrador
	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, "Acesso restrito a administradores")
		return
	}

	var information models.Information
	if err := c.ShouldBindJSON(information); err != nil {
		c.JSON(http.StatusBadRequest, "Dados inválidos")
		return
	}

	// Obtendo ID do autor (administrador)
	authorID, _ := c.Get("userID")
	information.AuthorID = authorID.(string)

	// Validação de campos obrigatórios
	if information.Title == "" || information.Content == "" {
		c.JSON(http.StatusBadRequest, "Título e conteúdo são obrigatórios")
		return
	}

	createdInfo, err := h.informationService.CreateInformation(c, information, information.AuthorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao criar informação")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Informação criada com sucesso",
		"data":    createdInfo,
	})
}

func (h *InformationHandler) GetByID(c *gin.Context) {
	infoID := c.Param("id")
	if infoID == "" {
		c.JSON(http.StatusBadRequest, "ID da informação não fornecido")
		return
	}

	info, err := h.informationService.GetInformation(c, infoID)
	if err != nil {
		c.JSON(http.StatusNotFound, "Informação não encontrada")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Informação obtida com sucesso",
		"data":    info,
	})
}

func (h *InformationHandler) Update(c *gin.Context) {
	// Verificando se é um usuário administrador
	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, "Acesso restrito a administradores")
		return
	}

	infoID := c.Param("id")
	if infoID == "" {
		c.JSON(http.StatusBadRequest, "ID da informação não fornecido")
		return
	}

	var updateData struct {
		Title    string   `json:"title"`
		Content  string   `json:"content"`
		ImageURL string   `json:"image_url"`
		Tags     []string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, "Dados inválidos")
		return
	}

	info, err := h.informationService.GetInformation(c, infoID)
	if err != nil {
		c.JSON(http.StatusNotFound, "Informação não encontrada")
		return
	}

	// Atualizar apenas os campos fornecidos
	if updateData.Title != "" {
		info.Title = updateData.Title
	}
	if updateData.Content != "" {
		info.Content = updateData.Content
	}

	// Registrar quem atualizou

	updatedInfo, err := h.informationService.UpdateInformation(c, string(info.ID), info)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao atualizar informação")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Informação atualizada com sucesso",
		"data":    updatedInfo,
	})
}

func (h *InformationHandler) Delete(c *gin.Context) {
	// Verificando se é um usuário administrador
	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, "Acesso restrito a administradores")
		return
	}

	infoID := c.Param("id")
	if infoID == "" {
		c.JSON(http.StatusBadRequest, "ID da informação não fornecido")
		return
	}

	err := h.informationService.DeleteInformation(c, infoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao excluir informação")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Informação excluída com sucesso",
	})
}

func (h *InformationHandler) GetAll(c *gin.Context) {
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

	infos, total, err := h.informationService.GetAllInformation(c, 0, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao buscar informações")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Informações obtidas com sucesso",
		"data": gin.H{
			"items":      infos,
			"total":      total,
			"page":       page,
			"limit":      limit,
			"totalPages": (total + limit - 1) / limit,
		},
	})
}
