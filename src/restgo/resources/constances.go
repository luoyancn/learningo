package resources

const user_create_json_schema = `
{
	"title": "User",
	"type": "object",
	"properties": {
		"name": {
			"type": "string",
			"minLength": 4,
			"maxLength": 36
		},
		"age": {
			"type": "integer",
			"minimux": 0,
			"maximux": 200
		},
		"sex": {
			"type": "string",
			"enum": ["men", "women"]
		}
	},
	"required": ["name", "sex", "age"],
	"additionalProperties": false
}
`

const user_update_json_schema = `
{
	"title": "User",
	"type": "object",
	"properties": {
		"name": {
			"type": "string",
			"minLength": 4,
			"maxLength": 36
		},
		"age": {
			"type": "integer",
			"minimux": 0,
			"maximux": 200
		}
	},
	"additionalProperties": false,
	"anyOf": [
		{
			"required": ["name"]
		},
		{
			"required": ["age"]
		}
	]
}
`

const role_create_json_schema = `
{
	"title": "Role",
	"type": "object",
	"properties": {
		"rolename": {
			"type": "string",
			"minLength": 4,
			"maxLength": 36
		}
	},
	"required": ["rolename"],
	"additionalProperties": false
}
`
