#!/bin/bash

argc=$#
if [ $argc -ne 5 ]; then
	echo "Usage: <starting lg(tps)> <stepping lg(tps)> <final lg(tps)> <starting run> <final run>"
	echo
	echo "       <starting lg(tps)> must be an integer greater than 0"
	echo "       <stepping lg(tps)> must be an integer"
	echo "       <final lg(tps)> must be an integer greater than <starting lg(tps)>"
	echo "       <starting run> must be an integer greater than 0"
	echo "       <final run> must be an integer greater than <starting run>"
	exit
fi

declare -a tests=(
	[0]="addOperators"
	[1]="requestPermits"
	[2]="logBeacons"
)
start=$1
step=$2
final=$3
startrun=$4
finalrun=$5

for test in ${tests[@]}
do
	for run in $(seq $startrun $finalrun)
	do
		for lg_tps in $(seq $start $step $final)
		do
			tps=$((2 ** $lg_tps))
			export WORKERS=$(($(nproc) * 4))
			export TPS=$tps
			# export TX_NUMBER=$(( ($TPS * 10 + $WORKERS - 1) / $WORKERS * $WORKERS ))
			export TX_DURATION=30

			echo "$test $run $TPS"

			fabric-samples/test-network/network.sh down &&
			fabric-samples/test-network/network.sh up createChannel -c abyssar &&
			fabric-samples/test-network/network.sh deployCC -c abyssar -ccn abyssarCC -ccl go -ccv 1.0 -ccs 1 -ccp "../../chaincode"

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
done
