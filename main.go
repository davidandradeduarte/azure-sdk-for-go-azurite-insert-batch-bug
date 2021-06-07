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
	tablename    = "TestTable"
)

func main() {
	insertBatch("1", "3")
	insert("1", "3")
	//query()
}

func insertBatch(partitionkey string, rowkey string) {
	client, err := storage.NewEmulatorClient()

	if err != nil {
		fmt.Printf("%s: \n", err)
	}

	ts = client.GetTableService()

	t := ts.GetTableReference(tablename)

	entity := t.GetEntityReference(partitionkey, rowkey)

	props := map[string]interface{}{
		"AmountDue":      200.23,
		"CustomerCode":   "123",
		"CustomerSince":  time.Date(1992, time.December, 20, 21, 55, 0, 0, time.UTC),
		"IsActive":       true,
		"NumberOfOrders": int64(255),
	}

	entity.Properties = props

	tb := t.NewBatch()

	tb.InsertOrMergeEntity(entity, true)

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
		return
	} else {
		fmt.Println("Inserted! ")
		return
	}
}

func insert(partitionkey string, rowkey string) {
	client, err := storage.NewEmulatorClient()

	if err != nil {
		fmt.Printf("%s: \n", err)
	}

	ts = client.GetTableService()

	t := ts.GetTableReference(tablename)

	entity := t.GetEntityReference(partitionkey, rowkey)

	props := map[string]interface{}{
		"AmountDue":      200.23,
		"CustomerCode":   "123",
		"CustomerSince":  time.Date(1992, time.December, 20, 21, 55, 0, 0, time.UTC),
		"IsActive":       true,
		"NumberOfOrders": int64(255),
	}

	entity.Properties = props

	err = entity.Insert(fullmetadata, nil)

	if cerr, ok := err.(storage.AzureStorageServiceError); ok {
		if cerr.Code == "TableNotFound" {
			if cerr := t.Create(uint(10), storage.FullMetadata, nil); cerr != nil {
				fmt.Printf("error creating table: %v.", cerr)
				return
			}
			// retry
			err = entity.Insert(fullmetadata, nil)
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

	ts = client.GetTableService()

	table := ts.GetTableReference(tablename)

	// timeout, metatadalevel, options
	entities, err := table.QueryEntities(30, fullmetadata, nil)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(entities.Entities[0])
	}
}
