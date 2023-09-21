package main

import (
	"encoding/json"
	"fmt"
	"log"
	"madmoh/hlf-uav/isinside"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"golang.org/x/exp/slices"
	"gonum.org/v1/gonum/spatial/kdtree"
)

type UAVContract struct {
	contractapi.Contract
}

// TODO: OperatorId, DroneId shouldn't be inside their respective structs

type Operator struct {
	OperatorId        string    `json:"OperatorId"`      // TODO: Unique accross operators
	CertificateTier   string    `json:"CertificateTier"` // TODO: Could also be an enum and possibly house multiple certificates
	CertificateExpiry time.Time `json:"CertificateExpiry"`
	DroneIds          []string  `json:"DroneIds"`
	FlightIds         []string  `json:"FlightIds"`
	ViolationIds      []string  `json:"ViolationIds"`
	Status            string    `json:"Status"`
}

type Drone struct {
	DroneId string    `json:"DroneId"`
	Expiry  time.Time `json:"Expiry"`
	// RemoteId string `json:"RemoteId"` // TODO: Unique accross drones
}

type Flight struct {
	FlightId         string       `json:"FlightId"`
	OperatorId       string       `json:"OperatorId"`
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
	// OperatorId       string    `json:"OperatorId"`
}

type Violation struct {
	ViolationId string    `json:"ViolationId"`
	OccuredAt   time.Time `json:"OccuredAt"`
	ReportedAt  time.Time `json:"ReportedAt"`
	Reason      string    `json:"Reason"`
	OperatorId  string    `json:"OperatorId"`
	FlightId    string    `json:"FlightId"`
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

// TODO: Make some link between the drone and operator
func (c *UAVContract) AddDrone(ctx contractapi.TransactionContextInterface, droneId string, expiry time.Time) error {
	// TODO: Similar checks to AddOperator
	exists, err := c.KeyExists(ctx, droneId)
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

// TODO: More business logic checks (caller identity)
func (c *UAVContract) AddCertificate(ctx contractapi.TransactionContextInterface, operatorId string, tier string, expiry time.Time) error {
	exists, err := c.KeyExists(ctx, operatorId)
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

func (c *UAVContract) RequestPermit(ctx contractapi.TransactionContextInterface, flightId string, droneId string, permitEffective time.Time, permitExpiry time.Time, vertices [][3]float64, facets [][3]uint64) error {
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
	exists = slices.Contains(operator.FlightIds, flightId)
	if exists {
		return fmt.Errorf("Operator already has flight %v", flightId)
	}
	if exists {
		return fmt.Errorf("Flight %v already exists", flightId)
	}
	tree := kdtree.Tree{Root: nil, Count: 0}
	flight := Flight{
		FlightId:         flightId,
		OperatorId:       operatorId,
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
	flightJSON, err := json.Marshal(flight)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(flightId, flightJSON)
	if err != nil {
		return err
	}
	operator.FlightIds = append(operator.FlightIds, flightId)
	operatorJSON, err = json.Marshal(operator)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(operatorId, operatorJSON)
}

func (c *UAVContract) EvaluatePermit(ctx contractapi.TransactionContextInterface, flightId string, decision string) error {
	exists, err := c.KeyExists(ctx, flightId)
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
	if decision == "APPROVED" {
		// vertices2D := isinside.ConvertFloat64To2D(flight.BoundaryVertices)
		// facets2D := isinside.ConvertUint64To2D(flight.BoundaryFacets)
		tree, _, _ := isinside.GenerateKDTreePlus(flight.BoundaryVertices, flight.BoundaryFacets)
		flight.Tree = tree
	}
	flightJSON, err = json.Marshal(flight)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(flightId, flightJSON)
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
	if !slices.Contains(operator.FlightIds, flightId) {
		return fmt.Errorf("flight %v is not associated with operator %v", flightId, operatorId)
	}

	exists, err = c.KeyExists(ctx, flightId)
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

	flight.Status = "ACTIVE"
	txTimestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get timestamp for receipt: %v", err)
	}
	flight.Takeoff = txTimestamp.AsTime()
	flightJSON, err = json.Marshal(flight)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(flightId, flightJSON)
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
	if !slices.Contains(operator.FlightIds, flightId) {
		return fmt.Errorf("flight %v is not associated with operator %v", flightId, operatorId)
	}

	exists, err = c.KeyExists(ctx, flightId)
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
		return fmt.Errorf("flight %v does not exist", flightId)
	}
	var flight Flight
	err = json.Unmarshal(flightJSON, &flight)
	if err != nil {
		return err
	}

	if flight.Status != "ACTIVE" {
		return fmt.Errorf("Flight %v status is not ACTIVE", flightId)
	}
	// TODO: Find if it's possible to get time based on the invocation time
	// TODO: Add margin of error
	if int(time.Since(flight.LastBeaconAt).Seconds()) != len(flight.Beacons) {
		return fmt.Errorf("mismatch in flight %v beacons and invocation time ", flightId)
	}
	flight.Beacons = append(flight.Beacons, newBeacons...)
	txTimestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get timestamp for receipt: %v", err)
	}
	flight.LastBeaconAt = txTimestamp.AsTime()
	flightJSON, err = json.Marshal(flight)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(flightId, flightJSON)
	if err != nil {
		return err
	}
	// TODO: Specify 1000 as a contract-wide parameter
	if len(flight.Beacons) > 1000 {
		err = c.AnalyzeBeacons(ctx, flightId)
	}
	if err != nil {
		return err
	}
	return nil
}

// TODO: Function design can be much improved
func (c *UAVContract) AnalyzeBeacons(ctx contractapi.TransactionContextInterface, flightId string) error {
	exists, err := c.KeyExists(ctx, flightId)
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

	inclusions := isinside.GetInclusions(flight.BoundaryVertices, flight.BoundaryFacets, flight.Beacons)
	for index, inclusion := range inclusions {
		if !inclusion {
			c.AddViolation(ctx, flight.OperatorId, flightId, index, "ESCAPE_PERMITTED_REGION")
			// TODO: Decide report first violation, last violation, all violations?
		}
	}
	flight.Beacons = nil
	flightJSON, err = json.Marshal(flight)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(flightId, flightJSON)
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
	if !slices.Contains(operator.FlightIds, flightId) {
		return fmt.Errorf("Flight %v is not associated with operator %v", flightId, operatorId)
	}

	exists, err = c.KeyExists(ctx, flightId)
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

	flight.Status = "PENDING"
	txTimestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get timestamp for receipt: %v", err)
	}
	flight.Landing = txTimestamp.AsTime()
	flightJSON, err = json.Marshal(flight)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(flightId, flightJSON)
}

