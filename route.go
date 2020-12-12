package main

import (
	"encoding/json"
	"net/http"
)

type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

var (
	posts []Post
)

func init() {
	posts = []Post{Post{ID: 1, Title: "Title 1", Text: "Text text text text text..."}}
}

func getAllPosts(w http.ResponseWriter, r *http.Request) {

	// specify headers
	w.Header().Set("Content-Type", "application/json")

	// encode data and send back
	// json.NewEncoder(w).Encode(posts)

	// or marshal
	result, err := json.Marshal(posts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
		w.Write([]byte(`{"error": "Error marshalling the post array"}`))
		return
	}

	w.WriteHeader(http.StatusOK) // 200

	w.Write(result)

}

func addPost(w http.ResponseWriter, r *http.Request) {

	// set header
	w.Header().Set("Content-Type", "application/json")

	// get data from the response and store in a temp struct
	var post Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
		w.Write([]byte(`{"error": "Error unmarshalling the post"}`))
		return
	}

	// add id
	post.ID = len(posts) + 1

	// append to data
	posts = append(posts, post)

	// status ok
	w.WriteHeader(http.StatusOK)

	// encode post using marshal before sending it back to client
	result, err := json.Marshal(post)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
		w.Write([]byte(`{"error": "Error marshalling the post"}`))
		return
	}

	w.Write(result)

}
