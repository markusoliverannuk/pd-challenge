package server

import (
	"challenge/internal/models"
	"challenge/internal/store"
	router "challenge/pkg"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
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
	s.Router.GET("/trackedusers", s.getTrackedUsers())
}

func (s *Server) getUserGists() http.HandlerFunc {
	type Response struct {
		OldGists []models.Gist `json:"old_gists"`
		NewGists []models.Gist `json:"new_gists"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.PathValue("id")

		// we check if user is already tracked
		if !s.GithubAPI.IsUserTracked(userId) {
			fmt.Println("Adding user")
			s.GithubAPI.AddUser(userId)
			// we load first gists to db
			time.Sleep(1000 * time.Millisecond)
		}

		old_gists, err := s.Store.Gists().GetUsersOld(userId)
		if err != nil {
			s.error(w, http.StatusBadRequest, err)
			return
		}
		new_gists, err := s.Store.Gists().GetUsersNew(userId)
		if err != nil {
			s.error(w, http.StatusBadRequest, err)
			return
		}

		// we call CreatePipedriveDeal for each new gist
		for _, gist := range new_gists {
			err := CreatePipedriveDeal(gist.Username, gist.Description, gist.Id)
			if err != nil {
				// w log the error and continue with to next gist
				fmt.Printf("Error creating Pipedrive deal for gist %s: %v\n", gist.Id, err)
			}
		}

		s.respond(w, http.StatusOK, &Response{OldGists: old_gists, NewGists: new_gists})
	}
}

func (s *Server) getTrackedUsers() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // quering the database to fetch unique values from usernames
        uniqueNames, err := s.Store.Gists().GetUniqueNames()
        if err != nil {
            s.error(w, http.StatusInternalServerError, err)
            return
        }

        
        s.respond(w, http.StatusOK, uniqueNames)
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
