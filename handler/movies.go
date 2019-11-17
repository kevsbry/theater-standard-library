package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

// Movies !
func Movies(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/movies")

	type Movie struct {
		Title     string
		IMDB      string
		Source    string
		DateAdded string
		Category  string
	}
	tempMovie := Movie{}
	movies := []Movie{}

	q := fmt.Sprintf("select m.Title, m.IMDB, m.MovieSource, m.Date_Added, c.CategoryName from movie as m, category as c where c.CategoryID = m.Category_FK")
	result, err := db.Query(q)
	if err != nil {
		log.Fatal(err)
	}

	for result.Next() {
		result.Scan(&tempMovie.Title, &tempMovie.IMDB, &tempMovie.Source, &tempMovie.DateAdded, &tempMovie.Category)
		movies = append(movies, tempMovie)
	}

	t := template.New("")
	t, _ = t.ParseFiles("static/theater/html/index.html")
	err = t.ExecuteTemplate(w, "index.html", movies)
	if err != nil {
		log.Fatal(err)
	}

}
