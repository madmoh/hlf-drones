#!/bin/bash

RESULTS_PATH=$1
_dir="$(cd "$(dirname "$0")" && pwd)"

for i in $(seq 25 25 300)
do
	${_dir}/restart-network.sh
	if [[ $i -gt 0 ]]; then
		${_dir}/send-transactions-parallel.sh 1 $i 1 $i 12
	fi
	sleep 5
	height=$(${_dir}/get-block-height.sh)
	echo -en "$i\t$height\t" >> $RESULTS_PATH
	${_dir}/request-snapshot.sh 0
	sleep 5
	echo -ne "$(docker exec $(docker ps -aqf "name=^peer0.org1.example.com") du -sb /var/hyperledger/production/ledgersData | cut -f1)\t" >> $RESULTS_PATH
	echo -ne "$(docker exec $(docker ps -aqf "name=^peer0.org1.example.com") du -sb /var/hyperledger/production/snapshots/completed/abyssar/$(($height - 1)) | cut -f1)\t" >> $RESULTS_PATH
	echo "" >> $RESULTS_PATH
done

