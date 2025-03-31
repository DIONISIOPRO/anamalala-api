package mongodb

// import (
// 	"context"
// 	"log"
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// // Estruturas para mapear os documentos do MongoDB
// type User struct {
// 	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	Name     string             `bson:"name" json:"name"`
// 	Province string             `bson:"province" json:"province"`
// 	Contact  string             `bson:"contact" json:"contact"`
// }

// type Post struct {
// 	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	Content   string             `bson:"content" json:"content"`
// 	AuthorID  primitive.ObjectID `bson:"author_id" json:"authorId"`
// 	Author    User               `bson:"-" json:"author"`
// 	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
// 	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
// 	Likes     int                `bson:"likes" json:"likes"`
// 	LikedBy   []primitive.ObjectID `bson:"liked_by" json:"likedBy"`
// 	Comments  []Comment          `bson:"-" json:"comments"`
// }

// type Comment struct {
// 	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	Content     string             `bson:"content" json:"content"`
// 	AuthorID    primitive.ObjectID `bson:"author_id" json:"authorId"`
// 	Author      User               `bson:"-" json:"author"`
// 	Reference   string             `bson:"reference" json:"reference"` // "post" ou "comment"
// 	ReferenceID primitive.ObjectID `bson:"reference_id" json:"referenceId"`
// 	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
// 	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
// 	Likes       int                `bson:"likes" json:"likes"`
// 	LikedBy     []primitive.ObjectID `bson:"liked_by" json:"likedBy"`
// 	Replies     []Comment          `bson:"-" json:"replies"`
// }

// // PostRepository representa o repositório de dados para postagens
// type PostRepository struct {
// 	db             *mongo.Database
// 	postsColl      *mongo.Collection
// 	commentsColl   *mongo.Collection
// 	usersColl      *mongo.Collection
// }

// // NewPostRepository cria uma nova instância do repositório de postagens
// func NewPostRepository(db *mongo.Database) *PostRepository {
// 	return &PostRepository{
// 		db:             db,
// 		postsColl:      db.Collection("posts"),
// 		commentsColl:   db.Collection("comments"),
// 		usersColl:      db.Collection("users"),
// 	}
// }

// // GetPosts busca postagens com paginação e todos os comentários aninhados
// func (r *PostRepository) GetPosts(ctx context.Context, limit, offset int) ([]Post, error) {
// 	// Definir opções de busca para paginação
// 	findOptions := options.Find().
// 		SetSort(bson.D{{Key: "created_at", Value: -1}}). // Mais recentes primeiro
// 		SetLimit(int64(limit)).
// 		SetSkip(int64(offset))

// 	// Buscar postagens
// 	cursor, err := r.postsColl.Find(ctx, bson.M{}, findOptions)
// 	if err != nil {
// 		log.Printf("Erro ao buscar postagens: %v", err)
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	// Decodificar resultados
// 	var posts []Post
// 	if err = cursor.All(ctx, &posts); err != nil {
// 		log.Printf("Erro ao decodificar postagens: %v", err)
// 		return nil, err
// 	}

// 	// Buscar autores das postagens
// 	err = r.populatePostAuthors(ctx, posts)
// 	if err != nil {
// 		log.Printf("Erro ao buscar autores das postagens: %v", err)
// 		return nil, err
// 	}

// 	// Buscar comentários para cada postagem
// 	for i := range posts {
// 		comments, err := r.getCommentsForReference(ctx, "post", posts[i].ID)
// 		if err != nil {
// 			log.Printf("Erro ao buscar comentários da postagem %s: %v", posts[i].ID.Hex(), err)
// 			continue
// 		}
// 		posts[i].Comments = comments
// 	}

// 	return posts, nil
// }

// // getCommentsForReference busca comentários recursivamente para um determinado tipo de referência e ID
// func (r *PostRepository) getCommentsForReference(ctx context.Context, referenceType string, referenceID primitive.ObjectID) ([]Comment, error) {
// 	// Buscar comentários para esta referência
// 	filter := bson.M{
// 		"reference":     referenceType,
// 		"reference_id":  referenceID,
// 	}
// 	findOptions := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}) // Mais antigos primeiro

// 	cursor, err := r.commentsColl.Find(ctx, filter, findOptions)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	// Decodificar comentários
// 	var comments []Comment
// 	if err = cursor.All(ctx, &comments); err != nil {
// 		return nil, err
// 	}

// 	// Se não houver comentários, retornar array vazio
// 	if len(comments) == 0 {
// 		return []Comment{}, nil
// 	}

// 	// Buscar informações dos autores dos comentários
// 	authorIDs := make([]primitive.ObjectID, 0, len(comments))
// 	authorMap := make(map[string]User)

// 	for _, comment := range comments {
// 		authorIDs = append(authorIDs, comment.AuthorID)
// 	}

