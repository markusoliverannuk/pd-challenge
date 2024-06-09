package server

import (
	"challenge/internal/models"
	"challenge/internal/store"
	"context"
	"time"

	"github.com/google/go-github/github"
	"gitlab.com/0x4149/logz"
	"golang.org/x/oauth2"
)

type GitHubAPP struct {
	AccessToken  string
	TrackedUsers []string
	Store        store.Store
}

type GistWorkerData struct {
	Username string
	Gists    []*github.Gist
}

func NewGithubAPP(access_token string, s store.Store) *GitHubAPP {
	return &GitHubAPP{
		AccessToken: access_token,
		Store:       s,
	}
}

func (g *GitHubAPP) AddUser(username string) {
	g.TrackedUsers = append(g.TrackedUsers, username)
}

func (g *GitHubAPP) Start() {
	logz.Info("Starting Github Scraper...")
	//Channel in order to get all the information from goroutines
	gistChannel := make(chan GistWorkerData)

	for _, user := range g.TrackedUsers {
		go g.gistWorker(user, gistChannel)
	}

	//strat tracking them every 5 seconds
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for gists := range gistChannel {
		for _, gist := range gists.Gists {
			Gist := models.Gist{
				Username:    gists.Username,
				Description: *gist.Description,
			}
			_, err := g.Store.Gists().CreateGist(Gist)
			if err != nil {
				logz.Error("Error: ", err)
			}

			for _, file := range gist.Files {
				File := models.File{
					Username: gists.Username,
					Path:     *file.RawURL,
				}
				_, err := g.Store.Gists().CreateFile(File)
				if err != nil {
					logz.Error("Fail error:", err)
				}
			}
		}
	}
}

func (g *GitHubAPP) gistWorker(username string, gistChannel chan<- GistWorkerData) {
	for {
		gists := g.getUserGists(username)
		data := GistWorkerData{
			Username: username,
			Gists:    gists,
		}
		gistChannel <- data
		time.Sleep(5 * time.Second)
	}
}

func (g *GitHubAPP) getUserGists(username string) []*github.Gist {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: g.AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	gists, _, err := client.Gists.List(ctx, username, nil)
	if err != nil {
		logz.Error("Failed to fetch gists: %v", err)
	}
	
	return gists
}
