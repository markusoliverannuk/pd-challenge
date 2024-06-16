package main

import (
	"bufio"
	server "challenge/internal/app"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"gitlab.com/0x4149/logz"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "./configs/config.json", "path to config")
}

func main() {
	pipedriveAPIKey := os.Getenv("PIPEDRIVE_API_KEY")
	githubAT := os.Getenv("GITHUB_AT")


	if pipedriveAPIKey == "" {
		fmt.Print("Enter your pipedrive api key: ")
		pipedriveAPIKey = readInput()
		os.Setenv("PIPEDRIVE_API_KEY", pipedriveAPIKey)
	}

	if githubAT == "" {
		fmt.Print("Enter your Github Access Token: ")
		githubAT = readInput()
		os.Setenv("GITHUB_AT", githubAT)
	}
	flag.Parse()
	logz.Verbos = true
	config := server.NewConfig()
	err := config.ReadConfig(configPath)
	if err != nil {
		logz.Error("Error reading config file: %s\n", err)
	}

	// mux := http.NewServeMux()
	// mux.HandleFunc("/", testHandler())
	// http.ListenAndServe(":8080", mux)
	log.Fatal(server.Start(config))
}

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
