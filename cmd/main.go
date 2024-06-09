package main

import (
	server "challenge/internal/app"
	"flag"
	"log"

	"gitlab.com/0x4149/logz"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "./configs/config.json", "path to config")
}

func main() {
	logz.Verbos = true
	flag.Parse()

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
