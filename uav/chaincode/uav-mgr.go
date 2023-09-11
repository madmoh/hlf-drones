package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type UAVContract struct {
	contractapi.Contract
}

type Operator struct {
	OperatorId string `json:"OperatorId"` // TODO: Unique accross operators
	// Other attributes
}

type Drone struct {
	DroneId  string    `json:"DroneId"` // TODO: Unique accross drones
	Expiary  time.Time `json:"Expiary"`
	RemoteId string    `json:"RemoteId"` // TODO: Unique accross drones
	// Other attributes
}

func (c *UAVContract) AddOperator(ctx contractapi.TransactionContextInterface, operatorId string) error {
	// TODO: Two options:
	// 1) Restrict call to ministry.
	// 2) Check claimed identity matches operator identity registered with the ministry. Two options:
	// 2.1) Operator claims he's registered with ministry + []providers. Each of them attests the claim.
	// 2.2) (Better) Incorporate it as part of the endorsing mechanism.
	exists, err := c.KeyExists(ctx, operatorId)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("Operator %v already exists", operatorId)
	}
	operator := Operator{
		OperatorId: operatorId,
	}
	operatorJSON, err := json.Marshal(operator)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(operatorId, operatorJSON)
}

func (c *UAVContract) AddDrone(ctx, droneId string, expiary time.Time, remoteId string) error {

}

func (c *UAVContract) KeyExists(ctx contractapi.TransactionContextInterface, key string) (bool, error) {
	valueJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return false, fmt.Errorf("Failed to read world state. Error: %v", err)
	}
	return valueJSON != nil, nil
}
