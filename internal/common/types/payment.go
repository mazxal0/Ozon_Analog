package types

type PaymentMethod string

const (
	PaymentMethodCard PaymentMethod = "bank_card"
	PaymentMethodSBP  PaymentMethod = "sbp"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSucceeded PaymentStatus = "succeeded"
	PaymentStatusCanceled  PaymentStatus = "canceled"
)
