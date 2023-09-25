'use strict'

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class LogBeaconsWorkload extends WorkloadModuleBase {
	constructor() {
		super();
	}

	async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
		await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
		// Create operators and flights
		for (let operatorIndex = 0; operatorIndex < this.roundArguments.operatorsCount; operatorIndex++) {
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

			for (let flightIndex = 0; flightIndex < this.roundArguments.flightsPerOperator; flightIndex++) {
				const flightId = `${this.workerIndex}_${operatorIndex}_${flightIndex}`
				console.log(`Worker ${this.workerIndex}: Adding flight ${flightId}`)
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
						"[[-1,0,-1],[-1,0,1],[-1,8,-1],[-1,8,1],[1,0,-1],[1,0,1],[1,8,-1],[1,8,1]]",
						"[[0,6,4],[0,2,6],[0,3,2],[0,1,3],[2,7,6],[2,3,7],[4,6,7],[4,7,5],[0,4,5],[0,5,1],[1,5,7],[1,7,3]]"
					],
					readOnly: false
				}
				await this.sutAdapter.sendRequests(requestRequestPermit)
				const requestLogTakeoff = {
					contractId: this.roundArguments.contractId,
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
		const operatorIndex = Math.floor(Math.random() * this.roundArguments.operatorsCount)
		const flightIndex = Math.floor(Math.random() * this.roundArguments.flightsPerOperator)
		const operatorId = `${this.workerIndex}_${operatorIndex}`
		const flightId = `${operatorId}_${flightIndex}`
		const request = {
			contractId: this.roundArguments.contractId,
			contractFunction: 'FlightsSC:LogBeacons',
			invokerIdentity: 'User1',
			contractArguments: [
				operatorId,
				flightId,
				"[" + `[${Math.random() * 2 - 1},${Math.random() * 8},${Math.random() * 2 - 1}],`.repeat(this.roundArguments.beaconsPerLog - 1) + `[${Math.random() * 2 - 1},${Math.random() * 8},${Math.random() * 2 - 1}]]`
			]
		}
		await this.sutAdapter.sendRequests(request)
	}

	async cleanupWorkloadModule() {
		for (let operatorIndex = 0; operatorIndex < this.roundArguments.operatorsCount; operatorIndex++) {
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

			for (let flightIndex = 0; flightIndex < this.roundArguments.flightsPerOperator; flightIndex++) {
				const flightId = `${this.workerIndex}_${operatorIndex}_${flightIndex}`
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
}

function createWorkloadModule() {
	return new LogBeaconsWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;