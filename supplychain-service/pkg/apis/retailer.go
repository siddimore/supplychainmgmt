package apis

import (
	"encoding/json"
	"fmt"
	"net/http"

	"supplychain-service/pkg/models"
	"supplychain-service/pkg/service"
	"supplychain-service/pkg/eventmgr"
)

// RetailerAPI represents the API for the retailer participant.
type RetailerAPI struct {
	DBService    *service.InMemoryDB
	EventManager *eventmgr.EventManager
}

func NewRetailerAPI(dbService *service.InMemoryDB, eventManager *eventmgr.EventManager) *RetailerAPI {
	return &RetailerAPI{
		DBService:    dbService,
		EventManager: eventManager,
	}
}

// SellHandler handles the "sell" endpoint for the retailer.
func (api *RetailerAPI) SellHandler(w http.ResponseWriter, r *http.Request) {
		// Create a new coffee product (e.g., received from the distributor)
		product := models.CreateCoffeeProduct("Arabica", "Sold Arabica beans", 13.99)
		event := models.Event{Name: "Sold", Payload: product}

		// Write to MongoDB
		if err := api.DBService.Write(product); err != nil {
			http.Error(w, "Failed to write product to MongoDB", http.StatusInternalServerError)
			return
		}

		api.EventManager.Advertise("Sold", event)

		// Respond with the sold product
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
		return
}

func (api *RetailerAPI) EventHandler(event models.Event) {
	// Handle Harvested event for the farmer here
	// Decide to Sell
	fmt.Printf("Retailer received Distributor event: %v\n", event)
}
