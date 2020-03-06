package main

import (
	"bytes"
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/lib/pq"
	"github.com/yansal/sql/build"
)

func main() {
	addr := ":" + os.Getenv("PORT")
	if addr == ":" {
		addr = ":8080"
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "sslmode=disable"
	}
	pqconn, err := pq.NewConnector(dsn)
	if err != nil {
		log.Fatal(err)
	}
	db := sql.OpenDB(pqconn)

	http.Handle("/", &handler{db: db})
	log.SetFlags(0)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

type handler struct{ db *sql.DB }

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.serveHTTP(w, r)
	if err == nil {
		return
	}
	log.Print(err)
	w.WriteHeader(http.StatusInternalServerError)
}

func (h *handler) serveHTTP(w http.ResponseWriter, r *http.Request) error {
	buf := new(bytes.Buffer)
	r.Header.Write(buf)
	header := buf.String()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	query, args := build.InsertInto("requests").
		Values(
			build.Value("method", build.Bind(r.Method)),
			build.Value("url", build.Bind(r.URL.String())),
			build.Value("header", build.Bind(header)),
			build.Value("body", build.Bind(string(body))),
		).Build()
	_, err = h.db.ExecContext(r.Context(), query, args...)
	return err
}
