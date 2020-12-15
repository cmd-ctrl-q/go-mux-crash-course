package repository

import (
	"database/sql"
	"log"
	"os"

	"github.com/cmd-ctrl-q/go-mux-crash-course/entity"
)

type sqliteRepo struct{}

// NewSQLiteRepository ...
func NewSQLiteRepository() PostRepository {
	os.Remove("./posts.db")

	db, err := sql.Open("sqlite3", "./posts.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	CREATE TABLE posts (id INTEGER NOT NULL PRIMARY KEY, title TEXT, txt TEXT);
	DELETE FROM posts; 
	`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}
	return &sqliteRepo{}
}

// Save ...
func (*sqliteRepo) Save(post *entity.Post) (*entity.Post, error) {
	db, err := sql.Open("sqlite3", "./posts.db")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	stmt, err := tx.Prepare("INSERT INTO posts(id, title, txt) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(post.ID, post.Title, post.Text)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return post, nil
}

// FindAll TODO (see postgres-repo.go)
func (*sqliteRepo) FindAll() ([]entity.Post, error) {

	db, err := sql.Open("sqlite3", "./posts.db")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	rows, err := db.Query("SELECT id, title, txt FROM posts")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	var posts []entity.Post = []entity.Post{}
	for rows.Next() {
		var id int64
		var title string
		var text string
		err = rows.Scan(&id, &title, &text)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		post := entity.Post{
			ID:    id,
			Title: title,
			Text:  text,
		}
		posts = append(posts, post)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return posts, nil
}

// FindOne ...
func (*sqliteRepo) FindOne(id int64) (*entity.Post, error) {

	db, err := sql.Open("sqlite3", "./posts.db")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	row := db.QueryRow("select id, title, txt from posts where id = ?", id)

	var post entity.Post
	if row != nil {
		var id int64
		var title string
		var text string
		err := row.Scan(&id, &title, &text)
		if err != nil {
			return nil, err
		} else {
			post = entity.Post{
				ID:    id,
				Title: title,
				Text:  text,
			}
		}
	}

	return &post, nil
}

func (*sqliteRepo) Delete(post *entity.Post) error {

	db, err := sql.Open("sqlite3", "./posts.db")
	if err != nil {
		log.Fatal(err)
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}
	stmt, err := tx.Prepare("DELETE FROM posts WHERE id = ?")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(post.ID)
	if err != nil {
		log.Fatal(err)
		return err
	}
	tx.Commit()
	return nil
}
