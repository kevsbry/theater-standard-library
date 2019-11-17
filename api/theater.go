package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//movieRating
type movieRating struct {
	Source string
	Value  string
}

// imdbJSON ! imdb data
type imdbJSON struct {
	Title      string
	Year       string
	Rated      string
	Released   string
	Runtime    string
	Genre      string
	Director   string
	Writer     string
	Actors     string
	Plot       string
	Language   string
	Country    string
	Awards     string
	Poster     string
	Ratings    []movieRating
	MetaScore  string
	imdbRating string
	imdbVotes  string
	imdbID     string
	Type       string
	Dvd        string
	BoxOffice  string
	Production string
	Website    string
	Response   string
}

// Theater !
func Theater(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
	// r.Header.Set("Access-Control-Allow-Origin", "*")

	db, err := sql.Open("mysql", "root:thisiscool@/theater_schema")
	if err != nil {
		log.Fatal(err)
	}

	paths := strings.Split(r.URL.Path, "/")

	switch true {
	case strings.HasPrefix(r.URL.Path, "/movie"):
		if paths[2] != "" {
			Movie(w, r, db, paths[2])
		}
		return
	case strings.HasPrefix(r.URL.Path, "/category"):
		if paths[2] != "" {
			Movies(w, r, db, paths[2])
		}
		return
	}
}

// Movie ! specific movie
func Movie(w http.ResponseWriter, r *http.Request, db *sql.DB, movieTitle string) {
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/movie")

	type Movie struct {
		MovieDetails imdbJSON
		MediaSource  string
	}

	//get movie

	var movieSource string
	var movieKey string

	if movieTitle != "" {
		wildcard := "%"
		result, err := db.Query(fmt.Sprintf("SELECT IMDB FROM movie WHERE Title LIKE '%s%s%s'", wildcard, movieTitle, wildcard))
		if err != nil {
			log.Fatal(err)
		}
		defer result.Close()
		for result.Next() {
			result.Scan(&movieKey)
		}

		// get Movie source
		res, err := db.Query(fmt.Sprintf("SELECT MovieSource FROM movie WHERE IMDB = '%s'", movieKey))
		if err != nil {
			log.Fatal(err)
		}
		defer res.Close()
		for res.Next() {
			res.Scan(&movieSource)
		}

		//get IMDB data
		const imdbKey = "92176bca"
		response, err := http.Get(fmt.Sprintf("https://www.omdbapi.com/?apikey=%s&i=%s", imdbKey, movieKey))
		if err != nil {
			log.Fatal(err)
		}
		respByte, _ := ioutil.ReadAll(response.Body)
		movieJSON := imdbJSON{}
		json.Unmarshal(respByte, &movieJSON)

		resp := map[string]interface{}{}
		resp["movie"] = Movie{
			MovieDetails: movieJSON,
			MediaSource:  movieSource,
		}

		json.NewEncoder(w).Encode(resp)
	}

}

// Movies ! multiple movies
func Movies(w http.ResponseWriter, r *http.Request, db *sql.DB, category string) {
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/category")

	// get movie key from DB
	type Movie struct {
		MovieKey    string
		MediaSource string
	}
	movie := []Movie{}
	tempMovie := Movie{}

	var q string

	if category == "all" {
		q = "SELECT IMDB, MovieSource FROM movie"
	} else {
		q = fmt.Sprintf("SELECT IMDB, MovieSource FROM movie WHERE Category_FK = (SELECT CategoryID FROM category WHERE CategoryName like '%s')", category)
	}

	result, err := db.Query(q)

	if err != nil {
		log.Fatal(err)
	}
	defer result.Close()

	for result.Next() {
		result.Scan(&tempMovie.MovieKey, &tempMovie.MediaSource)
		movie = append(movie, tempMovie)
	}

	const imdbKey = "92176bca"
	type Context struct {
		MovieKey    string
		MediaSource string
		MovieInfo   imdbJSON
	}
	context := []Context{}

	for _, m := range movie {
		reponse, err := http.Get(fmt.Sprintf("https://www.omdbapi.com/?apikey=%s&i=%s", imdbKey, m.MovieKey))
		if err != nil {
			log.Fatal(err)
		}
		responseByte, _ := ioutil.ReadAll(reponse.Body)
		imdb := imdbJSON{}
		_ = json.Unmarshal(responseByte, &imdb)
		context = append(context, Context{
			MovieKey:    m.MovieKey,
			MediaSource: m.MediaSource,
			MovieInfo:   imdb,
		})

	}

	json.NewEncoder(w).Encode(context)
}
