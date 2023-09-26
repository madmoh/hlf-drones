package contract

import (
	"encoding/json"
	"fmt"
	"madmoh/hlf-uav/isinside"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"golang.org/x/exp/slices"
)

type FlightsSC struct {
	contractapi.Contract
}

// TODO: Additional business logic checks (e.g., if not already active)
func (c *FlightsSC) LogTakeoff(ctx contractapi.TransactionContextInterface, operatorId string, flightId string) error {
	// operatorId, err := ctx.GetClientIdentity().GetID()
	// if err != nil {
	// 	return err
	// }
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
	if !slices.Contains(operator.FlightIds, flightId) {
		return fmt.Errorf("flight %v is not associated with operator %v", flightId, operatorId)
	}

	exists, err = KeyExists(ctx, flightId)
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
func (c *FlightsSC) LogBeacons(ctx contractapi.TransactionContextInterface, operatorId string, flightId string, newBeacons [][3]float64) error {
	// operatorId, err := ctx.GetClientIdentity().GetID()
	// if err != nil {
	// 	return err
	// }
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
	if !slices.Contains(operator.FlightIds, flightId) {
		return fmt.Errorf("flight %v is not associated with operator %v", flightId, operatorId)
	}

	exists, err = KeyExists(ctx, flightId)
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
	// TODO: Enable after testing
	// TODO: Add margin of error to secondsPassed != len(newBeacons)
	// var secondsPassed int
	// if !flight.LastBeaconAt.IsZero() {
	// 	secondsPassed = int(time.Since(flight.LastBeaconAt).Seconds())
	// } else {
	// 	secondsPassed = int(time.Since(flight.Takeoff).Seconds())
	// }
	// if secondsPassed != len(newBeacons) {
	// 	return fmt.Errorf("mismatch in flight %v beacons and invocation time, tried to add %v beacons, should have added %v beacons", flightId, len(newBeacons), secondsPassed)
	// }
	// TODO: Fix bug where analysis ignores last submitted beacons
	flight.Beacons = append(flight.Beacons, newBeacons...)
	txTimestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return fmt.Errorf("failed to get timestamp for receipt: %v", err)
	}
	flight.LastBeaconAt = txTimestamp.AsTime()
	// TODO: Specify threshold as a contract-wide parameter
	if len(flight.Beacons) > 120 {
		AnalyzeBeacons(ctx, &flight)
	}
	flightJSON, err = json.Marshal(flight)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(flightId, flightJSON)
	if err != nil {
		return err
	}
	return nil
}

// TODO: Function design can be much improved
func AnalyzeBeacons(ctx contractapi.TransactionContextInterface, flight *Flight) {
	boundariesCount := len(flight.BoundariesTypes)
	inclusions := make([][]bool, boundariesCount)
	for i := 0; i < boundariesCount; i++ {
		inclusions[i] = isinside.GetInclusions(flight.BoundariesVertices[i], flight.BoundariesFacets[i], flight.Beacons)
	}
	beaconsCount := len(flight.Beacons)
	compliance := make([]int, beaconsCount) // 0 good, -1 entered, 1 escaped
	for b := 0; b < beaconsCount; b++ {
		for i := 0; i < boundariesCount; i++ {
			if flight.BoundariesTypes[i] && !inclusions[i][b] {
				compliance[b] = 1
				break
			} else if !flight.BoundariesTypes[i] && inclusions[i][b] {
				compliance[b] = -1
				break
			}
		}
	}
	for i, comply := range compliance {
		if comply != 0 {
			reason := ""
			if comply == 1 {
				reason = "ESCAPE_PERMITTED"
			} else if comply == -1 {
				reason = "ENTER_RESTRICTED"
			}
			violationsSC := ViolationsSC{}
			violationsSC.AddViolation(ctx, flight.OperatorId, flight.FlightId, i, reason)
		}
	}
	flight.Beacons = make([][3]float64, 0)
}

// TODO: Additional business logic checks
func (c *FlightsSC) LogLanding(ctx contractapi.TransactionContextInterface, operatorId string, flightId string) error {
	// operatorId, err := ctx.GetClientIdentity().GetID()
	// if err != nil {
	// 	return err
	// }
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
	if !slices.Contains(operator.FlightIds, flightId) {
		return fmt.Errorf("Flight %v is not associated with operator %v", flightId, operatorId)
	}

	exists, err = KeyExists(ctx, flightId)
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

func (c *FlightsSC) GetFlight(ctx contractapi.TransactionContextInterface, key string) (*Flight, error) {
	objectJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("failed to read world state. Error: %v", err)
	}
	if objectJSON == nil {
		return nil, fmt.Errorf("object %s does not exist", key)
	}
	var object Flight
	err = json.Unmarshal(objectJSON, &object)
	if err != nil {
		return nil, err
	}
	return &object, nil
}
