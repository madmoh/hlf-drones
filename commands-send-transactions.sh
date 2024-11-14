#!/bin/bash

argc=$#
if [ $argc -ne 3 ]; then
	echo "Usage: <test name> <tps> <duration>"
	echo
	echo "       <test name> must direct to a valid benchmark at caliper-workspace/benchmarks/benchconfig-<test name>.yaml.tmpl"
	echo "       <tps> must be an integer greater than 0"
	echo "       <duration> must be an integer greater than 0"
	exit
fi

test=$1
start=$2
duration=$3

echo "$start"
export WORKERS=$(($(nproc) * 4))
export TPS=$start
export TX_DURATION=$duration
if [ $test = "requestPermits" ]; then
	export OPERATORS_PER_WORKER=$(($TPS*15/$WORKERS+1))
fi
if [ $test = "logBeacons" ]; then
	export OPERATORS_PER_WORKER=$(($TPS*120/$WORKERS+1))
fi

envsubst < caliper-workspace/benchmarks/benchconfig-${test}.yaml.tmpl > caliper-workspace/benchmarks/temp.yaml

bun --cwd "${PWD}/caliper-workspace" \
	caliper launch manager \
	--caliper-workspace ./ \
	--caliper-networkconfig networks/networkconfig.yaml \
	--caliper-benchconfig benchmarks/temp.yaml \
	--caliper-flow-only-test
