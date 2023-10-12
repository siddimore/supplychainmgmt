package apis

import (
	"encoding/json"
	"fmt"
	"net/http"

	"supplychain-service/pkg/models"
	"supplychain-service/pkg/service"
	"supplychain-service/pkg/eventmgr"
)

// ConsumerAPI represents the API for the consumer participant.
type ConsumerAPI struct {
	DBService    *service.InMemoryDB
	EventManager *eventmgr.EventManager
}

func NewConsumerAPI(dbService *service.InMemoryDB, eventManager *eventmgr.EventManager) *ConsumerAPI {
	return &ConsumerAPI{
		DBService:    dbService,
		EventManager: eventManager,
	}
}

// ConsumeHandler handles the "consume" endpoint for the consumer.
func (api *ConsumerAPI) ConsumeHandler(w http.ResponseWriter, r *http.Request) {
		// Create a new coffee product (e.g., purchased from the retailer)
		product := models.CreateCoffeeProduct("Arabica", "Consumed Arabica beans", 14.99)
		event := models.Event{Name: "Consumed", Payload: product}

		// Write to mongod db
		if err := api.DBService.Write(product); err != nil {
			http.Error(w, "Failed to write product to MongoDB", http.StatusInternalServerError)
			return
		}

		api.EventManager.Advertise("Consumed", event)

		// Respond with the consumed product
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
		return
}

func (api *ConsumerAPI) EventHandler(event models.Event) {
	// Handle Harvested event for the farmer here
	// Decide to Buy
	fmt.Printf("Consumer received Sold event: %v\n", event)
}