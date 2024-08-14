'use strict'

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class RequestPermitsWorkload extends WorkloadModuleBase {
	constructor() {
		super();
	}

	async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
		await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
		this.lastOperator = new Array(this.totalWorkers).fill(0)
		this.lastFlight = new Array(this.totalWorkers * this.roundArguments.operatorsPerWorker).fill(0)

		// Create operators
		for (let operatorIndex = 0; operatorIndex < this.roundArguments.operatorsPerWorker; operatorIndex++) {
			const operatorId = `${this.workerIndex}_${operatorIndex}`
			console.log(`Worker ${this.workerIndex}: Adding operator ${operatorId}`)
			const requestAddOperator = {
				contractId: this.roundArguments.contractId,
				contractFunction: 'RecordsSC:AddOperator',
				invokerIdentity: 'User1',
				contractArguments: [operatorId],
				readOnly: false
			}
			await this.sutAdapter.sendRequests(requestAddOperator)
		}
	}

	async submitTransaction() {
		const operatorIndex = this.lastOperator[this.workerIndex]
		const flightIndex = this.lastFlight[this.workerIndex * this.roundArguments.operatorsPerWorker + operatorIndex]
		const operatorId = `${this.workerIndex}_${operatorIndex}`
		const flightId = `${operatorId}_${flightIndex}`

		const requestRequestPermit = {
			contractId: this.roundArguments.contractId,
			contractFunction: 'RecordsSC:RequestPermit',
			invokerIdentity: 'User1',
			contractArguments: [
				operatorId,
				flightId,
				"",
				"2024-01-01T00:00:00Z",
				"2025-01-01T00:00:00Z",
				"[[[-1,0,-1],[-1,0,1],[-1,8,-1],[-1,8,1],[1,0,-1],[1,0,1],[1,8,-1],[1,8,1]]]",
				"[[[0,6,4],[0,2,6],[0,3,2],[0,1,3],[2,7,6],[2,3,7],[4,6,7],[4,7,5],[0,4,5],[0,5,1],[1,5,7],[1,7,3]]]",
				"[true]"
			],
			readOnly: false
		}

		this.lastOperator[this.workerIndex] = (this.lastOperator[this.workerIndex] + 1) % this.roundArguments.operatorsPerWorker
		this.lastFlight[this.workerIndex * this.roundArguments.operatorsPerWorker + operatorIndex] += 1
		console.log(`Worker ${this.workerIndex}: Requesting permit ${flightId}`)
		await this.sutAdapter.sendRequests(requestRequestPermit)


		// console.log(`Worker ${this.workerIndex}: Adding ${this.roundArguments.beaconsPerLog} logs for flightId ${flightId}`)
		// await this.sutAdapter.sendRequests(request)


		// const operatorId = `${this.workerIndex}_${this.lastFlight[this.workerIndex] % this.roundArguments.operatorsPerWorker}`
		// const flightId = `${operatorId}_${this.lastFlight[this.workerIndex]}`
		// const requestRequestPermit = {
		// 	contractId: this.roundArguments.contractId,
		// 	contractFunction: 'RecordsSC:RequestPermit',
		// 	invokerIdentity: 'User1',
		// 	contractArguments: [
		// 		operatorId,
		// 		flightId,
		// 		"",
		// 		"2024-01-01T00:00:00Z",
		// 		"2025-01-01T00:00:00Z",
		// 		"[[[-1,0,-1],[-1,0,1],[-1,8,-1],[-1,8,1],[1,0,-1],[1,0,1],[1,8,-1],[1,8,1]]]",
		// 		"[[[0,6,4],[0,2,6],[0,3,2],[0,1,3],[2,7,6],[2,3,7],[4,6,7],[4,7,5],[0,4,5],[0,5,1],[1,5,7],[1,7,3]]]",
		// 		"[true]"
		// 	],
		// 	readOnly: false
		// }
		// console.log(`Worker ${this.workerIndex}: Requesting permit ${flightId}`)
		// await this.sutAdapter.sendRequests(requestRequestPermit)
		// this.lastFlight[this.workerIndex]++
	}

	// async cleanupWorkloadModule() {
	// 	for (let operatorIndex = 0; operatorIndex < this.roundArguments.operatorsPerWorker; operatorIndex++) {
	// 		const operatorId = `${this.workerIndex}_${operatorIndex}`
	// 		console.log(`Worker ${this.workerIndex}: Deleting operator ${operatorId}`)
	// 		const requestDeleteOperator = {
	// 			contractId: this.roundArguments.contractId,
	// 			contractFunction: 'RecordsSC:DeleteOperator',
	// 			invokerIdentity: 'User1',
	// 			contractArguments: [operatorId],
	// 			readOnly: false
	// 		}
	// 		await this.sutAdapter.sendRequests(requestDeleteOperator)
	// 	}

	// 	for (let i = 0; i < this.lastFlight[this.workerIndex]; i++) {
	// 		const operatorId = `${this.workerIndex}_${i % this.roundArguments.operatorsPerWorker}`
	// 		const flightId = `${operatorId}_${i}`
	// 		console.log(`Worker ${this.workerIndex}: Deleting flight ${flightId}`)
	// 		const requestDeleteFlight = {
	// 			contractId: this.roundArguments.contractId,
	// 			contractFunction: 'RecordsSC:DeleteFlight',
	// 			invokerIdentity: 'User1',
	// 			contractArguments: [flightId],
	// 			readOnly: false
	// 		}
	// 		await this.sutAdapter.sendRequests(requestDeleteFlight)
	// 	}
	// }
}

function createWorkloadModule() {
	return new RequestPermitsWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
