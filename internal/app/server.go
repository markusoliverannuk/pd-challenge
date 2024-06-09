package server

import (
	"challenge/internal/store"
	router "challenge/pkg"
	"encoding/json"
	"log"
	"net/http"

	"gitlab.com/0x4149/logz"
)

type Server struct {
	Addr         string
	GithubAPI    *GitHubAPP
	PipedriveAPI string
	Store        store.Store
	Logger       *log.Logger
	Server       *http.Server
	Router       *router.Router
}

func newServer(srv *http.Server, g *GitHubAPP, store store.Store, router *router.Router) *Server {
	s := &Server{
		GithubAPI: g,
		Store:     store,
		Logger:    log.Default(),
		Server:    srv,
		Router:    router,
	}

	s.configureRouter()

	return s
}

func (s *Server) configureRouter() {
	s.Router.GET("/user/{id}", s.getUserGists())
}

func (s *Server) getUserGists() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gists, err := s.Store.Gists().GetUsers(r.PathValue("id"))
		if err != nil {
			s.error(w, http.StatusBadRequest, err)
			return
		}

		for _, gist := range gists {
			err := CreateDealForGist(gist)
			if err != nil {
				logz.Error("Error creating deal for gist: ", err)
			} else {
				logz.Info("Successfully created deal for gist: ", gist.Description)
			}
		}

		s.respond(w, http.StatusOK, gists)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

func (s *Server) error(w http.ResponseWriter, code int, err error) {
	s.respond(w, code, map[string]string{"error": err.Error()})
}

func (s *Server) respond(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
