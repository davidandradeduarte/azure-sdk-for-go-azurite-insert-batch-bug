package main

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/storage"
)

var (
	tableCli storage.TableServiceClient
)

const (
	fullmetadata = "application/json;odata=fullmetadata"
	tablename    = "TestTable"
)

func main() {
	insert("1", "3")
	query()
}

func insert(partitionkey string, rowkey string) {

	client, err := storage.NewEmulatorClient()

	if err != nil {
		fmt.Printf("%s: \n", err)
	}

	tableCli = client.GetTableService()

	fmt.Println(tableCli)

	table := tableCli.GetTableReference(tablename)

	entity := table.GetEntityReference(partitionkey, rowkey)

	props := map[string]interface{}{
		"AmountDue":      200.23,
		"CustomerCode":   "123",
		"CustomerSince":  time.Date(1992, time.December, 20, 21, 55, 0, 0, time.UTC),
		"IsActive":       true,
		"NumberOfOrders": int64(255),
	}

	entity.Properties = props

	tb := table.NewBatch()

	tb.InsertOrMergeEntity(entity, true)

	if err := tb.ExecuteBatch(); err != nil {
		if cerr, ok := err.(storage.AzureStorageServiceError); ok {
			if cerr.Code == "TableNotFound" {
				if cerr := table.Create(uint(10), storage.FullMetadata, nil); cerr != nil {
					fmt.Printf("error creating table: %v.", cerr)
					return
				}
				// retry
				err = tb.ExecuteBatch()
			}
		}
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// err = entity.Insert(fullmetadata, nil)

	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Inserted! ")
		return
	}
}

func query() {

	client, err := storage.NewEmulatorClient()

	if err != nil {
		fmt.Printf("%s: \n", err)
	}

	tableCli = client.GetTableService()

	fmt.Println(tableCli)

	table := tableCli.GetTableReference(tablename)

	// timeout, metatadalevel, options
	entities, err := table.QueryEntities(30, fullmetadata, nil)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(entities.Entities[0])
	}
}
