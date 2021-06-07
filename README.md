# Reproducing a bug in Azure Go's SDK with Azurite when inserting in batch mode to a non-existing table

[Github issue link](https://github.com/Azure/Azurite/issues/814)

For Go see [main.go](go-sdk/main.go):
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

For C# see [Program.cs](dotnet-sdk/Program.cs):
```csharp
private static void Insert(){
    // working - error is 404 TableNotFound, as expected
}

private static void InserBatcht(){
    // not working - returns System.IO.InvalidDataException: Invalid header line: HTTP/1.1 400 Bad Request
}
```