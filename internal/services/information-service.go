package services

import (
	"context"
	"time"

	"github.com/anamalala/internal/models"
	"github.com/anamalala/internal/repositories/interfaces"
)

type InformationService struct {
	infoRepo interfaces.InformationRepository
}

func NewInformationService(infoRepo interfaces.InformationRepository) InformationService {
	return InformationService{
		infoRepo: infoRepo,
	}
}

func (s *InformationService) CreateInformation(ctx context.Context, info models.Information, authorID string) (models.Information, error) {
	// Converter ID do autor para ObjectID
	// Configurar campos do artigo
	info.AuthorID = authorID
	info.CreatedAt = time.Now()
	info.UpdatedAt = time.Now()

	// Salvar artigo
	err := s.infoRepo.Create(ctx, info)
	if err != nil {
		return models.Information{}, err
	}

	return info, nil
}

func (s *InformationService) GetInformation(ctx context.Context, id string) (models.Information, error) {

	info, err := s.infoRepo.FindByID(ctx, id)
	if err != nil {
		return models.Information{}, err
	}

	return info, nil
}

func (s *InformationService) GetAllInformation(ctx context.Context, page, limit int) (models.Informations, int, error) {
	infoItems, total, err := s.infoRepo.List(ctx, int64(page), int64(limit), true)
	if err != nil {
		return nil, 0, err
	}

	return infoItems, int(total), nil
}

func (s *InformationService) UpdateInformation(ctx context.Context, id string, updateData models.Information) (models.Information, error) {
	// Obter informação atual
	info, err := s.infoRepo.FindByID(ctx, id)
	if err != nil {
		return models.Information{}, err
	}

	// Atualizar campos
	if updateData.Title != "" {
		info.Title = updateData.Title
	}

	if updateData.Content != "" {
		info.Content = updateData.Content
	}

	// Atualizar data de modificação
	info.UpdatedAt = time.Now()

	// Salvar alterações
	err = s.infoRepo.Update(ctx, info)
	if err != nil {
		return models.Information{}, err
	}

	return info, nil
}

func (s *InformationService) DeleteInformation(ctx context.Context, id string) error {
	return s.infoRepo.Delete(ctx, id)
}
