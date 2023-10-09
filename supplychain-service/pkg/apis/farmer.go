package apis

import (
	"encoding/json"
	"net/http"

	"supplychain-service/pkg/models"
	"supplychain-service/pkg/service"
	"supplychain-service/pkg/eventmgr"
)

// FarmerAPI represents the API for the farmer participant.
// type FarmerAPI struct {
// 	BlobService *service.AzureBlobService
// }

// MongodB and EventManager
type FarmerAPI struct {
	DBService    *service.InMemoryDB
	EventManager *eventmgr.EventManager
}

func NewFarmerAPI(dbService *service.InMemoryDB, eventManager *eventmgr.EventManager) *FarmerAPI {
	return &FarmerAPI{
		DBService:    dbService,
		EventManager: eventManager,
	}
}

//HarvestHandler handles the "harvest" endpoint for the farmer.
func (api *FarmerAPI) HarvestHandler(w http.ResponseWriter, r *http.Request) {
	// Authenticate the request and check the user's role
	allowedRoles := []string{"farmer"}
	if err := AuthMiddleware(allowedRoles, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a new coffee product
		product := models.CreateCoffeeProduct("Arabica", "High-quality Arabica beans", 8.99)
		event := models.Event{Name: "Harvested", Payload: product}

		// Write to MongoDB
		if err := api.DBService.Write(product); err != nil {
			http.Error(w, "Failed to write product to MongoDB", http.StatusInternalServerError)
			return
		}

		api.EventManager.Advertise("Harvested", event)
		// Respond with the harvested product
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	})); err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
}