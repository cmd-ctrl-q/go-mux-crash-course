package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/cmd-ctrl-q/go-mux-crash-course/cache"
	"github.com/cmd-ctrl-q/go-mux-crash-course/entity"
	"github.com/cmd-ctrl-q/go-mux-crash-course/repository"
	"github.com/cmd-ctrl-q/go-mux-crash-course/service"
	"github.com/stretchr/testify/assert"
)

const (
	ID    int64  = 123
	TITLE string = "Title 1"
	TEXT  string = "Text 1"
)

var (
	postRepo       repository.PostRepository = repository.NewPostgresRepository()
	postSrv        service.PostService       = service.NewPostService(postRepo)
	postCacheSrv   cache.PostCache           = cache.NewRedisCache("localhost:6379", 0, 10)
	postController PostController            = NewPostController(postSrv, postCacheSrv)
)

func setup() {
	var post entity.Post = entity.Post{
		ID:    ID,
		Title: TITLE,
		Text:  TEXT,
	}
	postRepo.Save(&post)
}

func tearDown(postID int64) {
	var post entity.Post = entity.Post{
		ID: postID,
	}
	postRepo.Delete(&post)
}

func TestAddPost(t *testing.T) {

	setup()

	// TODO (steps):
	// 1. Create a new HTTP Post request
	var jsonData = []byte(`{"title": "` + TITLE + `", "txt": "` + TEXT + `"}`)
	// create new post request
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonData))

	// 2. Assign HTTP Handler function (controller AddPost function)
	handler := http.HandlerFunc(postController.AddPost)

	// 3. Record HTTP Response (httptest)
	resp := httptest.NewRecorder()

	// 4. Dispatch the HTTP request
	handler.ServeHTTP(resp, req)

	// 5. Add Assertions on the HTTP Status Code and the response
	status := resp.Code

	if status != http.StatusOK {
		t.Errorf("Handler returned a wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Decode the HTTP response and transform the json into a Post struct
	var post entity.Post
	json.NewDecoder(io.Reader(resp.Body)).Decode(&post) // cast resp.Body to a Reader type

	// Assert HTTP response - make an assertion on the data of the post
	assert.NotNil(t, ID, post.ID) // check if identifier is not nil
	assert.Equal(t, TITLE, post.Title)
	assert.Equal(t, TEXT, post.Text)

	// clean up database
	tearDown(ID)
}

// Doesnt work
func TestGetAllPosts(t *testing.T) {

	// Insert new post
	setup()

	// Create new HTTP request
	req, _ := http.NewRequest("GET", "/posts", nil)

	// Assign HTTP Requesthandler function (controller function)
	handler := http.HandlerFunc(postController.GetAllPosts)
	// Record the HTTP Response
	response := httptest.NewRecorder()
	// Dispatch the HTTP Request
	handler.ServeHTTP(response, req)

	// Assert HTTP status
	status := response.Code
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Decode HTTP respones
	var posts []entity.Post
	json.NewDecoder(io.Reader(response.Body)).Decode(&posts)

	// Assert HTTP response
	assert.Equal(t, ID, posts[0].ID)
	assert.Equal(t, TITLE, posts[0].Title)
	assert.Equal(t, TEXT, posts[0].Text)
}

func TestGetPostByID(t *testing.T) {

	// insert new post
	setup()

	// Create new HTTP request
	req, _ := http.NewRequest("GET", "/posts/"+strconv.FormatInt(ID, 10), nil)

	// Assign HTTP Requesthandler function (controller function)
	handler := http.HandlerFunc(postController.GetPostByID)
	// Record the HTTP Response
	response := httptest.NewRecorder()
	// Dispatch the HTTP Request
	handler.ServeHTTP(response, req)

	// Assert HTTP status
	status := response.Code
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Decode HTTP respones
	var post entity.Post
	json.NewDecoder(io.Reader(response.Body)).Decode(&post)

	// Assert HTTP response
	assert.Equal(t, ID, post.ID)
	assert.Equal(t, TITLE, post.Title)
	assert.Equal(t, TEXT, post.Text)

	// Cleanup database
	tearDown(ID)
}
