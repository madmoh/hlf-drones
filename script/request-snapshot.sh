#!/bin/bash

argc=$#
if [ $argc -ne 1 ]; then
	echo "Usage: <height>"
	echo
	echo "       <height> must be an integer greater or equal to 0, or -1 for current block height"
	exit
fi

if [ $1 -eq "-1" ]; then
  _dir="$(cd "$(dirname "$0")" && pwd)"
	height=$(${_dir}/get-block-height.sh)
else
  height=$1
fi


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

peer snapshot submitrequest -c abyssar -b $height --peerAddress $CORE_PEER_ADDRESS --tlsRootCertFile $CORE_PEER_TLS_ROOTCERT_FILE


# Environment variables for Org2
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID=Org2MSP
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051
export FABRIC_LOGGING_SPEC=error

peer snapshot submitrequest -c abyssar -b $height --peerAddress $CORE_PEER_ADDRESS --tlsRootCertFile $CORE_PEER_TLS_ROOTCERT_FILE
