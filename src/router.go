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


func (api API) RouteGetBranches(ctx *iris.Context)  {
	branches, err := GetBranches(database)
	if err != nil {
		ctx.JSON(iris.StatusNotFound, iris.Map{"status": "error"})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{"branches": branches})
	}
}

