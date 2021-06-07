# Reproducing two bugs in Azurite when inserting in batch mode

## Batch insert into a non-existing table returns unexpected error

[**Github issue**](https://github.com/Azure/Azurite/issues/814)

For Go see [main.go](go-sdk/814/main.go):

```golang
func main(){
    insert()
    insertBatch()
}

func insert(){
    // working - error is 404 TableNotFound, as expected
}

func insertBatch(){
    // not working - returns a EOF string error
}
```

For C# see [Program.cs](dotnet-sdk/814/Program.cs):

```csharp
static void Main(string[] args)
{
    Insert();
    InsertBatch();
}

private static void Insert(){
    // working - error is 404 TableNotFound, as expected
}

private static void InsertBatch(){
    // not working - returns System.IO.InvalidDataException: Invalid header line: HTTP/1.1 400 Bad Request
}
```

## Batch insert with Go SDK (using Azurite) only inserts one record

[**Github issue**](https://github.com/Azure/azure-sdk-for-go/issues/14750)

For Go see [main.go](go-sdk/14750/main.go):

```golang
func main(){
    insertBatch()
}

func insertBatch(){
    // not working - inserts only one of the two entities added to the batch
}
```

For C# see [Program.cs](dotnet-sdk/14750/Program.cs):

```csharp
static void Main(string[] args)
{
    InsertBatch();
}

private static void InsertBatch(){
    // working - inserts multiple records in batch mode
}
```
