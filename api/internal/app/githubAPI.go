package server

import (
	"challenge/internal/models"
	"challenge/internal/store"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/go-github/github"
	"github.com/google/uuid"
	"gitlab.com/0x4149/logz"
	"golang.org/x/oauth2"
)

type GitHubAPP struct {
	AccessToken  string
	TrackedUsers []string
	Store        store.Store
	GistChannel  chan GistWorkerData
}

type GistWorkerData struct {
	Username string
	Gists    []*github.Gist
}

func NewGithubAPP(access_token string, s store.Store) *GitHubAPP {
	return &GitHubAPP{
		AccessToken: access_token,
		Store:       s,
		GistChannel: make(chan GistWorkerData),
	}
}

func (g *GitHubAPP) AddUser(username string) {
	g.TrackedUsers = append(g.TrackedUsers, username)
	go g.gistWorker(username)
}

func (g *GitHubAPP) IsUserTracked(username string) bool {
	for _, user := range g.TrackedUsers {
		if user == username {
			return true
		}
	}
	return false
}

func (g *GitHubAPP) Start() {
	logz.Info("Starting Github Scraper...")
	//strat tracking them every 5 seconds
	// ticker := time.NewTicker(100 * time.Second)
	// defer ticker.Stop()

	for gists := range g.GistChannel {
		for _, gist := range gists.Gists {
			Gist := models.Gist{
				Username:    gists.Username,
				Description: *gist.Description,
			}
			user_gist, err := g.Store.Gists().CreateGist(Gist)
			if err != nil {
				//logz.Error("Errora: ", err)
			}

			if user_gist != nil {
				for _, file := range gist.Files {
					logz.Info(file)
					File := models.File{
						Id:       user_gist.Id,
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
}

func (g *GitHubAPP) gistWorker(username string) {
	for {
		gists := g.getUserGists(username)
		data := GistWorkerData{
			Username: username,
			Gists:    gists,
		}
		g.GistChannel <- data
		time.Sleep(10800 * time.Second)
	}
}

func (g *GitHubAPP) getUserGists(username string) []*github.Gist {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: g.AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	var allGists []*github.Gist
	opts := &github.GistListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		gists, resp, err := client.Gists.List(ctx, username, opts)
		if err != nil {
			log.Fatalf("Failed to fetch gists: %v", err)
		}

		//Process each gist
		for _, gist := range gists {
			// Check if description is nil
			if gist.Description == nil || *gist.Description == "" {
				description := "untitled - " + uuid.New().String()
				gist.Description = &description
			}
			fmt.Println("Gist DESCRIPTION:", gist.Description)

		}

		allGists = append(allGists, gists...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allGists
}