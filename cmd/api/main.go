package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/anamalala/internal/config"
	"github.com/anamalala/internal/handlers"
	"github.com/anamalala/internal/middlewares"
	"github.com/anamalala/internal/repositories/mongodb"
	routers "github.com/anamalala/internal/router"
	"github.com/anamalala/internal/services"
	"github.com/anamalala/internal/utils"
	"github.com/anamalala/pkg/logger"
	"github.com/anamalala/pkg/sms"
)

var smsConfig = sms.SMSConfig{}

func main() {
	// Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	// Inicializar configuração
	cfg, _ := config.LoadConfig()

	// Inicializar logger
	appLogger := logger.NewLogger(cfg.Enviroment)
	appLogger.Info("Iniciando servidor da API ANAMALALA...")

	// Conectar ao MongoDB
	mongoClient, err := connectToMongoDB(cfg)
	if err != nil {
		appLogger.Fatal("Falha ao conectar ao MongoDB:", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := mongoClient.Close(ctx); err != nil {
			appLogger.Error("Erro ao desconectar do MongoDB:", err)
		}
	}()
	appLogger.Info("Conectado ao MongoDB com sucesso")

	// Inicializar repositórios
	var lock = &sync.RWMutex{}

	appLogger.Info("Iniciando repositorios")

	userRepo := mongodb.NewUserRepository(&mongoClient)
	postRepo := mongodb.NewPostRepository(&mongoClient)
	infoRepo := mongodb.NewInformationRepository(&mongoClient)
	commentRepo := mongodb.NewCommentRepository(&mongoClient)
	suggestionRepo := mongodb.NewSuggestionRepository(&mongoClient)

	// Inicializar utilitários

	appLogger.Info(" A Inicializar utilitários")

	tokenUtil := utils.NewTokenUtil(cfg.JWT.Secret, time.Duration(time.Hour.Abs()*5))
	validator := utils.NewValidator()

	// Inicializar serviço de SMS
	// smsConfig.APIKey = ""
	// smsConfig.ProviderType = ""
	// smsConfig.SenderID = ""
	smsService, _ := sms.NewService(&smsConfig, appLogger)

	// Inicializar serviços

	appLogger.Info(" A Inicializar serviços")

	authService := services.NewAuthService(userRepo, tokenUtil)
	userService := services.NewUserService(userRepo)
	infoService := services.NewInformationService(infoRepo)
	chatroomService := services.NewChatroomService(postRepo, commentRepo, userRepo)
	suggestionService := services.NewSuggestionService(suggestionRepo, userRepo)
	adminService := services.NewAdminService(userRepo, postRepo, commentRepo, smsService)

	// Inicializar handlers
	appLogger.Info(" A Inicializar handlers")

	authHandler := handlers.NewAuthHandler(authService, validator)
	userHandler := handlers.NewUserHandler(userService)
	infoHandler := handlers.NewInformationHandler(infoService)
	chatroomHandler := handlers.NewChatroomHandler(chatroomService, lock)
	suggestionHandler := handlers.NewSuggestionHandler(suggestionService)
	adminHandler := handlers.NewAdminHandler(adminService)

	// Inicializar middlewares
	appLogger.Info(" A Inicializar middlewares")

	authMiddleware := middlewares.NewAuthMiddleware(tokenUtil, userService)
	adminMiddleware := middlewares.NewAdminMiddleware(tokenUtil)
	//	loggerMiddleware := middlewares.NewLoggerMiddleware(appLogger)

	// Configurar router (Gin)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middlewares.CorsMiddleware())

	// Configurar rotas

	routers.SetupRoutes(
		router,
		authHandler,
		userHandler,
		infoHandler,
		chatroomHandler,
		suggestionHandler,
		adminHandler,
		authMiddleware,
		adminMiddleware,
	)

	// Configurar servidor HTTP
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Iniciar servidor em uma goroutine separada
	go func() {
		certFile := "/etc/letsencrypt/live/freesexy.net/fullchain.pem"
		keyFile := "/etc/letsencrypt/live/freesexy.net/privkey.pem"
		appLogger.Info("Servidor HTTP iniciado na porta 8080")
		if err := server.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Erro ao iniciar servidor:", err)
		}
	}()

	// Configurar encerramento gracioso
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	appLogger.Info("Encerrando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Fatal("Erro ao encerrar servidor:", err)
	}

	appLogger.Info("Servidor encerrado com sucesso")
}

func connectToMongoDB(cfg *config.Config) (mongodb.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongodb.Connect(ctx, cfg.Database.URI, cfg.Database.Name)
	if err != nil {
		return mongodb.Client{}, err
	}
	return *client, nil
}
