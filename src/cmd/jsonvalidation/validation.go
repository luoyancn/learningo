package main

import (
	"jsonschemavalidation"
)

func main() {
	jsonschemavalidation.ValidateJsonSchema()
	jsonschemavalidation.ValidateJsonSchema_complex()
}
