package db

import (
	"fmt"
	"frdocker/db/models"
	"frdocker/types"
	"testing"
)

func TestDB(t *testing.T) {
	container := &models.Container{
		Container: &types.Container{
			IP:      "12345678",
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
		NetWorkId: 1,
	}
	ContainerMgo.InsertOne(container)
	container = &models.Container{}
	ContainerMgo.FindOne("container.ip", "123456").Decode(container)
	fmt.Println(container)

}
