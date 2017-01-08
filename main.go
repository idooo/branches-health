package main

import (
	"log"
	"branches-health/core"
	"github.com/boltdb/bolt"
	"os"
	"encoding/json"
	"github.com/kataras/iris"
	"strconv"
	"flag"
)

type Configuration struct {
	Repositories	*[]string
	DatabasePath	*string
	ServerPort	*int
}

func readConfig (filename string) Configuration {
	file, _ := os.Open(filename)
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Printf("Can't read configuration file %s : %s", filename, err)
	}

	if configuration.DatabasePath == nil {
		currPath, _ := os.Getwd()
		defaultPath := currPath + "/branches-health.db"
		configuration.DatabasePath = &defaultPath
	}

	if configuration.Repositories == nil {
		defaultRepos := make([]string, 0)
		configuration.Repositories = &defaultRepos
	}

	if configuration.ServerPort == nil {
		defaultPort := 8080
		configuration.ServerPort = &defaultPort
	}

	return configuration
}

func getInfoAboutBranches (repositories []string, database *bolt.DB) {

	for _, repoName := range repositories {
		branches := core.GetInfoFromGit(repoName)

		for _, branch := range branches {
			log.Printf("Get information about: %s/%s", repoName, branch.Name)
			branch.Save(database)
		}
	}
}

func main() {

	configPathPtr := flag.String(
		"config",
		"/etc/branches-health/config.json",
		"path to a configuration file")

	flag.Parse()

	configuration := readConfig(*configPathPtr)

	database, err := bolt.Open(*configuration.DatabasePath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	go getInfoAboutBranches(*configuration.Repositories, database)

	router := core.NewRouter(database)

	iris.Get("/api/repositories", router.RouteGetRepositories)
	iris.Get("/api/branches", router.RouteGetBranches)
	iris.Listen(":" + strconv.Itoa(*configuration.ServerPort))
}