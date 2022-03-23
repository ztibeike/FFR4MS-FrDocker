package db

import (
	"context"
	"fmt"
	"frdocker/models"
	"frdocker/types"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestDB(t *testing.T) {
	container := &models.Container{
		Container: &types.Container{
			IP:      "123456",
			Port:    "123456",
			ID:      "123456",
			Group:   "123456",
			Gateway: "123456",
			Name:    "123456",
			Leaf:    true,
			Health:  true,
			States: []*types.State{{
				Id: &types.StateId{
					StartWith: &types.StateEndpointEvent{
						IP:       "123456",
						HttpType: "123456",
					},
					EndWith: &types.StateEndpointEvent{
						IP:       "123456",
						HttpType: "123456",
					},
				},
				Ecc:      1.0,
				Variance: &types.Vector{Data: []float64{1.0, 2.0}},
				Sigma:    1.0,
				K:        1,
				MaxTime:  1.0,
				MinTime:  1.0,
			}},
		},
		NetworkId: "23233",
	}
	// ContainerMgo.InsertOne(container)
	// container = &models.Container{}
	// result := ContainerMgo.FindOne("container.ip", "1234567")
	// if result != nil {
	// 	result.Decode(container)
	// }
	// fmt.Println(container)
	ContainerMgo.ReplaceOne(bson.D{{Key: "container.ip", Value: "123456"}}, container)
}

func TestFindMany(t *testing.T) {
	var dbContainers []*models.Container
	filter := bson.D{
		{Key: "networkId", Value: "b0e7978c-e44e-454c-946e-66f3232467e2"},
	}
	cursor := ContainerMgo.FindMany(filter)
	cursor.All(context.TODO(), &dbContainers)
	dbContainers[0].Container.States = append(dbContainers[0].Container.States, &types.State{})
	fmt.Println(dbContainers)
}

func TestFindOne(t *testing.T) {
	filter := bson.D{
		{Key: "name", Value: "br-46facbce8c7"},
	}
	var network *models.NetWork
	NetworkMgo.FindOne(filter).Decode(&network)
	fmt.Println(network)
}
