package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/Azure/go-autorest/autorest/to"
)

func main() {
	// insertBatch()
	// query()
	execWithConnectionString()
}

func insertBatch() {
	client, err := aztables.NewClientWithNoCredential("http://localhost:10002/TestTable", nil)
	handle(err)
	_, err = client.Create(context.Background(), nil)
	handle(err)

	entity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: "pencils",
			RowKey:       "id-003",
		},
		Properties: map[string]interface{}{
			"Product":      "Ticonderoga Pencils",
			"Price":        5.00,
			"Count":        aztables.EDMInt64(12345678901234),
			"ProductGUID":  aztables.EDMGUID("some-guid-value"),
			"DateReceived": aztables.EDMDateTime(time.Now()),
			"ProductCode":  aztables.EDMBinary([]byte("somebinaryvalue")),
		},
	}

	data, err := json.Marshal(entity)
	handle(err)

	_, err = client.AddEntity(context.Background(), data, nil)
	handle(err)
}

func query() {
	client, err := aztables.NewClientWithNoCredential("http://0.0.0.0:10002", nil)
	handle(err)
	filter := "PartitionKey eq 'markers' or RowKey eq 'id-003'"
	options := &aztables.ListEntitiesOptions{
		Filter: &filter,
		Select: to.StringPtr("RowKey,Value,Product,Available"),
		Top:    to.Int32Ptr(15),
	}

	pager := client.List(options)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		fmt.Printf("Received: %v entitiesn", len(resp.Entities))

		for _, entity := range resp.Entities {
			var myEntity aztables.EDMEntity
			err = json.Unmarshal(entity, &myEntity)
			handle(err)

			fmt.Printf("Received: %v, %v, %v, %vn", myEntity.Properties["RowKey"], myEntity.Properties["Value"], myEntity.Properties["Product"], myEntity.Properties["Available"])
		}
	}

	err = pager.Err()
	handle(err)
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}
