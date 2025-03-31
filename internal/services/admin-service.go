package services

import (
	"context"
	"errors"
	"time"

	"github.com/anamalala/internal/models"
	"github.com/anamalala/internal/repositories/interfaces"
	"github.com/anamalala/pkg/sms"
)

type AdminService struct {
	userRepo    interfaces.UserRepository
	postRepo    interfaces.PostRepository
	commentRepo interfaces.CommentRepository
	smsService  *sms.Service
}

func NewAdminService(
	userRepo interfaces.UserRepository,
	postRepo interfaces.PostRepository,
	commentRepo interfaces.CommentRepository,
	smsService *sms.Service,
) AdminService {
	return AdminService{
		userRepo:    userRepo,
		postRepo:    postRepo,
		commentRepo: commentRepo,
		smsService:  smsService,
	}
}

func (s *AdminService) BanUser(ctx context.Context, userID string) error {

	var user = models.User{}

	// Obter usuário
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verificar se já está banido
	if !user.Active {
		return errors.New("usuário já está banido")
	}

	// Banir usuário
	user.Active = false
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(ctx, user)
}

func (s *AdminService) UnbanUser(ctx context.Context, userID string) error {

	var user = models.User{}
	// Obter usuário
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verificar se está banido
	if user.Active {
		return errors.New("usuário não está banido")
	}

	// Desbanir usuário
	user.Active = true
	return s.userRepo.Update(ctx, user)
}

func (s *AdminService) PromoteToAdmin(ctx context.Context, userID string) error {
	var user = models.User{}

	// Obter usuário
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verificar se já é administrador
	if user.Role == "admin" {
		return errors.New("usuário já é administrador")
	}
	// Promover a administrador
	user.Role = "admin"
	return s.userRepo.Update(ctx, user)
}

func (s *AdminService) DemoteFromAdmin(ctx context.Context, userID string) error {

	var user = models.User{}

	// Obter usuário
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verificar se é administrador
	if user.Role != "admin" {
		return errors.New("usuário não é administrador")
	}

	// Rebaixar para usuário comum
	user.Role = "user"

	return s.userRepo.Update(ctx, user)
}

func (s *AdminService) GetBannedUsers(ctx context.Context, page, limit int) (models.Users, int, error) {
	var users = models.Users{}

	users, total, err := s.userRepo.InactiveUsers(ctx, int64(page), int64(limit))
	if err != nil {
		return nil, 0, err
	}

	// Remover informações sensíveis
	for _, user := range users {
		user.Password = ""
	}

	return users, int(total), nil
}

func (s *AdminService) GetAdminUsers(ctx context.Context, page, limit int) (models.Users, int, error) {
	var users = models.Users{}

	users, total, err := s.userRepo.ListByRole(ctx, "admin", int64(page), int64(limit))
	if err != nil {
		return models.Users{}, 0, err
	}

	// Remover informações sensíveis
	for i, user := range users {
		user.Password = ""
		users[i] = user
	}
	return users, int(total), nil
}

func (s *AdminService) SendSMS(ctx context.Context, message string, contacts []string) (int, error) {
	if len(contacts) == 0 {
		return 0, errors.New("nenhum contato fornecido")
	}

	// Contador de mensagens enviadas
	sentCount := 0

	// Enviar SMS para cada contato
	for _, contact := range contacts {
		err := s.smsService.Send(contact, message)
		if err == nil {
			sentCount++
		}
		// Continuar mesmo se houver erro para tentar enviar para todos os contatos
	}

	if sentCount == 0 {
		return 0, errors.New("falha ao enviar SMS para todos os contatos")
	}

	return sentCount, nil
}

func (s *AdminService) SendSMSToProvince(ctx context.Context, message string, province string) (int, error) {
	// Obter todos os usuários da província

	var users = models.Users{}
	users, _, err := s.userRepo.ListByProvince(ctx, province, 0, 0) // Sem paginação para obter todos
	if err != nil {
		return 0, err
	}

	if len(users) == 0 {
		return 0, errors.New("nenhum usuário encontrado na província")
	}

	// Extrair contatos
	contacts := make([]string, len(users))
	for i, user := range users {
		contacts[i] = user.Contact
	}

	// Enviar SMS para os contatos
	return s.SendSMS(ctx, message, contacts)
}

func (s *AdminService) SendSMSToAllUsers(ctx context.Context, message string) (int, error) {
	// Obter todos os usuários
	var users = models.Users{}
	users, _, err := s.userRepo.List(ctx, 0, 0) // Sem paginação para obter todos
	if err != nil {
		return 0, err
	}

	if len(users) == 0 {
		return 0, errors.New("nenhum usuário encontrado")
	}

	// Extrair contatos
	contacts := make([]string, len(users))
	for i, user := range users {
		contacts[i] = user.Contact
	}

	// Enviar SMS para os contatos
	return s.SendSMS(ctx, message, contacts)
}

func (s *AdminService) DeleteInappropriateContent(ctx context.Context, postID string, reason string) error {

	var post = models.Post{}
	// Obter a postagem
	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		return err
	}
	var user = models.User{}
	// Notificar o autor (opcional)
	user, err = s.userRepo.FindByID(ctx, post.UserID)
	if err == nil && user.Contact != "" {
		message := "Sua postagem foi removida por violar as diretrizes da comunidade. Motivo: " + reason
		s.smsService.Send(user.Contact, message)
	}

	// Excluir todos os comentários da postagem
	err = s.commentRepo.Delete(ctx, postID)
	if err != nil {
		return err
	}

	// Excluir a postagem
	return s.postRepo.Delete(ctx, postID)
}

func (s *AdminService) GetSystemStats(ctx context.Context) (map[string]interface{}, error) {
	// Total de usuários
	users, _, err := s.userRepo.List(ctx, 0, 0)

	totalUsers := len(users)

	if err != nil {
		return nil, err
	}

	// Total de usuários por província
	provinceStats := make(map[string]int)
	provinces := []string{"Maputo", "Gaza", "Inhambane", "Sofala", "Manica",
		"Tete", "Zambézia", "Nampula", "Cabo Delgado", "Niassa"}

	for index, province := range provinces {
		users, _, _ := s.userRepo.ListByProvince(ctx, provinces[index], 0, 0)

		provinceStats[province] = len(users)
	}

	// Total de postagens
	_, totalPosts, err := s.postRepo.List(ctx, 0, 0)
	if err != nil {
		return nil, err
	}

	// Usuários banidos
	bannedUsers, _, err := s.userRepo.InactiveUsers(ctx, 0, 0)
	if err != nil {
		return nil, err
	}

	// Usuários administradores
	adminUsers, _, err := s.userRepo.ListByRole(ctx, "admin", 0, 0)
	if err != nil {
		return nil, err
	}

	// Compilar estatísticas
	stats := map[string]interface{}{
		"totalUsers":    totalUsers,
		"totalPosts":    totalPosts,
		"bannedUsers":   bannedUsers,
		"adminUsers":    adminUsers,
		"provinceStats": provinceStats,
	}

	return stats, nil
}

func (s *AdminService) GetAdminLogs(page, limit int, adminID string) (map[string]string, error) {
	return map[string]string{}, nil
}

func (s *AdminService) DeleteUserAccount(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer func() {
		cancel()
	}()
	return s.userRepo.Delete(ctx, userID)
}
