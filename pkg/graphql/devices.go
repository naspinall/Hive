package graphql

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/naspinall/Hive/pkg/models"
)

type DeviceResolver struct {
	models.DeviceService
	DeviceType *graphql.Object
	QueryType  *graphql.Object
	Schema     graphql.Schema
}

func (dr *DeviceResolver) Device(q graphql.ResolveParams) (interface{}, error) {
	id, ok := q.Args["id"].(int)
	if ok {
		device, err := dr.ByID(uint(id), q.Context)
		if err != nil {
			return nil, err
		}
		return device, err
	}

	return nil, nil
}

func (dr *DeviceResolver) Devices(q graphql.ResolveParams) (interface{}, error) {
	return dr.Many(100, q.Context)
}

func (dr *DeviceResolver) DeviceCount(q graphql.ResolveParams) (interface{}, error) {
	count, err := dr.Count(q.Context)
	if err != nil {
		return 0, nil
	}
	return &Count{count}, nil
}

func NewDeviceResolver(deviceService models.DeviceService) *DeviceResolver {

	dr := &DeviceResolver{
		DeviceService: deviceService,
	}

	dr.DeviceType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Device",
			Fields: graphql.Fields{
				"ID": &graphql.Field{
					Type: graphql.Int,
				},
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"imei": &graphql.Field{
					Type: graphql.String,
				},
				"longitude": &graphql.Field{
					Type: graphql.Float,
				},
				"latitude": &graphql.Field{
					Type: graphql.Float,
				},
			},
		},
	)

	dr.QueryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"device": &graphql.Field{
					Type:        dr.DeviceType,
					Description: "Get device by id",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
					},
					Resolve: dr.Device,
				},
				"devices": &graphql.Field{
					Type:        graphql.NewList(dr.DeviceType),
					Description: "Get device by id",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
					},
					Resolve: dr.Devices,
				},
			},
		})

	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query: dr.QueryType,
		},
	)
	if err != nil {
		log.Println("Query Compilation Error", err)
	}

	dr.Schema = schema

	return dr
}

func (dr *DeviceResolver) HandleQuery(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	result := graphql.Do(graphql.Params{
		Schema:        dr.Schema,
		RequestString: query,
		Context:       r.Context(),
	})

	if result.HasErrors() {
		log.Printf("Graphql Error: %v", result.Errors)
	}

	json.NewEncoder(w).Encode(result)
}
