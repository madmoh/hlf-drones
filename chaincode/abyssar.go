package main

import (
	"log"

	"madmoh/hlf-uav/contract"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	abyssar, err := contractapi.NewChaincode(
		&contract.RecordsSC{},
		&contract.FlightsSC{},
		&contract.ViolationsSC{},
	)
	if err != nil {
		log.Panicf("Error creating Abyssar chaincode: %v", err)
	}
	err = abyssar.Start()
	if err != nil {
		log.Panicf("Error starting Abyssar chaincode: %v", err)
	}
}
