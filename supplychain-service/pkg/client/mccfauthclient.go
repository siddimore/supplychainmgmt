package mccfclient


import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Create Client with cerrt and private key
func Create() (*http.Client, error) {
	// Load client certificate and key
	// TODO: fetch this cert from keyvault
	// Load Sample Cert
	cert, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		return nil, fmt.Errorf("error loading client certificate and key: %v", err)
	}

	// This is a self signed cerrt
	caCert, err := ioutil.ReadFile("client.crt")
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
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}, nil
}

func AuthorizeAccess(client *http.Client, url string, requestBody []byte) ([]byte, error) {
	// Send the HTTP POST request
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error sending POST request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return responseBody, nil
}


