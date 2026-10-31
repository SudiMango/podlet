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
	JobTypeEnsureDocker        JobType = "ensure_docker"
	JobTypeInstallDocker       JobType = "install_docker"
	JobTypeEnsurePodletNetwork JobType = "ensure_podlet_network"
	JobTypeCreatePodletNetwork JobType = "create_podlet_network"
	JobTypeEnsureCfTunnel      JobType = "ensure_cf_tunnel"
	JobTypeCreateCfTunnel      JobType = "create_cf_tunnel"
	JobTypeEnsureCfContainer   JobType = "ensure_cf_container"
	JobTypeCreateCfContainer   JobType = "create_cf_container"
	JobTypeEnsureRailpack      JobType = "ensure_railpack"
	JobTypeInstallRailpack     JobType = "install_railpack"
	JobTypeEnsureBuildkit      JobType = "ensure_buildkit"
	JobTypeSetupBuildkit       JobType = "setup_buildkit"
	JobTypeCloneRepository     JobType = "clone_repository"
)

type JobStatus string

const (
	StatusPending JobStatus = "pending"
	StatusOngoing JobStatus = "ongoing"
	StatusSuccess JobStatus = "success"
	StatusFailed  JobStatus = "failed"
)