// 	// Deduplica os IDs de autores
// 	uniqueAuthorIDs := make(map[string]primitive.ObjectID)
// 	for _, id := range authorIDs {
// 		uniqueAuthorIDs[id.Hex()] = id
// 	}

// 	// Converte o mapa de volta para slice
// 	deduplicatedAuthorIDs := make([]primitive.ObjectID, 0, len(uniqueAuthorIDs))
// 	for _, id := range uniqueAuthorIDs {
// 		deduplicatedAuthorIDs = append(deduplicatedAuthorIDs, id)
// 	}

// 	// Busca os autores
// 	if len(deduplicatedAuthorIDs) > 0 {
// 		userCursor, err := r.usersColl.Find(ctx, bson.M{"_id": bson.M{"$in": deduplicatedAuthorIDs}})
// 		if err != nil {
// 			log.Printf("Erro ao buscar autores dos comentários: %v", err)
// 		} else {
// 			defer userCursor.Close(ctx)
// 			var users []User
// 			if err = userCursor.All(ctx, &users); err != nil {
// 				log.Printf("Erro ao decodificar autores dos comentários: %v", err)
// 			} else {
// 				for _, user := range users {
// 					authorMap[user.ID.Hex()] = user
// 				}
// 			}
// 		}
// 	}

// 	// Preencher informações dos autores e buscar respostas recursivamente
// 	for i := range comments {
// 		// Preencher autor
// 		if author, ok := authorMap[comments[i].AuthorID.Hex()]; ok {
// 			comments[i].Author = author
// 		}

// 		// Buscar respostas (comentários dos comentários) recursivamente
// 		replies, err := r.getCommentsForReference(ctx, "comment", comments[i].ID)
// 		if err != nil {
// 			log.Printf("Erro ao buscar respostas do comentário %s: %v", comments[i].ID.Hex(), err)
// 			continue
// 		}
// 		comments[i].Replies = replies
// 	}

// 	return comments, nil
// }

// // populatePostAuthors busca e preenche informações dos autores das postagens
// func (r *PostRepository) populatePostAuthors(ctx context.Context, posts []Post) error {
// 	if len(posts) == 0 {
// 		return nil
// 	}

// 	// Extrair IDs dos autores das postagens
// 	authorIDs := make([]primitive.ObjectID, 0, len(posts))
// 	for _, post := range posts {
// 		authorIDs = append(authorIDs, post.AuthorID)
// 	}

// 	// Deduplica os IDs de autores
// 	uniqueAuthorIDs := make(map[string]primitive.ObjectID)
// 	for _, id := range authorIDs {
// 		uniqueAuthorIDs[id.Hex()] = id
// 	}

// 	// Converte o mapa de volta para slice
// 	deduplicatedAuthorIDs := make([]primitive.ObjectID, 0, len(uniqueAuthorIDs))
// 	for _, id := range uniqueAuthorIDs {
// 		deduplicatedAuthorIDs = append(deduplicatedAuthorIDs, id)
// 	}

// 	// Buscar autores no banco de dados
// 	filter := bson.M{"_id": bson.M{"$in": deduplicatedAuthorIDs}}
// 	cursor, err := r.usersColl.Find(ctx, filter)
// 	if err != nil {
// 		return err
// 	}
// 	defer cursor.Close(ctx)

// 	// Mapear autores por ID
// 	var users []User
// 	if err = cursor.All(ctx, &users); err != nil {
// 		return err
// 	}

// 	authorMap := make(map[string]User, len(users))
// 	for _, user := range users {
// 		authorMap[user.ID.Hex()] = user
// 	}

// 	// Atribuir autores às postagens
// 	for i := range posts {
// 		if author, ok := authorMap[posts[i].AuthorID.Hex()]; ok {
// 			posts[i].Author = author
// 		}
// 	}

// 	return nil
// }

// // GetPostByID busca uma postagem específica por ID com todos os comentários aninhados
// func (r *PostRepository) GetPostByID(ctx context.Context, postID string) (*Post, error) {
// 	id, err := primitive.ObjectIDFromHex(postID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Buscar a postagem
// 	var post Post
// 	err = r.postsColl.FindOne(ctx, bson.M{"_id": id}).Decode(&post)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Buscar autor da postagem
// 	var author User
// 	err = r.usersColl.FindOne(ctx, bson.M{"_id": post.AuthorID}).Decode(&author)
// 	if err == nil {
// 		post.Author = author
// 	}

// 	// Buscar comentários
// 	comments, err := r.getCommentsForReference(ctx, "post", post.ID)
// 	if err != nil {
// 		log.Printf("Erro ao buscar comentários da postagem %s: %v", post.ID.Hex(), err)
// 	} else {
// 		post.Comments = comments
// 	}

// 	return &post, nil
// }
