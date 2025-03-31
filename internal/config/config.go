package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config contém todas as configurações do aplicativo
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	SMS      SMSConfig
	Enviroment string
}

// ServerConfig contém configurações relacionadas ao servidor HTTP
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig contém configurações relacionadas ao MongoDB
type DatabaseConfig struct {
	URI        string
	Name       string
	PoolSize   uint64
	MaxIdleTime time.Duration
}

// JWTConfig contém configurações relacionadas a autenticação JWT
type JWTConfig struct {
	Secret           string
	ExpirationHours  int
	RefreshSecret    string
	RefreshExpHours  int
}

// SMSConfig contém configurações para o serviço de SMS
type SMSConfig struct {
	APIKey       string
	APISecret    string
	ServiceURL   string
	SenderID     string
}

// LoadConfig carrega todas as configurações do ambiente
func LoadConfig() (*Config, error) {
	// Carrega variáveis de ambiente do arquivo .env se existir
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	// Configurações do servidor
	serverPort := getEnv("SERVER_PORT", "8080")
	readTimeout, _ := strconv.Atoi(getEnv("SERVER_READ_TIMEOUT", "15"))
	writeTimeout, _ := strconv.Atoi(getEnv("SERVER_WRITE_TIMEOUT", "15"))
	idleTimeout, _ := strconv.Atoi(getEnv("SERVER_IDLE_TIMEOUT", "60"))

	// Configurações do banco de dados
	dbURI := getEnv("MONGODB_URI", "mongodb+srv://NAMUETHO:knOef2hbvwvcqsJe@cluster0.dbgb4.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")
	dbName := getEnv("MONGODB_NAME", "anamalala")
	dbPoolSize, _ := strconv.ParseUint(getEnv("MONGODB_POOL_SIZE", "100"), 10, 64)
	dbMaxIdleTime, _ := strconv.Atoi(getEnv("MONGODB_MAX_IDLE_TIME", "30"))

	// Configurações JWT
	jwtSecret := getEnv("JWT_SECRET", "anamalala_secret_key")
	jwtExpHours, _ := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	jwtRefreshSecret := getEnv("JWT_REFRESH_SECRET", "anamalala_refresh_key")
	jwtRefreshExpHours, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRATION_HOURS", "168")) // 7 dias

	// Configurações SMS
	smsAPIKey := getEnv("SMS_API_KEY", "")
	smsAPISecret := getEnv("SMS_API_SECRET", "")
	smsServiceURL := getEnv("SMS_SERVICE_URL", "")
	smsSenderID := getEnv("SMS_SENDER_ID", "ANAMALALA")

	return &Config{
		Enviroment: "dev",
		Server: ServerConfig{
			Port:         serverPort,
			ReadTimeout:  time.Duration(readTimeout) * time.Second,
			WriteTimeout: time.Duration(writeTimeout) * time.Second,
			IdleTimeout:  time.Duration(idleTimeout) * time.Second,
		},
		Database: DatabaseConfig{
			URI:        dbURI,
			Name:       dbName,
			PoolSize:   dbPoolSize,
			MaxIdleTime: time.Duration(dbMaxIdleTime) * time.Minute,
		},
		JWT: JWTConfig{
			Secret:          jwtSecret,
			ExpirationHours: jwtExpHours,
			RefreshSecret:   jwtRefreshSecret,
			RefreshExpHours: jwtRefreshExpHours,
		},
		SMS: SMSConfig{
			APIKey:     smsAPIKey,
			APISecret:  smsAPISecret,
			ServiceURL: smsServiceURL,
			SenderID:   smsSenderID,
		},
	}, nil
}

// getEnv obtém uma variável de ambiente ou retorna um valor padrão se não estiver definida
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
