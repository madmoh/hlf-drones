# These commands are for setting up the network, channel, and chaincode deployment using the network.sh script in fabric-samples/test-network
# It is assumed all pre-requisites are installed, and the chaincode is developed and ready for deployment

# Execute all commands from test-network directory

cd <.../fabric-samples/test-network/>

export PATH=$PATH:${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/

./network.sh down

./network.sh up

./network.sh createChannel -c abyssar

./network.sh deployCC -c abyssar -ccn abyssarCC -ccl go -ccv 1.0 -ccs 1 -ccp "../../chaincode"

# To test manually:
## To become Admin@org1.example.com and go through peer0.org1.example.com
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

## To become Admin@org2.example.com and go through peer0.org2.example.com
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051


OPERATOR_ID=operator00
DRONE_ID=drone00
FLIGHT_ID=flight00

peer chaincode invoke \
	--channelID abyssar --name abyssarCC \
	--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
	--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
	--orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
	--ctor '{"function":"RecordsSC:AddOperator","Args":["'${OPERATOR_ID}'"]}'

peer chaincode query --channelID abyssar --name abyssarCC --ctor '{"Args":["RecordsSC:GetOperator","'${OPERATOR_ID}'"]}'

peer chaincode invoke \
	--channelID abyssar --name abyssarCC \
	--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
	--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
	--orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
	--ctor '{"function":"RecordsSC:AddDrone","Args":["'${DRONE_ID}'","2023-09-28T00:00:00Z"]}'

peer chaincode query --channelID abyssar --name abyssarCC --ctor '{"Args":["RecordsSC:GetDrone","'${DRONE_ID}'"]}'

peer chaincode invoke \
	--channelID abyssar --name abyssarCC \
	--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
	--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
	--orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
	--ctor '{"function":"RecordsSC:AddCertificate","Args":["'${OPERATOR_ID}'","COMMERCIAL","2023-10-28T00:00:00Z"]}'

peer chaincode query --channelID abyssar --name abyssarCC --ctor '{"Args":["RecordsSC:GetOperator","'${OPERATOR_ID}'"]}'

# [
# 	[-1,0,-1],
# 	[-1,0,1],
# 	[-1,8,-1],
# 	[-1,8,1],
# 	[1,0,-1],
# 	[1,0,1],
# 	[1,8,-1],
# 	[1,8,1],
# ]
# [[-1,0,-1],[-1,0,1],[-1,8,-1],[-1,8,1],[1,0,-1],[1,0,1],[1,8,-1],[1,8,1]]
# [
# 	[0,6,4],
# 	[0,2,6],
# 	[0,3,2],
# 	[0,1,3],
# 	[2,7,6],
# 	[2,3,7],
# 	[4,6,7],
# 	[4,7,5],
# 	[0,4,5],
# 	[0,5,1],
# 	[1,5,7],
# 	[1,7,3],
# ]
# [[0,6,4],[0,2,6],[0,3,2],[0,1,3],[2,7,6],[2,3,7],[4,6,7],[4,7,5],[0,4,5],[0,5,1],[1,5,7],[1,7,3]]

peer chaincode invoke \
	--channelID abyssar --name abyssarCC \
	--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
	--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
	--orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
	--ctor '{"function":"RecordsSC:RequestPermit","Args":["'${OPERATOR_ID}'","'${FLIGHT_ID}'","'${DRONE_ID}'","2023-09-21T00:00:00Z","2023-11-28T00:00:00Z","[[[-1,0,-1],[-1,0,1],[-1,8,-1],[-1,8,1],[1,0,-1],[1,0,1],[1,8,-1],[1,8,1]]]","[[[0,6,4],[0,2,6],[0,3,2],[0,1,3],[2,7,6],[2,3,7],[4,6,7],[4,7,5],[0,4,5],[0,5,1],[1,5,7],[1,7,3]]]","[true]"]}'

peer chaincode query --channelID abyssar --name abyssarCC --ctor '{"Args":["RecordsSC:GetOperator","'${OPERATOR_ID}'"]}'
peer chaincode query --channelID abyssar --name abyssarCC --ctor '{"Args":["FlightsSC:GetFlight","'${FLIGHT_ID}'"]}'

peer chaincode invoke \
	--channelID abyssar --name abyssarCC \
	--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
	--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
	--orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
	--ctor '{"function":"RecordsSC:EvaluatePermit","Args":["'${FLIGHT_ID}'","APPROVED"]}'

peer chaincode query --channelID abyssar --name abyssarCC --ctor '{"Args":["FlightsSC:GetFlight","'${FLIGHT_ID}'"]}'

peer chaincode invoke \
	--channelID abyssar --name abyssarCC \
	--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
	--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
	--orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
	--ctor '{"function":"FlightsSC:LogTakeoff","Args":["'${OPERATOR_ID}'","'${FLIGHT_ID}'"]}'

peer chaincode query --channelID abyssar --name abyssarCC --ctor '{"Args":["FlightsSC:GetFlight","'${FLIGHT_ID}'"]}'

peer chaincode invoke \
	--channelID abyssar --name abyssarCC \
	--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
	--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
	--orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
	--ctor '{"function":"FlightsSC:LogBeacons","Args":["'${OPERATOR_ID}'","'${FLIGHT_ID}'","[[50.1,10.2,20.3],[50.1,10.2,20.3],[50.1,10.2,20.3],[50.1,10.2,20.3],[50.1,10.2,20.3],[50.1,10.2,20.3],[50.1,10.2,20.3],[50.1,10.2,20.3],[50.1,10.2,20.3],[50.1,10.2,20.3],[50.1,10.2,20.3],[50.1,10.2,20.3],[50.1,10.2,20.3],[50.1,10.2,20.3],[50.1,10.2,20.3]]"]}'
peer chaincode invoke \
	--channelID abyssar --name abyssarCC \
	--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
	--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
	--orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
	--ctor '{"function":"FlightsSC:LogBeacons","Args":["'${OPERATOR_ID}'","'${FLIGHT_ID}'","[[100.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3]]"]}'

peer chaincode query --channelID abyssar --name abyssarCC --ctor '{"Args":["RecordsSC:GetOperator","'${OPERATOR_ID}'"]}'
peer chaincode query --channelID abyssar --name abyssarCC --ctor '{"Args":["RecordsSC:GetDrone","'${DRONE_ID}'"]}'
peer chaincode query --channelID abyssar --name abyssarCC --ctor '{"Args":["FlightsSC:GetFlight","'${FLIGHT_ID}'"]}'

VIOLATION_ID=$(peer chaincode query --channelID abyssar --name abyssarCC --ctor '{"Args":["RecordsSC:GetOperator","'${OPERATOR_ID}'"]}' | jq '.ViolationIds[-1]' -r)
peer chaincode query --channelID abyssar --name abyssarCC --ctor '{"Args":["ViolationsSC:GetViolation","'${VIOLATION_ID}'"]}'

peer chaincode invoke \
	--channelID abyssar --name abyssarCC \
	--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
	--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
	--orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
	--ctor '{"function":"UpdateOperatorStatus","Args":["'${OPERATOR_ID}'","TEMP_BAN_APPEAL"]}'

peer chaincode query --channelID abyssar --name abyssarCC --ctor '{"Args":["RecordsSC:GetOperator","'${OPERATOR_ID}'"]}'
