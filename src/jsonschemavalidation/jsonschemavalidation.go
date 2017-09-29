package jsonschemavalidation

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

var user_json_schema = `
	{"type": "string"}
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
