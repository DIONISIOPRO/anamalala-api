package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"slices"

	"github.com/anamalala/internal/models"
	"github.com/anamalala/internal/repositories/interfaces"
)

type ChatroomService struct {
	postRepo    interfaces.PostRepository
	commentRepo interfaces.CommentRepository
	userRepo    interfaces.UserRepository
}

// CommentJob representa um job para buscar comentários de um post
type CommentJob struct {
	Post *models.Post
}

// CommentWorker representa um worker para processar jobs de comentários
type CommentWorker struct {
	ID   int
	Jobs <-chan CommentJob
	wg   *sync.WaitGroup
}

// NewCommentWorker cria um novo worker para buscar comentários
func NewCommentWorker(id int, jobs <-chan CommentJob, wg *sync.WaitGroup) *CommentWorker {
	return &CommentWorker{
		ID:   id,
		Jobs: jobs,
		wg:   wg,
	}
}

func NewChatroomService(
	postRepo interfaces.PostRepository,
	commentRepo interfaces.CommentRepository,
	userRepo interfaces.UserRepository,
) ChatroomService {
	return ChatroomService{
		postRepo:    postRepo,
		commentRepo: commentRepo,
		userRepo:    userRepo,
	}
}

func (s *ChatroomService) CreatePost(ctx context.Context, post models.Post) (models.Post, error) {
	// Verificar se o autor existe
	user, err := s.userRepo.FindByID(ctx, post.UserID)
	if err != nil {
		return models.Post{}, errors.New("autor não encontrado")
	}
	// Configurar campos da postagem
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	post.Likes = 0
	post.Author = models.Author{
		Name: user.Name,
		ID:   user.ID,
	}
	// Salvar postagem
	post, err = s.postRepo.Create(ctx, post)
	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func (s *ChatroomService) GetAllPosts(ctx context.Context, page, limit int) (models.Posts, int, error) {
	posts, total, err := s.postRepo.List(ctx, int64(page), int64(limit))
	if err != nil {
		return nil, 0, err
	}
	posts = FetchPostsComments(posts, len(posts), s.commentRepo)
	return posts, int(total), nil

}

func (s *ChatroomService) GetPostByID(ctx context.Context, id string) (models.Post, error) {
	post, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return models.Post{}, err
	}
	posts := []models.Post{post}
	posts = FetchPostsComments(posts, 1, s.commentRepo)

	for _, p := range posts {
		if p.ID == post.ID {
			return p, nil
		}
	}
	return post, nil
}

func (s *ChatroomService) CommentPost(ctx context.Context, referenceID string, comment models.Comment, authorID string) (models.Comment, error) {
	_, err := s.postRepo.FindByID(ctx, referenceID)
	if err != nil {
		return models.Comment{}, errors.New("postagem não encontrada")
	}
	user, err := s.userRepo.FindByID(ctx, authorID)
	if err != nil {
		return models.Comment{}, errors.New("usuario que comenta não encontrada")
	}
	comment.ReferenceID = referenceID
	comment.UserID = authorID
	comment.CreatedAt = time.Now()
	comment.Likes = 0
	comment.Reference = "post"
	comment.Author = models.Author{
		Name: user.Name,
		ID:   user.ID,
	}
	err = s.commentRepo.Create(ctx, comment)
	if err != nil {
		return models.Comment{}, err
	}
	return comment, nil
}

func (s *ChatroomService) ReplayComment(ctx context.Context, referenceID string, comment models.Comment, authorID string) (models.Comment, error) {
	_, err := s.commentRepo.FindByID(ctx, referenceID)
	if err != nil {
		return models.Comment{}, errors.New("comentario não encontrado")
	}
	user, err := s.userRepo.FindByID(ctx, authorID)
	if err != nil {
		return models.Comment{}, errors.New("usuario que comenta não encontrada")
	}
	comment.ReferenceID = referenceID
	comment.UserID = authorID
	comment.CreatedAt = time.Now()
	comment.Likes = 0
	comment.Reference = "comment"
	comment.Author = models.Author{
		Name: user.Name,
		ID:   user.ID,
	}
	comment.Comments = []models.Comment{}
	err = s.commentRepo.Create(ctx, comment)
	if err != nil {
		return models.Comment{}, err
	}
	return comment, nil
}

func (s *ChatroomService) GetCommentsByPostID(ctx context.Context, postID string, page, limit int) (models.Comments, int, error) {
	comments, total, err := s.commentRepo.ListByPostID(ctx, postID, int64(page), int64(limit))
	if err != nil {
		return nil, 0, err
	}
	return comments, int(total), nil
}

