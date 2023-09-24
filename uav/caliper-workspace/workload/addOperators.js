'use strict'

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class MyWorkload extends WorkloadModuleBase {
	constructor() {
		super();
	}

	// async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
	// 	await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);

	// 	for (let i = 0; i < this.roundArguments.cars; i++) {
	// 		const vin = `${this.workerIndex}_${i}`
	// 		console.log(`Worker ${this.workerIndex}: Creating car ${vin}`)
	// 		const request = {
	// 			contractId: this.roundArguments.contractId,
	// 			contractFunction: 'CreateCar',
	// 			invokerIdentity: 'User1',
	// 			contractArguments: ['plink', '4', 'tayoto', 'mahmadto', '1000000', vin],
	// 			readOnly: false
	// 		}
	// 		await this.sutAdapter.sendRequests(request)
	// 	}
	// }

	async submitTransaction() {
		const operatorId = `${this.workerIndex}_${Date.now()}`
		const request = {
			contractId: this.roundArguments.contractId,
			contractFunction: 'AddOperator',
			invokerIdentity: 'User1',
			contractArguments: [operatorId],
			readOnly: false
		}
		// console.log(`Worker ${this.workerIndex}: Calling AddOperator(${operatorId})`)
		await this.sutAdapter.sendRequests(request)
	}

	// async cleanupWorkloadModule() {
	// 	for (let i = 0; i < this.roundArguments.cars; i++) {
	// 		const vin = `${this.workerIndex}_${i}`
	// 		console.log(`Worker ${this.workerIndex}: Deleting car ${vin}`)
	// 		const request = {
	// 			contractId: this.roundArguments.contractId,
	// 			contractFunction: 'DeleteCar',
	// 			invokerIdentity: 'User1',
	// 			contractArguments: [vin],
	// 			readOnly: false
	// 		}
	// 		await this.sutAdapter.sendRequests(request)
	// 	}
	// }
}

function createWorkloadModule() {
	return new MyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;