# Branches Health [![Build Status](https://travis-ci.org/idooo/branches-health.svg?branch=master)](https://travis-ci.org/idooo/branches-health)

Simple service that retrieves meta data from git repositories
about remote branches and provides an info which of them were merged and which not updated for a white

## Build

`glide` was used for dependency management here so this should work: 

```
glide install
go build
```

There is also a simple HTML page served from `/` that must be "compiled" to be a part of app's code. So use `./compile-template.sh` script that will update `core/template.go`. 

During development you can specify path to serve your assets folder directly:

```
./branches-health -dev-assets=/path/to/branches-health/assets
```

## Run

Pass a path to configuration file

```
./branches-health -config=./config/default.json
```

Properties in a configuration file are easy to understand. Check `config/example.json` for example.
Read [robfig/cron docs](https://godoc.org/github.com/robfig/cron) to know more about `UpdateSchedule` format

```json
{
  "Repositories": [
    "https://github.com/idooo/test"
  ],
  "DatabasePath": "/tmp/branches-health.db",
  "ServerPort": 8080,
  "UpdateSchedule": "@midnight" 
}
```

## Endpoints

#### GET: /api/repositories

Returns map `repository -> []branch`, example response:

```
{
    "https://github.com/idooo/test-repo": [
        {
            "Repository": "https://github.com/idooo/test-repo",
            "Name": "origin/0.6.0",
            "FullPath": "https://github.com/idooo/test-repo/origin/0.6.0",
            "IsMerged": true,
            "IsOutdated": true,
            "Author": "idooo",
            "LastUpdated": "2015-08-10T19:03:59-06:00"
        },
        ...
    ]
    ...
}
```

#### GET: /api/branches

Returns a list of branches `[]branch`, example response:

```
{
    "branches": [
        {
            "Repository": "https://github.com/idooo/test-repo",
            "Name": "origin/0.6.0",
            "FullPath": "https://github.com/idooo/test-repo/origin/0.6.0",
            "IsMerged": true,
            "IsOutdated": true,
            "Author": "idooo",
            "LastUpdated": "2015-08-10T19:03:59-06:00"
        },
        ...
    ]
}
```


# License

##### The MIT License (MIT)

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.


