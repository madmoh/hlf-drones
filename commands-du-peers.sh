#!/bin/bash

for container in $(docker ps --filter "ancestor=hyperledger/fabric-peer" --format "{{.ID}}"); do
  echo "Disk usage for container $container:"
  docker exec $container du -sh /var/hyperledger/production/
done
