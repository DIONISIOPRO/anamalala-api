package utils

import (
	"errors"
	"regexp"
	"unicode"
)
type  Validator struct{}

func  NewValidator() Validator {
	return Validator{}
}

// ValidateEmail valida o formato do email
func  (v *Validator)  ValidateEmail(email string) bool {
	// Expressão regular para validação básica de email
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}
// ValidateEmail valida o formato do email

// ValidatePhoneNumber valida número de telefone moçambicano
// Padrão: +258xxxxxxxxx ou 258xxxxxxxxx ou 8xxxxxxxx
func (v *Validator) ValidatePhoneNumber(phone string) bool {
	// Remover espaços
	re := regexp.MustCompile(`\s+`)
	phone = re.ReplaceAllString(phone, "")
	
	// Verificar padrões válidos
	pattern := `^(\+258|258)?8[234567]\d{7}$`
	
	re = regexp.MustCompile(pattern)
	return re.MatchString(phone)
}

// FormatPhoneNumber formata número de telefone para formato padrão
func (v *Validator) FormatPhoneNumber(phone string) string {
	// Remover espaços e caracteres não numéricos
	re := regexp.MustCompile(`\D+`)
	phone = re.ReplaceAllString(phone, "")
	
	// Se começar com 258, manter como está
	if len(phone) >= 12 && phone[:3] == "258" {
		return phone[3:]
	}
	return phone
}

// ValidatePassword verifica se a senha atende aos requisitos de segurança
func (v *Validator)  ValidatePassword(password string) error {
	if len(password) < 6 {
		return errors.New("a senha deve ter pelo menos 6 caracteres")
	}
	
	// Verificação mais rigorosa pode ser adicionada conforme necessário
	// Por exemplo, exigir letras maiúsculas, números, etc.
	
	hasLetter := false
	hasDigit := false
	
	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
	}
	
	if !hasLetter || !hasDigit {
		return errors.New("a senha deve conter pelo menos uma letra e um número")
	}
	
	return nil
}

// ValidateProvince verifica se a província é válida para Moçambique
func (v *Validator) ValidateProvince(province string) bool {
	validProvinces := map[string]bool{
		"Maputo Cidade":  true,
		"Maputo":         true,
		"Gaza":           true,
		"Inhambane":      true,
		"Manica":         true,
		"Sofala":         true,
		"Tete":           true,
		"Zambézia":       true,
		"Nampula":        true,
		"Cabo Delgado":   true,
		"Niassa":         true,
	}
	
	return validProvinces[province]
}
