package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var tracer trace.Tracer

func main() {
	// init tracer and ensure cleanup after app ends
	ctx := context.Background()
	exp := newExporter(ctx)

	tp := newTraceProvider(exp)
	defer func() { _ = tp.Shutdown(ctx) }()

	otel.SetTracerProvider(tp)
	tracer = tp.Tracer("otel-texample")

	r := gin.Default()
	r.Use(otelgin.Middleware("otlp-example"))
	r.GET("/reverse/:str", reverseStrHandler())

	r.Run() // listen on :8080
}

func newExporter(ctx context.Context) *otlptrace.Exporter {
	secureOption := otlptracegrpc.WithInsecure() // TODO Use WithTLSCredentials for prod

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint("localhost:4317"),
		),
	)

	if err != nil {
		log.Fatalf("failed to init exporter: %v", err)
	}

	return exporter
}

func newTraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	r, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("otlp-example"),
			semconv.ServiceNamespaceKey.String("localsvc"),
		),
	)

	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // TODO Set sane sampling rates when deployed
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}
