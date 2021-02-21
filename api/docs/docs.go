// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/fllws/{username}": {
            "post": {
                "description": "Eiter follows a user, unfollows a user or returns a list of users's followers",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Follow, unfollow or get followers",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Number of results returned",
                        "name": "no",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Something about latest",
                        "name": "latest",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object"
                        }
                    },
                    "204": {
                        "description": "No Content",
                        "schema": {
                            "type": "object"
                        }
                    },
                    "401": {
                        "description": "unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/latest": {
            "get": {
                "description": "Get the latest x",
                "produces": [
                    "application/json"
                ],
                "summary": "Get the latest x",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/msgs": {
            "get": {
                "description": "Gets the latest messages in descending order.",
                "produces": [
                    "application/json"
                ],
                "summary": "Gets the latest messages",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Number of results returned",
                        "name": "no",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object"
                        }
                    },
                    "401": {
                        "description": "unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/msgs/{username}": {
            "post": {
                "description": "Gets the latest messages per user",
                "produces": [
                    "application/json"
                ],
                "summary": "Gets the latest messages per user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Number of results returned",
                        "name": "no",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Something about latest",
                        "name": "latest",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object"
                        }
                    },
                    "401": {
                        "description": "unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "Registers a user, provided that the given info passes all checks.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Registers a user",
                "responses": {
                    "203": {
                        "description": ""
                    },
                    "400": {
                        "description": "unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:8001",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "Minitwit API",
	Description: "This API is comsumed by the simulator for the course DevOps, Software Evolution and Software Maintenance @ ITU Spring 2021",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
