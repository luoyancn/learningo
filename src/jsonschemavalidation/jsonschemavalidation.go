package jsonschemavalidation

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

/**
type uuid_format_checker struct{}

func (this uuid_format_checker) IsFormat(input interface{}) bool {
	uuid_string, ok := input.(string)
	if false == ok {
		return false
	}
	_, err := uuid.FromString(uuid_string)
	if nil != err {
		fmt.Printf("ERROR: %v\n", err)
		return false
	}
	return true
}
gojsonschema.FormatCheckers.Add("uuid", uuid_format_checker{})
*/

var user_json_schema = `
	{"type": "string"}
`

var complex_json_schema = `
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
	schema_loader := gojsonschema.NewStringLoader(user_json_schema)
	data_loader := gojsonschema.NewStringLoader(`"lucifer"`)
	result, err := gojsonschema.Validate(schema_loader, data_loader)
	if nil != err {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, err := range result.Errors() {
			fmt.Printf("- %s\n", err)
		}
		return
	}
	fmt.Printf("The document is valid\n")
}

func ValidateJsonSchema_complex() {
	schema_loader := gojsonschema.NewStringLoader(complex_json_schema)
	data_loader := gojsonschema.NewStringLoader(
		`{"id": "c9b83353-5005-4d31-b21f-e99f3ae33386",
			"name": "zhangjl", "sex": "men", "desc": "zhangjl",
			"married": true, "addresses": ["beijing", "chongqing", "wuhan"]}`)
	result, err := gojsonschema.Validate(schema_loader, data_loader)
	if nil != err {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, err := range result.Errors() {
			fmt.Printf("- %s\n", err)
		}
		return
	}
	fmt.Printf("The document is valid\n")
}
