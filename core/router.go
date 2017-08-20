package core

import (
	"fmt"
	"io/ioutil"

	"github.com/boltdb/bolt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

var database *bolt.DB
var pathToAssets string

type API struct {
	context.Context
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

func (api API) RouteGetRepositories(ctx context.Context) {
	branches, err := GetBranches(database)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(map[string]interface{}{"status": "error"})
	} else {
		repositories := make(map[string][]Branch)

		for _, branch := range branches {
			if _, present := repositories[branch.Repository]; !present {
				repositories[branch.Repository] = make([]Branch, 0)
			}
			repositories[branch.Repository] = append(repositories[branch.Repository], branch)
		}

		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(repositories)
	}
}

func (api API) RouteGetBranches(ctx context.Context) {
	branches, err := GetBranches(database)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(map[string]interface{}{"status": "error"})
	} else {
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(map[string]interface{}{"branches": branches})
	}
}

func (api API) RouteGetIndex(ctx context.Context) {
	if len(pathToAssets) > 0 {
		data, err := ioutil.ReadFile(pathToAssets + "/index.html")
		if err != nil {
			fmt.Printf("Can't open index.html: %s", err)
		}
		ctx.StatusCode(iris.StatusOK)
		ctx.HTML(string(data))
	} else {
		ctx.StatusCode(iris.StatusOK)
		ctx.HTML(IndexTemplate)
	}
}
