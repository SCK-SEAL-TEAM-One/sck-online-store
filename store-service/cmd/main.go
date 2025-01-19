package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"store-service/cmd/api"
	"store-service/internal/cart"
	"store-service/internal/healthcheck"
	"store-service/internal/order"
	"store-service/internal/payment"
	"store-service/internal/shipping"

	"store-service/internal/point"
	"store-service/internal/product"

	"github.com/Cyprinus12138/otelgin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"go.elastic.co/apm/module/apmgin"

	"github.com/penglongli/gin-metrics/ginmetrics"

	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	serviceName  = os.Getenv("SERVICE_NAME")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	insecure     = os.Getenv("INSECURE_MODE")
)

func getResource() *resource.Resource {
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Printf("Could not set resources: %s ", err)
	}
	return resources
}

func initMetric(resources *resource.Resource, ) func(context context.Context ) error {

	secureOption := otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if len(insecure) > 0 {
		secureOption = otlpmetricgrpc.WithInsecure()
	}

	exporter, err := otlpmetricgrpc.New(
		context.Background(), 
		secureOption, 
		otlpmetricgrpc.WithEndpoint(collectorURL),
	)
	if err != nil {
		log.Fatal(err)
	}
	otel.SetMeterProvider(
		sdkmetric.NewMeterProvider(
			sdkmetric.WithResource(resources),
			sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		),
	)
	return exporter.Shutdown
}

func initTracer(resources *resource.Resource) func(context.Context) error {

	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if len(insecure) > 0 {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptracegrpc.New(
		context.Background(), 
		secureOption, 
		otlptracegrpc.WithEndpoint(collectorURL),
	)
	if err != nil {
		log.Fatal(err)
	}
	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithResource(resources),
			sdktrace.WithBatcher(exporter),
		),
	)
	return exporter.Shutdown
}
	
func main() {

	bankGatewayEndpoint := "thirdparty:8882"
	shippingGatewayEndpoint := "thirdparty:8883"
	pointGatewayEndpoint := "point-service:8001"
	storeWebEndpoint := "http://localhost:3000"

	if os.Getenv("BANK_GATEWAY") != "" {
		bankGatewayEndpoint = os.Getenv("BANK_GATEWAY")
	}
	if os.Getenv("SHIPPING_GATEWAY") != "" {
		shippingGatewayEndpoint = os.Getenv("SHIPPING_GATEWAY")
	}
	if os.Getenv("POINT_GATEWAY") != "" {
		pointGatewayEndpoint = os.Getenv("POINT_GATEWAY")
	}
	if os.Getenv("STORE_WEB_HOST") != "" {
		storeWebEndpoint = os.Getenv("STORE_WEB_HOST")
	}

	dbConnection := "user:password@(db:3306)/store"
	if os.Getenv("DB_CONNECTION") != "" {
		dbConnection = os.Getenv("DB_CONNECTION")
	}

	connection, err := sqlx.Connect("mysql", dbConnection)
	if err != nil {
		log.Fatalln("cannot connect to database", err)
	}
	defer connection.Close()

	productRepository := product.ProductRepositoryMySQL{
		DBConnection: connection,
	}
	orderRepository := order.OrderRepositoryMySQL{
		DBConnection: connection,
	}
	cartRepository := cart.CartRepositoryMySQL{
		DBConnection: connection,
	}
	shippingRepository := shipping.ShippingRepositoryMySQL{
		DBConnection: connection,
	}

	bankGateway := payment.BankGateway{
		BankEndpoint: "http://" + bankGatewayEndpoint,
	}
	shippingGateway := shipping.ShippingGateway{
		ShippingEndpoint: "http://" + shippingGatewayEndpoint,
	}
	pointGateway := point.PointGateway{
		PointEndpoint: "http://" + pointGatewayEndpoint,
	}

	paymentService := payment.PaymentService{
		BankGateway:       &bankGateway,
		ShippingGateway:   &shippingGateway,
		OrderRepository:   &orderRepository,
		ProductRepository: productRepository,
	}
	pointService := point.PointService{
		PointGateway: &pointGateway,
	}
	cartService := cart.CartService{
		CartRepository: &cartRepository,
	}
	productService := product.ProductService{
		ProductRepository: &productRepository,
	}
	orderService := order.OrderService{
		CartRepository:     cartRepository,
		OrderRepository:    &orderRepository,
		PointService:       pointService,
		ProductRepository:  productRepository,
		ShippingRepository: shippingRepository,
	}

	cartAPI := api.CartAPI{
		CartService: &cartService,
	}
	orderAPI := api.OrderAPI{
		OrderService: &orderService,
	}
	paymentAPI := api.PaymentAPI{
		PaymentService: &paymentService,
	}
	productAPI := api.ProductAPI{
		ProductService: &productService,
	}
	pointAPI := api.PointAPI{
		PointService: pointService,
	}


	resource := getResource()
	cleanupTracer := initTracer(resource)
	defer cleanupTracer(context.Background())

	cleanupMetric := initMetric(resource)
	defer cleanupMetric(context.Background())

	route := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{storeWebEndpoint}
	route.Use(cors.New(config))
	route.Use(otelgin.Middleware(serviceName))


	// get global Monitor object
	m := ginmetrics.GetMonitor()

	// +optional set metric path, default /debug/metrics
	m.SetMetricPath("/metrics")
	// +optional set slow time, default 5s
	m.SetSlowTime(10)
	// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
	// used to p95, p99
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10, 50, 100, 500})

	// set middleware for gin
	m.Use(route)

	route.Use(apmgin.Middleware(route))
	route.GET("/api/v1/product", productAPI.SearchHandler)
	route.GET("/api/v1/product/:id", productAPI.GetProductHandler)
	route.GET("/api/v1/cart", cartAPI.GetCartHandler)
	route.PUT("/api/v1/addCart", cartAPI.AddCartHandler)
	route.PUT("/api/v1/updateCart", cartAPI.UpdateCartHandler)
	route.POST("/api/v1/order", orderAPI.SubmitOrderHandler)
	route.POST("/api/v1/confirmPayment", paymentAPI.ConfirmPaymentHandler)
	route.GET("/api/v1/point", pointAPI.TotalPointHandler)
	route.POST("/api/v1/point", pointAPI.DeductPointHandler)

	route.GET("/api/v1/health", func(context *gin.Context) {
		user, err := healthcheck.GetUserNameFromDB(connection)
		if err != nil {
			context.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		context.JSON(200, gin.H{
			"message": user,
		})
	})

	log.Fatal(route.Run(":8000"))
}

func GetUserNameFromDB(connection *sqlx.DB) User {
	user := User{}
	err := connection.Get(&user, "SELECT id,first_name,last_name FROM user WHERE id=1")
	if err != nil {
		fmt.Printf("Get user name from tearup get error : %s", err.Error())
		return User{}
	}
	return user
}

type User struct {
	ID        int    `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}
