package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"supplychain-service/pkg/apis"
	"supplychain-service/pkg/service"
)

func main() {
	// Initialize immutable blob service (replace with your actual blob service implementation)
	blobService, _ := service.NewAzureBlobService("your-account-name", "your-account-key", "your-container-name")

	// Initialize API instances for each participant
	farmerAPI := apis.NewFarmerAPI(blobService)
	distributorAPI := apis.NewDistributorAPI(blobService)
	retailerAPI := apis.NewRetailerAPI(blobService)
	consumerAPI := apis.NewConsumerAPI(blobService)

	// Create a new router instance
	router := mux.NewRouter()

	// Farmer Endpoints
	router.HandleFunc("/farmer/harvest", farmerAPI.HarvestHandler).Methods("POST")
	// Add other farmer endpoints here for processed, packed, for sale, etc.

	// Distributor Endpoints
	router.HandleFunc("/distributor/receive", distributorAPI.ReceiveHandler).Methods("POST")
	// Add other distributor endpoints here for shipped, etc.

	// Retailer Endpoints
	router.HandleFunc("/retailer/sell", retailerAPI.SellHandler).Methods("POST")
	// Add other retailer endpoints here for sold, etc.

	// Consumer Endpoints
	router.HandleFunc("/consumer/consume", consumerAPI.ConsumeHandler).Methods("POST")
	// Add other consumer endpoints here for consumed, etc.

	// Serve the API using the router
	http.Handle("/", router)

	// Start the server on port 8080
	http.ListenAndServe(":8080", nil)
}



// package main

// import (
// 	"io"
// 	"log"
// 	"net/http"
// )

// func main() {
// 	// Hello world, the web server

// 	helloHandler := func(w http.ResponseWriter, req *http.Request) {
// 		io.WriteString(w, "Hello, world!\n")
// 	}

// 	http.HandleFunc("/hello", helloHandler)
//     log.Println("Listing for requests at http://localhost:8000/hello")
// 	log.Fatal(http.ListenAndServe(":8000", nil))
// }