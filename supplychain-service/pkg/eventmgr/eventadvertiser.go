package eventmgr

import "sync"
import 	"supplychain-service/pkg/models"

type EventHandlerFunc func(event models.Event)

type EventManager struct {
	subscribers map[models.EventType][]EventHandlerFunc
	mu          sync.RWMutex
}

func NewEventManager() *EventManager {
	return &EventManager{
		subscribers: make(map[models.EventType][]EventHandlerFunc),
	}
}

func (em *EventManager) Advertise(eventType models.EventType, event models.Event) {
	em.mu.RLock()
	defer em.mu.RUnlock()

	if handlers, ok := em.subscribers[eventType]; ok {
		for _, handler := range handlers {
			go func(handler EventHandlerFunc) {
				handler(event)
			}(handler)
		}
	}
}

func (em *EventManager) Subscribe(eventType models.EventType, handler EventHandlerFunc) {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.subscribers[eventType] = append(em.subscribers[eventType], handler)
}