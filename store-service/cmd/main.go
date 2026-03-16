package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"store-service/cmd/api"
	"store-service/internal/auth"
	"store-service/internal/cart"
	"store-service/internal/healthcheck"
	"store-service/internal/middleware"
	storeOtel "store-service/internal/otel"
	"store-service/internal/order"
	"store-service/internal/payment"
	"store-service/internal/shipping"
	"store-service/internal/user"
	"time"

	"store-service/internal/point"
	"store-service/internal/product"

	"github.com/XSAM/otelsql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	_ "store-service/cmd/docs"

	_ "time/tzdata"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var (
	serviceName  = os.Getenv("OTEL_SERVICE_NAME")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	insecure     = os.Getenv("INSECURE_MODE")
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if collectorURL != "" {
		cleanup, err := storeOtel.InitOtel(ctx, serviceName, collectorURL, insecure)
		if err != nil {
			log.Fatalf("failed to initialize OpenTelemetry: %v", err)
		}
		defer cleanup()
	}

	http.DefaultTransport = otelhttp.NewTransport(http.DefaultTransport)

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

	dbConnection := "user:password@(localhost:3306)/store?parseTime=True"
	if os.Getenv("DB_CONNECTION") != "" {
		dbConnection = os.Getenv("DB_CONNECTION")
	}

	driverName, err := otelsql.Register("mysql", otelsql.WithAttributes(semconv.DBSystemMySQL))
	if err != nil {
		log.Fatalln("cannot register otelsql driver", err)
	}

	connection, err := sqlx.Connect(driverName, dbConnection)
	if err != nil {
		log.Fatalln("cannot connect to database", err)
	}
	defer connection.Close()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

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
	userRepository := auth.UserRepositoryMySQL{
		DBConnection: connection,
	}
	jwtManager := &auth.JWTTokenManager{
		SecretKey: jwtSecret,
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

	PDFHelper := order.OrderSummaryPDFGenerator{}
	orderHelper := order.OrderHelper{}

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
		UserRepository:     userRepository,
		PDFHelper:          PDFHelper,
		OrderHelper:        orderHelper,
		Clock:              time.Now,
	}
	authService := auth.AuthService{
		UserRepository:  userRepository,
		JWTTokenManager: jwtManager,
		PasswordHelper:  user.BcryptPasswordChecker{},
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
	authAPI := api.AuthAPI{
		AuthService: authService,
	}

	route := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{storeWebEndpoint}
	// allow uid in headers
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "uid", "Authorization"}
	// allow cookies
	config.AllowCredentials = true
	route.Use(cors.New(config))

	route.Use(otelgin.Middleware(serviceName))

	v1 := route.Group("/api/v1")
	// -------------------------------------------
	// Public /api/v1 endpoints
	// -------------------------------------------
	v1.GET("/health", func(context *gin.Context) {
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

	v1.GET("/refreshToken", authAPI.RefreshTokenHandler)
	v1.POST("/login", authAPI.LoginHandler)

	// -------------------------------------------
	// Protected /api/v1 endpoints
	// -------------------------------------------
	protected := v1.Group("/")
	protected.Use(middleware.AuthUser(jwtSecret))

	protected.GET("/product", productAPI.SearchHandler)
	protected.GET("/product/:id", productAPI.GetProductHandler)

	protected.GET("/cart", cartAPI.GetCartHandler)
	protected.PUT("/addCart", cartAPI.AddCartHandler)
	protected.PUT("/updateCart", cartAPI.UpdateCartHandler)

	protected.POST("/order", orderAPI.SubmitOrderHandler)
	protected.POST("/order/:id/summary", orderAPI.GetOrderSummaryHandler)
	protected.POST("/confirmPayment", paymentAPI.ConfirmPaymentHandler)

	protected.GET("/point", pointAPI.TotalPointHandler)
	protected.POST("/point", pointAPI.DeductPointHandler)

	//docs.SwaggerInfo.BasePath = "/api/v1"
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

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
