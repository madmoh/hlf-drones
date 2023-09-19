# These commands are for setting up the network, channel, and chaincode deployment using the network.sh script in fabric-samples/test-network
# It is assumed all pre-requisites are installed, and the chaincode is developed and ready for deployment

# Execute all commands from test-network directory

export PATH=$PATH:${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/

./network.sh down

./network.sh up -ca -s couchdb -verbose

./network.sh createChannel -ca -c carschannel -s couchdb -verbose

./network.sh deployCC -c carschannel -ccn cars -ccl go -ccv 1.0 -ccs 1 -ccp "../chaincode" # -ccep "OutOf(2,'Org1MSP.peer','Org2MSP.peer')"

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

    ## As peer0.org1.example.com
    peer chaincode invoke \
        --channelID carschannel --name cars \
        --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
        --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
        --orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
        --ctor '{"function":"InitLedger","Args":[]}'

    ## As peer0.org1.example.com
    peer chaincode query --channelID carschannel --name cars --ctor '{"Args":["QueryAllCars"]}'

    ## As peer0.org1.example.com
    peer chaincode invoke \
        --channelID carschannel --name cars \
        --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
        --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
        --orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
        --ctor '{"function":"TransferCar","Args":["ABCD1234","MMM"]}'

    ## As peer0.org2.example.com
    peer chaincode query --channelID carschannel --name cars --ctor '{"Args":["QueryCar","ABCD1234"]}'

    ## As peer0.org2.example.com
    peer chaincode invoke \
        --channelID carschannel --name cars \
        --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
        --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
        --orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
        --ctor '{"function":"CreateCar","Args":["Gold","4","Nissan","YAY","15000","IJKL4321"]}'

    ## As peer0.org2.example.com
    peer chaincode query --channelID carschannel --name cars --ctor '{"Args":["CarExists","IJKL4321"]}'

    ## As peer0.org2.example.com
    peer chaincode query --channelID carschannel --name cars --ctor '{"Args":["QueryAllCars"]}'

# To test automatically (Caliper):
