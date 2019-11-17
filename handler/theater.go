package handler

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
)

// Theater !
func Theater(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/theater")
	db, err := sql.Open("mysql", "root:thisiscool@/theater_schema")
	if err != nil {
		log.Fatal(err)
	}

	switch true {
	case strings.HasPrefix(r.URL.Path, "/movies"):
		Movies(w, r, db)
		break
	}
}
