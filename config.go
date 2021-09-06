package main

import (
	"car-svc/internal/app"
	"car-svc/internal/http"
	"car-svc/internal/lib/spanner"
	"context"

	lib_env "github.com/tomwangsvc/lib-svc/env"
	lib_log "github.com/tomwangsvc/lib-svc/log"
	lib_spanner "github.com/tomwangsvc/lib-svc/spanner"
)

type Config struct {
	App     app.Config
	Http    http.Config
	Spanner spanner.Config
}

func newConfig(ctx context.Context, env lib_env.Env) *Config {
	lib_log.Info(ctx, "Initializing config")

	config := Config{
		App: app.Config{
			Env: env,
		},
		Http: http.Config{
			Env: env,
		},

		Spanner: spanner.Config{
			ClientConfig: lib_spanner.ClientConfigWithMinOpened(env, 80),
			DatabaseId:   env.SpannerDatabaseId,
			Env:          env,
			InstanceId:   env.SpannerInstanceId,
			ProjectId:    env.GcpProjectId,
		},
	}

	lib_log.Info(ctx, "Initialized config")
	return &config
}
