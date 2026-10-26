package github

import "go.uber.org/fx"

var Module = fx.Module("github",
	fx.Provide(NewClient),
)
