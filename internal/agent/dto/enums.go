package dto

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
	JobTypeRemoveCfContainer   JobType = "remove_cf_container"
	JobTypeRemoveBuildkit      JobType = "remove_buildkit"
	JobTypeRemovePodletNetwork JobType = "remove_podlet_network"
	JobTypeUninstallRailpack   JobType = "uninstall_railpack"
	JobTypeCloneRepository     JobType = "clone_repository"
	JobTypeBuildRepository     JobType = "build_repository"
)
