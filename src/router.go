package core

import (
	"github.com/boltdb/bolt"
	"github.com/kataras/iris"
)

var database *bolt.DB

type API struct {
	*iris.Context
}

func NewRouter(db *bolt.DB) API {
	database = db
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


