package contract_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"madmoh/hlf-uav/contract"
	"madmoh/hlf-uav/contract/mocks"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/require"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o mocks/transaction.go -fake-name TransactionContext . transactionContext
type transactionContext interface {
	contractapi.TransactionContextInterface
}

//counterfeiter:generate -o mocks/chaincodestub.go -fake-name ChaincodeStub . chaincodeStub
type chaincodeStub interface {
	shim.ChaincodeStubInterface
}

//counterfeiter:generate -o mocks/statequeryiterator.go -fake-name StateQueryIterator . stateQueryIterator
type stateQueryIterator interface {
	shim.StateQueryIteratorInterface
}

//counterfeiter:generate -o mocks/clientidentity.go -fake-name ClientIdentity . clientIdentity
type clientIdentity interface {
	cid.ClientIdentity
}

func TestAddOperatorSuccess(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	chaincodeStub.GetStateReturns(nil, nil)
	err = recordsSC.AddOperator(transactionContext, "test-operator")
	require.NoError(t, err)
}

func TestAddOperatorCatchStateError(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = recordsSC.AddOperator(transactionContext, "test-operator")
	require.EqualError(t, err, "failed to read world state. Error: unable to retrieve asset")
}

func TestAddOperatorCatchDuplicateError(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	chaincodeStub.GetStateReturns([]byte{}, nil)
	err = recordsSC.AddOperator(transactionContext, "test-operator")
	require.EqualError(t, err, "Operator test-operator already exists")
}

func TestGetOperator(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	operatorExpected := &contract.Operator{}
	bytes, err := json.Marshal(operatorExpected)
	require.NoError(t, err)
	chaincodeStub.GetStateReturns(bytes, nil)
	operator, err := recordsSC.GetOperator(transactionContext, "test-operator")
	require.NoError(t, err)
	require.Equal(t, operatorExpected, operator)
}

func TestGetOperatorCatchStateError(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	operator, err := recordsSC.GetOperator(transactionContext, "test-operator")
	require.EqualError(t, err, "failed to read world state. Error: unable to retrieve asset")
	require.Nil(t, operator)
}

func TestGetOperatorCatchEmptyError(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	chaincodeStub.GetStateReturns(nil, nil)
	_, err = recordsSC.GetOperator(transactionContext, "test-operator")
	require.EqualError(t, err, "object test-operator does not exist")
}

func TestAddDrone(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	// Case: Key doesn't exist, no error in retrieving the state
	chaincodeStub.GetStateReturns(nil, nil)
	err = recordsSC.AddDrone(transactionContext, "test-drone", time.Now())
	require.NoError(t, err)

	// Case: Key doesn't exist, error in retrieving the state
	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = recordsSC.AddDrone(transactionContext, "test-drone", time.Now())
	require.EqualError(t, err, "failed to read world state. Error: unable to retrieve asset")

	// Case: Key exists, no error in retrieving the state
	chaincodeStub.GetStateReturns([]byte{}, nil)
	err = recordsSC.AddDrone(transactionContext, "test-drone", time.Now())
	require.EqualError(t, err, "Drone test-drone already exists")
}

func TestGetDrone(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	// Case: Key doesn't exist, no error in retrieving the state
	chaincodeStub.GetStateReturns(nil, nil)
	_, err = recordsSC.GetDrone(transactionContext, "test-drone")
	require.EqualError(t, err, "object test-drone does not exist")

	// Case: Key doesn't exist, error in retrieving the state
	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	drone, err := recordsSC.GetDrone(transactionContext, "test-drone")
	require.EqualError(t, err, "failed to read world state. Error: unable to retrieve asset")
	require.Nil(t, drone)

	// Case: Key exists, no error in retrieving the state
	droneExpected := &contract.Drone{}
	bytes, err := json.Marshal(droneExpected)
	require.NoError(t, err)
	chaincodeStub.GetStateReturns(bytes, nil)
	drone, err = recordsSC.GetDrone(transactionContext, "test-drone")
	require.NoError(t, err)
	require.Equal(t, droneExpected, drone)
}

func TestAddCertificate(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	// Case: Key doesn't exist, no error in retrieving the state
	chaincodeStub.GetStateReturns(nil, nil)
	err = recordsSC.AddCertificate(transactionContext, "test-operator", "CERTIFIED", time.Now())
	require.EqualError(t, err, "Operator test-operator does not exist")

	// Case: Key doesn't exist, error in retrieving the state
	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = recordsSC.AddCertificate(transactionContext, "test-operator", "CERTIFIED", time.Now())
	require.EqualError(t, err, "failed to read world state. Error: unable to retrieve asset")

	// Case: Key exists, no error in retrieving the state
	operatorExpected := &contract.Operator{}
	bytes, err := json.Marshal(operatorExpected)
	require.NoError(t, err)
	chaincodeStub.GetStateReturns(bytes, nil)
	err = recordsSC.AddCertificate(transactionContext, "test-operator", "CERTIFIED", time.Now())
	require.NoError(t, err)
}

