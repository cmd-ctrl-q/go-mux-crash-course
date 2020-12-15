package service

import (
	"testing"

	"github.com/cmd-ctrl-q/go-mux-crash-course/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mock repository structure will implement the PostRepository interface
type MockRepository struct {
	mock.Mock
}

func (mock *MockRepository) Save(post *entity.Post) (*entity.Post, error) {
	args := mock.Called() // return the arguments
	result := args.Get(0) // 0 for 0th index argument

	// type assertion result.(*entity.Post)
	return result.(*entity.Post), args.Error(1) // error is in 1st index
}

func (mock *MockRepository) FindAll() ([]entity.Post, error) {
	args := mock.Called()
	result := args.Get(0)

	return result.([]entity.Post), args.Error(1)
}

func TestFindAll(t *testing.T) {
	// create variable for mock repository
	mockRepo := new(MockRepository) // new reference to repo

	var identifier int64 = 1

	post := entity.Post{ID: 1, Title: "A", Text: "B"}

	// set up expectation
	mockRepo.On("FindAll").Return([]entity.Post{post}, nil)

	// create test service
	testService := NewPostService(mockRepo)

	// call the FindAll function from this service
	result, _ := testService.FindAll()

	// create Mock Assertion: Behavioral
	mockRepo.AssertExpectations(t) // assertion on the expectation created above

	// create Data Assertion on each attribute
	assert.Equal(t, identifier, result[0].ID)
	assert.Equal(t, "a", result[0].Title)
	assert.Equal(t, "b", result[0].Text)
}

func TestCreate(t *testing.T) {
	mockRepo := new(MockRepository)
	// var identifier int64 = 1
	post := entity.Post{Title: "A", Text: "B"}

	// Set expectations. when the save function of the mock repo is invoked, it will return a ref to the post
	mockRepo.On("Save").Return(&post, nil)

	// create service were testing
	testService := NewPostService(mockRepo)

	// call the create funcction of the service
	result, err := testService.Create(&post) //pass ref of post

	// add assertions
	// assertion for mockrepo
	mockRepo.AssertExpectations(t)

	// Data Assertion
	// asssert.Equal(t, identifier, result.ID)  // we expect identifier = 1, from result.ID, (doesnt work because random)
	assert.NotNil(t, result.ID) // we expect identifier = 1, from result.ID
	assert.Equal(t, "A", result.Title)
	assert.Equal(t, "B", result.Text)
	assert.Nil(t, err) // we expect this error to be nil if the test passes

	// two options for testing identifier since it is randomly generated
	// 1. mock the rand library and set the value that we expect when the Create method is called
	// 2. change the assertion and expect a NotNil value
}

// test validate message
// create two test cases, if post == nil, if post.Title == ""

func TestValidateEmptyPost(t *testing.T) {
	// make instance of the service
	testService := NewPostService(nil)

	err := testService.Validate(nil) // assuming post is nil
	// add assertion from testify library, use for assertion on data, and mocks

	assert.NotNil(t, err) // if post passes

	// Equal(t, expectedValue, actualValue)
	assert.Equal(t, "The post is empty", err.Error()) // Error() returns the error message
}

func TestValidateEmtpyTitle(t *testing.T) {
	testService := NewPostService(nil)
	post := entity.Post{ID: 1, Title: "", Text: "lorem ipsum"}

	err := testService.Validate(&post)

	assert.NotNil(t, err) // if post passes

	assert.Equal(t, "The post title is empty", err.Error()) // if post fails, check error
}
