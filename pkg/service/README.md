### Build Kclvm C Api Lib
Choose make targe by your OS and Computer Architecture. Eg darwin-amd64
```
$make darwin-amd64
$make clean
```

### Call Kclvm Service by C API
```
client := service.NewKclvmServiceClient()
client.IsNative = true //set native flag
..........
result ,err :=client.ExecProgram(args)
```