package main

import (
	"fmt"
	// "io/ioutil"
	"net/http"
	"github.com/gorilla/mux"
	"supplychain-service/pkg/apis"
	"supplychain-service/pkg/client"
	"supplychain-service/pkg/models"
	"supplychain-service/pkg/service"
	"supplychain-service/pkg/eventmgr"
)

type MongoDBConfig struct {
	ConnectionString string
	DatabaseName     string
}

func main() {
	fmt.Println("This is a log message.")
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

	// Create InMemoryDBService
	inMemoryDbService := service.NewInMemoryDB()

	// // Initialize the event manager
	eventManager := eventmgr.NewEventManager()

	// Create MCCF client for Decentralized Authz
	authzClient, err := mccfclient.Create()
	if err != nil {
		fmt.Println("Error creating HTTP client:", err)
		return
	}

	// TODO: Above can be injected using GoContainers 
	// ex: https://github.com/vardius/gocontainer?utm_campaign=awesomego&utm_medium=referral&utm_source=awesomego
	// Initialize API instances for each participant
	// This API works with mccfAuthz 
	farmerAPI := apis.NewFarmerAPI(inMemoryDbService, eventManager)

	// This API needs to be updated with mccfAuthz 
	distributorAPI := apis.NewDistributorAPI(inMemoryDbService, eventManager)
	// This API needs to be updated with mccfAuthz 
	retailerAPI := apis.NewRetailerAPI(inMemoryDbService, eventManager)
	// This API needs to be updated with mccfAuthz 	
	consumerAPI := apis.NewConsumerAPI(inMemoryDbService, eventManager)

	// // Create a new router instance
	router := mux.NewRouter()

	// // TODO: Update all endpoint handlefunc to use MCCF.Authz instead
	// // Farmer Endpoints
	
	// router.HandleFunc("/farmer/harvest", apis.AuthMiddleware([]string{"farmer"}, farmerAPI.HarvestHandler)).Methods("POST")
	router.HandleFunc("/farmer/harvest", authzClient.AuthorizeAccess("886ddc0bd4f89adbf1e6ef81d66163b4e86202464d8131530d1468094d373ea3/action/harvested", farmerAPI.HarvestHandler)).Methods("POST")

	// Distributor Endpoints
	router.HandleFunc("/distributor/receive", authzClient.AuthorizeAccess("58ac36072ec568adc72d8a04dd2ae99b4d810ab373c0eec4323cb6652e4776e2/action/received", distributorAPI.ReceiveHandler)).Methods("POST")
	// Subscribe to HaverstedEvent by Farmer and handle it
	eventManager.Subscribe(models.HarvestedEvent, distributorAPI.EventHandler)
	// Add other distributor endpoints here for shipped, etc.

	// Retailer Endpoints
	router.HandleFunc("/retailer/sell", authzClient.AuthorizeAccess("0ac065deee073918ff159282f521e4716aecf9b73f761dea925f21721a46ae81", retailerAPI.SellHandler)).Methods("POST")
	// Subscribe to ReceivedEvent by Distributor and handle it
	eventManager.Subscribe(models.ReceivedEvent, retailerAPI.EventHandler)

	// Consumer Endpoints
	router.HandleFunc("/consumer/consume", authzClient.AuthorizeAccess("eec867e122818210a3ef92409d0870d14ac345b39596e608fad0ca94f40d8ff5", consumerAPI.ConsumeHandler)).Methods("POST")
	// Subscribe to SoldEvent by Retailer and handle it
	eventManager.Subscribe(models.SoldEvent, consumerAPI.EventHandler)
	// Add other consumer endpoints here for consumed, etc.

	// // Serve the API using the router
	http.Handle("/", router)

	// // Start the server on port 8080
	http.ListenAndServe("127.0.0.1:8080", nil)
}
