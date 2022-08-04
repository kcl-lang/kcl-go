### Build Kclvm C Api Lib
Choose make targe by your OS and Computer Architecture. Eg darwin-amd64
```
$ cd kclvm_api
$ make darwin-amd64
$ make clean
```

### Call Kclvm Service by C API
enable Kclvm Service C API by environment variable
```
$ export KCLVM_SERVICE_CLIENT_HANDLER=native

```

enable Kclvm Service C API by source
```
client := service.NewKclvmServiceClient()
client.IsNative = true
..........
result ,err :=client.ExecProgram(args)
```