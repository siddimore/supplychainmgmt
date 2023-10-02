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
	BlobService *service.AzureBlobService
}

var consumerEventManager = eventmgr.NewEventManager()

// NewConsumerAPI creates a new ConsumerAPI instance.
func NewConsumerAPI(blobService *service.AzureBlobService) *ConsumerAPI {
	return &ConsumerAPI{
		BlobService: blobService,
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

		// Write the event as an immutable blob
		if err := api.BlobService.WriteBlob(product, "Consumed"); err != nil {
			http.Error(w, "Failed to write blob", http.StatusInternalServerError)
			return
		}

		consumerEventManager.Advertise("Consumed", event)

		// Respond with the consumed product
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	})); err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
}

// ... (
