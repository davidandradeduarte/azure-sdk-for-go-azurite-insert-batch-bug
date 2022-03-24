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
	insertBatch()
	query()
}

func add() {
	sc, err := aztables.NewServiceClientFromConnectionString("DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;TableEndpoint=http://127.0.0.1:10002/devstoreaccount1;", nil)
	handleErr(err)
	client := sc.NewClient("TestTable")

	_, err = client.Create(context.Background(), nil)
	if !tableExists(err) {
		handleErr(err)
	}

	for _, v := range generateEntities() {
		resp, err := client.AddEntity(context.Background(), v, nil)
		handleErr(err)
		fmt.Println(resp)
	}
}

func insertBatch() {
	sc, err := aztables.NewServiceClientFromConnectionString("DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;TableEndpoint=http://127.0.0.1:10002/devstoreaccount1;", nil)
	handleErr(err)
	client := sc.NewClient("TestTable")

	_, err = client.Create(context.Background(), nil)
	if !tableExists(err) {
		handleErr(err)
	}

	var batch []aztables.TransactionAction
	for _, v := range generateEntities() {
		batch = append(batch, aztables.TransactionAction{
			ActionType: aztables.Add,
			Entity:     v,
		})
	}

	resp, err := client.SubmitTransaction(context.Background(), batch, nil)
	handleErr(err)
	fmt.Println(resp)
}

func query() {
	sc, err := aztables.NewServiceClientFromConnectionString("DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;TableEndpoint=http://127.0.0.1:10002/devstoreaccount1;", nil)
	handleErr(err)
	client := sc.NewClient("TestTable")

	options := &aztables.ListEntitiesOptions{
		Select: to.StringPtr("RowKey,Value,Product,Available"),
		Top:    to.Int32Ptr(15),
	}

	pager := client.List(options)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		for _, entity := range resp.Entities {
			var myEntity aztables.EDMEntity
			err = json.Unmarshal(entity, &myEntity)
			handleErr(err)
			fmt.Println(myEntity)
		}
	}

	err = pager.Err()
	handleErr(err)
}

func generateEntities() [][]byte {
	uuid := uuid.New().String()

	var entities []aztables.EDMEntity
	entities = append(entities, aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: uuid,
			RowKey:       "rkey1",
		},
		Properties: map[string]interface{}{
			"product": "product1",
		},
	})

	entities = append(entities, aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: uuid,
			RowKey:       "rkey2",
		},
		Properties: map[string]interface{}{
			"product": "product2",
		},
	})

	var response [][]byte
	for _, v := range entities {
		e, err := json.Marshal(v)
		handleErr(err)
		response = append(response, e)
	}
	return response
}

func handleErr(err error) {
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
