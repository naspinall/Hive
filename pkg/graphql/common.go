package graphql

import "github.com/graphql-go/graphql"

var CountType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Count",
		Fields: graphql.Fields{
			"count": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

type Count struct {
	count uint `json:"count"`
}
