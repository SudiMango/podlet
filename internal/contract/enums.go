package contract

type ServerLifecycleState string

const (
	ServerLifecycleActive          ServerLifecycleState = "active"
	ServerLifecyclePendingDeletion ServerLifecycleState = "pending_deletion"
)

type HeartbeatAction string

const (
	ActionUninstall HeartbeatAction = "uninstall"
)

type JobType string

const (
	JobTypeDeployNewRepository JobType = "deploy_new_repository"
)

type JobStatus string

const (
	StatusPending JobStatus = "pending"
	StatusOngoing JobStatus = "ongoing"
	StatusSuccess JobStatus = "success"
	StatusFailed  JobStatus = "failed"
)
