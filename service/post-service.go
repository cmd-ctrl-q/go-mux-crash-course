package service

import (
	"errors"
	"math/rand"

	"github.com/cmd-ctrl-q/go-mux-crash-course/entity"
	"github.com/cmd-ctrl-q/go-mux-crash-course/repository"
)

// PostService interface
type PostService interface {
	Validate(post *entity.Post) error
	Create(post *entity.Post) (*entity.Post, error)
	FindAll() ([]entity.Post, error)
	FindByID(id int64) (*entity.Post, error)
}

// implement interface
type service struct{}

var (
	// new reference to Firestore Repository
	// repo repository.PostRepository = repository.NewFirestoreRepository()
	postRepo repository.PostRepository
)

// NewPostService implements the PostService interface and returns a PostService struct
func NewPostService(repo repository.PostRepository) PostService {
	postRepo = repo
	return &service{}
}

func (s *service) Validate(post *entity.Post) error {
	// if post is emtpy
	if post == nil {
		err := errors.New("The post is empty")
		return err
	}
	if post.Title == "" {
		err := errors.New("The post title is empty")
		return err
	}
	return nil
}

func (s *service) Create(post *entity.Post) (*entity.Post, error) {
	post.ID = rand.Int63()
	return postRepo.Save(post)
}
func (s *service) FindAll() ([]entity.Post, error) {
	return postRepo.FindAll()
}

func (s *service) FindByID(id int64) (*entity.Post, error) {
	return postRepo.FindOne(id)
}
