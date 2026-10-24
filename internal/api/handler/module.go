package handler

import "go.uber.org/fx"

var Module = fx.Module("service",
	fx.Provide(
		NewHealthHandler,
		NewProfileHandler,
		NewServerHandler,
	),
)
