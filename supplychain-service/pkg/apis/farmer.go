package apis

import (
	"encoding/json"
	"net/http"

	"supplychain-service/pkg/models"
	"supplychain-service/pkg/service"
	"supplychain-service/pkg/eventmgr"
)

// FarmerAPI represents the API for the farmer participant.
type FarmerAPI struct {
	BlobService *service.AzureBlobService
}

var farmerEventManager = eventmgr.NewEventManager()
// NewFarmerAPI creates a new FarmerAPI instance.
func NewFarmerAPI(blobService *service.AzureBlobService) *FarmerAPI {
	//eventManager := eventmgr.NewEventManager()
	return &FarmerAPI{
		BlobService: blobService,
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

		// Write the event as an immutable blob
		if err := api.BlobService.WriteBlob(product, "Harvested"); err != nil {
			http.Error(w, "Failed to write blob", http.StatusInternalServerError)
			return
		}

		farmerEventManager.Advertise("Harvested", event)
		// Respond with the harvested product
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	})); err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
}