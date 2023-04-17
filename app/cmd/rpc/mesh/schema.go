package mesh

var fallbackNodeFileSchema = `
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "type": "array",
  "minItems": 1,
  "uniqueItems": true,
  "items": [
    {
      "type": "object",
      "properties": {
        "priv_key": {
          "type": "string"
        },
        "servicer_url": {
          "type": "string"
        }
      },
      "additionalProperties": true,
      "required": [
        "priv_key",
        "servicer_url"
      ]
    }
  ]
}
`

var nodeFileSchema = `
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "type": "array",
  "minItems": 1,
  "uniqueItems": true,
  "items": [
    {
      "type": "object",
      "required": ["url", "keys"],
	  "additionalProperties": true,
      "properties": {
        "name": {
          "type": "string"
		},
        "url": {
          "type": "string"
        },
        "keys": {
          "type": "array",
          "minItems": 1,
          "uniqueItems": true,
          "items": [
            {
              "type": "string"
            }
          ]
        }
      }
    }
  ]
}
`

var plainChainsMapSchema = `
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "type": "object",
  "patternProperties": {
    ".*": {
      "type": "string"
    }
  }
}
`

var richChainsMapSchema = `
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "type": "object",
  "patternProperties": {
    ".*": {
      "type": "object",
	  "required": ["label"],
	  "additionalProperties": true,
      "properties": {
        "label": {
          "type": "string"
		}
      }
    }
  }
}
`
