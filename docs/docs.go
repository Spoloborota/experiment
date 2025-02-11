// docs/docs.go
// Package docs GENERATED BY SWAG; DO NOT EDIT
package docs

import (
	"github.com/swaggo/swag"
)

var doc = `{
	"swagger": "2.0",
	"info": {
		"description": "This is a sample server Petstore server.",
		"version": "1.0",
		"title": "Swagger Example API",
		"contact": {
			"email": "support@example.com"
		},
		"license": {
			"name": "Apache 2.0",
			"url": "http://www.apache.org/licenses/LICENSE-2.0.html"
		}
	},
	"host": "localhost:8080",
	"basePath": "/",
	"paths": {
		"/register": {
			"post": {
				"summary": "Register a new user",
				"parameters": [
					{
						"in": "body",
						"name": "body",
						"description": "User info",
						"required": true,
						"schema": {
							"$ref": "#/definitions/models.User"
						}
					}
				],
				"responses": {
					"201": {
						"description": "Created",
						"schema": {
							"$ref": "#/definitions/models.User"
						}
					},
					"400": {
						"description": "Bad Request"
					},
					"500": {
						"description": "Internal Server Error"
					}
				}
			}
		},
		"/login": {
			"post": {
				"summary": "Login a user",
				"parameters": [
					{
						"in": "body",
						"name": "body",
						"description": "User credentials",
						"required": true,
						"schema": {
							"$ref": "#/definitions/models.User"
						}
					}
				],
				"responses": {
					"200": {
						"description": "OK",
						"schema": {
							"type": "object",
							"properties": {
								"token": {
									"type": "string"
								}
							}
						}
					},
					"400": {
						"description": "Bad Request"
					},
					"401": {
						"description": "Unauthorized"
					},
					"500": {
						"description": "Internal Server Error"
					}
				}
			}
		},
		"/profile": {
			"get": {
				"summary": "View user profile",
				"responses": {
					"200": {
						"description": "OK",
						"schema": {
							"$ref": "#/definitions/models.User"
						}
					},
					"401": {
						"description": "Unauthorized"
					},
					"500": {
						"description": "Internal Server Error"
					}
				}
			}
		}
	},
	"definitions": {
		"models.User": {
			"type": "object",
			"required": [
				"firstName",
				"lastName",
				"passwordHash"
			],
			"properties": {
				"id": {
					"type": "integer",
					"example": 1
				},
				"firstName": {
					"type": "string",
					"example": "John"
				},
				"lastName": {
					"type": "string",
					"example": "Doe"
				},
				"age": {
					"type": "integer",
					"example": 30
				},
				"gender": {
					"type": "string",
					"example": "Male"
				},
				"interests": {
					"type": "string",
					"example": "Reading"
				},
				"city": {
					"type": "string",
					"example": "New York"
				},
				"passwordHash": {
					"type": "string",
					"example": "hashed_password"
				}
			}
		}
	}
}`

func init() {
	swag.Register("swagger", &s{doc: doc})
}

type s struct {
	doc string
}

func (s *s) ReadDoc() string {
	return s.doc
}
