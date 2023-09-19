# https://hyperledger.github.io/caliper/v0.5.0/fabric-tutorial/tutorials-fabric-existing

cd <.../cars/>

curl -sSL https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/bootstrap.sh| bash -s -- 2.4.9 1.5.5

cd <.../cars/fabric-samples/test-network>

./network down

./network.sh up -s couchdb

./network.sh createChannel -c carschannel -s couchdb

./network.sh deployCC -c carschannel -ccn cars -ccl go -ccv 1.0 -ccs 1 -ccp "../../chaincode"

cd <.../cars/caliper-workspace>

bun caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/networkconfig.yaml --caliper-benchconfig benchmarks/benchconfig-cars.yaml --caliper-flow-only-test