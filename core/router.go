package core

import (
	"github.com/boltdb/bolt"
	"github.com/kataras/iris"
	"io/ioutil"
	"fmt"
)

var database *bolt.DB
var pathToAssets string;

type API struct {
	*iris.Context
}

func NewRouter(db *bolt.DB, assetsPath string) API {
	database = db
	pathToAssets = assetsPath

	if len(pathToAssets) > 0 {
		fmt.Print("Using development version of assets")
	} else {
		fmt.Print("Using compiled version of assets")
	}
	return API{}
}


func (api API) RouteGetRepositories(ctx *iris.Context)  {
	branches, err := GetBranches(database)
	if err != nil {
		ctx.JSON(iris.StatusNotFound, iris.Map{"status": "error"})
	} else {
		repositories := make(map[string][]Branch)

		for _, branch := range branches {
			if _, present := repositories[branch.Repository]; !present {
				repositories[branch.Repository] = make([]Branch, 0)
			}
			repositories[branch.Repository] = append(repositories[branch.Repository], branch)
		}

		ctx.JSON(iris.StatusOK, repositories)
	}
}

func (api API) RouteGetBranches(ctx *iris.Context)  {
	branches, err := GetBranches(database)
	if err != nil {
		ctx.JSON(iris.StatusNotFound, iris.Map{"status": "error"})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"branches": branches})
	}
}

func (api API) RouteGetIndex(ctx *iris.Context)  {
	if len(pathToAssets) > 0 {
		data, err := ioutil.ReadFile(pathToAssets + "/index.html")
		if err != nil {
			fmt.Printf("Can't open index.html: %s", err)
		}
		ctx.HTML(iris.StatusOK, string(data))
	} else {
		ctx.HTML(iris.StatusOK, IndexTemplate)
	}
}
