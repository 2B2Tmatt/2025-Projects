package helper

import (
	"database/sql"
	"net/http"
	"text/template"
	"time"
)

type User struct {
	User string
}

type Post struct {
	Username   string
	Title      string
	Body       string
	Created_at time.Time
}

type PageData struct {
	User  *User
	Posts []Post
}

func GetUsername(db *sql.DB, uid int64) (*User, error) {
	user := User{""}
	err := db.QueryRow(`SELECT display_name FROM users WHERE id=$1`, uid).Scan(&user.User)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func GetPost(db *sql.DB) ([]Post, error) {
	var posts []Post
	rows, err := db.Query(`
	SELECT users.display_name, posts.title, posts.body, posts.created_at 
	FROM users 
	INNER JOIN posts 
	ON users.id = posts.user_id
	ORDER BY created_at DESC`)
	if err != nil {
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.Username, &p.Title, &p.Body, &p.Created_at)
		if err != nil {
			return posts, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func GetUserPost(db *sql.DB, username string) ([]Post, error) {
	var posts []Post
	rows, err := db.Query(`
	SELECT users.display_name, posts.title, posts.body, posts.created_at 
	FROM users 
	INNER JOIN posts 
	ON users.id = posts.user_id
	WHERE users.display_name = $1
	ORDER BY created_at DESC`, username)
	if err != nil {
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.Username, &p.Title, &p.Body, &p.Created_at)
		if err != nil {
			return posts, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func Render(w http.ResponseWriter, data any, filepath ...string) error {
	tmpl, err := template.ParseFiles(filepath...)
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return err
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "render error", http.StatusInternalServerError)
		return err
	}
	return nil
}
