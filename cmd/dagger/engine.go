package main

import (
	"context"
	"os"

	"dagger.io/dagger/telemetry"
	"github.com/dagger/dagger/engine"
	"github.com/dagger/dagger/engine/client"
	"github.com/dagger/dagger/engine/slog"
	enginetel "github.com/dagger/dagger/engine/telemetry"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type runClientCallback func(context.Context, *client.Client) error

func withEngine(
	ctx context.Context,
	params client.Params,
	fn runClientCallback,
) error {
	return Frontend.Run(ctx, opts, func(ctx context.Context) (rerr error) {
		// Init tracing as early as possible and shutdown after the command
		// completes, ensuring progress is fully flushed to the frontend.
		cleanupTelemetry := initEngineTelemetry(ctx)
		defer cleanupTelemetry(rerr)

		if debug {
			params.LogLevel = slog.LevelDebug
		}

		if params.RunnerHost == "" {
			var err error
			params.RunnerHost, err = engine.RunnerHost()
			if err != nil {
				return err
			}
		}

		params.DisableHostRW = disableHostRW

		params.EngineCallback = Frontend.ConnectedToEngine
		params.CloudURLCallback = Frontend.SetCloudURL

		params.EngineTrace = telemetry.SpanForwarder{
			Processors: telemetry.SpanProcessors,
		}
		params.EngineLogs = telemetry.LogForwarder{
			Processors: telemetry.LogProcessors,
		}
		params.WithTerminal = withTerminal
		params.Interactive = interactive

		// Connect to and run with the engine
		sess, ctx, err := client.Connect(ctx, params)
		if err != nil {
			return err
		}
		defer sess.Close()

		return fn(ctx, sess)
	})
}

func initEngineTelemetry(ctx context.Context) func(error) {
	// Setup telemetry config
	telemetryCfg := telemetry.Config{
		Detect:   true,
		Resource: Resource(),

		LiveTraceExporters: []sdktrace.SpanExporter{Frontend.SpanExporter()},
		LiveLogExporters:   []sdklog.Exporter{Frontend.LogExporter()},
	}
	if spans, logs, ok := enginetel.ConfiguredCloudExporters(ctx); ok {
		telemetryCfg.LiveTraceExporters = append(telemetryCfg.LiveTraceExporters, spans)
		telemetryCfg.LiveLogExporters = append(telemetryCfg.LiveLogExporters, logs)
	}
	ctx = telemetry.Init(ctx, telemetryCfg)

	// Set the full command string as the name of the root span.
	//
	// If you pass credentials in plaintext, yes, they will be leaked; don't do
	// that, since they will also be leaked in various other places (like the
	// process tree). Use Secret arguments instead.
	ctx, span := Tracer().Start(ctx, spanName(os.Args))

	// Set up global slog to log to the primary span output.
	slog.SetDefault(slog.SpanLogger(ctx, InstrumentationLibrary))

	// Set the span as the primary span for the frontend.
	Frontend.SetPrimary(span.SpanContext().SpanID())

	// Direct command stdout/stderr to span stdio via OpenTelemetry.
	stdio := telemetry.SpanStdio(ctx, InstrumentationLibrary)
	rootCmd.SetOut(stdio.Stdout)
	rootCmd.SetErr(stdio.Stderr)

	return func(rerr error) {
		stdio.Close()
		telemetry.End(span, func() error { return rerr })
		telemetry.Close()
	}
}
