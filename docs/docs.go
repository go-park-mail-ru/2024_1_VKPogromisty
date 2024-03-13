// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Petr Mitin",
            "url": "https://github.com/Petr09Mitin",
            "email": "petr09mitin@mail.ru"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/is-authorized": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "check if user is authorized",
                "operationId": "auth/is-authorized",
                "parameters": [
                    {
                        "type": "string",
                        "description": "session_id=some_session",
                        "name": "Cookie",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.JSONResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "body": {
                                            "$ref": "#/definitions/services.IsAuthorizedResponse"
                                        }
                                    }
                                }
                            ]
                        },
                        "headers": {
                            "Set-Cookie": {
                                "type": "string",
                                "description": "session_id=some_session_id; Path=/; HttpOnly; Expires=Thu, 01 Jan 1970 00:00:00 GMT;"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/auth/login/": {
            "post": {
                "description": "login user by email and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "handle user's login",
                "operationId": "auth/login",
                "parameters": [
                    {
                        "description": "Email of the user",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Password of the user",
                        "name": "password",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.JSONResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "body": {
                                            "$ref": "#/definitions/services.LoginResponse"
                                        }
                                    }
                                }
                            ]
                        },
                        "headers": {
                            "Set-Cookie": {
                                "type": "string",
                                "description": "session_id=some_session_id; Path=/; Max-Age=36000; HttpOnly;"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/errors.HTTPError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/errors.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errors.HTTPError"
                        }
                    }
                }
            }
        },
        "/auth/logout/": {
            "delete": {
                "description": "logout user that is authorized",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "handle user's logout",
                "operationId": "auth/logout",
                "parameters": [
                    {
                        "type": "string",
                        "description": "session_id=some_session",
                        "name": "Cookie",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "headers": {
                            "Set-Cookie": {
                                "type": "string",
                                "description": "session_id=some_session_id; Path=/; HttpOnly; Expires=Thu, 01 Jan 1970 00:00:00 GMT;"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/errors.HTTPError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/errors.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errors.HTTPError"
                        }
                    }
                }
            }
        },
        "/auth/signup/": {
            "post": {
                "description": "registrate user by his data",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "handle user's registration flow",
                "operationId": "auth/signup",
                "parameters": [
                    {
                        "type": "string",
                        "description": "First name of the user",
                        "name": "firstName",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Last name of the user",
                        "name": "lastName",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Email of the user",
                        "name": "email",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "minLength": 6,
                        "type": "string",
                        "description": "Password of the user",
                        "name": "password",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "minLength": 6,
                        "type": "string",
                        "description": "Repeat password of the user",
                        "name": "repeatPassword",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "format": "date",
                        "example": "2021-01-01",
                        "description": "Date of birth of the user",
                        "name": "dateOfBirth",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "Avatar of the user",
                        "name": "avatar",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.JSONResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "body": {
                                            "$ref": "#/definitions/services.User"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/errors.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errors.HTTPError"
                        }
                    }
                }
            }
        },
        "/posts/": {
            "get": {
                "description": "list posts to authorized user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "posts"
                ],
                "summary": "list all posts",
                "operationId": "posts/",
                "parameters": [
                    {
                        "type": "string",
                        "description": "session_id=some_session",
                        "name": "Cookie",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.JSONResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "body": {
                                            "$ref": "#/definitions/services.ListPostsResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/errors.HTTPError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/errors.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errors.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "errors.HTTPError": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "services.IsAuthorizedResponse": {
            "type": "object",
            "properties": {
                "isAuthorized": {
                    "type": "boolean"
                }
            }
        },
        "services.ListPostsResponse": {
            "type": "object",
            "properties": {
                "posts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/services.PostWithAuthor"
                    }
                }
            }
        },
        "services.LoginResponse": {
            "type": "object",
            "properties": {
                "sessionId": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/services.User"
                }
            }
        },
        "services.Post": {
            "type": "object",
            "properties": {
                "attachments": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "authorId": {
                    "type": "integer"
                },
                "creationDate": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2021-01-01T00:00:00Z"
                },
                "postId": {
                    "type": "integer"
                },
                "text": {
                    "type": "string"
                }
            }
        },
        "services.PostWithAuthor": {
            "type": "object",
            "properties": {
                "author": {
                    "$ref": "#/definitions/services.User"
                },
                "post": {
                    "$ref": "#/definitions/services.Post"
                }
            }
        },
        "services.User": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string",
                    "example": "default_avatar.png"
                },
                "dateOfBirth": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2021-01-01T00:00:00Z"
                },
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                },
                "registrationDate": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2021-01-01T00:00:00Z"
                },
                "userId": {
                    "type": "integer"
                }
            }
        },
        "utils.JSONResponse": {
            "type": "object",
            "properties": {
                "body": {}
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Socio API",
	Description:      "First version of Socio API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
