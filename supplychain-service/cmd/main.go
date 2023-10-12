package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"supplychain-service/pkg/apis"
	"supplychain-service/pkg/service"
	"supplychain-service/pkg/eventmgr"
)

type MongoDBConfig struct {
	ConnectionString string
	DatabaseName     string
}

func main() {
	// Initialize immutable blob service (replace with your actual blob service implementation)
	// TODO: Remove Blob Service
	// blobService, _ := service.NewAzureBlobService("your-account-name", "your-account-key", "your-container-name")

	// Load MongoDB configuration (you can load this from a configuration file or environment variables)
	// Init MongoDB
	// mongoDBConfig := MongoDBConfig{
	// 	ConnectionString: "mongodb://username:password@localhost:27017", // Replace with your MongoDB connection string
	// 	DatabaseName:     "your-database-name",                          // Replace with your database name
	// }

	// Initialize MongoDB service
	// mongoDBService, err := service.NewMongoDBService(mongoDBConfig.ConnectionString, mongoDBConfig.DatabaseName)
	// if err != nil {
	// 	fmt.Printf("Failed to connect to MongoDB:", err)
	// 	return
	// }

	inMemoryDbService := service.NewInMemoryDB()

	// Initialize the event manager
	eventManager := eventmgr.NewEventManager()

	// TODO: Above can be injected using GoContainers 
	// ex: https://github.com/vardius/gocontainer?utm_campaign=awesomego&utm_medium=referral&utm_source=awesomego


	// Initialize API instances for each participant
	farmerAPI := apis.NewFarmerAPI(inMemoryDbService, eventManager)
	distributorAPI := apis.NewDistributorAPI(inMemoryDbService, eventManager)
	retailerAPI := apis.NewRetailerAPI(inMemoryDbService, eventManager)
	consumerAPI := apis.NewConsumerAPI(inMemoryDbService, eventManager)

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
