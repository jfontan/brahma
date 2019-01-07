package brahma

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "gopkg.in/src-d/go-log.v1"
)

type Repository struct {
	ID  string `json:"name"`
	URL string `json:"url"`
}

type Server struct {
	repos   []Repository
	current int
}

func (s *Server) router() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/repository", s.repository)

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
	}

	s.current++
}

func StartServer() error {
	server := &Server{
		repos: []Repository{
			{"0", "https://github.com/jfontan/cangallo"},
			{"1", "https://github.com/jfontan/borges"},
			{"2", "https://github.com/src-d/borges"},
		},
	}

	return http.ListenAndServe(":8765", server.router())
}
