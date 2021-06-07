package main

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/storage"
)

var (
	ts storage.TableServiceClient
)

const (
	fullmetadata = "application/json;odata=fullmetadata"
)

func main() {
	insertBatch("1", "3")
	query()
}

func insertBatch(partitionkey string, rowkey string) {
	client, err := storage.NewEmulatorClient()

	if err != nil {
		fmt.Printf("%s: \n", err)
	}

	ts = client.GetTableService()

	t := ts.GetTableReference("InsertBatchTestTable")

	entity := t.GetEntityReference(partitionkey, rowkey)

	props := map[string]interface{}{
		"AmountDue":      200.23,
		"CustomerCode":   "123",
		"CustomerSince":  time.Date(1992, time.December, 20, 21, 55, 0, 0, time.UTC),
		"IsActive":       true,
		"NumberOfOrders": int64(255),
	}

	entity.Properties = props

	entity2 := t.GetEntityReference("foo", "bar")

	props2 := map[string]interface{}{
		"AmountDue":      100,
		"CustomerCode":   "111",
		"CustomerSince":  time.Date(2020, time.December, 20, 21, 55, 0, 0, time.UTC),
		"IsActive":       false,
		"NumberOfOrders": int64(20),
	}

	entity2.Properties = props2

	// create table beforehand
	if cerr := t.Create(uint(10), storage.FullMetadata, nil); cerr != nil {
		fmt.Printf("error creating table: %v.", cerr)
		return
	}

	tb := t.NewBatch()

	tb.InsertOrMergeEntity(entity, true)
	tb.InsertOrMergeEntity(entity2, true)

	if err := tb.ExecuteBatch(); err != nil {
		if cerr, ok := err.(storage.AzureStorageServiceError); ok {
			if cerr.Code == "TableNotFound" {
				if cerr := t.Create(uint(10), storage.FullMetadata, nil); cerr != nil {
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

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Inserted!")
	}
}

func query() {
	client, err := storage.NewEmulatorClient()

	if err != nil {
		fmt.Printf("%s: \n", err)
	}

	t := client.GetTableService()

	table := t.GetTableReference("InsertBatchTestTable")

	// timeout, metatadalevel, options
	entities, err := table.QueryEntities(30, fullmetadata, nil)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Number of records: %d\n", len(entities.Entities))
		for _, v := range entities.Entities {
			fmt.Println(v)
		}
	}

	err = table.Delete(10, &storage.TableOptions{})
	if err != nil {
		fmt.Println(err)
	}
}
