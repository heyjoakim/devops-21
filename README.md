# DevOps, Software Evolution and Software Maintenance

[![Build & Deploy](https://github.com/heyjoakim/devops-21/actions/workflows/deploy.yml/badge.svg)](https://github.com/heyjoakim/devops-21/actions/workflows/deploy.yml)
[![Tests](https://github.com/heyjoakim/devops-21/actions/workflows/build_test.yml/badge.svg)](https://github.com/heyjoakim/devops-21/actions/workflows/build_test.yml)

> This project revolves around a forum application called minitwit. The functionalities includes signing up, logging in, posting messages, following other users. The forum has a public timeline where all messages are displayed. Furthermore, if a user is signed in, a personal timeline exists that displays a users own messages aswell as messages of followed users.

## Setting up env variable

Run the following command `cp .env.example .env` to create a .env file

Populate the .env file with your db connection string

## Running the application

To run the application locally `go run minitwit.go`

Server running on port http://localhost:8000

## Test the application

To execute unit tests `go test -v`

## Remote access

The latest release is running in the cloud with Digital Ocean at http://206.189.14.172:8000

## Dependencies

### Libraries

> Table only lists direct dependencies. Verbose dependency graph can be found [here](assets/dep_app_simple.png).

| **Dependency**              | **Version**                        | **Description**                                    |
| --------------------------- | ---------------------------------- | -------------------------------------------------- |
| github.com/gorilla/mux      | 1.8.0                              | Framework for HTTP request handling.               |
| github.com/gorilla/sessions | 1.2.1                              | Provides access to read and write session cookies. |
| github.com/mattn/go-sqlite3 | 1.14.6                             | Database driver for SQLite3.                       |
| gorm.io/gorm                | 1.20.12                            | ORM for Go.                                        |
| golang.org/x/crypto/bcrypt  | v0.0.0-20201221181555-eec23a3978ad | Used to hash passwords and verify password hashes. |

### Cloud dependencies

> These services are responsible for cloud hosting.

| **Service**      | **Provider**    | **Description**                |
| ---------------- | --------------- | ------------------------------ |
| Droplet     | Digital Ocean | Hosting of web application     |
| Docker Container | Docker          | Containerizing of applications |

## API

### Docs

Swagger is used for API documentation. Documentation can be found on endpoint `/api/swagger`.

#### Update docs

When annotations have been added or updated, run the command `swag init -g minitwit.go`

#### Monitoring tools

Prometheus server `http://142.93.103.26:9090`

Minitwit metrics from prometheus + custom `206.189.14.172:8000//metrics`

Grafana consuming prometheus `http://164.90.165.111:3000/d/JJQvP88Mz/prometheus-2-0-stats`, requires credentials

##### Requirements

Install the following packages **_OUTSIDE_** this repository (devops-21) - otherwise unnecessary packages will be added to go mod file.

```
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/http-swagger
go get -u github.com/alecthomas/template
```

Also make sure your `$GOPATH/bin` is added to your $PATH to be able to run `swag`, else you will be prompted with something like `zsh: command not found: swag`.

#### License

Licensed by the [MIT License ](LICENSE)


## Authors

_Joakim Hey Hinnerskov (jhhi), Ask Harup Sejsbo (asse), Kasper Olsen (kols), Petya Buchkova (pebu) and Thomas Tyge Andersen (thta)._
