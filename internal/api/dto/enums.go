package dto

type ServerLifecycleState string

const (
	ServerLifecycleActive          ServerLifecycleState = "active"
	ServerLifecyclePendingDeletion ServerLifecycleState = "pending_deletion"
)
