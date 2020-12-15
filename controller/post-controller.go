package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/cmd-ctrl-q/go-mux-crash-course/cache"
	"github.com/cmd-ctrl-q/go-mux-crash-course/entity"
	"github.com/cmd-ctrl-q/go-mux-crash-course/errors"
	"github.com/cmd-ctrl-q/go-mux-crash-course/service"
)

type controller struct{}

var (
	// postService service.PostService = service.NewPostService() // old
	postService service.PostService // new, now added to NewPostController
	postCache   cache.PostCache     // new post cache object
)

// PostController ...
type PostController interface {
	GetAllPosts(w http.ResponseWriter, r *http.Request)
	AddPost(w http.ResponseWriter, r *http.Request)
	GetPostByID(w http.ResponseWriter, r *http.Request)
}

// NewPostController ...
func NewPostController(service service.PostService, cache cache.PostCache) PostController {
	postService = service // set global service with postService in param
	postCache = cache     // set global cache with postCache in param
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

// GetPostByID from service and store value into cache. (faster to get post from cache if possible)
func (*controller) GetPostByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get the id from the url
	postID := strings.Split(r.URL.Path, "/")[2]

	// cache: check/get data from the postID
	var post *entity.Post = postCache.Get(postID)
	if post == nil {
		n, _ := strconv.ParseInt(postID, 10, 64)
		post, err := postService.FindByID(n)
		if err != nil {
			// cannot find id in db
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errors.ServiceError{Message: "No posts found!"})
			return
		}
		// store value into cache. Set(key, value)
		postCache.Set(postID, post)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(post)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(post)
	}
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
