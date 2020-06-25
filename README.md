# eserveless-platform

A simple serverless platform implemented for learning. 
Play with it, no kidding it can run cloud functions.

We still need some clean up & documentation but it works so Yeah!!!


## some notes along the way

https://crontab.guru/every-5-minutes
https://www.json2yaml.com/
https://docs.gofiber.io/ctx#context
https://www.sohamkamani.com/blog/2017/10/18/golang-adding-database-to-web-application/
https://docs.docker.com/engine/api/sdk/examples/#run-a-container

To build `go build ./cmd/web`

[post, put, get, delete] /projects
[any] /endpoints/:project/:function
[interal-triger] /ticks/:project/:function
[get] /logs/:project/:function
[get/post/put/delete] /users

commandline-tool managemen-app 
- projects [init list push delete]
- functions [list logs delete]
- users [list create change delete]
- init server username pwd


