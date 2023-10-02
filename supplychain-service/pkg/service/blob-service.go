package service

import (
    "context"
    "encoding/json"
    "fmt"
    //"log"
    "net/url"
    "github.com/Azure/azure-storage-blob-go/azblob"

    "supplychain-service/pkg/models"
)

// AzureBlobService handles interactions with Azure Blob Storage for immutability.
type AzureBlobService struct {
    containerURL azblob.ContainerURL
}

// NewAzureBlobService creates a new AzureBlobService instance.
func NewAzureBlobService(accountName, accountKey, containerName string) (*AzureBlobService, error) {
    credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
    if err != nil {
        return nil, err
    }
    // The URL typically looks like this:
    u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))

    p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
    serviceUrl:= azblob.NewServiceURL(*u, p)
    containerURL := serviceUrl.NewContainerURL(containerName)

    return &AzureBlobService{
        containerURL: containerURL,
    }, nil
}

// WriteBlob writes a coffee product event as an immutable JSON blob to Azure Blob Storage.
func (s *AzureBlobService) WriteBlob(product *models.CoffeeProduct, event string) error {
    blobName := fmt.Sprintf("%d_%s.json", product.ID, event)
    data, err := json.Marshal(product)
    if err != nil {
        return err
    }

    blobURL := s.containerURL.NewBlockBlobURL(blobName)
    ctx := context.Background()

    _, err = azblob.UploadBufferToBlockBlob(ctx, data, blobURL, azblob.UploadToBlockBlobOptions{
        BlobHTTPHeaders: azblob.BlobHTTPHeaders{ContentType: "application/json"},
    })
    if err != nil {
        return err
    }

    return nil
}
