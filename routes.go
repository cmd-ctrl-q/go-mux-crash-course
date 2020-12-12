package main

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/cmd-ctrl-q/go-mux-crash-course/entity"
	"github.com/cmd-ctrl-q/go-mux-crash-course/repository"
)

var (
	// posts []Post
	repo repository.PostRepository = repository.NewPostRepository()
)

// func init() {
// 	posts = []Post{Post{ID: 1, Title: "Title 1", Text: "Text text text text text..."}}
// }

func getAllPosts(w http.ResponseWriter, r *http.Request) {

	// specify headers
	w.Header().Set("Content-Type", "application/json")

	// get posts from firebase database
	posts, err := repo.FindAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
		w.Write([]byte(`{"error": "Error getting the posts from firebase"}`))
		return
	}

	// encode data (or marshal) and send back
	// result, err := json.Marshal(posts)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError) // 500
	// 	w.Write([]byte(`{"error": "Error marshalling the post array"}`))
	// 	return
	// }

	w.WriteHeader(http.StatusOK) // 200

	// encode the data and send it back to client
	json.NewEncoder(w).Encode(posts)

	// w.Write(result)

}

func addPost(w http.ResponseWriter, r *http.Request) {

	// set header
	w.Header().Set("Content-Type", "application/json")

	// get data from the response and store in a temp struct
	var post entity.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
		w.Write([]byte(`{"error": "Error unmarshalling the post"}`))
		return
	}

	// add id
	// post.ID = len(posts) + 1 // old
	post.ID = rand.Int63() // new

	// append to data
	// posts = append(posts, post) // old
	repo.Save(&post)

	// status ok
	w.WriteHeader(http.StatusOK)

	// encode data and send back to client
	json.NewEncoder(w).Encode(post)

	// encode post using marshal before sending it back to client
	// result, err := json.Marshal(post)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError) // 500
	// 	w.Write([]byte(`{"error": "Error marshalling the post"}`))
	// 	return
	// }

	// w.Write(result)

}
