package services

import (
	"context"
	"errors"
	"time"

	"github.com/anamalala/internal/models"
	"github.com/anamalala/internal/repositories/interfaces"
)

type SuggestionService struct {
	suggestionRepo interfaces.SuggestionRepository
	userRepo       interfaces.UserRepository
}

func NewSuggestionService(
	suggestionRepo interfaces.SuggestionRepository,
	userRepo interfaces.UserRepository,
) SuggestionService {
	return SuggestionService{
		suggestionRepo: suggestionRepo,
		userRepo:       userRepo,
	}
}

func (s *SuggestionService) CreateSuggestion(ctx context.Context, suggestion models.Suggestion, userID string) (models.Suggestion, error) {
	// Converter ID do usuário para ObjectID
	// Verificar se o usuário existe
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return models.Suggestion{}, errors.New("usuário não encontrado")
	}
	
	// Configurar campos da sugestão
	suggestion.UserID = userID
	suggestion.CreatedAt = time.Now()
	suggestion.Status = "pending" // pending, reviewed, implemented, rejected
	
	// Salvar sugestão
	err = s.suggestionRepo.Create(ctx, suggestion)
	if err != nil {
		return models.Suggestion{}, err
	}
	
	return suggestion, nil
}

func (s *SuggestionService) GetAllSuggestions(ctx context.Context, page, limit int) (models.Suggestions, int, error) {
	suggestions, total, err := s.suggestionRepo.List(ctx, int64(page), int64(limit), models.SuggestionStatusAll)
	if err != nil {
		return nil, 0, err
	}
	
	return suggestions, int(total), nil
}

func (s *SuggestionService) GetSuggestionByID(ctx context.Context, id string) (models.Suggestion, error) {

	suggestion, err := s.suggestionRepo.FindByID(ctx, id)
	if err != nil {
		return models.Suggestion{}, err
	}
	
	return suggestion, nil
}

func (s *SuggestionService) GetSuggestionsByUserID(ctx context.Context, userID string, page, limit int) (models.Suggestions, int, error) {
	suggestions, total, err := s.suggestionRepo.ListByUserID(ctx, userID, int64(page), int64(limit))
	if err != nil {
		return nil, 0, err
	}
	
	return suggestions, int(total), nil
}

func (s *SuggestionService) UpdateSuggestionStatus(ctx context.Context, id string, status string, adminResponse string) (models.Suggestion, error) {
	// Validar status
	validStatuses := map[string]bool{
		"pending":     true,
		"reviewed":    true,
		"implemented": true,
		"rejected":    true,
	}
	
	if !validStatuses[status] {
		return models.Suggestion{}, errors.New("status inválido")
	}
	
	// Obter sugestão atual
	suggestion, err := s.suggestionRepo.FindByID(ctx, id)
	if err != nil {
		return models.Suggestion{}, err
	}
	
	// Atualizar status e resposta


	if status == "reviewed"{
		suggestion.Status = models.SuggestionStatusReviewed
	}
	if status == "approved"{
		suggestion.Status = models.SuggestionStatusApproved
	}
	if status == "rejected"{
		suggestion.Status = models.SuggestionStatusRejected
	}

	suggestion.UpdatedAt = time.Now()
	
	// Salvar alterações
	err = s.suggestionRepo.Update(ctx, suggestion)
	if err != nil {
		return models.Suggestion{}, err
	}
	
	return models.Suggestion{}, nil
}

func (s *SuggestionService) DeleteSuggestion(ctx context.Context, id string) error {

	return s.suggestionRepo.Delete(ctx,id)
}

func (s *SuggestionService) GetSuggestionsByStatus(ctx context.Context, status string, page, limit int) (models.Suggestions, int, error) {
	// Validar status
	validStatuses := map[string]bool{
		"pending":     true,
		"reviewed":    true,
		"implemented": true,
		"rejected":    true,
	}
	
	if !validStatuses[status] {
		return nil, 0, errors.New("status inválido")
	}
	
	suggestions, total, err := s.suggestionRepo.GetByStatus(ctx, status, 0, 0)
	if err != nil {
		return nil, 0, err
	}
	
	return suggestions, int(total), nil
}
