package middlewares

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracerProvider *sdktrace.TracerProvider
	loggerProvider *sdklog.LoggerProvider
	tracer         trace.Tracer
)

// InitOpenTelemetry initializes OpenTelemetry tracer and logger providers
func InitOpenTelemetry() error {
	ctx := context.Background()

	// Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("gfp-api"),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	// Setup trace exporter (console/stdout)
	traceExporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Setup tracer provider
	tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	otel.SetTracerProvider(tracerProvider)

	// Setup propagators for distributed tracing
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	// Setup log exporter (console/stdout)
	logExporter, err := stdoutlog.New(
		stdoutlog.WithPrettyPrint(),
	)
	if err != nil {
		return fmt.Errorf("failed to create log exporter: %w", err)
	}

	// Setup logger provider
	loggerProvider = sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
		sdklog.WithResource(res),
	)

	// Create global tracer
	tracer = otel.Tracer("gfp-api")

	log.Println("✅ OpenTelemetry initialized successfully")
	return nil
}

// Shutdown gracefully shuts down the OpenTelemetry providers
func Shutdown(ctx context.Context) error {
	log.Println("🔄 Shutting down OpenTelemetry...")

	// Create a timeout context for shutdown
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var errs []error

	// Shutdown tracer provider
	if tracerProvider != nil {
		if err := tracerProvider.Shutdown(shutdownCtx); err != nil {
			errs = append(errs, fmt.Errorf("tracer provider shutdown error: %w", err))
		}
	}

	// Shutdown logger provider
	if loggerProvider != nil {
		if err := loggerProvider.Shutdown(shutdownCtx); err != nil {
			errs = append(errs, fmt.Errorf("logger provider shutdown error: %w", err))
		}
	}

	if len(errs) > 0 {
		for _, err := range errs {
			log.Printf("❌ %v\n", err)
		}
		return fmt.Errorf("shutdown completed with %d error(s)", len(errs))
	}

	log.Println("✅ OpenTelemetry shutdown complete")
	return nil
}

// Tracer returns the global tracer instance
func Tracer() trace.Tracer {
	return tracer
}

// LoggerProvider returns the global logger provider instance
func LoggerProvider() *sdklog.LoggerProvider {
	return loggerProvider
}

// LogInfo logs an informational message with trace context
func LogInfo(ctx context.Context, message string, attrs ...interface{}) {
	spanCtx := trace.SpanFromContext(ctx).SpanContext()
	if spanCtx.IsValid() {
		log.Printf("[INFO] [trace_id=%s span_id=%s] %s %v\n",
			spanCtx.TraceID().String(),
			spanCtx.SpanID().String(),
			message,
			attrs,
		)
	} else {
		log.Printf("[INFO] %s %v\n", message, attrs)
	}
}

// LogError logs an error message with trace context
func LogError(ctx context.Context, message string, err error, attrs ...interface{}) {
	spanCtx := trace.SpanFromContext(ctx).SpanContext()
	if spanCtx.IsValid() {
		fmt.Fprintf(os.Stderr, "[ERROR] [trace_id=%s span_id=%s] %s: %v %v\n",
			spanCtx.TraceID().String(),
			spanCtx.SpanID().String(),
			message,
			err,
			attrs,
		)
	} else {
		fmt.Fprintf(os.Stderr, "[ERROR] %s: %v %v\n", message, err, attrs)
	}
}
