### Build Kclvm C Api Lib
Choose make targe by your OS and Computer Architecture. Eg darwin-amd64
```
$ cd ../../kclvm_capi
$ make darwin-amd64
$ make clean
```

### Call Kclvm Service by C API
First, make sure CGO is enabled and use build tag kclvm_service_capi
```
$ export CGO_ENABLED=1
$ go build -tags kclvm_service_capi
$ go test -tags kclvm_service_capi ./...
```


enable Kclvm Service C API by environment variable
```
$ export KCLVM_SERVICE_CLIENT_HANDLER=native

```

enable Kclvm Service C API by source
```
service.Default_IsNative=true
client := service.NewKclvmServiceClient()
..........
result ,err :=client.ExecProgram(args)
```