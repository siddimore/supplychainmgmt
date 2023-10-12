
package apis

import (
	"encoding/json"
	"fmt"
	"net/http"

	"supplychain-service/pkg/models"
	"supplychain-service/pkg/service"
	"supplychain-service/pkg/eventmgr"
)

// DistributorAPI represents the API for the distributor participant.
// MongodB and EventManager
type DistributorAPI struct {
	DBService    *service.InMemoryDB
	EventManager *eventmgr.EventManager
}

func NewDistributorAPI(dbService *service.InMemoryDB, eventManager *eventmgr.EventManager) *DistributorAPI {
	return &DistributorAPI{
		DBService:    dbService,
		EventManager: eventManager,
	}
}

// ReceiveHandler handles the "receive" endpoint for the distributor.
func (api *DistributorAPI) ReceiveHandler(w http.ResponseWriter, r *http.Request) {
	product := models.CreateCoffeeProduct("Arabica", "Received Arabica beans", 12.99)
	event := models.Event{Name: "Received", Payload: product}

	// Write to MongoDB
	if err := api.DBService.Write(product); err != nil {
		http.Error(w, "Failed to write product to MongoDB", http.StatusInternalServerError)
		return
	}

	api.EventManager.Advertise("Received", event)

	// Respond with the received product
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
	return
}

func (api *DistributorAPI) EventHandler(event models.Event) {
	// Handle Harvested event for the farmer here
	// Decide to Distribute
	fmt.Printf("Distributor received Harvested event: %v\n", event)
}
