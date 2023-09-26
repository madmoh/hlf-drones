package contract

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"golang.org/x/exp/slices"
)

type RecordsSC struct {
	contractapi.Contract
}

func (c *RecordsSC) GetClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {
	return ctx.GetClientIdentity().GetID()
}

func (c *RecordsSC) AddOperator(ctx contractapi.TransactionContextInterface, operatorId string) error {
	// TODO: Two options:
	// 1) Restrict call to ministry.
	// 2) Check claimed identity matches operator identity registered with the ministry. Two options:
	// 2.1) Operator claims he's registered with ministry + []providers. Each of them attests the claim.
	// 2.2) (Better) Incorporate it as part of the endorsing mechanism.
	exists, err := KeyExists(ctx, operatorId)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("Operator %v already exists", operatorId)
	}
	operator := Operator{
		OperatorId:      operatorId,
		CertificateTier: "NO_CERTIFICATE",
		DroneIds:        make([]string, 0),
		FlightIds:       make([]string, 0),
		ViolationIds:    make([]string, 0),
		Status:          "NORMAL", // NORMAL, (TERM/PERM)_BAN[_APP] for temporary/permanent bans with possible appeal
	}
	operatorJSON, err := json.Marshal(operator)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(operatorId, operatorJSON)
}

func (c *RecordsSC) DeleteOperator(ctx contractapi.TransactionContextInterface, operatorId string) error {
	exists, err := KeyExists(ctx, operatorId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Operator %v does not exist", operatorId)
	}
	return ctx.GetStub().DelState(operatorId)
}

// TODO: Make some link between the drone and operator
func (c *RecordsSC) AddDrone(ctx contractapi.TransactionContextInterface, droneId string, expiry time.Time) error {
	// TODO: Similar checks to AddOperator
	exists, err := KeyExists(ctx, droneId)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("Drone %v already exists", droneId)
	}
	drone := Drone{
		DroneId: droneId,
		Expiry:  expiry,
		// RemoteId: remoteId,
	}
	droneJSON, err := json.Marshal(drone)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(droneId, droneJSON)
}

func (c *RecordsSC) DeleteDrone(ctx contractapi.TransactionContextInterface, dronerId string) error {
	exists, err := KeyExists(ctx, dronerId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Drone %v does not exist", dronerId)
	}
	return ctx.GetStub().DelState(dronerId)
}

// TODO: More business logic checks (caller identity)
func (c *RecordsSC) AddCertificate(ctx contractapi.TransactionContextInterface, operatorId string, tier string, expiry time.Time) error {
	exists, err := KeyExists(ctx, operatorId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Operator %v does not exist", operatorId)
	}
	operatorJSON, err := ctx.GetStub().GetState(operatorId)
	if err != nil {
		return fmt.Errorf("failed to read from state. Error: %v", err)
	}
	if operatorJSON == nil {
		return fmt.Errorf("Operator %v does not exist", operatorId)
	}
	var operator Operator
	err = json.Unmarshal(operatorJSON, &operator)
	if err != nil {
		return err
	}
	operator.CertificateTier = tier
	operator.CertificateExpiry = expiry
	operatorJSON, err = json.Marshal(operator)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(operatorId, operatorJSON)
}

