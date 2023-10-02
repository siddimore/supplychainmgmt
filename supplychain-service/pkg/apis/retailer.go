package apis

import (
	"encoding/json"
	"net/http"

	"supplychain-service/pkg/models"
	"supplychain-service/pkg/service"
	"supplychain-service/pkg/eventmgr"
)

// RetailerAPI represents the API for the retailer participant.
type RetailerAPI struct {
	BlobService *service.AzureBlobService
}

var retailerEventManager = eventmgr.NewEventManager()

// NewRetailerAPI creates a new RetailerAPI instance.
func NewRetailerAPI(blobService *service.AzureBlobService) *RetailerAPI {
	return &RetailerAPI{
		BlobService: blobService,
	}
}

// SellHandler handles the "sell" endpoint for the retailer.
func (api *RetailerAPI) SellHandler(w http.ResponseWriter, r *http.Request) {
	// Authenticate the request and check the user's role
	allowedRoles := []string{"retailer"}
	if err := AuthMiddleware(allowedRoles, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a new coffee product (e.g., received from the distributor)
		product := models.CreateCoffeeProduct("Arabica", "Sold Arabica beans", 13.99)
		event := models.Event{Name: "Sold", Payload: product}

		// Write the event as an immutable blob
		if err := api.BlobService.WriteBlob(product, "Sold"); err != nil {
			http.Error(w, "Failed to write blob", http.StatusInternalServerError)
			return
		}

		retailerEventManager.Advertise("Sold", event)

		// Respond with the sold product
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	})); err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
}
