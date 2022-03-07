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
	insertBatchWithConnectionString()
	queryBatchWithConnectionString()

}
func insertBatchWithConnectionString() {
	sc, err := aztables.NewServiceClientFromConnectionString("DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;TableEndpoint=http://127.0.0.1:10002/devstoreaccount1;", nil)
	handle(err)
	client := sc.NewClient("TestTable")

	// assuming table is not created
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

	entity2 := aztables.EDMEntity{
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

	var batch []aztables.TransactionAction

	e1, err := json.Marshal(entity)
	handle(err)

	e2, err := json.Marshal(entity2)
	handle(err)

	batch = append(batch, aztables.TransactionAction{
		ActionType: aztables.InsertMerge,
		Entity:     e1,
	})

	batch = append(batch, aztables.TransactionAction{
		ActionType: aztables.InsertMerge,
		Entity:     e2,
	})

	resp, err := client.SubmitTransaction(context.Background(), batch, nil)
	handle(err)
	fmt.Println(resp)
}

func queryBatchWithConnectionString() {
	sc, err := aztables.NewServiceClientFromConnectionString("DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;TableEndpoint=http://127.0.0.1:10002/devstoreaccount1;", nil)
	handle(err)
	client := sc.NewClient("TestTable")

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