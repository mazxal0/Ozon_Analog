package types

type OrderStatus string

const (
	inProgress OrderStatus = "in_progress"
	paid       OrderStatus = "paid"
	completed  OrderStatus = "completed"
	failed     OrderStatus = "failed"
)
