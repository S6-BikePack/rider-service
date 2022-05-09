// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
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
        "/api/riders": {
            "get": {
                "description": "gets all riders in the system",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "get all riders",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.ridersResponse"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "creates a new rider",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "create rider",
                "parameters": [
                    {
                        "description": "Add rider",
                        "name": "rider",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.BodyCreateRider"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.RiderResponse"
                        }
                    }
                }
            }
        },
        "/api/riders/{id}": {
            "get": {
                "description": "gets a rider from the system by its ID",
                "produces": [
                    "application/json"
                ],
                "summary": "get rider",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Rider id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.RiderResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "updates a rider's information",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "update rider",
                "parameters": [
                    {
                        "description": "Update rider",
                        "name": "rider",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.BodyCreateRider"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Rider id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.RiderResponse"
                        }
                    }
                }
            }
        },
        "/api/riders/{id}/location": {
            "put": {
                "description": "updates a rider's location",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "update rider location",
                "parameters": [
                    {
                        "description": "Update rider",
                        "name": "rider",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.BodyLocation"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Rider id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.RiderResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.BodyCreateRider": {
            "type": "object",
            "properties": {
                "capacity": {
                    "$ref": "#/definitions/dto.CreateDimensions"
                },
                "id": {
                    "type": "string"
                },
                "serviceArea": {
                    "type": "integer"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "dto.BodyLocation": {
            "type": "object",
            "properties": {
                "latitude": {
                    "type": "number"
                },
                "longitude": {
                    "type": "number"
                }
            }
        },
        "dto.CreateDimensions": {
            "type": "object",
            "properties": {
                "depth": {
                    "type": "integer"
                },
                "height": {
                    "type": "integer"
                },
                "width": {
                    "type": "integer"
                }
            }
        },
        "dto.RiderResponse": {
            "type": "object",
            "properties": {
                "capacity": {
                    "$ref": "#/definitions/dto.riderResponseCapacity"
                },
                "id": {
                    "type": "string"
                },
                "location": {
                    "$ref": "#/definitions/dto.riderResponseLocation"
                },
                "serviceArea": {
                    "$ref": "#/definitions/dto.riderResponseArea"
                },
                "status": {
                    "type": "integer"
                },
                "user": {
                    "$ref": "#/definitions/dto.riderResponseUser"
                }
            }
        },
        "dto.riderResponseArea": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "identifier": {
                    "type": "string"
                }
            }
        },
        "dto.riderResponseCapacity": {
            "type": "object",
            "properties": {
                "depth": {
                    "type": "integer"
                },
                "height": {
                    "type": "integer"
                },
                "width": {
                    "type": "integer"
                }
            }
        },
        "dto.riderResponseLocation": {
            "type": "object",
            "properties": {
                "latitude": {
                    "type": "number"
                },
                "longitude": {
                    "type": "number"
                }
            }
        },
        "dto.riderResponseUser": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "dto.ridersResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "serviceArea": {
                    "type": "integer"
                },
                "status": {
                    "type": "integer"
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
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
