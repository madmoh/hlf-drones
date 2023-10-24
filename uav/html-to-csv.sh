#!/bin/bash

argc=$#
if [ $argc -ne 3 ]; then
	echo "Usage: <test name> <starting run> <final run>"
	echo
	echo "       <test name> must direct to a valid Caliper report at caliper-workspace/reports/<test name>-tps-run.html"
	echo "       <starting run> must be an integer greater than 0"
	echo "       <final run> must be an integer greater than <starting run>"
	exit
fi

test=$1
startrun=$2
finalrun=$3
output="caliper-workspace/reports/$test-rs$startrun-rf$finalrun.csv"

touch $output
echo "Test,TargetSendRate,Success,Fail,SendRate,MaxLatency,MinLatency,AvgLatency,Throughput" > $output

for run in $(seq $startrun $finalrun)
do
	for f in caliper-workspace/reports/$test-*-$run.html
	do
		xq $f -q "#benchmarksummary table tr td" | sed 's/-/,/g;:a;N;$!ba;s/\n/,/g' >> $output
	done
done


