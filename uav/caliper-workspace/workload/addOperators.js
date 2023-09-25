'use strict'

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class AddOperatorWorkload extends WorkloadModuleBase {
	constructor() {
		super();
	}

	async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
		await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
		this.I = new Array(this.totalWorkers).fill(0)
	}

	async submitTransaction() {
		const operatorId = `${this.workerIndex}_${this.I[this.workerIndex]}`
		this.I[this.workerIndex]++
		const request = {
			contractId: this.roundArguments.contractId,
			contractFunction: 'RecordsSC:AddOperator',
			invokerIdentity: 'User1',
			contractArguments: [operatorId],
			readOnly: false
		}
		await this.sutAdapter.sendRequests(request)
	}

	async cleanupWorkloadModule() {
		for (let w = 0; w < this.totalWorkers; w++) {
			for (let i = 0; i < this.I[w]; i++) {
				const operatorId = `${w}_${i}`
				console.log(`Worker ${w}: Deleting operator ${operatorId}`)
				const requestDeleteOperator = {
					contractId: this.roundArguments.contractId,
					contractFunction: 'RecordsSC:DeleteOperator',
					invokerIdentity: 'User1',
					contractArguments: [operatorId],
					readOnly: false
				}
				await this.sutAdapter.sendRequests(requestDeleteOperator)
			}
		}
	}
}

function createWorkloadModule() {
	return new AddOperatorWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;