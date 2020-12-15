package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/cmd-ctrl-q/go-mux-crash-course/entity"
	_ "github.com/lib/pq"
)

type postgresRepo struct{}

var (
	db  *sql.DB
	err error
)

func init() {
	db, err = sql.Open("postgres", "postgres://postgres:password@localhost/muxdb?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("Database connection successful")
}

// NewPostgresRepository ...
func NewPostgresRepository() PostRepository {

	// db, err = sql.Open("postgres", "./posts.db")
	// connect to db
	db, err = sql.Open("postgres", "postgres://postgres:password@localhost/muxdb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// delete table in database
	_, err := db.Exec("DROP TABLE posts;")

	sqlStmt := `CREATE TABLE posts (id BIGINT PRIMARY KEY NOT NULL, title TEXT NOT NULL, txt TEXT); DELETE FROM posts;`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}
	return &postgresRepo{}
}

func (*postgresRepo) Save(post *entity.Post) (*entity.Post, error) {

	db, err = sql.Open("postgres", "postgres://postgres:password@localhost/muxdb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer db.Close()

	_, err := db.Exec("INSERT INTO posts (id, title, txt) VALUES ($1, $2, $3);", post.ID, post.Title, post.Text)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// NOTE: Tx transaction doesnt work
	// tx, err := db.Begin()
	// if err != nil {
	// 	log.Fatal(err)
	// 	return nil, err
	// }

	// stmt, err := tx.Prepare("INSERT INTO posts (id, title, txt) VALUES ($1, $2, $3);")
	// if err != nil {
	// 	log.Fatal(err)
	// 	return nil, err
	// }
	// defer stmt.Close()

	// _, err = stmt.Exec(post.ID, post.Title, post.Text)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return nil, err
	// }

	return post, nil
}

// FindAll TODO
func (*postgresRepo) FindAll() ([]entity.Post, error) {

	rows, err := db.Query("SELECT * FROM posts;")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		post := entity.Post{}
		err = rows.Scan(&post.ID, &post.Title, &post.Text)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return posts, nil
}
