package schemas

const USER_CREATE_JSON_SCHEMA = `
{
	"title": "user",
	"type": "object",
	"properties": {
		"user": {
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
				},
				"age":{
					"type": "integer",
					"minimum": 0,
					"maximum": 156
				},
				"sex":{
					"type": "string",
					"enum": ["men", "women"]
				}
			},
			"required": ["name", "password", "age", "sex"],
			"additionalProperties": false
		}
	},
	"required": ["user"],
	"additionalProperties": false
}
`
