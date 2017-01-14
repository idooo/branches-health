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
	"github.com/robfig/cron"
)

type Configuration struct {
	Repositories	*[]string
	DatabasePath	*string
	ServerPort	*int
	UpdateSchedule	*string
}

// Reads configuration file from the specified location and
// applies the default values if needed
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

	if configuration.UpdateSchedule == nil {
		defaultSchedule := "@midnight"
		configuration.UpdateSchedule = &defaultSchedule
	}

	return configuration
}

func main() {

	configPathPtr := flag.String(
		"config",
		"/etc/branches-health/config.json",
		"path to a configuration file")

	assetsPathPtr := flag.String(
		"dev-assets",
		"",
		"path to a assets folder")

	flag.Parse()

	configuration := readConfig(*configPathPtr)

	database, err := bolt.Open(*configuration.DatabasePath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Schedule job to run regularly
	c := cron.New()
	c.AddFunc(*configuration.UpdateSchedule, func() {
		core.GetBranchesInfoForRepos(*configuration.Repositories, database)
	})
	c.Start()

	// Execute our first job
	go core.GetBranchesInfoForRepos(*configuration.Repositories, database)

	// Setup Iris to serve HTTP requests
	router := core.NewRouter(database, *assetsPathPtr)
	iris.Get("/api/repositories", router.RouteGetRepositories)
	iris.Get("/api/branches", router.RouteGetBranches)
	iris.Get("/", router.RouteGetIndex)
	iris.Listen(":" + strconv.Itoa(*configuration.ServerPort))
}