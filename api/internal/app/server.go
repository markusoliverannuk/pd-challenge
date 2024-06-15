package server

import (
	"challenge/internal/models"
	"challenge/internal/store"
	router "challenge/pkg"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

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
			logz.Info("Adding user")
			s.GithubAPI.AddUser(userId)
			// we load first gists to db
			
			time.Sleep(15 * time.Second)
			
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
			go func() {
				originalID, err := extractOriginalID(gist.Files[0])
				if err != nil {
					return
				}
				logz.Info("Original ID", originalID)
				error := CreatePipedriveDeal(gist.Username, gist.Description, gist.Id, originalID)

				if error != nil {
					// w log the error and continue with to next gist
					fmt.Printf("Error creating Pipedrive deal for gist %s: %v\n", gist.Id, err)
				}
			}()
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
// what this function does is it extracts the Gist ID from the first file (there might be more) in the gist. there has to be at least one so its safe
// to extract from ...File[0]. from the URL we can get the ID of the gist.
func extractOriginalID(url string) (string, error) {
	startIndex := strings.Index(url, "/gist.githubusercontent.com/") + len("/gist.githubusercontent.com/")
	if startIndex == -1 {
		return "", fmt.Errorf("start marker not found in URL")
	}

	endIndex := strings.Index(url[startIndex:], "/raw/")
	if endIndex == -1 {
		return "", fmt.Errorf("end marker not found in URL")
	}

	return url[startIndex : startIndex+endIndex], nil
}

func calculateSleepDuration(numGists int) time.Duration {
	
	// bit of messing around, this will probably be gone soon
	const baseSleepTime = 100 * time.Millisecond 
	maxSleepTime := 20 * time.Second             

	sleepTime := time.Duration(numGists) * baseSleepTime
	if sleepTime > maxSleepTime {
		return maxSleepTime
	}
	return sleepTime
}


