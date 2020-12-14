package controller

import (
	"encoding/json"
	"net/http"

	"github.com/cmd-ctrl-q/go-mux-crash-course/entity"
	"github.com/cmd-ctrl-q/go-mux-crash-course/errors"
	"github.com/cmd-ctrl-q/go-mux-crash-course/service"
)

type controller struct{}

var (
	// postService service.PostService = service.NewPostService() // old
	postService service.PostService // new, now added to NewPostController
)

type PostController interface {
	GetAllPosts(w http.ResponseWriter, r *http.Request)
	AddPost(w http.ResponseWriter, r *http.Request)
}

func NewPostController(service service.PostService) PostController {
	postService = service
	return &controller{}
}

func (*controller) GetAllPosts(w http.ResponseWriter, r *http.Request) {

	// specify headers
	w.Header().Set("Content-Type", "application/json")

	// get posts from firebase database
	posts, err := postService.FindAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
		json.NewEncoder(w).Encode(errors.ServiceError{Message: "Error getting the posts"})
		return
	}

	w.WriteHeader(http.StatusOK) // 200

	// encode the data and send it back to client
	json.NewEncoder(w).Encode(posts)

}

func (*controller) AddPost(w http.ResponseWriter, r *http.Request) {

	// set header
	w.Header().Set("Content-Type", "application/json")

	// get data from the response and store in a temp struct
	var post entity.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
		json.NewEncoder(w).Encode(errors.ServiceError{Message: "Error adding the post"})
		return
	}

	// validate post
	err1 := postService.Validate(&post)
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
		json.NewEncoder(w).Encode(errors.ServiceError{Message: err1.Error()})
		return
	}

	result, err2 := postService.Create(&post)
	if err2 != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
		json.NewEncoder(w).Encode(errors.ServiceError{Message: "Error saving the post"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
