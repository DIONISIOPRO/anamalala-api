package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/anamalala/internal/models"
	"github.com/anamalala/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Define WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin for now. In production, restrict this.
	},
}

// Define message types for WebSocket communication
const (
	MsgTypeNewPost    = "new_post"
	MsgTypeDeletePost = "delete_post"
	MsgTypeNewComment = "new_comment"
	MsgTypeDelComment = "delete_comment"
	MsgTypeLikePost   = "like_post"
	MsgTypeLikeCmnt   = "like_comment"
)

// WebSocketMessage represents the structure of messages sent over WebSocket
type WebSocketMessage struct {
	Type    string      `json:"type"`
	Payload any `json:"payload"`
}

// ChatroomHandler handles HTTP and WebSocket requests for the chatroom
type ChatroomHandler struct {
	chatroomService services.ChatroomService
	// WebSocket connection management
	clients    map[string][]*websocket.Conn // Map of userID to connections (a user can have multiple connections)
	clientsMux *sync.RWMutex                 // Mutex for thread-safe access to clients map
}

// NewChatroomHandler creates a new instance of ChatroomHandler
func NewChatroomHandler(chatroomService services.ChatroomService, clientsMux *sync.RWMutex  ) ChatroomHandler {
	return ChatroomHandler{
		chatroomService: chatroomService,
		clients:         make(map[string][]*websocket.Conn),
		clientsMux: clientsMux ,
	}
}

// HandleWebSocket handles WebSocket connections
func (h *ChatroomHandler) HandleWebSocket(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, "Usuário não autenticado")
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("Failed to upgrade connection: %v\n", err)
		return
	}

	// Register the client
	h.registerClient(userID.(string), conn)

	// Handle disconnection when the function returns
	defer h.unregisterClient(userID.(string), conn)

	// Handle incoming messages (if needed)
	h.handleMessages(conn, userID.(string))
}

// registerClient adds a new WebSocket connection to the clients map
func (h *ChatroomHandler) registerClient(userID string, conn *websocket.Conn) {
	h.clientsMux.Lock()
	defer h.clientsMux.Unlock()

	h.clients[userID] = append(h.clients[userID], conn)
	fmt.Printf("Client registered: %s (total connections: %d)\n", userID, len(h.clients[userID]))
}
// unregisterClient removes a WebSocket connection from the clients map
func (h *ChatroomHandler) unregisterClient(userID string, conn *websocket.Conn) {
	h.clientsMux.Lock()
	defer h.clientsMux.Unlock()

	// Find and remove the connection
	connections := h.clients[userID]
	for i, c := range connections {
		if c == conn {
			// Close the connection
			conn.Close()
			// Remove from slice
			h.clients[userID] = append(connections[:i], connections[i+1:]...)
			break
		}
	}

	// If no more connections for this user, remove the user entry
	if len(h.clients[userID]) == 0 {
		delete(h.clients, userID)
	}

	fmt.Printf("Client unregistered: %s (remaining connections: %d)\n", userID, len(h.clients[userID]))
}

// handleMessages processes incoming WebSocket messages
func (h *ChatroomHandler) handleMessages(conn *websocket.Conn, userID string) {
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("WebSocket error: %v\n", err)
			}
			break // Exit the loop if there's an error
		}
		// We're not processing incoming messages for now, just keeping the connection alive
	}
}

// broadcastMessage sends a message to all connected clients
func (h *ChatroomHandler) broadcastMessage(message WebSocketMessage) {
	h.clientsMux.RLock()
	defer h.clientsMux.RUnlock()

	messageJSON, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("Error marshaling message: %v\n", err)
		return
	}

	for userID, connections := range h.clients {
		for _, conn := range connections {
			err := conn.WriteMessage(websocket.TextMessage, messageJSON)
			if err != nil {
				fmt.Printf("Error sending message to %s: %v\n", userID, err)
				// We don't remove the connection here; it will be handled on next read operation
			}
		}
	}
}

