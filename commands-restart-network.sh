#!/bin/bash

fabric-samples/test-network/network.sh down &&
fabric-samples/test-network/network.sh up createChannel -c abyssar &&
# (cd fabric-samples/test-network/addOrg3 ; ./addOrg3.sh up -c abyssar) &&
fabric-samples/test-network/network.sh deployCC -c abyssar -ccn abyssarCC -ccl go -ccv 1.0 -ccs 1 -ccp "../../chaincode"