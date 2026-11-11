package contract

var ChainDeploySteps = []ChainStep{
	StepCloneRepository,
	StepBuildImage,
	StepCreateDnsRecord,
	StepRemovePrevDeploy,
}
