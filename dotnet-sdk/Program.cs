using System;
using System.Collections.Generic;
using System.Linq;
using Azure;
using Azure.Data.Tables;

namespace play
{
    class Program
    {
        static void Main(string[] args)
        {
            Insert();
            InsertBatch();
        }
        
        private static void Insert()
        {
            const string tableName = "InsertTestTable";
            var client = new TableClient(
                "UseDevelopmentStorage=true;",
                tableName);

            const string partitionKey = "InsertSample";
            var entity = new TableEntity(partitionKey, "01")
            {
                { "Product", "Marker Set" },
                { "Price", 5.00 },
                { "Quantity", 21 }
            };

            try
            {
                client.AddEntity(entity);
            }
            catch (RequestFailedException e)
            when(e.Status == 404 && e.ErrorCode == "TableNotFound")
            {
                // create table 
                client.Create();
                
                // retry
                client.AddEntity(entity);
            }
            
            Console.WriteLine("Inserted! ");
        }
        
        private static void InsertBatch()
        {
            const string tableName = "InsertBatchTestTable";
            var client = new TableClient(
                "UseDevelopmentStorage=true;",
                tableName);

            const string partitionKey = "BatchInsertSample";
            var entityList = new List<TableEntity>
            {
                new(partitionKey, "01")
                {
                    {"Product", "Marker"},
                    {"Price", 5.00},
                    {"Brand", "Premium"}
                },
                new(partitionKey, "02")
                {
                    {"Product", "Pen"},
                    {"Price", 3.00},
                    {"Brand", "Premium"}
                },
                new(partitionKey, "03")
                {
                    {"Product", "Paper"},
                    {"Price", 0.10},
                    {"Brand", "Premium"}
                },
                new(partitionKey, "04")
                {
                    {"Product", "Glue"},
                    {"Price", 1.00},
                    {"Brand", "Generic"}
                },
            };

            // Create the batch.
            var addEntitiesBatch = new List<TableTransactionAction>();

            // Add the entities to be added to the batch.
            addEntitiesBatch.AddRange(entityList.Select(e =>
                new TableTransactionAction(TableTransactionActionType.Add, e)));

            try
            {
                // Submit the batch.
                client.SubmitTransaction(addEntitiesBatch);
            }
            catch (RequestFailedException e)
                when(e.Status == 404 && e.ErrorCode == "TableNotFound")
            {
                // create table 
                client.Create();
                
                // retry
                client.SubmitTransaction(addEntitiesBatch);
            }
            
            Console.WriteLine("Inserted! ");
        }
    }
}