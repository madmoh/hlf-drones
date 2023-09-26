package contract

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type ViolationsSC struct {
	contractapi.Contract
}

func (c *ViolationsSC) AddViolation(ctx contractapi.TransactionContextInterface, operatorId string, flightId string, beaconIndex int, reason string) (string, error) {
	exists, err := KeyExists(ctx, flightId)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", fmt.Errorf("Flight %v does not exist", flightId)
	}
	flightJSON, err := ctx.GetStub().GetState(flightId)
	if err != nil {
		return "", fmt.Errorf("failed to read from state. Error: %v", err)
	}
	if flightJSON == nil {
		return "", fmt.Errorf("Flight %v does not exist", flightId)
	}
	var flight Flight
	err = json.Unmarshal(flightJSON, &flight)
	if err != nil {
		return "", err
	}

	occuredAt := flight.LastBeaconAt.Add(time.Second * time.Duration(-len(flight.Beacons)+beaconIndex+1))
	violationId := fmt.Sprintf("%v%v%v", operatorId, flightId, occuredAt.Unix())
	exists, err = KeyExists(ctx, violationId)
	if err != nil {
		return "", err
	}
	if exists {
		return "", fmt.Errorf("Violation %v already exists", violationId)
	}
	txTimestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "", fmt.Errorf("failed to get timestamp for receipt: %v", err)
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
		return "", err
	}
	err = ctx.GetStub().PutState(violationId, violationJSON)
	if err != nil {
		return violationId, err
	}

	exists, err = KeyExists(ctx, operatorId)
	if err != nil {
		return violationId, err
	}
	if !exists {
		return violationId, fmt.Errorf("operator %v does not exist", operatorId)
	}
	operatorJSON, err := ctx.GetStub().GetState(operatorId)
	if err != nil {
		return violationId, fmt.Errorf("failed to read from state. Error: %v", err)
	}
	if operatorJSON == nil {
		return violationId, fmt.Errorf("operator %v does not exist", operatorId)
	}
	var operator Operator
	err = json.Unmarshal(operatorJSON, &operator)
	if err != nil {
		return violationId, err
	}
	operator.ViolationIds = append(operator.ViolationIds, violationId)
	operatorJSON, err = json.Marshal(operator)
	if err != nil {
		return violationId, err
	}
	err = ctx.GetStub().PutState(operatorId, operatorJSON)
	if err != nil {
		return violationId, err
	}
	return violationId, nil
}

func (c *ViolationsSC) GetViolation(ctx contractapi.TransactionContextInterface, key string) (*Violation, error) {
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

func (c *ViolationsSC) DeleteViolation(ctx contractapi.TransactionContextInterface, violationId string) error {
	exists, err := KeyExists(ctx, violationId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Violation %v does not exist", violationId)
	}
	return ctx.GetStub().DelState(violationId)
}
