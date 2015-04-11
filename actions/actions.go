package actions

type orderId string

type Order struct {
	Id orderId
}

func NewOrder() Order {
	return Order{
		Id: "asdflkj",
	}
}