func (c *UAVContract) AddViolation(ctx contractapi.TransactionContextInterface, operatorId string, flightId string, beaconIndex int, reason string) error {
	exists, err := c.KeyExists(ctx, flightId)
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

	occuredAt := flight.LastBeaconAt.Add(time.Second * time.Duration(-len(flight.Beacons)+beaconIndex+1))
	violationId := fmt.Sprintf("%v%v%v", operatorId, flightId, occuredAt)
	exists, err = c.KeyExists(ctx, violationId)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("Violation %v already exists", violationId)
	}
	txTimestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get timestamp for receipt: %v", err)
	}
	violation := Violation{
		OccuredAt:  occuredAt,
		ReportedAt: txTimestamp.AsTime(),
		Reason:     reason,
		OperatorId: operatorId,
		FlightId:   flightId,
	}
	violationJSON, err := json.Marshal(violation)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(violationId, violationJSON)
}

func (c *UAVContract) UpdateOperatorStatus(ctx contractapi.TransactionContextInterface, operatorId string, status string) error {
	exists, err := c.KeyExists(ctx, operatorId)
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

func (c *UAVContract) KeyExists(ctx contractapi.TransactionContextInterface, key string) (bool, error) {
	valueJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return false, fmt.Errorf("failed to read world state. Error: %v", err)
	}
	return valueJSON != nil, nil
}

func (c *UAVContract) GetOperator(ctx contractapi.TransactionContextInterface, key string) (*Operator, error) {
	objectJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("failed to read world state. Error: %v", err)
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

func (c *UAVContract) GetDrone(ctx contractapi.TransactionContextInterface, key string) (*Drone, error) {
	objectJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("failed to read world state. Error: %v", err)
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

// func (c *UAVContract) GetFlight(ctx contractapi.TransactionContextInterface, key string) (*Flight, error) {
// 	objectJSON, err := ctx.GetStub().GetState(key)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read world state. Error: %v", err)
// 	}
// 	if objectJSON == nil {
// 		return nil, fmt.Errorf("object %s does not exist", key)
// 	}
// 	var object Flight
// 	err = json.Unmarshal(objectJSON, &object)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &object, nil
// }

func (c *UAVContract) GetViolation(ctx contractapi.TransactionContextInterface, key string) (*Violation, error) {
	objectJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("failed to read world state. Error: %v", err)
	}
	if objectJSON == nil {
		return nil, fmt.Errorf("object %s does not exist", key)
	}
	var object Violation
	err = json.Unmarshal(objectJSON, &object)
	if err != nil {
		return nil, err
	}
	return &object, nil
}

func (c *UAVContract) GetClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {
	return ctx.GetClientIdentity().GetID()
}

func main() {
	uavChaincode, err := contractapi.NewChaincode(&UAVContract{})
	if err != nil {
		log.Panicf("Error creating UAV chaincode: %v", err)
	}
	err = uavChaincode.Start()
	if err != nil {
		log.Panicf("Error starting UAV chaincode: %v", err)
	}
}
