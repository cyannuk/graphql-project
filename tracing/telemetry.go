package tracing

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/savsgio/gotils/strconv"
	otelcontrib "go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"

	"graphql-project/config"
)

type TracerProvider struct {
	trace.TracerProvider
}

type noopTracer struct {
	noop.Tracer
}

var tracer trace.Tracer

func (t noopTracer) Start(ctx context.Context, _ string, _ ...trace.SpanStartOption) (context.Context, trace.Span) {
	return ctx, noop.Span{}
}

func (tp TracerProvider) Tracer(name string, options ...trace.TracerOption) trace.Tracer {
	if tp.TracerProvider == nil {
		return noopTracer{}
	}
	return tp.TracerProvider.Tracer(name, options...)
}

func (tp TracerProvider) Shutdown(ctx context.Context) error {
	if tp.TracerProvider != nil {
		if provider, ok := tp.TracerProvider.(*sdktrace.TracerProvider); ok {
			return provider.Shutdown(ctx)
		}
	}
	return nil
}

func InitProviders(cfg *config.Config) (TracerProvider, error) {
	if !cfg.EnableTracing() {
		tracerProvider := TracerProvider{}
		otel.SetTracerProvider(tracerProvider)
		tracer = tracerProvider.Tracer("")
		return tracerProvider, nil
	}
	traceExporter, err := jaeger.New(jaeger.WithAgentEndpoint())
	if err != nil {
		return TracerProvider{}, err
	}

	// resource.Merge(resource.Default(), resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String("graphql-service")))
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // trace.TraceIDRatioBased(0.1)
		sdktrace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String("graphql-service"))),
	)

	// otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)
	tracer = tracerProvider.Tracer("", oteltrace.WithInstrumentationVersion(otelcontrib.Version()))

	return TracerProvider{tracerProvider}, nil
}

func InitSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return tracer.Start(ctx, name)
}

func Middleware(tracerProvider *sdktrace.TracerProvider) fiber.Handler {
	tracer := tracerProvider.Tracer("", oteltrace.WithInstrumentationVersion(otelcontrib.Version()))
	propagators := otel.GetTextMapPropagator()
	return func(ctx *fiber.Ctx) error {
		headerCarrier := make(propagation.HeaderCarrier, 16)
		ctx.Request().Header.VisitAll(func(k, v []byte) {
			headerCarrier.Set(strconv.B2S(k), strconv.B2S(v))
		})
		ctx.SetUserContext(propagators.Extract(ctx.UserContext(), headerCarrier))
		ctx.Locals("tracer", tracer)
		return ctx.Next()
	}
}
