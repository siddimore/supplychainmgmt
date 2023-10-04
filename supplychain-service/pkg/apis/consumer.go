package apis

import (
	"encoding/json"
	"net/http"

	"supplychain-service/pkg/models"
	"supplychain-service/pkg/service"
	"supplychain-service/pkg/eventmgr"
)

// ConsumerAPI represents the API for the consumer participant.
type ConsumerAPI struct {
	DBService    *service.MongoDBService
	EventManager *eventmgr.EventManager
}

func NewConsumerAPI(dbService *service.MongoDBService, eventManager *eventmgr.EventManager) *ConsumerAPI {
	return &ConsumerAPI{
		DBService:    dbService,
		EventManager: eventManager,
	}
}

// ConsumeHandler handles the "consume" endpoint for the consumer.
func (api *ConsumerAPI) ConsumeHandler(w http.ResponseWriter, r *http.Request) {
	// Authenticate the request and check the user's role
	allowedRoles := []string{"consumer"}
	if err := AuthMiddleware(allowedRoles, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})); err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
}

// ... (
