# eserveless-platform

A simple serverless platform implemented for learning. 
Play with it, no kidding it can run cloud functions.
We still need some clean up but it works so Yeah!!!


# some notes along the way

https://crontab.guru/every-5-minutes
https://www.json2yaml.com/
https://docs.gofiber.io/ctx#context
https://www.sohamkamani.com/blog/2017/10/18/golang-adding-database-to-web-application/
https://docs.docker.com/engine/api/sdk/examples/#run-a-container


 go build ./cmd/web

[go, rust]
{
	id: uuid,
	repo: git,
	type: go|node,
	functions: [
		{
			name: 'echo',
			type: http|event|cron,
			meta: {interval: 4}, 
		}
	]

}


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


fetch gitrepo
parse .eserveless
create/update db
rebuild all functions


projects:
	id
	name:
	repo_url:
	last_commit:

functions:
	id
	name
	runtime
	type
	interval
	meta
	project_id
	

