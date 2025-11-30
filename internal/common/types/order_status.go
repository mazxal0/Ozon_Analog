package types

type OrderStatus string

const (
	InProgress OrderStatus = "in_progress"
	Paid       OrderStatus = "paid"
	Completed  OrderStatus = "completed"
	Failed     OrderStatus = "failed"
	Cancelled  OrderStatus = "cancelled"
)
