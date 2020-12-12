package repository

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/cmd-ctrl-q/go-mux-crash-course/entity"
)

type PostRepository interface {
	Save(post *entity.Post) (*entity.Post, error)
	FindAll() ([]entity.Post, error)
}

type repo struct{}

// NewPostRepository
func NewPostRepository() PostRepository {
	return &repo{}
}

const (
	projectID      string = "go-mux-crash-course"
	collectionName string = "posts"
)

func (r *repo) Save(post *entity.Post) (*entity.Post, error) {
	// access the context of the app
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create firestore client: %v", err)
		return nil, err
	}

	defer client.Close()

	// if no error, add new element to collection
	_, _, err = client.Collection(collectionName).Add(ctx, map[string]interface{}{ // map[key]
		"ID":    post.ID,
		"Title": post.Title,
		"Text":  post.Text,
	})
	if err != nil {
		log.Fatalf("Failed to add a new post (likely because of unlikely event key matches another key): %v", err)
		return nil, err
	}

	return post, nil
}

func (r *repo) FindAll() ([]entity.Post, error) {

	// create an empty context object
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create a Firestore Client: %v", err)
		return nil, err
	}

	defer client.Close()

	var posts []entity.Post
	// iterate post document from firestore and add the elements into a posts array
	iterator := client.Collection(collectionName).Documents(ctx)
	for {
		// access each element using Next()
		doc, err := iterator.Next()
		if err != nil {
			log.Fatalf("Failed to iterate the list of posts: %v", err)
			return nil, err
		}
		post := entity.Post{
			ID:    doc.Data()["ID"].(int64), // .(int) is a type assertion, throws error if not int
			Title: doc.Data()["Title"].(string),
			Text:  doc.Data()["Text"].(string),
		}
		// add to posts collection
		posts = append(posts, post)
	}
	return posts, nil
}
