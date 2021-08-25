package main

import (
	"car-svc/internal/app"
	"car-svc/internal/http"
	"car-svc/internal/lib/schema"
	"car-svc/internal/lib/spanner"
	"log"

	lib_context "github.com/tomwangsvc/lib-svc/context"
	lib_countries "github.com/tomwangsvc/lib-svc/countries"
	lib_env "github.com/tomwangsvc/lib-svc/env"
	lib_log "github.com/tomwangsvc/lib-svc/log"
	lib_os "github.com/tomwangsvc/lib-svc/os"
	lib_schema "github.com/tomwangsvc/lib-svc/schema"
	lib_storage "github.com/tomwangsvc/lib-svc/storage"
)

//revive:disable:cyclomatic
func main() {
	env, err := lib_env.New("car-svc")
	if err != nil {
		log.Fatal("Failed initializing env: ", err)
	}
	env.SpannerDatabaseId = "car-svc"
	env.SpannerInstanceId = "tom-wang-dev"
	env.GcpProjectId = "data-fabric-323905"
	env.GcpProjectNumber = "518937487179"
	ctx := lib_context.NewStartUpContext()

	if err = lib_log.Init(ctx, *env); err != nil {
		log.Fatal("Failed initializing logger: ", err)
	}

	lib_log.Info(ctx, "Initializing config")
	config := newConfig(ctx, *env)
	lib_log.Info(ctx, "Initialized config")

	schemaClient, err := lib_schema.NewClient(ctx, schema.SupportedSchema())
	if err != nil {
		lib_log.Fatal(ctx, "Failed initializing schema client", lib_log.FmtError(err))
	}

	libStorageClient, err := lib_storage.NewClient(ctx)
	if err != nil {
		lib_log.Fatal(ctx, "Failed initializing storage client", lib_log.FmtError(err))
	}
	defer func() {
		if err := libStorageClient.Close(); err != nil {
			lib_log.Error(ctx, "Failed closing storage client", lib_log.FmtError(err))
		}
	}()

	spannerClient, err := spanner.NewClient(ctx, config.Spanner)
	if err != nil {
		lib_log.Fatal(ctx, "Failed initializing spanner client", lib_log.FmtError(err))
	}
	defer spannerClient.Close()

	appClient := app.NewClient(config.App, spannerClient)

	countriesMetadata, err := lib_countries.NewMetadata(ctx)
	if err != nil {
		lib_log.Fatal(ctx, "Failed initializing country metadata", lib_log.FmtError(err))
	}

	lib_log.Info(ctx, "Initializing http client")
	httpClient, err := http.NewClient(config.Http, appClient, schemaClient, *countriesMetadata)
	if err != nil {
		lib_log.Fatal(ctx, "Failed initializing http client", lib_log.FmtError(err))
	}
	lib_log.Info(ctx, "Initialized http client")

	lib_os.CleanUpAndExitOnInterrupt(ctx, []lib_os.Closer{spannerClient}, []lib_os.CloserWithError{libStorageClient}, []lib_os.Flusher{})

	lib_log.Info(ctx, "Listening and serving HTTP client", lib_log.FmtInt("config.Http.Env.Port", config.Http.Env.Port))
	if err := httpClient.ListenAndServe(); err != nil {
		lib_log.Fatal(ctx, "HTTP client unexpectedly returned with error, terminating...", lib_log.FmtError(err))
	}
	lib_log.Fatal(ctx, "HTTP client unexpectedly returned, terminating...")
	//revive:enable:cyclomatic
}
