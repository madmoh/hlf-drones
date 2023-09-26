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

func TestAddOperatorCatchGetError(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf(""))
	err = recordsSC.AddOperator(transactionContext, "test-operator")
	require.EqualError(t, err, "")
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

func TestGetOperatorSuccess(t *testing.T) {
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

func TestGetOperatorCatchGetError(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf(""))
	operator, err := recordsSC.GetOperator(transactionContext, "test-operator")
	require.EqualError(t, err, "")
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

func TestAddDroneSuccess(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	chaincodeStub.GetStateReturns(nil, nil)
	err = recordsSC.AddDrone(transactionContext, "test-drone", time.Now())
	require.NoError(t, err)
}

func TestAddDroneCatchGetError(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf(""))
	err = recordsSC.AddDrone(transactionContext, "test-drone", time.Now())
	require.EqualError(t, err, "")
}

func TestAddDroneCatchDuplicateError(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	chaincodeStub.GetStateReturns([]byte{}, nil)
	err = recordsSC.AddDrone(transactionContext, "test-drone", time.Now())
	require.EqualError(t, err, "Drone test-drone already exists")
}

func TestGetDroneSuccess(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	droneExpected := &contract.Drone{}
	bytes, err := json.Marshal(droneExpected)
	require.NoError(t, err)
	chaincodeStub.GetStateReturns(bytes, nil)
	drone, err := recordsSC.GetDrone(transactionContext, "test-drone")
	require.NoError(t, err)
	require.Equal(t, droneExpected, drone)
}

func TestGetDroneCatchGetError(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf(""))
	drone, err := recordsSC.GetDrone(transactionContext, "test-drone")
	require.EqualError(t, err, "")
	require.Nil(t, drone)
}

func TestGetDroneCatchEmptyError(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	chaincodeStub.GetStateReturns(nil, nil)
	_, err = recordsSC.GetDrone(transactionContext, "test-drone")
	require.EqualError(t, err, "object test-drone does not exist")
}

func TestAddCertificateSuccess(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	operatorExpected := &contract.Operator{}
	bytes, err := json.Marshal(operatorExpected)
	require.NoError(t, err)
	chaincodeStub.GetStateReturns(bytes, nil)
	err = recordsSC.AddCertificate(transactionContext, "test-operator", "CERTIFIED", time.Now())
	require.NoError(t, err)
}

func TestAddCertificateCatchGetOperatorError(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf(""))
	err = recordsSC.AddCertificate(transactionContext, "test-operator", "CERTIFIED", time.Now())
	require.EqualError(t, err, "")
}

func TestAddCertificateCatchEmptyOperatorError(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	chaincodeStub.GetStateReturns(nil, nil)
	err = recordsSC.AddCertificate(transactionContext, "test-operator", "CERTIFIED", time.Now())
	require.EqualError(t, err, "Operator test-operator does not exist")
}

func TestRequestPermitCatchGetOperatorError(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	chaincodeStub.GetStateReturnsOnCall(0, nil, fmt.Errorf(""))
	err = recordsSC.RequestPermit(transactionContext, "test-operator", "test-flight", "test-drone", time.Now(), time.Now(), [][][3]float64{}, [][][3]uint64{}, []bool{})
	require.EqualError(t, err, "")
}

func TestRequestPermitCatchEmptyOperatorError(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	chaincodeStub.GetStateReturnsOnCall(0, nil, nil)
	err = recordsSC.RequestPermit(transactionContext, "test-operator", "test-flight", "test-drone", time.Now(), time.Now(), [][][3]float64{}, [][][3]uint64{}, []bool{})
	require.EqualError(t, err, "operator test-operator does not exist")
}

func TestRequestPermitCatchCorruptOperatorError(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	chaincodeStub.GetStateReturnsOnCall(0, []byte{}, nil)
	err = recordsSC.RequestPermit(transactionContext, "test-operator", "test-flight", "test-drone", time.Now(), time.Now(), [][][3]float64{}, [][][3]uint64{}, []bool{})
	require.EqualError(t, err, "failed to unmarshal. Error: unexpected end of JSON input")
}

func TestRequestPermitCatchDuplicateFlightError(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	operator := contract.Operator{OperatorId: "test-operator", FlightIds: []string{"test-flight"}}
	operatorJSON, err := json.Marshal(operator)
	require.NoError(t, err)
	chaincodeStub.GetStateReturnsOnCall(0, operatorJSON, nil)
	err = recordsSC.RequestPermit(transactionContext, "test-operator", "test-flight", "test-drone", time.Now(), time.Now(), [][][3]float64{}, [][][3]uint64{}, []bool{})
	require.EqualError(t, err, "operator already has flight test-flight")
}

func TestRequestPermitCatchPutFlightError(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	operator := contract.Operator{OperatorId: "test-operator", FlightIds: []string{}}
	operatorJSON, err := json.Marshal(operator)
	require.NoError(t, err)
	chaincodeStub.GetStateReturnsOnCall(0, operatorJSON, nil)
	chaincodeStub.PutStateReturnsOnCall(0, fmt.Errorf(""))
	err = recordsSC.RequestPermit(transactionContext, "test-operator", "test-flight", "test-drone", time.Now(), time.Now(), [][][3]float64{}, [][][3]uint64{}, []bool{})
	require.EqualError(t, err, "")
}

func TestRequestPermitSuccess(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	recordsSC := contract.RecordsSC{}
	err := error(nil)

	operator := contract.Operator{OperatorId: "test-operator", FlightIds: []string{}}
	operatorJSON, err := json.Marshal(operator)
	require.NoError(t, err)
	chaincodeStub.GetStateReturnsOnCall(0, operatorJSON, nil)
	chaincodeStub.PutStateReturnsOnCall(0, nil)
	err = recordsSC.RequestPermit(transactionContext, "test-operator", "test-flight", "test-drone", time.Now(), time.Now(), [][][3]float64{}, [][][3]uint64{}, []bool{})
	require.NoError(t, err)
}
