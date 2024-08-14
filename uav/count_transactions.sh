#!/bin/bash

cd fabric-samples/test-network

export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/

# Environment variables for Org1
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export FABRIC_LOGGING_SPEC=error

CHANNEL_NAME="abyssar"

# Get the blockchain height
HEIGHT=$(peer channel getinfo -c $CHANNEL_NAME | cut -c 18- | jq .height)

# Initialize transaction count
TOTAL_TRANSACTIONS=0

# Loop through each block
for (( i=0; i<$HEIGHT; i++ ))
do
    # Get block data
    BLOCK_DATA=$(peer channel fetch $i -c $CHANNEL_NAME --orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem")

    # Extract number of transactions in the block
    TX_COUNT=$(configtxgen -inspectBlock "${CHANNEL_NAME}_${i}.block" | jq .data.data | jq length)
		# echo "${i}: ${TX_COUNT}"

    # Add to total transactions
    TOTAL_TRANSACTIONS=$(($TOTAL_TRANSACTIONS + $TX_COUNT))

    # Clean up
    rm "${CHANNEL_NAME}_${i}.block"
done

echo "Total number of transactions: $TOTAL_TRANSACTIONS"