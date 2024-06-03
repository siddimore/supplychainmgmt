package client


import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

type CustomHTTPClient struct {
	client *http.Client
}

type ResponseData struct {
	Allowed bool `json:"allowed"`
}

// Create Client with cerrt and private key
func Create() (*CustomHTTPClient, error) {

	// Load client certificate and key
	// TODO: fetch this cert from keyvault
	var clientCertPath = "./pkg/client/client.pem"
	var clientKeyPath = "./pkg/client/client_privk.pem"
	var caCertPath = "./pkg/client/servicecert.pem"

	clientCertPath, err := filepath.Abs(clientCertPath)
	if err != nil {
		return nil, err
	}

	clientKeyPath, err = filepath.Abs(clientKeyPath)
	if err != nil {
		return nil, err
	}

	caCertPath, err = filepath.Abs(caCertPath)
	if err != nil {
		return nil, err
	}

	cert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error loading client certificate and key: %v", err)
	}

	// This is a self signed cerrt
	caCert, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		return nil, fmt.Errorf("error loading Self certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create a TLS configuration with client certificate and CA certificate
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	// Create and return the HTTP client with the TLS configuration
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

return &CustomHTTPClient{client: client}, nil
}


func (c *CustomHTTPClient) AuthorizeAccess(path string, next http.HandlerFunc) http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {
	const baseUrl = "https://decentralized-rbac-1.com/app/"
	response, err := c.client.Get(baseUrl + path)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	if err != nil {
		return
	}

	// Read and use the response body as needed
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

		// Create a struct variable to hold the parsed JSON data
		var responseData ResponseData

		// Unmarshal the JSON data into the struct
		err = json.Unmarshal(body, &responseData)
		if err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			return
		}
	
		// Access the parsed data using the struct field
		fmt.Printf("Client Response Allowed: %v\n", responseData.Allowed)

		if responseData.Allowed {
			next.ServeHTTP(w, r)
		}else{
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
		}

	}
}

