package dto

import "github.com/SudiMango/podlet/internal/contract"

type ChainStep struct {
	Name    contract.ChainStep
	JobType JobType
}
