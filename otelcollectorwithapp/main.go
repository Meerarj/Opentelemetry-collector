package main

import (
	"context"
	"log"

	"github.com/SigNoz/sample-golang-app/controllers"
	// "github.com/SigNoz/sample-golang-app/metrics"
	"github.com/SigNoz/sample-golang-app/models"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	serviceName  = "my-app"
	collectorURL = "localhost:4317"
	insecure     = "true"
)

func initTracer() func(context.Context) error {

	var secureOption otlptracegrpc.Option

	// if strings.ToLower(insecure) == "true" || insecure == "1" || strings.ToLower(insecure) == "t" {
	// 	secureOption = otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	// } else {
	// 	secureOption = otlptracegrpc.WithInsecure()
	// }
	if len(insecure) > 0 {
		secureOption = otlptracegrpc.WithInsecure()
	}
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(collectorURL),
		),
	)

	if err != nil {
		log.Fatalf("Failed to create exporter: %v", err)
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Fatalf("Could not set resources: %v", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)
	return exporter.Shutdown
}

func main() {

	cleanup := initTracer()
	defer cleanup(context.Background())

	// provider := metrics.InitMeter()
	// defer provider.Shutdown(context.Background())

	// meter := provider.Meter("sample-golang-app")
	// metrics.GenerateMetrics(meter)

	r := gin.Default()
	r.Use(otelgin.Middleware(serviceName))
	// Connect to database
	models.ConnectDatabase()

	// Routes
	r.GET("/books", controllers.FindBooks)
	r.GET("/books/:id", controllers.FindBook)
	r.POST("/books", controllers.CreateBook)
	r.PATCH("/books/:id", controllers.UpdateBook)
	r.DELETE("/books/:id", controllers.DeleteBook)

	// Run the server
	r.Run(":8092")

}
