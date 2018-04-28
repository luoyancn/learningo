package schemas

const AUTH_JSON_SCHEMA = `
{
	"title": "Auth",
	"type": "object",
	"properties": {
		"auth": {
			"type": "object",
			"properties":{
				"name":{
					"type": "string",
					"minLength": 4,
					"maxLength": 36
				},
				"password":{
					"type": "string",
					"minLength": 4,
					"maxLength": 36
				}
			},
			"required": ["name", "password"],
			"additionalProperties": false
		}
	},
	"required": ["auth"],
	"additionalProperties": false
}
`
