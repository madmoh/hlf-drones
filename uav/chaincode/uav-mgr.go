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
	OperatorId         string    `json:"OperatorId"`      // TODO: Unique accross operators
	CertificateTier    string    `json:"CertificateTier"` // TODO: Could also be an enum
	CertificateExpiary time.Time `json:CertificateExpiary`
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
		OperatorId:         operatorId,
		CertificateTier:    "NO_CERTIFICATE",
		CertificateExpiary: time.Now(),
	}
	operatorJSON, err := json.Marshal(operator)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(operatorId, operatorJSON)
}

func (c *UAVContract) AddDrone(ctx contractapi.TransactionContextInterface, droneId string, expiary time.Time, remoteId string) error {
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
		Expiary:  expiary,
		RemoteId: remoteId,
	}
	droneJSON, err := json.Marshal(drone)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(droneId, droneJSON)
}

func (c *UAVContract) AddCertificate(ctx contractapi.TransactionContextInterface, operatorId string, tier string, expiary time.Time) error {
	exists, err := c.KeyExists(ctx, operatorId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Operator %v does not exist", operatorId)
	}
	operator := Operator{
		OperatorId:         operatorId,
		CertificateTier:    tier,
		CertificateExpiary: expiary,
	}
	operatorJSON, err := json.Marshal(operator)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(operatorId, operatorJSON)
}

// RecordsSC.RequestPermit(UASId, PermitStartDateTime, PermitEndDateTime, {Geo bounds})
// - Requires manual approval by the Provider ("forward" to the operator's providers)

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
