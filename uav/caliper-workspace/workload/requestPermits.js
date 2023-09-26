'use strict'

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class LogBeaconsWorkload extends WorkloadModuleBase {
	constructor() {
		super();
	}

	async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
		await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
		this.I = new Array(this.totalWorkers).fill(0)
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
		const operatorId = `${this.workerIndex}_${this.I[this.workerIndex] % this.roundArguments.operatorsPerWorker}`
		const flightId = `${operatorId}_${this.I[this.workerIndex]}`
		const requestRequestPermit = {
			contractId: this.roundArguments.contractId,
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
		this.I[this.workerIndex]++
	}

	async cleanupWorkloadModule() {
		for (let operatorIndex = 0; operatorIndex < this.roundArguments.operatorsPerWorker; operatorIndex++) {
			const operatorId = `${this.workerIndex}_${operatorIndex}`
			console.log(`Worker ${this.workerIndex}: Deleting operator ${operatorId}`)
			const requestDeleteOperator = {
				contractId: this.roundArguments.contractId,
				contractFunction: 'RecordsSC:DeleteOperator',
				invokerIdentity: 'User1',
				contractArguments: [operatorId],
				readOnly: false
			}
			await this.sutAdapter.sendRequests(requestDeleteOperator)
		}

		for (let i = 0; i < this.I[this.workerIndex]; i++) {
			const operatorId = `${this.workerIndex}_${this.I[this.workerIndex] % this.roundArguments.operatorsPerWorker}`
			const flightId = `${operatorId}_${this.I[this.workerIndex]}`
			console.log(`Worker ${this.workerIndex}: Deleting flight ${flightId}`)
				const requestDeleteFlight = {
					contractId: this.roundArguments.contractId,
					contractFunction: 'RecordsSC:DeleteFlight',
					invokerIdentity: 'User1',
					contractArguments: [flightId],
					readOnly: false
				}
				await this.sutAdapter.sendRequests(requestDeleteFlight)
		}
	}
}

function createWorkloadModule() {
	return new LogBeaconsWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;