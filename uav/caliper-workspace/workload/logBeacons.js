'use strict'

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class LogBeaconsWorkload extends WorkloadModuleBase {
	constructor() {
		super();
	}

	async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
		await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
		this.lastOperator = new Array(this.totalWorkers).fill(0)
		this.lastFlight = new Array(this.totalWorkers * this.roundArguments.operatorsPerWorker).fill(0)

		// Create operators and flights
		for (let operatorIndex = 0; operatorIndex < this.roundArguments.operatorsPerWorker; operatorIndex++) {
			const operatorId = `${this.workerIndex}_${operatorIndex}`
			// console.log(`Worker ${this.workerIndex}: Adding operator ${operatorId}`)
			const requestAddOperator = {
				contractId: this.roundArguments.chaincodeName,
				contractFunction: 'RecordsSC:AddOperator',
				invokerIdentity: 'User1',
				contractArguments: [operatorId],
				readOnly: false
			}
			await this.sutAdapter.sendRequests(requestAddOperator)

			for (let flightIndex = 0; flightIndex < this.roundArguments.flightsPerOperator; flightIndex++) {
				const flightId = `${this.workerIndex}_${operatorIndex}_${flightIndex}`
				// console.log(`Worker ${this.workerIndex}: Adding flight ${flightId}`)
				const requestRequestPermit = {
					contractId: this.roundArguments.chaincodeName,
					contractFunction: 'RecordsSC:RequestPermit',
					invokerIdentity: 'User1',
					contractArguments: [
						operatorId,
						flightId, 
						"", 
						"2023-09-21T00:00:00Z", 
						"2023-11-28T00:00:00Z", 
						"[[[-1,0,-1],[-1,0,1],[-1,8,-1],[-1,8,1],[1,0,-1],[1,0,1],[1,8,-1],[1,8,1]]]",
						"[[[0,6,4],[0,2,6],[0,3,2],[0,1,3],[2,7,6],[2,3,7],[4,6,7],[4,7,5],[0,4,5],[0,5,1],[1,5,7],[1,7,3]]]",
						"[true]"
					],
					readOnly: false
				}
				await this.sutAdapter.sendRequests(requestRequestPermit)
				const requestEvaluatePermit = {
					contractId: this.roundArguments.chaincodeName,
					contractFunction: 'RecordsSC:EvaluatePermit',
					invokerIdentity: 'User1',
					contractArguments: [flightId, "APPROVED"],
					readOnly: false
				}
				await this.sutAdapter.sendRequests(requestEvaluatePermit)
				const requestLogTakeoff = {
					contractId: this.roundArguments.chaincodeName,
					contractFunction: 'FlightsSC:LogTakeoff',
					invokerIdentity: 'User1',
					contractArguments: [operatorId, flightId],
					readOnly: false
				}
				await this.sutAdapter.sendRequests(requestLogTakeoff)
			}
		}
	}

	async submitTransaction() {
		const operatorIndex = this.lastOperator[this.workerIndex]
		const flightIndex = this.lastFlight[this.workerIndex * this.roundArguments.operatorsPerWorker + operatorIndex]
		const operatorId = `${this.workerIndex}_${operatorIndex}`
		const flightId = `${operatorId}_${flightIndex}`
		const request = {
			contractId: this.roundArguments.chaincodeName,
			contractFunction: 'FlightsSC:LogBeacons',
			invokerIdentity: 'User1',
			contractArguments: [
				operatorId,
				flightId,
				"[" + `[${Math.random() * 2 - 1},${Math.random() * 8},${Math.random() * 2 - 1}],`.repeat(this.roundArguments.beaconsPerLog - 1) + `[${Math.random() * 2 - 1},${Math.random() * 8},${Math.random() * 2 - 1}]]`
			]
		}
		this.lastOperator[this.workerIndex] = (this.lastOperator[this.workerIndex] + 1) % this.roundArguments.operatorsPerWorker
		this.lastFlight[this.workerIndex * this.roundArguments.operatorsPerWorker + operatorIndex] = (this.lastFlight[this.workerIndex * this.roundArguments.operatorsPerWorker + operatorIndex] + 1) % this.roundArguments.flightsPerOperator
		// console.log(`Worker ${this.workerIndex}: Adding logs for flightId ${flightId}`)
		await this.sutAdapter.sendRequests(request)
	}

	// async cleanupWorkloadModule() {
	// 	for (let operatorIndex = 0; operatorIndex < this.roundArguments.operatorsPerWorker; operatorIndex++) {
	// 		const operatorId = `${this.workerIndex}_${operatorIndex}`
	// 		console.log(`Worker ${this.workerIndex}: Deleting operator ${operatorId}`)
	// 		const requestDeleteOperator = {
	// 			contractId: this.roundArguments.chaincodeName,
	// 			contractFunction: 'RecordsSC:DeleteOperator',
	// 			invokerIdentity: 'User1',
	// 			contractArguments: [operatorId],
	// 			readOnly: false
	// 		}
	// 		await this.sutAdapter.sendRequests(requestDeleteOperator)

	// 		for (let flightIndex = 0; flightIndex < this.roundArguments.flightsPerOperator; flightIndex++) {
	// 			const flightId = `${this.workerIndex}_${operatorIndex}_${flightIndex}`
	// 			console.log(`Worker ${this.workerIndex}: Deleting flight ${flightId}`)
	// 			const requestDeleteFlight = {
	// 				contractId: this.roundArguments.chaincodeName,
	// 				contractFunction: 'RecordsSC:DeleteFlight',
	// 				invokerIdentity: 'User1',
	// 				contractArguments: [flightId],
	// 				readOnly: false
	// 			}
	// 			await this.sutAdapter.sendRequests(requestDeleteFlight)
	// 		}
	// 	}
	// }
}

function createWorkloadModule() {
	return new LogBeaconsWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
