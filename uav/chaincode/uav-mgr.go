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
	OperatorId        string            `json:"OperatorId"`      // TODO: Unique accross operators
	CertificateTier   string            `json:"CertificateTier"` // TODO: Could also be an enum
	CertificateExpiry time.Time         `json:CertificateExpiry`
	ActivePermits     map[string]Permit `json:"ActivePermits"`
}

type Drone struct {
	DroneId  string    `json:"DroneId"` // TODO: Unique accross drones
	Expiry   time.Time `json:"Expiry"`
	RemoteId string    `json:"RemoteId"` // TODO: Unique accross drones
}

type Permit struct {
	// PermitId         string    `json:"PermitId"`
	// OperatorId       string    `json:"OperatorId"`
	DroneId          string    `json:"DroneId"`
	PermitEffective  time.Time `json:"PermitEffective"`
	PermitExpiry     time.Time `json:"PermitExpiry"`
	BoundaryVertices []float64 `json:"BoundaryVertices"`
	BoundaryFacets   []float64 `json:"BoundaryFacets"`
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
		OperatorId:        operatorId,
		CertificateTier:   "NO_CERTIFICATE",
		CertificateExpiry: time.Now(),
		ActivePermits:     make(map[string]Permit),
	}
	operatorJSON, err := json.Marshal(operator)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(operatorId, operatorJSON)
}

func (c *UAVContract) AddDrone(ctx contractapi.TransactionContextInterface, droneId string, expiry time.Time, remoteId string) error {
	// TODO: Similar checks to AddOperator
	exists, err := c.KeyExists(ctx, droneId)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("Drone %v already exists", droneId)
	}
	drone := Drone{
		DroneId:  droneId,
		Expiry:   expiry,
		RemoteId: remoteId,
	}
	droneJSON, err := json.Marshal(drone)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(droneId, droneJSON)
}

func (c *UAVContract) AddCertificate(ctx contractapi.TransactionContextInterface, operatorId string, tier string, expiry time.Time) error {
	exists, err := c.KeyExists(ctx, operatorId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Operator %v does not exist", operatorId)
	}
	operator := Operator{
		OperatorId:        operatorId,
		CertificateTier:   tier,
		CertificateExpiry: expiry,
	}
	operatorJSON, err := json.Marshal(operator)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(operatorId, operatorJSON)
}

func (c *UAVContract) RequestPermit(ctx contractapi.TransactionContextInterface, droneId string, permitEffective time.Time, permitExpiry time.Time, vertices []float64, facets []float64) error {
	operatorId, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return err
	}
	exists, err := c.KeyExists(ctx, operatorId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Operator %v does not exist", operatorId)
	}
	operatorJSON, err := ctx.GetStub().GetState(operatorId)
	if err != nil {
		return fmt.Errorf("Failed to read from state. Error: %v", err)
	}
	if operatorJSON == nil {
		return fmt.Errorf("Operator %v does not exist", operatorId)
	}
	var operator Operator
	err = json.Unmarshal(operatorJSON, operator)
	permitId := fmt.Sprintf("%v", time.Now().UnixNano())
	_, exists = operator.ActivePermits[permitId]
	if exists {
		return fmt.Errorf("Permit %v already exists", permitId)
	}
	permit := Permit{
		DroneId:          droneId,
		PermitEffective:  permitEffective,
		PermitExpiry:     permitExpiry,
		BoundaryVertices: vertices,
		BoundaryFacets:   facets,
	}
	operator.ActivePermits[permitId] = permit
	operatorJSON, err = json.Marshal(operator)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(operatorId, operatorJSON)
}

// RecordsSC.RejectPermit(PermitId)
// - Log the error and invalidate the request

// RecordsSC.AcceptPermit(PermitId)
// - Builds the kdtree

// Flights.LogTakeoff(UASId).
// - Record the time of takingoff to start expecting n RemoteId messages every n seconds (require n < ?)

// Flights.LogBeacons(...?)
// - Wait until N (cumulative n) is > ...? then execute Flights.AnalyzeBeacons

// Flights.AnalyzeBeacons(...?)
// - Execute isinside
// - Invokes Warnings.ReportFlight(...)

// Flights.LogLanding(UASId)

// ?.BanOperator or ?.BanDrone

func (c *UAVContract) KeyExists(ctx contractapi.TransactionContextInterface, key string) (bool, error) {
	valueJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return false, fmt.Errorf("Failed to read world state. Error: %v", err)
	}
	return valueJSON != nil, nil
}
