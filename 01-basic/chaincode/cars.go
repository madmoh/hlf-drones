package main

import (
	"fmt"
	"encoding/json"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Defining smart contract
type SmartContract struct {
	contractapi.Contract
}

// Defining cars
type Car struct {
	Color string `json:"Color"`
	Doors int `json:"Doors"`
	Make string `json:"Make"`
	Owner string `json:"Owner"`
	Value int `json:"Value"`
	VIN string `json:"VIN"`
}

// Inititalizing smart contract with fixed set of cars
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// Array of initial cars
	cars := []Car{
		{Color: "Silver", Doors: 4, Make: "Lexus", Owner: "CLM", Value: 10000, VIN: "ABCD1234"},
		{Color: "Black", Doors: 2, Make: "Mercedes-Benz", Owner: "MMM", Value: 5000, VIN: "EFGH5678"},
	}
	// Add the car using PutState()
	for _, car := range cars {
		carJSON, err := json.Marshal(car)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(car.VIN, carJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}
	return nil
}

// Add car dynamically
func (s *SmartContract) CreateCar(ctx contractapi.TransactionContextInterface, color string, doors int, make string, owner string, value int, vin string) error {
	// Make sure the car does not already exists
	exists, err := s.CarExists(ctx, vin)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the car %s already exists", vin)
	}
	// Create the car struct and JSON
	car := Car {
		Color: color,
		Doors: doors,
		Make: make,
		Owner: owner,
		Value: value,
		VIN: vin,
	}
	carJSON, err := json.Marshal(car)
	if err != nil {
		return err
	}
	// Add the car using PutState
	return ctx.GetStub().PutState(vin, carJSON)
}

// Query car by VIN
func (s *SmartContract) QueryCar(ctx contractapi.TransactionContextInterface, vin string) (*Car, error) {
	carJSON, err := ctx.GetStub().GetState(vin)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if carJSON == nil {
		return nil, fmt.Errorf("the car %s does not exist", vin)
	}
	var car Car
	err = json.Unmarshal(carJSON, &car)
	if err != nil {
		return nil, err
	}
	return &car, nil
}

func (s *SmartContract) UpdateCar(ctx contractapi.TransactionContextInterface, color string, doors int, make string, owner string, value int, vin string) error {
	// Make sure the car does already exist
	exists, err := s.CarExists(ctx, vin)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the car %s does not exist", vin)
	}
	// Overwrite original car
	car := Car {
		Color: color,
		Doors: doors,
		Make: make,
		Owner: owner,
		Value: value,
		VIN: vin,
	}
	carJSON, err := json.Marshal(car)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(vin, carJSON)
}

func (s *SmartContract) DeleteCar(ctx contractapi.TransactionContextInterface, vin string) error {
	// Make sure the car does already exist
	exists, err := s.CarExists(ctx, vin)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the car %s does not exist", vin)
	}
	return ctx.GetStub().DelState(vin)
}

func (s *SmartContract) CarExists(ctx contractapi.TransactionContextInterface, vin string) (bool, error) {
	carJSON, err := ctx.GetStub().GetState(vin)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return carJSON != nil, nil
}

func (s *SmartContract) TransferCar(ctx contractapi.TransactionContextInterface, vin string, newOwner string) error {
	car, err := s.QueryCar(ctx, vin)
	if err != nil {
		return err
	}
	car.Owner = newOwner
	carJSON, err := json.Marshal(car)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(vin, carJSON)
}

func (s *SmartContract) QueryAllCars(ctx contractapi.TransactionContextInterface) ([]*Car, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	var cars []*Car
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var car Car
		err = json.Unmarshal(queryResponse.Value, &car)
		if err != nil {
			return nil, err
		}
		cars = append(cars, &car)
	}
	return cars, nil
}

func main() {
	carChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating cars chaincode: %v", err)
	}
	err = carChaincode.Start()
	if err != nil {
		log.Panicf("Error starting cars chaincode: %v", err)
	}
}
