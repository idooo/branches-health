package main

import (
	"log"
	"branches/src"
	"github.com/boltdb/bolt"
	"os"
	"encoding/json"
	"github.com/kataras/iris"
	"strconv"
)

type Configuration struct {
	Repositories	[]string
	DatabasePath	string
	ServerPort	int
}

func readConfig (filename string) Configuration {
	file, _ := os.Open(filename)
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Fatalf("Can't read configuration file %s : %s", filename, err)
	}
	return configuration
}

func getInfoAboutBranches (repositories []string, database *bolt.DB) {

	for _, repoName := range repositories {
		branches := core.GetInfoFromGit(repoName)

		for _, branch := range branches {
			log.Printf("branch: %s", branch.Name)
			branch.Save(database)
		}
	}
}

func main() {

	configuration := readConfig("config/default.json")

	database, err := bolt.Open(configuration.DatabasePath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	go getInfoAboutBranches(configuration.Repositories, database)

	router := core.NewRouter(database)

	iris.Get("/api/repositories", router.RouteGetBranches)
	iris.Listen(":" + strconv.Itoa(configuration.ServerPort))
}