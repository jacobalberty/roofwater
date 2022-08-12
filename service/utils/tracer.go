package utils

import (
	"context"
	"os"
	"os/user"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var Tracer trace.Tracer

type TracerConfig struct {
	ServiceName string `envconfig:"SERVICE_NAME" default:"roofwaterd" desc:"Service name to use for tracing"`
}

func InitializeTracer(ctx context.Context, tc TracerConfig) func() {
	exp, err := newExporter(ctx, tc)
	if err != nil {
		os.Exit(1)
	}
	tp := newTraceProvider(exp, tc)
	otel.SetTracerProvider(tp)
	Tracer = tp.Tracer(tc.ServiceName)
	return func() { _ = tp.Shutdown(ctx) }
}
func newExporter(ctx context.Context, tc TracerConfig) (sdktrace.SpanExporter, error) {
	switch {
	case os.Getenv("OTEL_EXPORTER_JAEGER_ENDPOINT") != "":
		return jaeger.New(jaeger.WithCollectorEndpoint())
	default:
		return stdouttrace.New(stdouttrace.WithPrettyPrint())

	}
}

func newTraceProvider(exp sdktrace.SpanExporter, tc TracerConfig) *sdktrace.TracerProvider {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.Int("process.pid", os.Getpid()),
			attribute.String("process.executable", os.Args[0]),
			attribute.StringSlice("process.command_args", os.Args[1:]),
			attribute.String("process.owner", user.Username),
			attribute.String("service.name", tc.ServiceName),
			attribute.String("library.language", "go"),
		),
	)

	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}
