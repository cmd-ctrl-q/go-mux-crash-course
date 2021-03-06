package repository

import (
	"log"

	"github.com/cmd-ctrl-q/go-mux-crash-course/config"
	"github.com/cmd-ctrl-q/go-mux-crash-course/entity"
	_ "github.com/lib/pq"
)

type postgresRepo struct{}

// NewPostgresRepository ...
func NewPostgresRepository() PostRepository {

	_, err := config.DB.Exec("DROP TABLE posts;")

	sqlStmt := `CREATE TABLE posts (id BIGINT PRIMARY KEY NOT NULL, title TEXT NOT NULL, txt TEXT); DELETE FROM posts;`

	_, err = config.DB.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}
	return &postgresRepo{}
}

func (*postgresRepo) Save(post *entity.Post) (*entity.Post, error) {

	_, err := config.DB.Exec("INSERT INTO posts (id, title, txt) VALUES ($1, $2, $3);", post.ID, post.Title, post.Text)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Tx (transaction) doesnt work
	// tx, err := DB.Begin()
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

	rows, err := config.DB.Query("SELECT * FROM posts;")
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

func (*postgresRepo) FindOne(id string) (*entity.Post, error) {

	// find the row with matching id as the param id
	row := config.DB.QueryRow("SELECT * FROM posts WHERE id = $1", id)

	// create post object
	var post entity.Post

	// scan the columns of the row and store in post attributes
	err := row.Scan(&post.ID, &post.Title, &post.Text)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (*postgresRepo) Delete(post *entity.Post) error {

	// _, err := config.DB.Exec("INSERT INTO posts (id, title, txt) VALUES ($1, $2, $3);", post.ID, post.Title, post.Text)
	_, err := config.DB.Exec("DELETE FROM posts WHERE id = $1", post.ID)
	if err != nil {
		return err
	}
	return nil
}
