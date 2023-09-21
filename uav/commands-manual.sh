# These commands are for setting up the network, channel, and chaincode deployment using the network.sh script in fabric-samples/test-network
# It is assumed all pre-requisites are installed, and the chaincode is developed and ready for deployment

# Execute all commands from test-network directory

export PATH=$PATH:${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/

./network.sh down

./network.sh up -verbose

./network.sh createChannel -c uavchannel -verbose

./network.sh deployCC -c uavchannel -ccn uav -ccl go -ccv 1.0 -ccs 1 -ccp "../../chaincode" -verbose

# To test manually:
## To become peer0.org1.example.com
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

## To become peer0.org2.example.com
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051


OPERATOR_ID=$(peer chaincode query --channelID uavchannel --name uav --ctor '{"Args":["GetClientIdentity"]}')
DRONE_ID=${OPERATOR_ID}-drone
FLIGHT_ID=${DRONE_ID}-flight01

peer chaincode invoke \
	--channelID uavchannel --name uav \
	--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
	--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
	--orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
	--ctor '{"function":"AddOperator","Args":["'${OPERATOR_ID}'"]}'

peer chaincode query --channelID uavchannel --name uav --ctor '{"Args":["GetOperator","'${OPERATOR_ID}'"]}'

peer chaincode invoke \
	--channelID uavchannel --name uav \
	--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
	--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
	--orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
	--ctor '{"function":"AddDrone","Args":["'${OPERATOR_ID}'-drone","2023-09-28T00:00:00Z"]}'

peer chaincode query --channelID uavchannel --name uav --ctor '{"Args":["GetDrone","'${OPERATOR_ID}'-drone"]}'

peer chaincode invoke \
	--channelID uavchannel --name uav \
	--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
	--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
	--orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
	--ctor '{"function":"AddCertificate","Args":["'${OPERATOR_ID}'","COMMERCIAL","2023-10-28T00:00:00Z"]}'

peer chaincode query --channelID uavchannel --name uav --ctor '{"Args":["GetOperator","'${OPERATOR_ID}'"]}'

# [
# 	[-1.0,0.0,1.0],
# 	[1.0,0.0,1.0],
# 	[1.0,0.0,-1.0],
# 	[-1.0,0.0,-1.0],
# 	[-1.0,5.0,1.0],
# 	[1.0,5.0,1.0],
# 	[1.0,5.0,-1.0],
# 	[-1.0,5.0,-1.0],
# ]
# [
# 	[0,1,3],
# 	[1,2,3],
# 	[2,5,6],
# 	[1,5,2],
# 	[0,4,1],
# 	[1,4,5],
# 	[0,7,4],
# 	[0,3,7],
# 	[3,6,7],
# 	[2,6,3],
# 	[4,6,5],
# 	[4,7,6],
# ]
# [[0,1,3],[1,2,3],[2,5,6],[1,5,2],[0,4,1],[1,4,5],[0,7,4],[0,3,7],[3,6,7],[2,6,3],[4,6,5],[4,7,6]]

peer chaincode invoke \
	--channelID uavchannel --name uav \
	--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
	--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
	--orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
	--ctor '{"function":"RequestPermit","Args":["'${OPERATOR_ID}'-drone-flight01","'${OPERATOR_ID}'-drone","2023-09-21T00:00:00Z","2023-11-28T00:00:00Z","[[-1.0,0.0,1.0],[1.0,0.0,1.0],[1.0,0.0,-1.0],[-1.0,0.0,-1.0],[-1.0,5.0,1.0],[1.0,5.0,1.0],[1.0,5.0,-1.0],[-1.0,5.0,-1.0]]","[[0,1,3],[1,2,3],[2,5,6],[1,5,2],[0,4,1],[1,4,5],[0,7,4],[0,3,7],[3,6,7],[2,6,3],[4,6,5],[4,7,6]]"]}'

peer chaincode query --channelID uavchannel --name uav --ctor '{"Args":["GetOperator","'${OPERATOR_ID}'"]}'
peer chaincode query --channelID uavchannel --name uav --ctor '{"Args":["GetFlight","'${OPERATOR_ID}'-drone-flight01"]}'

peer chaincode invoke \
	--channelID uavchannel --name uav \
	--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
	--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
	--orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
	--ctor '{"function":"EvaluatePermit","Args":["'${OPERATOR_ID}'-drone-flight01","APPROVED"]}'

peer chaincode invoke \
	--channelID uavchannel --name uav \
	--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
	--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
	--orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
	--ctor '{"function":"LogTakeoff","Args":["'${OPERATOR_ID}'-drone-flight01"]}'

# [
# 	[0.1,0.2,0.3],
# 	[0.1,0.2,0.3],
# 	[0.1,0.2,0.3],
# 	[0.1,0.2,0.3],
# 	[0.1,0.2,0.3],
# 	[0.1,0.2,0.3],
# 	[0.1,0.2,0.3],
# 	[0.1,0.2,0.3],
# 	[0.1,0.2,0.3],
# ]
# [[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3]]

peer chaincode invoke \
	--channelID uavchannel --name uav \
	--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
	--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
	--orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
	--ctor '{"function":"LogBeacons","Args":["'${OPERATOR_ID}'-drone-flight01","[[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[0.1,0.2,0.3],[50.1,10.2,20.3]]"]}'
peer chaincode invoke \
	--channelID uavchannel --name uav \
	--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
	--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
	--orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
	--ctor '{"function":"LogBeacons","Args":["'${OPERATOR_ID}'-drone-flight01","[[0.1,0.2,0.3]]"]}'

peer chaincode query --channelID uavchannel --name uav --ctor '{"Args":["GetOperator","'${OPERATOR_ID}'"]}'
peer chaincode query --channelID uavchannel --name uav --ctor '{"Args":["GetDrone","'${DRONE_ID}'"]}'
peer chaincode query --channelID uavchannel --name uav --ctor '{"Args":["GetFlight","'${FLIGHT_ID}'"]}'
peer chaincode query --channelID uavchannel --name uav --ctor '{"Args":["GetViolation","'${NO_IDEA}'"]}'