// CreatePost handles the creation of a new post
func (h *ChatroomHandler) CreatePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, "Usuário não autenticado")
		return
	}
	var content struct{
		Content string `json:"content"`
	}

	var post = models.Post{}
	if err := c.ShouldBindJSON(&content); err != nil {
		c.JSON(http.StatusBadRequest, "Dados inválidos")
		return
	}

	post.Content = content.Content

	// Validação de campos obrigatórios
	if post.Content == "" {
		c.JSON(http.StatusBadRequest, "Conteúdo é obrigatório")
		return
	}

	post.UserID = userID.(string)

	createdPost, err := h.chatroomService.CreatePost(c, post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao criar postagem")
		return
	}

	// Broadcast new post to all clients
	h.broadcastMessage(WebSocketMessage{
		Type:    MsgTypeNewPost,
		Payload: createdPost,
	})

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Postagem criada com sucesso",
		"data":    createdPost,
	})
}

// GetPosts retrieves posts with pagination
func (h *ChatroomHandler) GetPosts(c *gin.Context) {
	// Parâmetros para paginação
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 20
	}

	posts, total, err := h.chatroomService.GetAllPosts(c, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao buscar postagens")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Postagens obtidas com sucesso",
		"data": gin.H{
			"posts":      posts,
			"total":      total,
			"page":       page,
			"limit":      limit,
			"totalPages": (total + limit - 1) / limit,
		},
	})
}

// GetRecentPostsTotal retrieves the count of recent posts and comments
func (h *ChatroomHandler) GetRecentPostsTotal(c *gin.Context) {
	var totalComments int = 0
	var totalPosts int = 0
	var recentPosts []models.Post = []models.Post{}

	posts, _, err := h.chatroomService.GetAllPosts(c, 0, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao buscar postagens")
		return
	}

	for _, p := range posts {
		if time.Now().After(p.CreatedAt.Add(time.Hour * 48)) {
			continue
		}
		recentPosts = append(recentPosts, p)
	}
	totalPosts = len(recentPosts)

	for _, p := range recentPosts {
		n := len(p.Comments)
		totalComments = totalComments + n
	}

	totalPosts = totalPosts + totalComments

	c.JSON(http.StatusOK, gin.H{
		"total": totalPosts,
	})
}

// GetPostByID retrieves a post by its ID
func (h *ChatroomHandler) GetPostByID(c *gin.Context) {
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, "ID da postagem não fornecido")
		return
	}

	post, err := h.chatroomService.GetPostByID(c, postID)
	if err != nil {
		c.JSON(http.StatusNotFound, "Postagem não encontrada")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Postagem obtida com sucesso",
		"data":    post,
	})
}

// DeletePost handles deletion of a post
func (h *ChatroomHandler) DeletePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, "Usuário não autenticado")
		return
	}

	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, "ID da postagem não fornecido")
		return
	}

	err := h.chatroomService.DeletePost(c, postID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao excluir postagem")
		return
	}

	// Broadcast post deletion to all clients
	h.broadcastMessage(WebSocketMessage{
		Type: MsgTypeDeletePost,
		Payload: gin.H{
			"postID": postID,
		},
	})

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Postagem excluída com sucesso",
	})
}

// CreateComment handles the creation of a new comment
func (h *ChatroomHandler) CommentPost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, "Usuário não autenticado")
		return
	}

	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, "ID da postagem não fornecido")
		return
	}

	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, "Dados inválidos")
		return
	}

	// Validação de campos obrigatórios
	if comment.Content == "" {
		c.JSON(http.StatusBadRequest, "Conteúdo é obrigatório")
		return
	}

	comment.UserID = userID.(string)
	comment.ReferenceID = postID
	comment.Reference = "post"

	createdComment, err := h.chatroomService.CommentPost(c, postID, comment, comment.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao criar comentário")
		return
	}

	// Broadcast new comment to all clients
	h.broadcastMessage(WebSocketMessage{
		Type: MsgTypeNewComment,
		Payload: gin.H{
			"comment": createdComment,
			"postID":  postID,
		},
	})

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Comentário criado com sucesso",
		"data":    createdComment,
	})
}


