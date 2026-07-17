// Package docs provides Swagger documentation metadata for DevOps Command Center.
// Generate full docs with: swag init -g cmd/server/main.go -o docs
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
        "/api/v1/auth/login": {
            "post": {
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["auth"],
                "summary": "Login",
                "parameters": [{
                    "in": "body",
                    "name": "body",
                    "required": true,
                    "schema": {"type": "object"}
                }],
                "responses": {"200": {"description": "OK"}}
            }
        },
        "/api/v1/dashboard/stats": {
            "get": {
                "security": [{"BearerAuth": []}],
                "tags": ["dashboard"],
                "summary": "Dashboard statistics",
                "responses": {"200": {"description": "OK"}}
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8095",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "DevOps Command Center API",
	Description:      "Enterprise DevOps Dashboard REST API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
