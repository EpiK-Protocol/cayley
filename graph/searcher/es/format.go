package es

const (
	index_node_mapping = `{
		"settings": {
			"number_of_shards": 1
		},
		"mappings": {
			"properties": {
				"id": { "type": "integer" },
				"content": { "type": "text" },
				"type": { "type": "text" }
			}
		}
	}`

	index_meta_mapping = `{
		"settings": {
			"number_of_shards": 1
		},
		"mappings": {
			"properties": {
				"maxid": { "type": "integer" }
			}
		}
	}`

	search_node = `{
		"query":{
			"bool":{
				"must":{
					"match":{
						"content":{
							"query": "%s",
							"operator": "and"
						}
					}
				},
				"filter":{
					"term":{
						"type":"string"
					}
				}
			}
		}
	}`
)
