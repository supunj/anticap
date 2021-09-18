// Package api GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package api

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "soberkoder@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/register": {
            "post": {
                "description": "Register handles the register node functionality",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "register"
                ],
                "summary": "Register handles the register node functionality",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.Node"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "types.Key": {
            "type": "object",
            "properties": {
                "privatekey": {
                    "type": "string"
                },
                "publickey": {
                    "type": "string"
                }
            }
        },
        "types.Location": {
            "type": "object",
            "properties": {
                "lat": {
                    "type": "number"
                },
                "lon": {
                    "type": "number"
                }
            }
        },
        "types.Node": {
            "type": "object",
            "properties": {
                "active": {
                    "type": "boolean"
                },
                "availability": {
                    "type": "boolean"
                },
                "bday": {
                    "type": "string"
                },
                "gender": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "key": {
                    "$ref": "#/definitions/types.Key"
                },
                "location": {
                    "$ref": "#/definitions/types.Location"
                },
                "mobile": {
                    "type": "string"
                },
                "rating": {
                    "$ref": "#/definitions/types.Rating"
                },
                "subscription": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/types.Subscription"
                    }
                },
                "vcode": {
                    "type": "string"
                }
            }
        },
        "types.Rating": {
            "type": "object",
            "properties": {
                "courtesy": {
                    "description": "Both",
                    "type": "number"
                },
                "noofCancelledRequests": {
                    "description": "Consumer",
                    "type": "integer"
                },
                "noofServicesDelivered": {
                    "description": "Provider",
                    "type": "integer"
                },
                "offersAccepted": {
                    "description": "Consumer",
                    "type": "integer"
                },
                "offersMade": {
                    "description": "Provider",
                    "type": "integer"
                },
                "offersRejected": {
                    "description": "Consumer",
                    "type": "integer"
                },
                "price": {
                    "description": "Provider",
                    "type": "number"
                },
                "promptPayment": {
                    "description": "Consumer",
                    "type": "number"
                },
                "quality": {
                    "description": "Provider",
                    "type": "number"
                },
                "recommendNo": {
                    "description": "Provider",
                    "type": "integer"
                },
                "recommendYes": {
                    "description": "Provider",
                    "type": "integer"
                },
                "speed": {
                    "description": "Provider",
                    "type": "number"
                }
            }
        },
        "types.Subscription": {
            "type": "object",
            "properties": {
                "as": {
                    "type": "string"
                },
                "channel": {
                    "type": "string"
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
	Host:        "localhost:8080",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "ANTICAP API",
	Description: "This is a sample serice for managing orders",
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
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
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