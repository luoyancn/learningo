package jsonschemavalidation

import (
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

const user_json_schema = `
	{"type": "string"}
`

const complex_json_schema = `
{
	"title": "Person",
	"type": "object",
	"properties":{
		"id": {
			"type": "string",
			"minLength": 32,
			"maxLength": 64,
			"format": "uuid"
		},
		"name": {
			"type": "string",
			"minLength": 4,
			"maxLength": 36
		},
		"email": {
			"type": "string",
			"maxLength": 64,
			"format": "email"
		},
		"sex": {
			"type": "string",
			"enum": ["men", "women"],
			"default": "men"
		},
		"desc": {
			"type": ["integer", "string", "null"],
			"default": null
		},
		"married": {
			"type": "boolean",
			"default": false
		},
		"addresses": {
			"type": "array",
			"items": {
				"type": "string",
				"minLength": 4,
				"maxLength": 36
			}
		}
	},
	"required" : ["name", "sex", "id"]
}
`

func ValidateJsonSchema() {
	schema := gojsonschema.NewStringLoader(user_json_schema)
	docs := gojsonschema.NewStringLoader(`"lucifer"`)
	res, err := gojsonschema.Validate(schema, docs)
	if err != nil {
		fmt.Printf("ERROR:%v\n", err)
		return
	}

	if !res.Valid() {
		var errs []string
		errs = append(errs, fmt.Sprintln(
			"Input is invalid, following errors found:"))
		for _, desc := range res.Errors() {
			errs = append(errs, fmt.Sprintf("- %s", desc))
		}
		fmt.Println(strings.Join(errs, "\n"))
		return
	}
	fmt.Println("The docs is valid!")
}

func Validate_complext_schema() {
	schema := gojsonschema.NewStringLoader(complex_json_schema)
	docs := gojsonschema.NewStringLoader(
		`{"name": "zhangjl", "sex": "men",
		"id": "299dbc8d-0bca-4189-9450-2f5bdf3aaf80",
		"address": ["beijing", "chongqing", "wuhan"],
		"married": true, "email": "zhangjl@hasp.io"}`)
	res, err := gojsonschema.Validate(schema, docs)
	if err != nil {
		fmt.Printf("ERROR:%v\n", err)
		return
	}

	if !res.Valid() {
		var errs []string
		errs = append(errs, fmt.Sprintln(
			"Input is invalid, following errors found:"))
		for _, desc := range res.Errors() {
			errs = append(errs, fmt.Sprintf("- %s", desc))
		}
		fmt.Println(strings.Join(errs, "\n"))
		return
	}
	fmt.Println("The docs is valid!")
}
