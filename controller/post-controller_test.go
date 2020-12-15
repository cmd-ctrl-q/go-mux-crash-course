package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cmd-ctrl-q/go-mux-crash-course/entity"
	"github.com/cmd-ctrl-q/go-mux-crash-course/repository"
	"github.com/cmd-ctrl-q/go-mux-crash-course/service"
	"github.com/stretchr/testify/assert"
)

// type MockRepository struct {
// 	mock.Mock
// }

// func (mock *MockRepository) GetAllPosts(w http.ResponseWriter, r *http.Request) {
// 	args := mock.Called()
// 	result := args.Get(0)

// 	// return result.([]entity.Post), args.Error(1)
// }

// func (mock *MockRepository) AddPost(w http.ResponseWriter, r *http.Request) {
// 	args := mock.Called()
// 	result := args.Get(0)

// 	// return result.([]entity.Post), args.Error(1)
// }

const (
	TITLE string = "Title 1"
	TEXT  string = "Text 1"
)

var (
	postRepo       repository.PostRepository = repository.NewSQLiteRepository()
	postSrv        service.PostService       = service.NewPostService(postRepo)
	postController PostController            = NewPostController(postSrv)
)

func TestAddPost(t *testing.T) {

	// TODO (steps):
	// 1. Create a new HTTP Post request
	var jsonData = []byte(`{"title": "` + TITLE + `", "text": "` + TEXT + `"}`)
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

	// Assert HTTP response
	assert.NotNil(t, post.ID) // check if identifier is not nil
	assert.Equal(t, TITLE, post.Title)
	assert.Equal(t, TEXT, post.Text)
}

func TestGetAllPosts(t *testing.T) {

}
