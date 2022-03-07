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
            InsertBatch();
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
                new(partitionKey, "rkey1")
                {
                    {"product", "product1"},
                },
                new(partitionKey, "rkey2")
                {
                    {"product", "product2"},
                }
            };

            // Create the batch.
            var addEntitiesBatch = new List<TableTransactionAction>();

            // Add the entities to be added to the batch.
            addEntitiesBatch.AddRange(entityList.Select(e =>
                new TableTransactionAction(TableTransactionActionType.Add, e)));

            // Create table beforehand
            client.Create();

            try
            {
                // Submit the batch.
                client.SubmitTransaction(addEntitiesBatch);
            }
            catch (RequestFailedException e)
                when (e.Status == 404 && e.ErrorCode == "TableNotFound")
            {
                // create table 
                client.Create();

                // retry
                client.SubmitTransaction(addEntitiesBatch);
            }

            Console.WriteLine("Inserted!");

            var entities = client.Query<TableEntity>();

            Console.WriteLine($"Number of records: {entities.Count()}");
            foreach (var entity in entities)
            {
                Console.WriteLine($"{entity.PartitionKey}:{entity.RowKey}");
            }
        }
    }
}