func (h *ChatroomHandler) ReplayComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, "Usuário não autenticado")
		return
	}

	commentID := c.Param("id")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, "ID da postagem não fornecido")
		return
	}

	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, "Dados inválidos")
		return
	}

	// Validação de campos obrigatórios
	if comment.Content == "" {
		c.JSON(http.StatusBadRequest, "Conteúdo é obrigatório")
		return
	}

	comment.UserID = userID.(string)
	comment.ReferenceID = commentID
	comment.Reference = "comment"

	createdComment, err := h.chatroomService.ReplayComment(c, commentID, comment, comment.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao criar comentário")
		return
	}

	// Broadcast new comment to all clients
	h.broadcastMessage(WebSocketMessage{
		Type: MsgTypeNewComment,
		Payload: gin.H{
			"comment": createdComment,
			"referenceId":  commentID,
		},
	})

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Comentário criado com sucesso",
		"data":    createdComment,
	})
}

// GetCommentsByPostID retrieves comments for a specific post
func (h *ChatroomHandler) GetCommentsByPostID(c *gin.Context) {
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, "ID da postagem não fornecido")
		return
	}

	// Parâmetros para paginação
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 20
	}

	comments, total, err := h.chatroomService.GetCommentsByPostID(c, postID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao buscar comentários")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Comentários obtidos com sucesso",
		"data": gin.H{
			"comments":   comments,
			"total":      total,
			"page":       page,
			"limit":      limit,
			"totalPages": (total + limit - 1) / limit,
		},
	})
}

// DeleteComment handles deletion of a comment
func (h *ChatroomHandler) DeleteComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, "Usuário não autenticado")
		return
	}

	commentID := c.Param("id")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, "ID do comentário não fornecido")
		return
	}

	// Get comment information before deletion to know which post it belongs to
	comment, err := h.chatroomService.GetCommentID(c, commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao buscar comentário")
		return
	}

	err = h.chatroomService.DeleteComment(c, commentID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao excluir comentário")
		return
	}

	// Broadcast comment deletion to all clients
	h.broadcastMessage(WebSocketMessage{
		Type: MsgTypeDelComment,
		Payload: gin.H{
			"commentID": commentID,
			"reference_id":    comment.ReferenceID,
		},
	})

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Comentário excluído com sucesso",
	})
}

// LikeComment handles liking a comment
func (h *ChatroomHandler) LikeComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, "Usuário não autenticado")
		return
	}

	commentID := c.Param("id")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, "ID do comentario não fornecido")
		return
	}

	comment, err := h.chatroomService.GetCommentID(c, commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao buscar comentário")
		return
	}

	err = h.chatroomService.LikeComment(c, commentID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao curtir o comentario")
		return
	}

	// Broadcast comment like to all clients
	h.broadcastMessage(WebSocketMessage{
		Type: MsgTypeLikeCmnt,
		Payload: gin.H{
			"commentID": commentID,
			"reference": comment.Reference,
			"reference_id":    comment.ReferenceID,
			"userID":    userID.(string),
		},
	})

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Comentario curtido com sucesso",
	})
}

// LikePost handles liking a post
func (h *ChatroomHandler) LikePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, "Usuário não autenticado")
		return
	}

	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, "ID da postagem não fornecido")
		return
	}

	post, err := h.chatroomService.LikePost(c, postID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao curtir postagem")
		return
	}

	// Broadcast post like to all clients
	h.broadcastMessage(WebSocketMessage{
		Type: MsgTypeLikePost,
		Payload: gin.H{
			"post": post,
			"postID": postID,
			"userID": userID.(string),
		},
	})

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Postagem curtida com sucesso",
	})
}

// UnlikePost handles unliking a post
func (h *ChatroomHandler) UnlikePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, "Usuário não autenticado")
		return
	}

	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, "ID da postagem não fornecido")
		return
	}

	post, err := h.chatroomService.LikePost(c, postID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Falha ao descurtir postagem")
		return
	}

	// Broadcast post unlike to all clients
	h.broadcastMessage(WebSocketMessage{
		Type: MsgTypeLikePost, // Re-use the same message type, clients can determine action based on payload
		Payload: gin.H{
			"postID":  postID,
			"userID":  userID.(string),
			"post": post,
		},
	})

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Postagem descurtida com sucesso",
	})
}
