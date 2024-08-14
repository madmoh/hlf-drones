#!/bin/bash

RESULTS_PATH=$1

for i in $(seq 225 25 300)
do
	./commands-restart-network.sh
	if [[ $i -gt 0 ]]; then
		./commands-send-transactions-parallel.sh 1 $i 1 $i 12
	fi
	sleep 5
	height=$(./commands-get-height.sh)
	echo -en "$i\t$height\t" >> $RESULTS_PATH
	./commands-request-snapshot.sh 0
	sleep 5
	echo -ne "$(docker exec $(docker ps -aqf "name=^peer0.org1.example.com") du -sb /var/hyperledger/production/ledgersData | cut -f1)\t" >> $RESULTS_PATH
	echo -ne "$(docker exec $(docker ps -aqf "name=^peer0.org1.example.com") du -sb /var/hyperledger/production/snapshots/completed/abyssar/$(($height - 1)) | cut -f1)\t" >> $RESULTS_PATH
	echo "" >> $RESULTS_PATH
done

