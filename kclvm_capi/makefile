KCLVM_CAPI_LIB_PATH = ../pkg/service/c_api_call/packaged/lib

default:
	cargo build --release
	
darwin-amd64:
	cargo build --release 
	-rm ${KCLVM_CAPI_LIB_PATH}/darwin-amd64/*.dylib
	cp ./target/release/libkclvm_capi.dylib ${KCLVM_CAPI_LIB_PATH}/darwin-amd64
	
darwin-arm64:
	cargo build --release 
	-rm ${KCLVM_CAPI_LIB_PATH}/darwin-arm64/*.dylib
	cp ./target/release/libkclvm_capi.dylib ${KCLVM_CAPI_LIB_PATH}/darwin-arm64

linux-amd64:
	cargo build --release 
	-rm ${KCLVM_CAPI_LIB_PATH}/linux-amd64/*.so
	cp ./target/release/libkclvm_capi.so ${KCLVM_CAPI_LIB_PATH}/linux-amd64

linux-arm64:
	cargo build --release 
	-rm ${KCLVM_CAPI_LIB_PATH}/linux-arm64/*.so
	cp ./target/release/libkclvm_capi.so ${KCLVM_CAPI_LIB_PATH}/linux-arm64

clean:
	rm -rf ./target