func (s *ChatroomService) LikePost(ctx context.Context, postID string, userID string) (models.Post, error) {
	// Obter postagem
	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		return models.Post{}, err
	}
	// Verificar se usuário já curtiu
	if slices.Contains(post.LikedUserId, userID) {
		// Usuário já curtiu, remover curtida (toggle)
		err = s.postRepo.RemoveLike(ctx, postID, userID)
		post.Likes = post.Likes - 1
	} else {
		err = s.postRepo.AddLike(ctx, postID, userID)
		post.Likes = post.Likes + 1

	}
	// Adicionar curtida

	if err != nil {
		return models.Post{}, err
	}

	posts := FetchPostsComments([]models.Post{post}, 2, s.commentRepo)

	for _, p  := range posts {
		if p.ID == post.ID{
			post = p
		}
	}

	return post, nil
}

func (s *ChatroomService) LikeComment(ctx context.Context, commentID string, userID string) error {
	// Obter comentário
	comment, err := s.commentRepo.FindByID(ctx, commentID)
	if err != nil {
		return err
	}

	// Verificar se usuário já curtiu
	if slices.Contains(comment.LikedUserId, userID) {
		// Usuário já curtiu, remover curtida (toggle)
		return s.commentRepo.RemoveLike(ctx, commentID, userID)
	}

	// Adicionar curtida
	return s.commentRepo.AddLike(ctx, commentID, userID)
}

func (s *ChatroomService) DeletePost(ctx context.Context, postID string, userID string) error {
	// Obter postagem
	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		return err
	}

	// Verificar se usuário é o autor ou administrador
	user, _ := s.userRepo.FindByID(ctx, userID)

	// Verificar se o usuário é o autor da postagem ou um administrador
	if post.UserID != user.ID && user.Role != "admin" {
		return errors.New("não autorizado a excluir esta postagem")
	}
	// Excluir todos os comentários da postagem
	err = s.commentRepo.Delete(ctx, postID)
	if err != nil {
		return err
	}

	// Excluir postagem
	return s.postRepo.Delete(ctx, postID)
}

func (s *ChatroomService) DeleteComment(ctx context.Context, commentID string, userID string) error {

	// Obter comentário
	comment, err := s.commentRepo.FindByID(ctx, commentID)
	if err != nil {
		return err
	}

	user, _ := s.userRepo.FindByID(ctx, userID)

	// Verificar se usuário é o autor ou administrador
	if user.Role != "admin" && comment.UserID != userID {
		return errors.New("não autorizado a excluir este comentário")
	}

	// Excluir comentário
	return s.commentRepo.Delete(ctx, commentID)
}

func (s *ChatroomService) GetCommentID(ctx context.Context, id string) (models.Comment, error) {
	comment, err := s.commentRepo.FindByID(ctx, id)
	if err != nil {
		return models.Comment{}, err
	}
	return comment, nil
}

// Start inicia o worker para processar jobs de comentários
func (w *CommentWorker) Start(commentRepo interfaces.CommentRepository) {
	fmt.Println("comecando um worker para fazer job de de comentario")

	go func() {
		defer w.wg.Done()
		for job := range w.Jobs {
			w.fetchCommentsWithReplies(job.Post, commentRepo)
		}
	}()
}

// fetchCommentsWithReplies busca comentários para um post específico com todas as respostas
func (w *CommentWorker) fetchCommentsWithReplies(post *models.Post, commentRepo interfaces.CommentRepository) {
	fmt.Println("comecando a procurar  comentarios")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	comments, _, _ := commentRepo.ListByPostID(ctx, post.ID, 0, 0)
	for i := range comments {
		comments[i].Comments = w.fetchNestedReplies(comments[i].ID, commentRepo)
	}
	post.Comments = comments
}

// fetchNestedReplies busca respostas aninhadas para um comentário
func (w *CommentWorker) fetchNestedReplies(commentID string, commentRepo interfaces.CommentRepository) []models.Comment {
	fmt.Println("comecando a procurar respostas de comentarios")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	replies, _, _ := commentRepo.ListByCommentID(ctx, commentID, 0, 0)
	return replies
}

// FetchPostsComments processa posts concorrentemente para buscar comentários com respostas
func FetchPostsComments(posts []models.Post, numWorkers int, commentRepo interfaces.CommentRepository) []models.Post {
	jobs := make(chan CommentJob, len(posts))
	var wg sync.WaitGroup
	// Cria workers
	for w := 1; w <= numWorkers; w++ {
		worker := NewCommentWorker(w, jobs, &wg)
		wg.Add(1)
		worker.Start(commentRepo)
	}
	for i := range posts {
		jobs <- CommentJob{Post: &posts[i]}
	}
	close(jobs)
	wg.Wait()
	return posts
}
