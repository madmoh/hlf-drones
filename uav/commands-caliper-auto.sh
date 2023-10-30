#!/bin/bash

argc=$#
if [ $argc -ne 6 ]; then
	echo "Usage: <test name> <starting tps> <stepping tps> <final tps> <starting run> <final run>"
	echo
	echo "       <test name> must direct to a valid benchmark at caliper-workspace/benchmarks/benchconfig-<test name>.yaml.tmpl"
	echo "       <starting tps> must be an integer greater than 0"
	echo "       <stepping tps> must be an integer"
	echo "       <final tps> must be an integer greater than <starting tps>"
	echo "       <starting run> must be an integer greater than 0"
	echo "       <final run> must be an integer greater than <starting run>"
	exit
fi

test=$1
start=$2
step=$3
final=$4
startrun=$5
finalrun=$6

for run in $(seq $startrun $finalrun)
do
	for tps in $(seq $start $step $final)
	do
		echo "$tps"
		fabric-samples/test-network/network.sh down &&
		fabric-samples/test-network/network.sh up createChannel -c abyssar &&
		fabric-samples/test-network/network.sh deployCC -c abyssar -ccn abyssarCC -ccl go -ccv 1.0 -ccs 1 -ccp "../../chaincode"

		export WORKERS=$(($(nproc) * 4))
		export TPS=$tps
		# export TX_NUMBER=$(( ($TPS * 10 + $WORKERS - 1) / $WORKERS * $WORKERS ))
		export TX_DURATION=30

		envsubst < caliper-workspace/benchmarks/benchconfig-${test}.yaml.tmpl > caliper-workspace/benchmarks/temp.yaml

		bun --cwd "${PWD}/caliper-workspace" \
			caliper launch manager \
			--caliper-workspace ./ \
			--caliper-networkconfig networks/networkconfig.yaml \
			--caliper-benchconfig benchmarks/temp.yaml \
			--caliper-flow-only-test \
			--caliper-report-path reports/${test}-${tps}-${run}.html \
			--caliper-report-precision 6 

	done
done
