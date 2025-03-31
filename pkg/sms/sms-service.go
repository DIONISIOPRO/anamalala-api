package sms

import (
	"errors"
	"fmt"
	"time"

	"github.com/anamalala/pkg/logger"
)

// Provider é uma interface para provedores de serviço de SMS
type Provider interface {
	Send(recipient, message string) error
}

// SMSConfig contém a configuração do serviço de SMS
type SMSConfig struct {
	ProviderType string
	APIKey       string
	SenderID     string
	// Outros campos específicos do provedor podem ser adicionados
}

// Service gerencia o envio de SMS
type Service struct {
	provider Provider
	logger   *logger.Logger
	config   *SMSConfig
}

// NewService cria uma nova instância de Service
func NewService(config *SMSConfig, logger *logger.Logger) (*Service, error) {
	// Criar provedor com base na configuração
	var provider Provider
	
	switch config.ProviderType {
	case "mock":
		provider = &MockProvider{logger: logger}
	case "africastalking": // Exemplo de provedor comum em África
		provider = &AfricasTalkingProvider{
			APIKey:   config.APIKey,
			SenderID: config.SenderID,
			logger:   logger,
		}
	case "twilio":
		provider = &TwilioProvider{
			APIKey:   config.APIKey,
			SenderID: config.SenderID,
			logger:   logger,
		}
	default:
		return nil, errors.New("provedor de SMS não suportado")
	}
	
	return &Service{
		provider: provider,
		logger:   logger,
		config:   config,
	}, nil
}

// Send envia uma mensagem SMS
func (s *Service) Send(recipient, message string) error {
	start := time.Now()
	err := s.provider.Send(recipient, message)
	duration := time.Since(start)
	
	if err != nil {
		s.logger.Error("sms_send_failed",
			"recipient", recipient,
			"error", err.Error(),
			"duration_ms", duration.Milliseconds(),
		)
		return err
	}
	
	s.logger.Info("sms_sent",
		"recipient", recipient,
		"message_length", len(message),
		"duration_ms", duration.Milliseconds(),
	)
	
	return nil
}

// SendBulk envia SMS em massa para múltiplos destinatários
func (s *Service) SendBulk(recipients []string, message string) (map[string]error, error) {
	results := make(map[string]error)
	
	for _, recipient := range recipients {
		err := s.Send(recipient, message)
		results[recipient] = err
	}
	
	return results, nil
}

// MockProvider é um provedor de SMS simulado para testes
type MockProvider struct {
	logger *logger.Logger
}

func (p *MockProvider) Send(recipient, message string) error {
	p.logger.Info("mock_sms_sent",
		"recipient", recipient,
		"message", message,
	)
	return nil
}

// AfricasTalkingProvider implementa o provedor AfricasTalking
type AfricasTalkingProvider struct {
	APIKey   string
	SenderID string
	logger   *logger.Logger
}

func (p *AfricasTalkingProvider) Send(recipient, message string) error {
	// Aqui seria implementada a integração com a API da AfricasTalking
	// Por enquanto, apenas registramos a tentativa
	p.logger.Info("africas_talking_send_attempt",
		"recipient", recipient,
		"api_key", fmt.Sprintf("%s****", p.APIKey[:4]),
		"sender_id", p.SenderID,
	)
	
	// Para implementação real, usar biblioteca HTTP para chamar a API
	return nil
}

// TwilioProvider implementa o provedor Twilio
type TwilioProvider struct {
	APIKey   string
	SenderID string
	logger   *logger.Logger
}

func (p *TwilioProvider) Send(recipient, message string) error {
	// Aqui seria implementada a integração com a API da Twilio
	// Por enquanto, apenas registramos a tentativa
	p.logger.Info("twilio_send_attempt",
		"recipient", recipient,
		"api_key", fmt.Sprintf("%s****", p.APIKey[:4]),
		"sender_id", p.SenderID,
	)
	
	// Para implementação real, usar biblioteca HTTP para chamar a API
	return nil
}
