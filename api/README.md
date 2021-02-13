# Minitwit API

> This API is comsumed by the simulator for the course DevOps, Software Evolution and Software Maintenance @ ITU Spring 2021

## How to run

```
go run api.go
```

Will serve on `http://localhost:8001`.

## Docs

Swagger is used for API documentation. Documentation can be found on endpoint `/swagger`.

### Update docs

When annotations have been added or updated, run the command ```swag init -g api.go```

#### Requirements
Install the following packages ___OUTSIDE___ this repository (devops-21) - otherwise unnecessary packages will be added to go mod file.
```
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/http-swagger
go get -u github.com/alecthomas/template
```
Also make sure your `$GOPATH/bin` is added to your $PATH to be able to run `swag`, else you will be prompted with something like `zsh: command not found: swag`.

<img src="https://media.giphy.com/media/L12g7V0J62bf2/giphy.gif" width="400" />
