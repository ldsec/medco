// Code generated by go-swagger; DO NOT EDIT.

package server

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "API of the MedCo connector, that orchestrates the query at the MedCo node and provides information about the MedCo network.",
    "title": "MedCo Connector",
    "contact": {
      "email": "medco-dev@listes.epfl.ch"
    },
    "license": {
      "name": "EULA",
      "url": "https://raw.githubusercontent.com/ldsec/medco-connector/master/LICENSE"
    },
    "version": "1.0.0"
  },
  "basePath": "/medco",
  "paths": {
    "/genomic-annotations/{annotation}": {
      "get": {
        "security": [
          {
            "medco-jwt": [
              "medco-genomic-annotations"
            ]
          }
        ],
        "tags": [
          "genomic-annotations"
        ],
        "summary": "Get genomic annotations values.",
        "operationId": "getValues",
        "parameters": [
          {
            "type": "string",
            "description": "Genomic annotation name.",
            "name": "annotation",
            "in": "path",
            "required": true
          },
          {
            "type": "string",
            "description": "Genomic annotation value.",
            "name": "value",
            "in": "query",
            "required": true
          },
          {
            "type": "integer",
            "default": 10,
            "description": "Limits the number of records retrieved.",
            "name": "limit",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "Queried annotation values.",
            "schema": {
              "type": "array",
              "items": {
                "type": "string"
              }
            }
          },
          "404": {
            "description": "Annotation not found."
          },
          "default": {
            "$ref": "#/responses/errorResponse"
          }
        }
      }
    },
    "/genomic-annotations/{annotation}/{value}": {
      "get": {
        "security": [
          {
            "medco-jwt": [
              "medco-genomic-annotations"
            ]
          }
        ],
        "tags": [
          "genomic-annotations"
        ],
        "summary": "Get variants corresponding to a genomic annotation value.",
        "operationId": "getVariants",
        "parameters": [
          {
            "type": "string",
            "description": "Genomic annotation name.",
            "name": "annotation",
            "in": "path",
            "required": true
          },
          {
            "type": "string",
            "description": "Genomic annotation value.",
            "name": "value",
            "in": "path",
            "required": true
          },
          {
            "type": "array",
            "items": {
              "enum": [
                "heterozygous",
                "homozygous",
                "unknown"
              ],
              "type": "string"
            },
            "default": [
              "heterozygous",
              "homozygous",
              "unknown"
            ],
            "description": "Genomic annotation zygosity.",
            "name": "zygosity",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "Queried variants.",
            "schema": {
              "type": "array",
              "items": {
                "type": "string"
              }
            }
          },
          "404": {
            "description": "Annotation or annotation value not found."
          },
          "default": {
            "$ref": "#/responses/errorResponse"
          }
        }
      }
    },
    "/network": {
      "get": {
        "security": [
          {
            "medco-jwt": [
              "medco-network"
            ]
          }
        ],
        "tags": [
          "medco-network"
        ],
        "summary": "Get network metadata.",
        "operationId": "getMetadata",
        "responses": {
          "200": {
            "$ref": "#/responses/networkMetadataResponse"
          },
          "default": {
            "$ref": "#/responses/errorResponse"
          }
        }
      }
    },
    "/node/explore/query": {
      "post": {
        "security": [
          {
            "medco-jwt": [
              "medco-explore"
            ]
          }
        ],
        "tags": [
          "medco-node"
        ],
        "summary": "MedCo-Explore query to the node.",
        "operationId": "exploreQuery",
        "parameters": [
          {
            "type": "boolean",
            "default": true,
            "description": "Request synchronous query (defaults to true).",
            "name": "sync",
            "in": "query"
          },
          {
            "$ref": "#/parameters/exploreQueryRequest"
          }
        ],
        "responses": {
          "200": {
            "$ref": "#/responses/exploreQueryResponse"
          },
          "default": {
            "$ref": "#/responses/errorResponse"
          }
        }
      }
    },
    "/node/explore/query/{queryId}": {
      "get": {
        "security": [
          {
            "medco-jwt": [
              "medco-explore"
            ]
          }
        ],
        "tags": [
          "medco-node"
        ],
        "summary": "Get status and result of a MedCo-Explore query.",
        "operationId": "getExploreQuery",
        "parameters": [
          {
            "type": "string",
            "description": "Query ID",
            "name": "queryId",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "$ref": "#/responses/exploreQueryResponse"
          },
          "404": {
            "description": "Query ID not found."
          },
          "default": {
            "$ref": "#/responses/errorResponse"
          }
        }
      }
    },
    "/node/explore/search": {
      "post": {
        "security": [
          {
            "medco-jwt": [
              "medco-explore"
            ]
          }
        ],
        "tags": [
          "medco-node"
        ],
        "summary": "Search through the ontology for MedCo-Explore query terms.",
        "operationId": "exploreSearch",
        "parameters": [
          {
            "$ref": "#/parameters/exploreSearchRequest"
          }
        ],
        "responses": {
          "200": {
            "$ref": "#/responses/exploreSearchResponse"
          },
          "default": {
            "$ref": "#/responses/errorResponse"
          }
        }
      }
    }
  },
  "definitions": {
    "exploreQuery": {
      "description": "MedCo-Explore query",
      "properties": {
        "differentialPrivacy": {
          "description": "differential privacy query parameters (todo)",
          "type": "object",
          "properties": {
            "queryBudget": {
              "type": "number"
            }
          }
        },
        "panels": {
          "description": "i2b2 panels (linked by an AND)",
          "type": "array",
          "items": {
            "type": "object",
            "required": [
              "not"
            ],
            "properties": {
              "items": {
                "description": "i2b2 items (linked by an OR)",
                "type": "array",
                "items": {
                  "type": "object",
                  "required": [
                    "encrypted"
                  ],
                  "properties": {
                    "encrypted": {
                      "type": "boolean"
                    },
                    "operator": {
                      "type": "string",
                      "enum": [
                        "exists",
                        "equals"
                      ]
                    },
                    "queryTerm": {
                      "type": "string"
                    },
                    "value": {
                      "type": "string"
                    }
                  }
                }
              },
              "not": {
                "description": "exclude the i2b2 panel",
                "type": "boolean"
              }
            }
          }
        },
        "type": {
          "$ref": "#/definitions/exploreQueryType"
        },
        "userPublicKey": {
          "type": "string"
        }
      }
    },
    "exploreQueryResultElement": {
      "type": "object",
      "properties": {
        "encryptedCount": {
          "type": "string"
        },
        "encryptedPatientList": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "status": {
          "type": "string",
          "enum": [
            "queued",
            "pending",
            "error",
            "available"
          ]
        },
        "timers": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "milliseconds": {
                "type": "integer",
                "format": "int64"
              },
              "name": {
                "type": "string"
              }
            }
          }
        }
      }
    },
    "exploreQueryType": {
      "type": "string",
      "enum": [
        "patient_list",
        "count_per_site",
        "count_per_site_obfuscated",
        "count_per_site_shuffled",
        "count_per_site_shuffled_obfuscated",
        "count_global",
        "count_global_obfuscated"
      ]
    },
    "exploreSearch": {
      "type": "object",
      "properties": {
        "path": {
          "type": "string"
        },
        "type": {
          "type": "string",
          "enum": [
            "children",
            "metadata"
          ]
        }
      }
    },
    "exploreSearchResultElement": {
      "type": "object",
      "required": [
        "leaf"
      ],
      "properties": {
        "code": {
          "type": "string"
        },
        "displayName": {
          "type": "string"
        },
        "leaf": {
          "type": "boolean"
        },
        "medcoEncryption": {
          "type": "object",
          "required": [
            "encrypted"
          ],
          "properties": {
            "childrenIds": {
              "type": "array",
              "items": {
                "type": "integer",
                "format": "int64"
              }
            },
            "encrypted": {
              "type": "boolean"
            },
            "id": {
              "type": "integer",
              "format": "int64"
            }
          }
        },
        "metadata": {
          "type": "object"
        },
        "name": {
          "type": "string"
        },
        "path": {
          "type": "string"
        },
        "type": {
          "type": "string",
          "enum": [
            "container",
            "concept",
            "concept_numeric",
            "concept_enum",
            "concept_text",
            "genomic_annotation"
          ]
        }
      }
    },
    "restApiAuthorization": {
      "type": "string",
      "enum": [
        "medco-network",
        "medco-explore",
        "medco-genomic-annotations"
      ]
    },
    "user": {
      "type": "object",
      "properties": {
        "authorizations": {
          "type": "object",
          "properties": {
            "exploreQuery": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/exploreQueryType"
              }
            },
            "restApi": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/restApiAuthorization"
              }
            }
          }
        },
        "id": {
          "type": "string"
        },
        "token": {
          "type": "string"
        }
      }
    }
  },
  "parameters": {
    "exploreQueryRequest": {
      "description": "MedCo-Explore query request.",
      "name": "queryRequest",
      "in": "body",
      "required": true,
      "schema": {
        "type": "object",
        "properties": {
          "id": {
            "type": "string"
          },
          "query": {
            "$ref": "#/definitions/exploreQuery"
          }
        }
      }
    },
    "exploreSearchRequest": {
      "description": "MedCo-Explore ontology search request.",
      "name": "searchRequest",
      "in": "body",
      "required": true,
      "schema": {
        "$ref": "#/definitions/exploreSearch"
      }
    }
  },
  "responses": {
    "errorResponse": {
      "description": "Error response.",
      "schema": {
        "type": "object",
        "properties": {
          "message": {
            "type": "string"
          }
        }
      }
    },
    "exploreQueryResponse": {
      "description": "MedCo-Explore query response.",
      "schema": {
        "type": "object",
        "properties": {
          "id": {
            "type": "string"
          },
          "query": {
            "$ref": "#/definitions/exploreQuery"
          },
          "result": {
            "$ref": "#/definitions/exploreQueryResultElement"
          }
        }
      }
    },
    "exploreSearchResponse": {
      "description": "MedCo-Explore search query response.",
      "schema": {
        "type": "object",
        "properties": {
          "results": {
            "type": "array",
            "items": {
              "$ref": "#/definitions/exploreSearchResultElement"
            }
          },
          "search": {
            "$ref": "#/definitions/exploreSearch"
          }
        }
      }
    },
    "networkMetadataResponse": {
      "description": "Network metadata (public key and nodes list).",
      "schema": {
        "type": "object",
        "properties": {
          "nodeIndex": {
            "type": "integer"
          },
          "nodes": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "index": {
                  "type": "integer"
                },
                "url": {
                  "type": "string"
                }
              }
            }
          },
          "public-key": {
            "description": "Aggregated public key of the collective authority.",
            "type": "string"
          }
        }
      }
    }
  },
  "securityDefinitions": {
    "medco-jwt": {
      "description": "MedCo JWT token.",
      "type": "oauth2",
      "flow": "application",
      "tokenUrl": "https://medco-demo.epfl.ch/auth"
    }
  },
  "security": [
    {
      "medco-jwt": []
    }
  ],
  "tags": [
    {
      "description": "MedCo Network API",
      "name": "medco-network"
    },
    {
      "description": "MedCo Node API",
      "name": "medco-node"
    },
    {
      "description": "Genomic Annotations Query API",
      "name": "genomic-annotations"
    }
  ],
  "externalDocs": {
    "description": "MedCo Technical Documentation",
    "url": "https://medco.epfl.ch/documentation"
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "API of the MedCo connector, that orchestrates the query at the MedCo node and provides information about the MedCo network.",
    "title": "MedCo Connector",
    "contact": {
      "email": "medco-dev@listes.epfl.ch"
    },
    "license": {
      "name": "EULA",
      "url": "https://raw.githubusercontent.com/ldsec/medco-connector/master/LICENSE"
    },
    "version": "1.0.0"
  },
  "basePath": "/medco",
  "paths": {
    "/genomic-annotations/{annotation}": {
      "get": {
        "security": [
          {
            "medco-jwt": [
              "medco-genomic-annotations"
            ]
          }
        ],
        "tags": [
          "genomic-annotations"
        ],
        "summary": "Get genomic annotations values.",
        "operationId": "getValues",
        "parameters": [
          {
            "type": "string",
            "description": "Genomic annotation name.",
            "name": "annotation",
            "in": "path",
            "required": true
          },
          {
            "type": "string",
            "description": "Genomic annotation value.",
            "name": "value",
            "in": "query",
            "required": true
          },
          {
            "type": "integer",
            "default": 10,
            "description": "Limits the number of records retrieved.",
            "name": "limit",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "Queried annotation values.",
            "schema": {
              "type": "array",
              "items": {
                "type": "string"
              }
            }
          },
          "404": {
            "description": "Annotation not found."
          },
          "default": {
            "description": "Error response.",
            "schema": {
              "type": "object",
              "properties": {
                "message": {
                  "type": "string"
                }
              }
            }
          }
        }
      }
    },
    "/genomic-annotations/{annotation}/{value}": {
      "get": {
        "security": [
          {
            "medco-jwt": [
              "medco-genomic-annotations"
            ]
          }
        ],
        "tags": [
          "genomic-annotations"
        ],
        "summary": "Get variants corresponding to a genomic annotation value.",
        "operationId": "getVariants",
        "parameters": [
          {
            "type": "string",
            "description": "Genomic annotation name.",
            "name": "annotation",
            "in": "path",
            "required": true
          },
          {
            "type": "string",
            "description": "Genomic annotation value.",
            "name": "value",
            "in": "path",
            "required": true
          },
          {
            "type": "array",
            "items": {
              "enum": [
                "heterozygous",
                "homozygous",
                "unknown"
              ],
              "type": "string"
            },
            "default": [
              "heterozygous",
              "homozygous",
              "unknown"
            ],
            "description": "Genomic annotation zygosity.",
            "name": "zygosity",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "Queried variants.",
            "schema": {
              "type": "array",
              "items": {
                "type": "string"
              }
            }
          },
          "404": {
            "description": "Annotation or annotation value not found."
          },
          "default": {
            "description": "Error response.",
            "schema": {
              "type": "object",
              "properties": {
                "message": {
                  "type": "string"
                }
              }
            }
          }
        }
      }
    },
    "/network": {
      "get": {
        "security": [
          {
            "medco-jwt": [
              "medco-network"
            ]
          }
        ],
        "tags": [
          "medco-network"
        ],
        "summary": "Get network metadata.",
        "operationId": "getMetadata",
        "responses": {
          "200": {
            "description": "Network metadata (public key and nodes list).",
            "schema": {
              "type": "object",
              "properties": {
                "nodeIndex": {
                  "type": "integer"
                },
                "nodes": {
                  "type": "array",
                  "items": {
                    "type": "object",
                    "properties": {
                      "index": {
                        "type": "integer"
                      },
                      "url": {
                        "type": "string"
                      }
                    }
                  }
                },
                "public-key": {
                  "description": "Aggregated public key of the collective authority.",
                  "type": "string"
                }
              }
            }
          },
          "default": {
            "description": "Error response.",
            "schema": {
              "type": "object",
              "properties": {
                "message": {
                  "type": "string"
                }
              }
            }
          }
        }
      }
    },
    "/node/explore/query": {
      "post": {
        "security": [
          {
            "medco-jwt": [
              "medco-explore"
            ]
          }
        ],
        "tags": [
          "medco-node"
        ],
        "summary": "MedCo-Explore query to the node.",
        "operationId": "exploreQuery",
        "parameters": [
          {
            "type": "boolean",
            "default": true,
            "description": "Request synchronous query (defaults to true).",
            "name": "sync",
            "in": "query"
          },
          {
            "description": "MedCo-Explore query request.",
            "name": "queryRequest",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "properties": {
                "id": {
                  "type": "string"
                },
                "query": {
                  "$ref": "#/definitions/exploreQuery"
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "MedCo-Explore query response.",
            "schema": {
              "type": "object",
              "properties": {
                "id": {
                  "type": "string"
                },
                "query": {
                  "$ref": "#/definitions/exploreQuery"
                },
                "result": {
                  "$ref": "#/definitions/exploreQueryResultElement"
                }
              }
            }
          },
          "default": {
            "description": "Error response.",
            "schema": {
              "type": "object",
              "properties": {
                "message": {
                  "type": "string"
                }
              }
            }
          }
        }
      }
    },
    "/node/explore/query/{queryId}": {
      "get": {
        "security": [
          {
            "medco-jwt": [
              "medco-explore"
            ]
          }
        ],
        "tags": [
          "medco-node"
        ],
        "summary": "Get status and result of a MedCo-Explore query.",
        "operationId": "getExploreQuery",
        "parameters": [
          {
            "type": "string",
            "description": "Query ID",
            "name": "queryId",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "MedCo-Explore query response.",
            "schema": {
              "type": "object",
              "properties": {
                "id": {
                  "type": "string"
                },
                "query": {
                  "$ref": "#/definitions/exploreQuery"
                },
                "result": {
                  "$ref": "#/definitions/exploreQueryResultElement"
                }
              }
            }
          },
          "404": {
            "description": "Query ID not found."
          },
          "default": {
            "description": "Error response.",
            "schema": {
              "type": "object",
              "properties": {
                "message": {
                  "type": "string"
                }
              }
            }
          }
        }
      }
    },
    "/node/explore/search": {
      "post": {
        "security": [
          {
            "medco-jwt": [
              "medco-explore"
            ]
          }
        ],
        "tags": [
          "medco-node"
        ],
        "summary": "Search through the ontology for MedCo-Explore query terms.",
        "operationId": "exploreSearch",
        "parameters": [
          {
            "description": "MedCo-Explore ontology search request.",
            "name": "searchRequest",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/exploreSearch"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "MedCo-Explore search query response.",
            "schema": {
              "type": "object",
              "properties": {
                "results": {
                  "type": "array",
                  "items": {
                    "$ref": "#/definitions/exploreSearchResultElement"
                  }
                },
                "search": {
                  "$ref": "#/definitions/exploreSearch"
                }
              }
            }
          },
          "default": {
            "description": "Error response.",
            "schema": {
              "type": "object",
              "properties": {
                "message": {
                  "type": "string"
                }
              }
            }
          }
        }
      }
    }
  },
  "definitions": {
    "exploreQuery": {
      "description": "MedCo-Explore query",
      "properties": {
        "differentialPrivacy": {
          "description": "differential privacy query parameters (todo)",
          "type": "object",
          "properties": {
            "queryBudget": {
              "type": "number"
            }
          }
        },
        "panels": {
          "description": "i2b2 panels (linked by an AND)",
          "type": "array",
          "items": {
            "type": "object",
            "required": [
              "not"
            ],
            "properties": {
              "items": {
                "description": "i2b2 items (linked by an OR)",
                "type": "array",
                "items": {
                  "type": "object",
                  "required": [
                    "encrypted"
                  ],
                  "properties": {
                    "encrypted": {
                      "type": "boolean"
                    },
                    "operator": {
                      "type": "string",
                      "enum": [
                        "exists",
                        "equals"
                      ]
                    },
                    "queryTerm": {
                      "type": "string"
                    },
                    "value": {
                      "type": "string"
                    }
                  }
                }
              },
              "not": {
                "description": "exclude the i2b2 panel",
                "type": "boolean"
              }
            }
          }
        },
        "type": {
          "$ref": "#/definitions/exploreQueryType"
        },
        "userPublicKey": {
          "type": "string"
        }
      }
    },
    "exploreQueryResultElement": {
      "type": "object",
      "properties": {
        "encryptedCount": {
          "type": "string"
        },
        "encryptedPatientList": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "status": {
          "type": "string",
          "enum": [
            "queued",
            "pending",
            "error",
            "available"
          ]
        },
        "timers": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "milliseconds": {
                "type": "integer",
                "format": "int64"
              },
              "name": {
                "type": "string"
              }
            }
          }
        }
      }
    },
    "exploreQueryType": {
      "type": "string",
      "enum": [
        "patient_list",
        "count_per_site",
        "count_per_site_obfuscated",
        "count_per_site_shuffled",
        "count_per_site_shuffled_obfuscated",
        "count_global",
        "count_global_obfuscated"
      ]
    },
    "exploreSearch": {
      "type": "object",
      "properties": {
        "path": {
          "type": "string"
        },
        "type": {
          "type": "string",
          "enum": [
            "children",
            "metadata"
          ]
        }
      }
    },
    "exploreSearchResultElement": {
      "type": "object",
      "required": [
        "leaf"
      ],
      "properties": {
        "code": {
          "type": "string"
        },
        "displayName": {
          "type": "string"
        },
        "leaf": {
          "type": "boolean"
        },
        "medcoEncryption": {
          "type": "object",
          "required": [
            "encrypted"
          ],
          "properties": {
            "childrenIds": {
              "type": "array",
              "items": {
                "type": "integer",
                "format": "int64"
              }
            },
            "encrypted": {
              "type": "boolean"
            },
            "id": {
              "type": "integer",
              "format": "int64"
            }
          }
        },
        "metadata": {
          "type": "object"
        },
        "name": {
          "type": "string"
        },
        "path": {
          "type": "string"
        },
        "type": {
          "type": "string",
          "enum": [
            "container",
            "concept",
            "concept_numeric",
            "concept_enum",
            "concept_text",
            "genomic_annotation"
          ]
        }
      }
    },
    "restApiAuthorization": {
      "type": "string",
      "enum": [
        "medco-network",
        "medco-explore",
        "medco-genomic-annotations"
      ]
    },
    "user": {
      "type": "object",
      "properties": {
        "authorizations": {
          "type": "object",
          "properties": {
            "exploreQuery": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/exploreQueryType"
              }
            },
            "restApi": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/restApiAuthorization"
              }
            }
          }
        },
        "id": {
          "type": "string"
        },
        "token": {
          "type": "string"
        }
      }
    }
  },
  "parameters": {
    "exploreQueryRequest": {
      "description": "MedCo-Explore query request.",
      "name": "queryRequest",
      "in": "body",
      "required": true,
      "schema": {
        "type": "object",
        "properties": {
          "id": {
            "type": "string"
          },
          "query": {
            "$ref": "#/definitions/exploreQuery"
          }
        }
      }
    },
    "exploreSearchRequest": {
      "description": "MedCo-Explore ontology search request.",
      "name": "searchRequest",
      "in": "body",
      "required": true,
      "schema": {
        "$ref": "#/definitions/exploreSearch"
      }
    }
  },
  "responses": {
    "errorResponse": {
      "description": "Error response.",
      "schema": {
        "type": "object",
        "properties": {
          "message": {
            "type": "string"
          }
        }
      }
    },
    "exploreQueryResponse": {
      "description": "MedCo-Explore query response.",
      "schema": {
        "type": "object",
        "properties": {
          "id": {
            "type": "string"
          },
          "query": {
            "$ref": "#/definitions/exploreQuery"
          },
          "result": {
            "$ref": "#/definitions/exploreQueryResultElement"
          }
        }
      }
    },
    "exploreSearchResponse": {
      "description": "MedCo-Explore search query response.",
      "schema": {
        "type": "object",
        "properties": {
          "results": {
            "type": "array",
            "items": {
              "$ref": "#/definitions/exploreSearchResultElement"
            }
          },
          "search": {
            "$ref": "#/definitions/exploreSearch"
          }
        }
      }
    },
    "networkMetadataResponse": {
      "description": "Network metadata (public key and nodes list).",
      "schema": {
        "type": "object",
        "properties": {
          "nodeIndex": {
            "type": "integer"
          },
          "nodes": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "index": {
                  "type": "integer"
                },
                "url": {
                  "type": "string"
                }
              }
            }
          },
          "public-key": {
            "description": "Aggregated public key of the collective authority.",
            "type": "string"
          }
        }
      }
    }
  },
  "securityDefinitions": {
    "medco-jwt": {
      "description": "MedCo JWT token.",
      "type": "oauth2",
      "flow": "application",
      "tokenUrl": "https://medco-demo.epfl.ch/auth"
    }
  },
  "security": [
    {
      "medco-jwt": []
    }
  ],
  "tags": [
    {
      "description": "MedCo Network API",
      "name": "medco-network"
    },
    {
      "description": "MedCo Node API",
      "name": "medco-node"
    },
    {
      "description": "Genomic Annotations Query API",
      "name": "genomic-annotations"
    }
  ],
  "externalDocs": {
    "description": "MedCo Technical Documentation",
    "url": "https://medco.epfl.ch/documentation"
  }
}`))
}
