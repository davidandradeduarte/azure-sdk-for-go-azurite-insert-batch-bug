# Reproducing a bug in Azure Go's SDK with Azurite when inserting in batch mode to a non-existing table

[Github issue link](https://github.com/Azure/azure-sdk-for-go/issues/14746)

See [main.go](main.go):

```golang
func main(){
    insert()
    insertBatch()
}

func insert(){
    // working, status code is TableNotFound as expected
}

func insertBatch(){
    // not working, returns a EOF string error
}
```