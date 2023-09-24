# https://hyperledger.github.io/caliper/v0.5.0/fabric-tutorial/tutorials-fabric-existing

cd <.../uav/>

# Do the following or simply copy the fabric-samples folder (containing bin + config + test-network)
curl -sSL https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/bootstrap.sh| bash -s -- 2.4.9 1.5.5

cd <.../uav/fabric-samples/test-network/>

export PATH=$PATH:${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/

./network.sh down

./network.sh up

./network.sh createChannel -c uavchannel

./network.sh deployCC -c uavchannel -ccn uav -ccl go -ccv 1.0 -ccs 1 -ccp "../../chaincode"

cd <.../cars/caliper-workspace>

bun caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/networkconfig.yaml --caliper-benchconfig benchmarks/benchconfig-uav.yaml --caliper-flow-only-test