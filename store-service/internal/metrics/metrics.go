package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var (
	OrdersCreated    metric.Int64Counter
	OrderRevenue     metric.Float64Counter
	OrderItemsCount  metric.Int64Histogram
	OrderTotalPrice  metric.Float64Histogram
	PaymentAttempts  metric.Int64Counter
	PaymentDuration  metric.Float64Histogram
)

func Init() {
	meter := otel.Meter("store-service")

	OrdersCreated, _ = meter.Int64Counter("orders.created",
		metric.WithDescription("Number of orders created"),
	)

	OrderRevenue, _ = meter.Float64Counter("order.revenue",
		metric.WithDescription("Total revenue from orders in THB"),
	)

	OrderItemsCount, _ = meter.Int64Histogram("order.items.count",
		metric.WithDescription("Number of items per order"),
	)

	OrderTotalPrice, _ = meter.Float64Histogram("order.total.thb",
		metric.WithDescription("Total price per order in THB"),
	)

	PaymentAttempts, _ = meter.Int64Counter("payment.attempts",
		metric.WithDescription("Number of payment attempts"),
	)

	PaymentDuration, _ = meter.Float64Histogram("payment.duration",
		metric.WithDescription("Duration of bank gateway payment calls in seconds"),
	)
}
