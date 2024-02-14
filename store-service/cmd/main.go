package main

import (
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

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"go.elastic.co/apm/module/apmgin"

	"github.com/penglongli/gin-metrics/ginmetrics"
)

func main() {

	bankGatewayEndpoint := "bank-gateway:8882"
	shippingGatewayEndpoint := "shipping-gateway:8882"
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

	dbConnecton := "user:password@(store-db:3306)/store"
	if os.Getenv("DBCONNECTION") != "" {
		dbConnecton = os.Getenv("DBCONNECTION")
	}
	connection, err := sqlx.Connect("mysql", dbConnecton)
	if err != nil {
		log.Fatalln("cannot connect to database", err)
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

	route := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{storeWebEndpoint}
	route.Use(cors.New(config))

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
