package models

type EventType string

const (
	HarvestedEvent EventType = "Harvested"
	ProcessedEvent EventType = "Processed"
	PackedEvent    EventType = "Packed"
	ForSaleEvent   EventType = "ForSale"
	SoldEvent      EventType = "Sold"
	ShippedEvent   EventType = "Shipped"
	ReceivedEvent  EventType = "Received"
	PurchasedEvent EventType = "Purchased"
)

// type Event struct {
// 	Type    EventType
// 	Payload CoffeeProduct `json:"payload"`
// }

// // States
const (
	HarvestedState  = "Harvested"
	ProcessedState  = "Processed"
	PackedState     = "Packed"
	ForSaleState    = "ForSale"
	SoldState       = "Sold"
	ShippedState    = "Shipped"
	ReceivedState   = "Received"
	PurchasedState  = "Purchased"
)

// Event represents an immutable event in the supply chain.
type Event struct {
    Name    string        `json:"name"`
    Payload *CoffeeProduct `json:"payload"`
}