func (c *RecordsSC) RequestPermit(ctx contractapi.TransactionContextInterface, operatorId string, flightId string, droneId string, permitEffective time.Time, permitExpiry time.Time, vertices [][][3]float64, facets [][][3]uint64, types []bool) error {
	// operatorId, err := ctx.GetClientIdentity().GetID()
	// if err != nil {
	// 	return fmt.Errorf("cannot get client identity. Error: %v", err)
	// }
	// exists, err := KeyExists(ctx, operatorId)
	// if err != nil {
	// 	return err
	// }
	// if !exists {
	// 	return fmt.Errorf("operator %v does not exist", operatorId)
	// }
	operatorJSON, err := ctx.GetStub().GetState(operatorId)
	if err != nil {
		return err
	}
	if operatorJSON == nil {
		return fmt.Errorf("operator %v does not exist", operatorId)
	}
	var operator Operator
	err = json.Unmarshal(operatorJSON, &operator)
	if err != nil {
		return fmt.Errorf("failed to unmarshal. Error: %v", err)
	}
	exists := slices.Contains(operator.FlightIds, flightId)
	if exists {
		return fmt.Errorf("operator already has flight %v", flightId)
	}
	// tree := kdtree.Tree{Root: nil, Count: 0}
	flight := Flight{
		FlightId:           flightId,
		OperatorId:         operatorId,
		DroneId:            droneId,
		PermitEffective:    permitEffective,
		PermitExpiry:       permitExpiry,
		BoundariesVertices: vertices,
		BoundariesFacets:   facets,
		BoundariesTypes:    types,
		// Tree:             &tree,
		Status:  "PENDING",
		Beacons: make([][3]float64, 0),
	}
	flightJSON, err := json.Marshal(flight)
	if err != nil {
		return fmt.Errorf("failed to marshal flight %v. Error: %v", flightJSON, err)
	}
	err = ctx.GetStub().PutState(flightId, flightJSON)
	if err != nil {
		return err
	}
	operator.FlightIds = append(operator.FlightIds, flightId)
	operatorJSON, err = json.Marshal(operator)
	if err != nil {
		return fmt.Errorf("failed to marshal operator %v. Error: %v", operatorJSON, err)
	}
	return ctx.GetStub().PutState(operatorId, operatorJSON)
}

func (c *RecordsSC) EvaluatePermit(ctx contractapi.TransactionContextInterface, flightId string, decision string) error {
	exists, err := KeyExists(ctx, flightId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Flight %v does not exist", flightId)
	}
	flightJSON, err := ctx.GetStub().GetState(flightId)
	if err != nil {
		return fmt.Errorf("failed to read from state. Error: %v", err)
	}
	if flightJSON == nil {
		return fmt.Errorf("Flight %v does not exist", flightId)
	}
	var flight Flight
	err = json.Unmarshal(flightJSON, &flight)
	if err != nil {
		return err
	}
	flight.Status = decision
	// TODO: Supposed to build kdtree here (once for the boundary) but cannot since the tree struct is not allowed in ctx
	// if decision == "APPROVED" {
	// 	// vertices2D := isinside.ConvertFloat64To2D(flight.BoundaryVertices)
	// 	// facets2D := isinside.ConvertUint64To2D(flight.BoundaryFacets)
	// 	tree, _, _ := isinside.GenerateKDTreePlus(flight.BoundaryVertices, flight.BoundaryFacets)
	// 	flight.Tree = tree
	// }
	flightJSON, err = json.Marshal(flight)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(flightId, flightJSON)
}

func (c *RecordsSC) DeleteFlight(ctx contractapi.TransactionContextInterface, flightId string) error {
	exists, err := KeyExists(ctx, flightId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Flight %v does not exist", flightId)
	}
	return ctx.GetStub().DelState(flightId)
}

func (c *RecordsSC) UpdateOperatorStatus(ctx contractapi.TransactionContextInterface, operatorId string, status string) error {
	exists, err := KeyExists(ctx, operatorId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Operator %v does not exist", operatorId)
	}
	operatorJSON, err := ctx.GetStub().GetState(operatorId)
	if err != nil {
		return fmt.Errorf("failed to read from state. Error: %v", err)
	}
	if operatorJSON == nil {
		return fmt.Errorf("Operator %v does not exist", operatorId)
	}
	var operator Operator
	err = json.Unmarshal(operatorJSON, &operator)
	if err != nil {
		return err
	}
	operator.Status = status
	operatorJSON, err = json.Marshal(operator)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(operatorId, operatorJSON)
}

func (c *RecordsSC) GetOperator(ctx contractapi.TransactionContextInterface, key string) (*Operator, error) {
	objectJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, err
	}
	if objectJSON == nil {
		return nil, fmt.Errorf("object %s does not exist", key)
	}
	var object Operator
	err = json.Unmarshal(objectJSON, &object)
	if err != nil {
		return nil, err
	}
	return &object, nil
}

func (c *RecordsSC) GetDrone(ctx contractapi.TransactionContextInterface, key string) (*Drone, error) {
	objectJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, err
	}
	if objectJSON == nil {
		return nil, fmt.Errorf("object %s does not exist", key)
	}
	var object Drone
	err = json.Unmarshal(objectJSON, &object)
	if err != nil {
		return nil, err
	}
	return &object, nil
}
