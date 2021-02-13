# Minitwit API

> This API is comsumed by the simulator for the course DevOps, Software Evolution and Software Maintenance @ ITU Spring 2021

## How to run

```
go run api.go
```

## Docs

Swagger is used for API documentation. Documentation can be found on endpoint `/swagger`.

### Update docs

When annotations have been added or updated, run the command ```swag init -g api.go```

#### Requirements
Install the following packages ___OUTSIDE___ this repository (devops-21) - otherwise unnecessary packages will be added to go mod file. Also make sure your `$GOPATH/bin` is added to your $PATH to be able to run `swag`, else you will be prompted with something like `zsh: command not found: swag`.
```
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/http-swagger
go get -u github.com/alecthomas/template
```
![alt text](https://giphy.com/embed/L12g7V0J62bf2)
<iframe src="https://giphy.com/embed/L12g7V0J62bf2" width="480" height="360" frameBorder="0" class="giphy-embed" allowFullScreen></iframe><p><a href="https://giphy.com/gifs/dancing-wtf-swag-L12g7V0J62bf2">via GIPHY</a></p>
