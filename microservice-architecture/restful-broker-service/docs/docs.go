// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/handle": {
            "post": {
                "description": "Handle request",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "api"
                ],
                "summary": "Handle request",
                "parameters": [
                    {
                        "description": "Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.CreateStatesInput": {
            "type": "object",
            "properties": {
                "color": {
                    "type": "string"
                },
                "machine": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "time": {
                    "type": "string"
                },
                "value": {
                    "type": "integer"
                }
            }
        },
        "main.Log": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "main.Request": {
            "type": "object",
            "properties": {
                "action": {
                    "type": "string"
                },
                "lastState": {
                    "$ref": "#/definitions/main.State"
                },
                "log": {
                    "$ref": "#/definitions/main.Log"
                },
                "machine": {
                    "type": "string"
                },
                "states": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.CreateStatesInput"
                    }
                }
            }
        },
        "main.State": {
            "type": "object",
            "properties": {
                "color": {
                    "type": "string"
                },
                "machine": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "time": {
                    "type": "string"
                },
                "value": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Broker Service represented as RESTful API",
	Description:      "This is a RESTful API for the Broker Service as a single point of entry for all api caLLs.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
