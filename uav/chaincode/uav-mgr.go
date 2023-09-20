package main

import (
	"encoding/json"
	"fmt"
	"log"
	"madmoh/hlf-uav/isinside"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"gonum.org/v1/gonum/spatial/kdtree"
)

type UAVContract struct {
	contractapi.Contract
}

// TODO: OperatorId, DroneId shouldn't be inside their respective structs

type Operator struct {
	OperatorId        string            `json:"OperatorId"`      // TODO: Unique accross operators
	CertificateTier   string            `json:"CertificateTier"` // TODO: Could also be an enum
	CertificateExpiry time.Time         `json:"CertificateExpiry"`
	Flights           map[string]Flight `json:"Flights"`
}

type Drone struct {
	DroneId  string    `json:"DroneId"` // TODO: Unique accross drones
	Expiry   time.Time `json:"Expiry"`
	RemoteId string    `json:"RemoteId"` // TODO: Unique accross drones
}

// TODO: Maybe a custom TimeOptional{bool,time.Time} is better
type Flight struct {
	// FlightId         string    `json:"FlightId"`
	// OperatorId       string    `json:"OperatorId"`
	DroneId          string       `json:"DroneId"`
	PermitEffective  time.Time    `json:"PermitEffective"`
	PermitExpiry     time.Time    `json:"PermitExpiry"`
	BoundaryVertices [][3]float64 `json:"BoundaryVertices"` // TODO: Expand to [][][3]
	BoundaryFacets   [][3]uint64  `json:"BoundaryFacets"`   // TODO: Expand to [][][3]
	Tree             *kdtree.Tree `json:"Tree"`
	Status           string       `json:"Status"` // PENDING, REJECTED, APPROVED, ACTIVE
	Takeoff          time.Time    `json:"Takeoff"`
	Landing          time.Time    `json:"Landing"`
	Beacons          [][3]float64 `json:"Beacons"`
	LastBeaconAt     time.Time    `json:"LastBeaconAt"`
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
		Flights:           make(map[string]Flight),
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

func (c *UAVContract) RequestPermit(ctx contractapi.TransactionContextInterface, droneId string, permitEffective time.Time, permitExpiry time.Time, vertices [][3]float64, facets [][3]uint64) error {
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
	flightId := fmt.Sprintf("%v", time.Now().UnixNano()) // TODO: Better way to generate a flightId
	_, exists = operator.Flights[flightId]
	if exists {
		return fmt.Errorf("Flight %v already exists", flightId)
	}
	tree := kdtree.Tree{Root: nil, Count: 0}
	flight := Flight{
		DroneId:          droneId,
		PermitEffective:  permitEffective,
		PermitExpiry:     permitExpiry,
		BoundaryVertices: vertices,
		BoundaryFacets:   facets,
		Tree:             &tree,
		Status:           "PENDING",
		Takeoff:          time.Now(),
		Landing:          time.Now(),
		Beacons:          make([][3]float64, 0),
		LastBeaconAt:     time.Now(),
	}
	operator.Flights[flightId] = flight
	operatorJSON, err = json.Marshal(operator)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(operatorId, operatorJSON)
}

func (c *UAVContract) EvaluatePermit(ctx contractapi.TransactionContextInterface, operatorId string, flightId string, decision string) error {
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
	flight, exists := operator.Flights[flightId]
	if !exists {
		return fmt.Errorf("Flight %v does not exist", flightId)
	}
	flight.Status = decision
	// TODO: Supposed to build kdtree here (once for the boundary) but cannot since the tree struct is not allowed in ctx
	if decision == "APPROVED" {
		// vertices2D := isinside.ConvertFloat64To2D(flight.BoundaryVertices)
		// facets2D := isinside.ConvertUint64To2D(flight.BoundaryFacets)
		tree, _, _ := isinside.GenerateKDTreePlus(flight.BoundaryVertices, flight.BoundaryFacets)
		flight.Tree = tree
	}
	operator.Flights[flightId] = flight
	operatorJSON, err = json.Marshal(operator)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(operatorId, operatorJSON)
}

// TODO: Additional business logic checks (e.g., if not already active)
func (c *UAVContract) LogTakeoff(ctx contractapi.TransactionContextInterface, flightId string) error {
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
	flight, exists := operator.Flights[flightId]
	if !exists {
		return fmt.Errorf("Flight %v does not exist", flightId)
	}
	flight.Status = "ACTIVE"
	flight.Takeoff = time.Now()
	operator.Flights[flightId] = flight
	operatorJSON, err = json.Marshal(operator)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(operatorId, operatorJSON)
}

// TODO: Additional business logic checks (e.g., if still within permit period)
func (c *UAVContract) LogBeacons(ctx contractapi.TransactionContextInterface, flightId string, newBeacons [][3]float64) error {
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
	flight, exists := operator.Flights[flightId]
	if !exists {
		return fmt.Errorf("Flight %v does not exist", flightId)
	}
	if flight.Status != "ACTIVE" {
		return fmt.Errorf("Flight %v status is not ACTIVE", flightId)
	}
	// TODO: Find if it's possible to get time based on the invocation time
	// TODO: Add margin of error
	if int(time.Since(flight.LastBeaconAt).Seconds()) != len(flight.Beacons) {
		return fmt.Errorf("Mismatch in flight %v beacons and invocation time ", flightId)
	}
	flight.Beacons = append(flight.Beacons, newBeacons...)
	flight.LastBeaconAt = time.Now()
	operator.Flights[flightId] = flight
	operatorJSON, err = json.Marshal(operator)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(operatorId, operatorJSON)
	if err != nil {
		return err
	}
	// TODO: Specify 1000 as a contract-wide parameter
	if len(flight.Beacons) > 1000 {
		err = c.AnalyzeBeacons(ctx, operator, operatorId, flight, flightId)
	}
	if err != nil {
		return err
	}
	return nil
}

// TODO: Function design can be much improved
func (c *UAVContract) AnalyzeBeacons(ctx contractapi.TransactionContextInterface, operator Operator, operatorId string, flight Flight, flightId string) error {
	inclusions := isinside.GetInclusions(flight.BoundaryVertices, flight.BoundaryFacets, flight.Beacons)
	for _, inclusion := range inclusions {
		if !inclusion {
			return fmt.Errorf("Drone %v has breached its permitted region", flight.DroneId)
		}
	}
	flight.Beacons = nil
	operator.Flights[flightId] = flight
	operatorJSON, err := json.Marshal(operator)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(operatorId, operatorJSON)
	if err != nil {
		return err
	}
	return nil
}

// TODO: Additional business logic checks
func (c *UAVContract) LogLanding(ctx contractapi.TransactionContextInterface, flightId string) error {
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
	flight, exists := operator.Flights[flightId]
	if !exists {
		return fmt.Errorf("Flight %v does not exist", flightId)
	}
	if flight.Status != "ACTIVE" {
		return fmt.Errorf("Flight %v status is not ACTIVE", flightId)
	}
	flight.Status = "PENDING"
	flight.Landing = time.Now()
	operator.Flights[flightId] = flight
	operatorJSON, err = json.Marshal(operator)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(operatorId, operatorJSON)
}

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

func main() {
	carChaincode, err := contractapi.NewChaincode(&UAVContract{})
	if err != nil {
		log.Panicf("Error creating cars chaincode: %v", err)
	}
	err = carChaincode.Start()
	if err != nil {
		log.Panicf("Error starting cars chaincode: %v", err)
	}
}
