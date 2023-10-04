package eventmgr

import "sync"
import 	"supplychain-service/pkg/models"

//

// const (
// 	HarvestedEvent EventType = "Harvested"
// 	ProcessedEvent EventType = "Processed"
// 	PackedEvent    EventType = "Packed"
// 	ForSaleEvent   EventType = "ForSale"
// 	SoldEvent      EventType = "Sold"
// 	ShippedEvent   EventType = "Shipped"
// 	ReceivedEvent  EventType = "Received"
// 	PurchasedEvent EventType = "Purchased"
// )

// type Event struct {
// 	Type    EventType
// 	Payload interface{}
// }

type EventManager struct {
	subscribers map[models.EventType][]chan models.Event
	mu          sync.RWMutex
}

func NewEventManager() *EventManager {
	return &EventManager{
		subscribers: make(map[models.EventType][]chan models.Event),
	}
}

func (em *EventManager) Advertise(eventType models.EventType, event models.Event) {
	em.mu.RLock()
	defer em.mu.RUnlock()

	if chans, ok := em.subscribers[eventType]; ok {
		for _, ch := range chans {
			go func(ch chan models.Event) {
				ch <- event
			}(ch)
		}
	}
}

func (em *EventManager) Subscribe(eventType models.EventType) chan models.Event {
	em.mu.Lock()
	defer em.mu.Unlock()

	ch := make(chan models.Event)
	em.subscribers[eventType] = append(em.subscribers[eventType], ch)
	return ch
}