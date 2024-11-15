# UAV Drones Compliance Monitoring using Hyperledger Fabric
This repository contains chaincode that performs real-time monitoring for drone geofencing compliance, and performance testing scripts that mainly use Hyperledger Caliper.

## What are these files? (repo structure)
### Chaincode files
- `chaincode/abyssar.go`: The main function, defines the chaincode wrapping three smart contracts (RecordsSC, FlightsSC, and ViolationsSC).
- `chaincode/contract/*`: Defines the smart contracts (attributes and functions) and unit test scripts.
- `chaincode/isinside/*`: Defines the algorithm for checking if a point in 3D space lies inside a closed 3D object.

### Hyperledger Caliper files
- `caliper-workspace/networks/*`: Configuration file to describe the type of the blockchain, the chaincode name, and organization architecture.
- `caliper-workspace/benchmarks/*`: Templates for configuration files to tell Caliper the details of the benchmark (number of parallel processes, duration of test, send rate type, send rate, workload path, and workload arguments).
- `caliper-workspace/workload/*`: JavaScript functions that define what exactly the peers must execute (prepare and send transactions).

### Bash and Python scripts
#### Helper scripts
- `script/restart-network.sh`: Stop the network. Start the network. Deploy the chaincode.
- `script/get-transactions-count.sh`: Retrieve the total number of transactions in the ledger. Very slow implementation (luckily none of the test scripts have to run this). Assumes the network is already running.
- `script/get-block-height.sh`: Retrieve the block height of the ledger. Assumes the network is already running.
- `script/request-snapshot.sh`: Create a snapshot of the ledger based on a specific block height (or last block height). Assumes the network is already running.
- `script/get-peer-disk-usage.sh`: Retrieve the size of the ledger. Assumes the network is already running.
- `script/compile-caliper-reports.py`: Combines all Caliper reports of the same benchmark into a .csv table. Reports of matching transaction send rate are merged as their average and standard deviation.

#### Test scripts (default scenario)
- `script/integration-test-manual.sh`: A series of commands that helps in confirming the correctness of the chaincode based on the default scenario (as described in the paper). The commands are meant to be ran manually. Running the entire script at once will certainly result in an error.
- `script/send-transactions-parallel.sh`: Send transactions based on the default scenario (as described in the paper), for a specific number of operators, flights, logs. The transactions are issued in parallel. Assumes the network is already running.
- `script/test-snapshots.sh`: Run the default scenario (as described in the paper) with iteratively increasing number of operators and logs while taking snapshots. Outputs the block height, ledger size, and snapshot size.

#### Test scripts (Caliper benchmarks)
- `script/send-transactions.sh`: Continuously send transactions based on one of Caliper benchmarks, at a fixed rate for a specific duration. Assumes the network is already running.
- `script/caliper-auto.sh`: Repeatedly restarts the network to run a sequence of Caliper benchmarks while linearly increasing the transaction send rate. Outputs Caliper reports.
- `script/caliper-log.sh`: Repeatedly restarts the network to run a sequence of Caliper benchmarks while exponentially increasing the transaction send rate. Outputs Caliper reports.

## How can I run these files? (instructions to run scripts)
1. Install pre-requisite packages. 
	1. Recommended packages in (https://hyperledger-fabric.readthedocs.io/en/latest/prereqs.html).
	1. JavaScript runtime (Bun) and npm. Relying purely on Node.js also works (must replace bun command with node in all scripts) but that comes with performance degradation guarantees.
	1. As the scripts are running, some command may rely on uninstalled packages. Install the missing packages on the fly.
	1. Optional quality of life packages: lazydocker, k9s.
1. Populate `fabric-samples` by following the installation steps in (https://hyperledger-fabric.readthedocs.io/en/latest/install.html). This code is tested with Fabric binaries version 2.5.9, CA binaries version 1.5.12, and the test network configuration provided in [fabric-samples](https://github.com/hyperledger/fabric-samples). Important subfolders are `bin`, `config`, and `test-network`.
1. Install node packages by moving to `caliper-workspace` and running `bun install`.
1. Run the script of liking (while at the root of the repo). Make sure to supply the expected arguments.
1. Recommended: If planning to stress test, close all unnecessary processes, including the desktop environment.