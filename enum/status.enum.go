package enum

type PaymentStatus string

const (
	Partial PaymentStatus = "partial"
	Unpaid  PaymentStatus = "unpaid"
	Paid    PaymentStatus = "paid"
)

type BookingStatus string

const (
	Pending   BookingStatus = "pending"
	Confirmed BookingStatus = "confirmed"
	Cancelled BookingStatus = "cancelled"
	Completed BookingStatus = "completed"
)
