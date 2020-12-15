package repository

import "github.com/cmd-ctrl-q/go-mux-crash-course/entity"

// PostRepository ...
type PostRepository interface {
	Save(post *entity.Post) (*entity.Post, error)
	FindAll() ([]entity.Post, error)
}
