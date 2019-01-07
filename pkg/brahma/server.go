package brahma

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	log "gopkg.in/src-d/go-log.v1"
)

type Repository struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type Server struct {
	repos   []Repository
	current int
	storage string
}

func (s *Server) router() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/repository", s.repository).Methods(http.MethodGet)
	router.HandleFunc("/upload/{id}", s.upload).Methods(http.MethodPost)

	return router
}

func (s *Server) repository(w http.ResponseWriter, r *http.Request) {
	log.Infof("new connection from %s: %s %s",
		r.RemoteAddr, r.Method, r.RequestURI)

	if s.current >= len(s.repos) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.NewEncoder(w).Encode(s.repos[s.current])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.current++
}

func (s *Server) upload(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	name := fmt.Sprintf("%s.siva", id)
	siva := filepath.Join(s.storage, name)

	l := log.With(log.Fields{
		"id":   id,
		"siva": siva,
	})
	l.Infof("uploading file")

	err := os.MkdirAll(s.storage, 0770)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f, err := os.Create(siva)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	n, err := io.Copy(f, r.Body)
	if err != nil {
		f.Close()
		os.Remove(siva)
		l.Errorf(err, "could not save file")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = f.Close()
	if err != nil {
		os.Remove(siva)
		l.Errorf(err, "could not save file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	l.With(log.Fields{"size": n}).Infof("finished uploading file")
	w.WriteHeader(http.StatusOK)
}

func StartServer() error {
	server := &Server{
		repos: []Repository{
			{"0", "https://github.com/jfontan/cangallo"},
			{"1", "https://github.com/jfontan/borges"},
			{"2", "https://github.com/src-d/borges"},
		},
		storage: "server",
	}

	return http.ListenAndServe(":8765", server.router())
}
