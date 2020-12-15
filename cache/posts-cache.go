package cache

import "github.com/cmd-ctrl-q/go-mux-crash-course/entity"

type PostCache interface {
	Set(key string, value *entity.Post)
	Get(key string) *entity.Post
}
