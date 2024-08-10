// Package swaggerHandlers Code generated by swaggo/swag. DO NOT EDIT
package swaggerHandlers

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
        "/api/v1/user/login": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Allows a user to refresh the access token",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Refresh the access token",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/userHandlers.AuthResult"
                        }
                    }
                }
            }
        },
        "/api/v1/user/me": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Allows a user to get their own information",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get the user's information",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/userHandlers.MeResult"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "userHandlers.AuthResult": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "expiration": {
                    "type": "integer"
                },
                "full_name": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "userHandlers.LoginData": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "userHandlers.LoginResult": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "expiration": {
                    "type": "integer"
                },
                "full_name": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "userHandlers.MeResult": {
            "type": "object",
            "properties": {
                "full_name": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
