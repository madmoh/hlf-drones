package contract

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

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
}

type Flight struct {
	FlightId           string         `json:"FlightId"`
	OperatorId         string         `json:"OperatorId"`
	DroneId            string         `json:"DroneId"`
	PermitEffective    time.Time      `json:"PermitEffective"`
	PermitExpiry       time.Time      `json:"PermitExpiry"`
	BoundariesVertices [][][3]float64 `json:"BoundariesVertices"` // TODO: Expand to [][][3]
	BoundariesFacets   [][][3]uint64  `json:"BoundariesFacets"`   // TODO: Expand to [][][3]
	BoundariesTypes    []bool         `json:"BoundariesTypes"`    // true: permitted/must be inside, false: restricted/must stay outside
	// Tree             *kdtree.Tree `json:"Tree"`
	Status       string       `json:"Status"` // PENDING, REJECTED, APPROVED, ACTIVE
	Takeoff      time.Time    `json:"Takeoff"`
	Landing      time.Time    `json:"Landing"`
	Beacons      [][3]float64 `json:"Beacons"`
	LastBeaconAt time.Time    `json:"LastBeaconAt"`
}

type Violation struct {
	ViolationId string    `json:"ViolationId"`
	OccuredAt   time.Time `json:"OccuredAt"`
	ReportedAt  time.Time `json:"ReportedAt"`
	Reason      string    `json:"Reason"`
	OperatorId  string    `json:"OperatorId"`
	FlightId    string    `json:"FlightId"`
}

func KeyExists(ctx contractapi.TransactionContextInterface, key string) (bool, error) {
	valueJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return false, fmt.Errorf("failed to read world state. Error: %v", err)
	}
	return valueJSON != nil, nil
}
