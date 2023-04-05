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
      "additionalProperties": false,
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
      "properties": {
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
      },
      "additionalProperties": false,
      "required": [
        "url",
        "keys"
      ]
    }
  ]
}
`
