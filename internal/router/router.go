package routers

import (
	"net/http"

	"github.com/anamalala/internal/handlers"
	"github.com/anamalala/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	authHandler handlers.AuthHandler,
	userHandler handlers.UserHandler,
	infoHandler handlers.InformationHandler,
	chatroomHandler handlers.ChatroomHandler,
	suggestionHandler handlers.SuggestionHandler,
	adminHandler handlers.AdminHandler,
	authMiddleware middlewares.AuthMiddlewares,
	adminMiddleware middlewares.AdminMiddlewares,
) {
	// Endpoints públicos
	public := r.Group("/api/v1")
	{
		// Health check
		public.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "online"})
		})

		// Autenticação
		auth := public.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/reset-password-request", authHandler.RequestPasswordReset)
			auth.POST("/reset-password", authHandler.ResetPassword)
			auth.POST("/reset-password_code_confirm", authHandler.VerifyResetPasswordToken)

		}

		// Informações (apenas leitura pública)
		info := public.Group("/info").Use(authMiddleware.AuthMiddleware())
		{
			info.GET("", infoHandler.GetAll)
			info.GET("/:id", infoHandler.GetByID)
		}
	}

	// Endpoints que exigem autenticação
	authenticated := r.Group("/api/v1")
	authenticated.Use(authMiddleware.AuthMiddleware())
	{
		// Perfil de usuário
		user := authenticated.Group("/user")
		{
			user.GET("/profile", userHandler.GetProfile)
			user.GET("/online_total", userHandler.GetTotalOnline)
			user.PUT("/profile", userHandler.UpdateProfile)
		}

		// Sala de bate-papo
		chatroom := authenticated.Group("/chatroom")
		{
			chatroom.GET("/ws",chatroomHandler.HandleWebSocket)

			chatroom.POST("/post", chatroomHandler.CreatePost)
			chatroom.GET("/posts", chatroomHandler.GetPosts)
			chatroom.GET("/recent_post_total", chatroomHandler.GetRecentPostsTotal)
			chatroom.GET("/post/:id", chatroomHandler.GetPostByID)
			chatroom.POST("/post/:id/comment", chatroomHandler.CommentPost)
			chatroom.POST("/comment/:id/comment", chatroomHandler.ReplayComment)
			chatroom.GET("/post/:id/comments", chatroomHandler.GetCommentsByPostID)
			chatroom.POST("/post/:id/like", chatroomHandler.LikePost)
			chatroom.POST("/comment/:id/like", chatroomHandler.LikeComment)
			chatroom.DELETE("/post/:id", chatroomHandler.DeletePost)
			chatroom.DELETE("/comment/:id", chatroomHandler.DeleteComment)
		}

		// Sugestões
		suggestion := authenticated.Group("/suggestions")
		{
			suggestion.POST("", suggestionHandler.CreateSuggestion)
		}

	}

	// Endpoints que exigem privilégios de administrador
	admin := r.Group("/api/v1/admin")
	admin.Use(authMiddleware.AuthMiddleware(), adminMiddleware.AdminMiddleware())
	{
		// Gestão de usuários
		admin.GET("/users", userHandler.GetAllUsers)
		admin.GET("/users/province/:province", userHandler.GetUsersByProvince)
		admin.POST("/users/:id/ban", adminHandler.BanUser)
		admin.POST("/users/:id/unban", adminHandler.UnbanUser)
		admin.POST("/users/:id/promote", adminHandler.PromoteToAdmin)
		admin.DELETE("/users/:id", adminHandler.BanUser)

		// Gestão de conteúdo (informações)
		admin.POST("/info", infoHandler.Create)
		admin.PUT("/info/:id", infoHandler.Update)
		admin.DELETE("/info/:id", infoHandler.Delete)

		// Moderação de conteúdo (posts e comentários)
		admin.DELETE("/posts/:id", chatroomHandler.DeletePost)
		admin.DELETE("/comments/:id", chatroomHandler.DeleteComment)

		// Gestão de sugestões
		admin.GET("/suggestions", suggestionHandler.GetAllSuggestions)
		admin.GET("/suggestions/:id", suggestionHandler.GetSuggestionByID)
		admin.PUT("/suggestions/:id/status", suggestionHandler.UpdateSuggestionStatus)

		// Estatísticas e dashboards
		admin.GET("/stats/users", adminHandler.GetDashboardStats)
	}
}
