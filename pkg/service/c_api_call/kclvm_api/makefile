default:
	cargo build --release
	
darwin-amd64:
	cargo build --release 
	cp ./target/release/libkclvm_api.dylib ../packaged/lib/darwin-amd64
	
darwin-arm64:
	cargo build --release 
	cp ./target/release/libkclvm_api.dylib ../packaged/lib/darwin-arm64

linux-amd64:
	cargo build --release 
	cp ./target/release/libkclvm_api.so ../packaged/lib/linux-amd64

linux-arm64:
	cargo build --release 
	cp ./target/release/libkclvm_api.so ../packaged/lib/linux-arm64

clean:
	rm -rf ./target