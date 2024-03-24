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
                                    "$ref": "#/definitions/json.JSONResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "body": {
                                            "$ref": "#/definitions/auth.IsAuthorizedResponse"
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
                                    "$ref": "#/definitions/json.JSONResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "body": {
                                            "$ref": "#/definitions/auth.LoginResponse"
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
                                    "$ref": "#/definitions/json.JSONResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "body": {
                                            "$ref": "#/definitions/domain.User"
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
                                    "$ref": "#/definitions/json.JSONResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "body": {
                                            "$ref": "#/definitions/posts.ListPostsResponse"
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
        },
        "/profile/{userID}": {
            "get": {
                "description": "get user profile with subscriptions info",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "profile"
                ],
                "summary": "get user profile with subscriptions info",
                "operationId": "profile/get",
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
                                    "$ref": "#/definitions/json.JSONResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "body": {
                                            "$ref": "#/definitions/profile.UserWithSubsInfo"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/errors.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
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
        "/subscriptions/": {
            "post": {
                "description": "subscribe to user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscriptions"
                ],
                "summary": "handle user's subscription flow",
                "operationId": "subscriptions/subscribe",
                "parameters": [
                    {
                        "description": "Subscribed to ID",
                        "name": "subscribedTo",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "type": "string",
                        "description": "session_id=some_session",
                        "name": "Cookie",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/json.JSONResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "body": {
                                            "$ref": "#/definitions/domain.Subscription"
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
            },
            "delete": {
                "description": "unsubscribe from user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscriptions"
                ],
                "summary": "handle user's unsubscription flow",
                "operationId": "subscriptions/unsubscribe",
                "parameters": [
                    {
                        "description": "User to unsubscribe from",
                        "name": "subscribedTo",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "type": "string",
                        "description": "session_id=some_session",
                        "name": "Cookie",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
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
                    "404": {
                        "description": "Not Found",
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
        "/subscriptions/friends/": {
            "get": {
                "description": "get user's friends",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscriptions"
                ],
                "summary": "get user's friends",
                "operationId": "subscriptions/friends",
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
                                    "$ref": "#/definitions/json.JSONResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "body": {
                                            "$ref": "#/definitions/subscriptions.GetFriendsResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "404": {
                        "description": "Not Found",
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
        "/subscriptions/subscribers/": {
            "get": {
                "description": "get user's subscribers",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscriptions"
                ],
                "summary": "get user's subscribers",
                "operationId": "subscriptions/subscribers",
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
                                    "$ref": "#/definitions/json.JSONResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "body": {
                                            "$ref": "#/definitions/subscriptions.GetSubscribersResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "404": {
                        "description": "Not Found",
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
        "/subscriptions/subscriptions/": {
            "get": {
                "description": "get user's subscriptions",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscriptions"
                ],
                "summary": "get user's subscriptions",
                "operationId": "subscriptions/subscriptions",
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
                                    "$ref": "#/definitions/json.JSONResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "body": {
                                            "$ref": "#/definitions/subscriptions.GetSubscriptionsResponse"
                                        }
                                    }
                                }
                            ]
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
        "auth.IsAuthorizedResponse": {
            "type": "object",
            "properties": {
                "isAuthorized": {
                    "type": "boolean"
                }
            }
        },
        "auth.LoginResponse": {
            "type": "object",
            "properties": {
                "sessionId": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/domain.User"
                }
            }
        },
        "domain.Post": {
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
                "content": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2021-01-01T00:00:00Z"
                },
                "postId": {
                    "type": "integer"
                },
                "updatedAt": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2021-01-01T00:00:00Z"
                }
            }
        },
        "domain.Subscription": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2021-01-01T00:00:00Z"
                },
                "subscribedTo": {
                    "type": "integer"
                },
                "subscriber": {
                    "type": "integer"
                },
                "subscriptionId": {
                    "type": "integer"
                },
                "updatedAt": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2021-01-01T00:00:00Z"
                }
            }
        },
        "domain.User": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string",
                    "example": "default_avatar.png"
                },
                "createdAt": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2021-01-01T00:00:00Z"
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
                "updatedAt": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2021-01-01T00:00:00Z"
                },
                "userId": {
                    "type": "integer"
                }
            }
        },
        "errors.HTTPError": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "json.JSONResponse": {
            "type": "object",
            "properties": {
                "body": {}
            }
        },
        "posts.ListPostsResponse": {
            "type": "object",
            "properties": {
                "posts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/posts.PostWithAuthor"
                    }
                }
            }
        },
        "posts.PostWithAuthor": {
            "type": "object",
            "properties": {
                "author": {
                    "$ref": "#/definitions/domain.User"
                },
                "post": {
                    "$ref": "#/definitions/domain.Post"
                }
            }
        },
        "profile.UserWithSubsInfo": {
            "type": "object",
            "properties": {
                "is_subscribed_to": {
                    "type": "boolean"
                },
                "is_subscriber": {
                    "type": "boolean"
                },
                "user": {
                    "$ref": "#/definitions/domain.User"
                }
            }
        },
        "subscriptions.GetFriendsResponse": {
            "type": "object",
            "properties": {
                "friends": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.User"
                    }
                }
            }
        },
        "subscriptions.GetSubscribersResponse": {
            "type": "object",
            "properties": {
                "subscribers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.User"
                    }
                }
            }
        },
        "subscriptions.GetSubscriptionsResponse": {
            "type": "object",
            "properties": {
                "subscriptions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.User"
                    }
                }
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
