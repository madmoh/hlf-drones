#!/bin/bash

# These commands are for setting up the network, channel, and chaincode deployment using the network.sh script in fabric-samples/test-network
# It is assumed all pre-requisites are installed, and the chaincode is developed and ready for deployment

# Execute all commands from test-network directory

cd fabric-samples/test-network/

export PATH=$PATH:${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/


## To become Admin@org1.example.com and go through peer0.org1.example.com
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export FABRIC_LOGGING_SPEC=error

OPERATOR_ID_STARTING=$1
OPERATOR_ID_ENDING=$2
NUM_FLIGHTS=$3
NUM_LOGS=$4
NUM_PROCS=$5

AddOperator() {
	local i=$1
	local OPERATOR_ID="operator_${i}"
	echo "Adding operator ${OPERATOR_ID}"
  peer chaincode invoke \
		--channelID abyssar --name abyssarCC \
		--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
		--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
		--orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
		--ctor '{"function":"RecordsSC:AddOperator","Args":["'${OPERATOR_ID}'"]}'
}
export -f AddOperator

AddCertificate(){
	local i=$1
	local OPERATOR_ID="operator_${i}"
	echo "Adding certificate for ${OPERATOR_ID}"
	peer chaincode invoke \
		--channelID abyssar --name abyssarCC \
		--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
		--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
		--orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
		--ctor '{"function":"RecordsSC:AddCertificate","Args":["'${OPERATOR_ID}'","COMMERCIAL","'$(date --iso-8601=seconds -d "$(date --iso-8601) +1 year")'"]}'
}
export -f AddCertificate

AddDrone(){
	local i=$1
	local OPERATOR_ID="operator_${i}"
	local DRONE_ID="drone_${OPERATOR_ID}"
	echo "Adding drone ${DRONE_ID}"
	peer chaincode invoke \
		--channelID abyssar --name abyssarCC \
		--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
		--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
		--orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
		--ctor '{"function":"RecordsSC:AddDrone","Args":["'${DRONE_ID}'","'$(date --iso-8601=seconds -d "$(date --iso-8601)")'"]}'
}
export -f AddDrone

RequestPermit(){
	local i=$1
	local OPERATOR_ID="operator_${i}"
	local DRONE_ID="drone_${OPERATOR_ID}"
	local FLIGHT_ID="flight_${DRONE_ID}"
	echo "Adding flight ${FLIGHT_ID}"
	peer chaincode invoke \
		--channelID abyssar --name abyssarCC \
		--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
		--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
		--orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
		--ctor '{"function":"RecordsSC:RequestPermit","Args":["'${OPERATOR_ID}'","'${FLIGHT_ID}'","'${DRONE_ID}'","'$(date --iso-8601=seconds -d "$(date --iso-8601)")'","'$(date --iso-8601=seconds -d "$(date --iso-8601) +90 day")'","[[[-1,0,-1],[-1,0,1],[-1,8,-1],[-1,8,1],[1,0,-1],[1,0,1],[1,8,-1],[1,8,1]]]","[[[0,6,4],[0,2,6],[0,3,2],[0,1,3],[2,7,6],[2,3,7],[4,6,7],[4,7,5],[0,4,5],[0,5,1],[1,5,7],[1,7,3]]]","[true]"]}'
}
export -f RequestPermit

EvaluatePermit(){
	local i=$1
	local OPERATOR_ID="operator_${i}"
	local DRONE_ID="drone_${OPERATOR_ID}"
	local FLIGHT_ID="flight_${DRONE_ID}"
	echo "Adding permit for ${FLIGHT_ID}"
	peer chaincode invoke \
		--channelID abyssar --name abyssarCC \
		--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
		--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
		--orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
		--ctor '{"function":"RecordsSC:EvaluatePermit","Args":["'${FLIGHT_ID}'","APPROVED"]}'
}
export -f EvaluatePermit

LogTakeoff(){
	local i=$1
	local OPERATOR_ID="operator_${i}"
	local DRONE_ID="drone_${OPERATOR_ID}"
	local FLIGHT_ID="flight_${DRONE_ID}"
	echo "Taking off ${FLIGHT_ID}"
	peer chaincode invoke \
		--channelID abyssar --name abyssarCC \
		--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
		--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
		--orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
		--ctor '{"function":"FlightsSC:LogTakeoff","Args":["'${OPERATOR_ID}'","'${FLIGHT_ID}'"]}'
}
export -f LogTakeoff

LogBeacons(){
	local i=$1
	local OPERATOR_ID="operator_${i}"
	local DRONE_ID="drone_${OPERATOR_ID}"
	local FLIGHT_ID="flight_${DRONE_ID}"
	echo "Logging beacons for ${FLIGHT_ID}"
	peer chaincode invoke \
		--channelID abyssar --name abyssarCC \
		--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
		--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
		--orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
		--ctor '{"function":"FlightsSC:LogBeacons","Args":["'${OPERATOR_ID}'","'${FLIGHT_ID}'","[[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3]]"]}'
}
export -f LogBeacons

LogLanding(){
	local i=$1
	local OPERATOR_ID="operator_${i}"
	local DRONE_ID="drone_${OPERATOR_ID}"
	local FLIGHT_ID="flight_${DRONE_ID}"
	echo "Landing ${FLIGHT_ID}"
	peer chaincode invoke \
		--channelID abyssar --name abyssarCC \
		--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
		--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
		--orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
		--ctor '{"function":"FlightsSC:LogLanding","Args":["'${OPERATOR_ID}'","'${FLIGHT_ID}'"]}'
}
export -f LogLanding


operators_i=$(seq $OPERATOR_ID_STARTING $OPERATOR_ID_ENDING)
echo "$operators_i" | parallel -j $NUM_PROCS AddOperator
echo "$operators_i" | parallel -j $NUM_PROCS AddCertificate
echo "$operators_i" | parallel -j $NUM_PROCS AddDrone
echo "$operators_i" | parallel -j $NUM_PROCS RequestPermit
sleep 2
echo "$operators_i" | parallel -j $NUM_PROCS EvaluatePermit
sleep 2
echo "$operators_i" | parallel -j $NUM_PROCS LogTakeoff
sleep 2
for k in $(seq $NUM_LOGS)
do
	echo "Logging beacons ${k}"
	echo "$operators_i" | parallel -j $NUM_PROCS LogBeacons
	sleep 2
done
echo "$operators_i" | parallel -j $NUM_PROCS LogLanding
sleep 2
