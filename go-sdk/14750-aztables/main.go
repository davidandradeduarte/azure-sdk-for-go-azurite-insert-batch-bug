package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/google/uuid"
)

func main() {
	// add()
	// insertBatch()
	query()
}

func add() {
	sc, err := aztables.NewServiceClientFromConnectionString("DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;TableEndpoint=http://127.0.0.1:10002/devstoreaccount1;", nil)
	handle(err)
	client := sc.NewClient("TestTable")

	// assuming table is not created
	_, err = client.Create(context.Background(), nil)
	if !tableExists(err) {
		handle(err)
	}

	uuid := generateUuid()

	entity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: uuid,
			RowKey:       "rkey1",
		},
		Properties: map[string]interface{}{
			"product": "product1",
			"price":   5.00,
		},
	}

	entity2 := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: uuid,
			RowKey:       "rkey2",
		},
		Properties: map[string]interface{}{
			"product": "product1",
			"price":   5.00,
		},
	}

	e1, err := json.Marshal(entity)
	handle(err)

	resp, err := client.AddEntity(context.Background(), e1, nil)
	handle(err)
	fmt.Println(resp)

	e2, err := json.Marshal(entity2)
	handle(err)

	resp, err = client.AddEntity(context.Background(), e2, nil)
	handle(err)
	fmt.Println(resp)
}

func insertBatch() {
	sc, err := aztables.NewServiceClientFromConnectionString("DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;TableEndpoint=http://127.0.0.1:10002/devstoreaccount1;", nil)
	handle(err)
	client := sc.NewClient("TestTable")

	// assuming table is not created
	_, err = client.Create(context.Background(), nil)
	if !tableExists(err) {
		handle(err)
	}

	uuid := generateUuid()

	entity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: uuid,
			RowKey:       "rkey1",
		},
		Properties: map[string]interface{}{
			"product": "product1",
			"price":   5.00,
		},
	}

	entity2 := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: uuid,
			RowKey:       "rkey2",
		},
		Properties: map[string]interface{}{
			"product": "product2",
			"price":   10.00,
		},
	}

	var batch []aztables.TransactionAction

	e1, err := json.Marshal(entity)
	handle(err)

	e2, err := json.Marshal(entity2)
	handle(err)

	batch = append(batch, aztables.TransactionAction{
		ActionType: aztables.Add,
		Entity:     e1,
	})

	batch = append(batch, aztables.TransactionAction{
		ActionType: aztables.Add,
		Entity:     e2,
	})

	resp, err := client.SubmitTransaction(context.Background(), batch, nil)
	handle(err)
	fmt.Println(resp)
}

func query() {
	sc, err := aztables.NewServiceClientFromConnectionString("DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;TableEndpoint=http://127.0.0.1:10002/devstoreaccount1;", nil)
	handle(err)
	client := sc.NewClient("TestTable")

	options := &aztables.ListEntitiesOptions{
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
			fmt.Println(myEntity)
		}
	}

	err = pager.Err()
	handle(err)
}

func generateUuid() string {
	id := uuid.New()
	return id.String()
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func tableExists(err error) bool {
	if err == nil {
		return false
	}
	var azErr *azcore.ResponseError
	if errors.As(err, &azErr) {
		return azErr.StatusCode == http.StatusConflict
	}
	return false
}