func TestRequestPermit(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)
	clientIdentity := &mocks.ClientIdentity{
		GetIDStub: func() (string, error) { return "test-operator", nil },
	}
	transactionContext.GetClientIdentityReturns(clientIdentity)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	// Case: Fail on checking KeyExists(operator) (err != nil)
	chaincodeStub.GetStateReturnsOnCall(0, nil, fmt.Errorf("cannot get client identity"))
	err = recordsSC.RequestPermit(transactionContext, "test-operator", "test-flight", "test-drone", time.Now(), time.Now(), [][][3]float64{}, [][][3]uint64{}, []bool{})
	require.EqualError(t, err, "failed to read world state. Error: cannot get client identity")

	// Case: Fail on checking KeyExists(opereator) (!exists)
	chaincodeStub.GetStateReturnsOnCall(1, nil, nil)
	err = recordsSC.RequestPermit(transactionContext, "test-operator", "test-flight", "test-drone", time.Now(), time.Now(), [][][3]float64{}, [][][3]uint64{}, []bool{})
	require.EqualError(t, err, "operator (test-operator) does not exist")

	// Case: Fail to read operator from state
	chaincodeStub.GetStateReturnsOnCall(2, []byte{}, nil)
	chaincodeStub.GetStateReturnsOnCall(3, nil, fmt.Errorf("failed to read from state"))
	err = recordsSC.RequestPermit(transactionContext, "test-operator", "test-flight", "test-drone", time.Now(), time.Now(), [][][3]float64{}, [][][3]uint64{}, []bool{})
	require.EqualError(t, err, "failed to read from state")

	// Case: Fail to get operator from state (empty)
	chaincodeStub.GetStateReturnsOnCall(4, []byte{}, nil)
	chaincodeStub.GetStateReturnsOnCall(5, nil, nil)
	err = recordsSC.RequestPermit(transactionContext, "test-operator", "test-flight", "test-drone", time.Now(), time.Now(), [][][3]float64{}, [][][3]uint64{}, []bool{})
	require.EqualError(t, err, "operator test-operator is empty")

	// Case: Fail to get operator from state (corrupt)
	chaincodeStub.GetStateReturnsOnCall(6, []byte{}, nil)
	chaincodeStub.GetStateReturnsOnCall(7, []byte{}, nil)
	err = recordsSC.RequestPermit(transactionContext, "test-operator", "test-flight", "test-drone", time.Now(), time.Now(), [][][3]float64{}, [][][3]uint64{}, []bool{})
	require.EqualError(t, err, "failed to unmarshal. Error: unexpected end of JSON input")

	// Case: Fail on checking flightId uniqueness
	operator := contract.Operator{OperatorId: "test-operator", FlightIds: []string{"test-flight"}}
	operatorJSON, err := json.Marshal(operator)
	require.NoError(t, err)
	chaincodeStub.GetStateReturnsOnCall(8, []byte{}, nil)
	chaincodeStub.GetStateReturnsOnCall(9, operatorJSON, nil)
	err = recordsSC.RequestPermit(transactionContext, "test-operator", "test-flight", "test-drone", time.Now(), time.Now(), [][][3]float64{}, [][][3]uint64{}, []bool{})
	require.EqualError(t, err, "operator already has flight test-flight")

	// Case: Fail on adding flight
	operator.FlightIds = []string{}
	operatorJSON, err = json.Marshal(operator)
	require.NoError(t, err)
	chaincodeStub.GetStateReturnsOnCall(10, []byte{}, nil)
	chaincodeStub.GetStateReturnsOnCall(11, operatorJSON, nil)
	chaincodeStub.PutStateReturnsOnCall(0, fmt.Errorf(""))
	err = recordsSC.RequestPermit(transactionContext, "test-operator", "test-flight", "test-drone", time.Now(), time.Now(), [][][3]float64{}, [][][3]uint64{}, []bool{})
	require.EqualError(t, err, "failed to write to world state")

	// Case: Success
	operator.FlightIds = []string{}
	operatorJSON, err = json.Marshal(operator)
	require.NoError(t, err)
	chaincodeStub.GetStateReturnsOnCall(12, []byte{}, nil)
	chaincodeStub.GetStateReturnsOnCall(13, operatorJSON, nil)
	chaincodeStub.PutStateReturnsOnCall(1, nil)
	err = recordsSC.RequestPermit(transactionContext, "test-operator", "test-flight", "test-drone", time.Now(), time.Now(), [][][3]float64{}, [][][3]uint64{}, []bool{})
	require.NoError(t, err)
}
