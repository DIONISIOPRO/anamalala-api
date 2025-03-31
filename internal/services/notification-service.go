package services

import (
	"context"
	"errors"
	"time"

	"github.com/anamalala/internal/models"
	"github.com/anamalala/internal/repositories/interfaces"
)

type NotificationService struct {
	notificationRepo interfaces.NotificationRepository
	userRepo         interfaces.UserRepository
}

func NewNotificationService(
	notificationRepo interfaces.NotificationRepository,
	userRepo interfaces.UserRepository,
) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
	}
}

func (s *NotificationService) CreateNotification(ctx context.Context, notification models.Notification) (models.Notification, error) {
	// Verificar se o usuário existe
	_, err := s.userRepo.FindByID(ctx, notification.UserID)
	if err != nil {
		return models.Notification{}, errors.New("usuário não encontrado")
	}

	// Configurar campos da notificação
	notification.CreatedAt = time.Now()
	notification.Read = false

	// Salvar notificação
	err = s.notificationRepo.Create(ctx, notification)
	if err != nil {
		return models.Notification{}, err
	}

	return notification, nil
}

func (s *NotificationService) GetUserNotifications(ctx context.Context, userID string, page, limit int) (models.Notifications, int, error) {
	notifications, total, err := s.notificationRepo.ListByUserID(ctx, userID, int64(page), int64(limit), false)
	if err != nil {
		return nil, 0, err
	}

	return notifications, int(total), nil
}

func (s *NotificationService) GetUnreadNotificationsCount(ctx context.Context, userID string) (int, error) {
	count, err := s.notificationRepo.CountUnread(ctx, userID)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID string, userID string) error {
	// Obter notificação
	notification, _, err := s.notificationRepo.ListByUserID(ctx, notificationID, 0, 0, true)
	if err != nil {
		return err
	}

	for _, n := range notification {
		s.notificationRepo.MarkAsRead(ctx, n.UserID)
	}

	return nil

}

func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID string) (int, error) {

	err := s.notificationRepo.MarkAllAsRead(ctx, userID)
	if err != nil {
		return 0, err
	}

	return 0, nil
}

func (s *NotificationService) DeleteNotification(ctx context.Context, notificationID string, userID string) error {

	// Obter notificação
	notification, err := s.notificationRepo.FindByID(ctx, notificationID)
	if err != nil {
		return err
	}

	// Verificar se a notificação pertence ao usuário
	if notification.UserID != userID {
		return errors.New("notificação não pertence a este usuário")
	}

	return s.notificationRepo.Delete(ctx, notificationID)
}

// Métodos para criar notificações específicas

func (s *NotificationService) NotifyNewComment(ctx context.Context, postAuthorID string, postID string, commenterID string) error {
	// Verificar se o autor do post é o mesmo que comentou
	if postAuthorID == commenterID {
		return nil // Não notificar o próprio usuário
	}

	// Obter informações do usuário que comentou
	commenter, err := s.userRepo.FindByID(ctx, commenterID)
	if err != nil {
		return err
	}

	// Criar notificação
	notification := models.Notification{
		UserID:    postAuthorID,
		Type:      "comment",
		Message:   commenter.Name + " comentou na sua postagem",
		Reference: postID,
		CreatedAt: time.Now(),
		Read:      false,
	}

	err = s.notificationRepo.Create(ctx, notification)
	return err
}

func (s *NotificationService) NotifyNewLike(ctx context.Context, contentAuthorID string, contentID string, likerID string, contentType string) error {
	// Verificar se o autor do conteúdo é o mesmo que curtiu
	if contentAuthorID == likerID {
		return nil // Não notificar o próprio usuário
	}

	// Obter informações do usuário que curtiu
	liker, err := s.userRepo.FindByID(ctx, likerID)
	if err != nil {
		return err
	}

	// Preparar mensagem com base no tipo de conteúdo
	var message string
	if contentType == "post" {
		message = liker.Name + " curtiu sua postagem"
	} else if contentType == "comment" {
		message = liker.Name + " curtiu seu comentário"
	} else {
		return errors.New("tipo de conteúdo inválido")
	}

	// Criar notificação
	notification := models.Notification{
		UserID:    contentAuthorID,
		Type:      "like",
		Message:   message,
		Reference: contentID,
		CreatedAt: time.Now(),
		Read:      false,
	}

	err = s.notificationRepo.Create(ctx, notification)
	return err
}

func (s *NotificationService) NotifyAdminAction(ctx context.Context, userID string, action string, reason string) error {
	// Criar notificação
	notification := models.Notification{
		UserID:    userID,
		Type:      "admin",
		Message:   action + ": " + reason,
		CreatedAt: time.Now(),
		Read:      false,
	}

	err := s.notificationRepo.Create(ctx, notification)
	return err
}

func (s *NotificationService) NotifyAllUsers(ctx context.Context, message string, notificationType string) (int, error) {
	// Obter todos os usuários
	users, _, err := s.userRepo.List(ctx, 0, 0) // Sem paginação para obter todos
	if err != nil {
		return 0, err
	}

	if len(users) == 0 {
		return 0, errors.New("nenhum usuário encontrado")
	}

	// Contador de notificações criadas
	createdCount := 0

	// Criar notificação para cada usuário
	for _, user := range users {
		notification := models.Notification{
			UserID:    user.ID,
			Type:      models.NotificationType(notificationType),
			Message:   message,
			CreatedAt: time.Now(),
			Read:      false,
		}

		err := s.notificationRepo.Create(ctx, notification)
		if err == nil {
			createdCount++
		}
		// Continuar mesmo se houver erro para tentar enviar para todos os usuários
	}

	if createdCount == 0 {
		return 0, errors.New("falha ao criar notificações para todos os usuários")
	}

	return createdCount, nil
}
