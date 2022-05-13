default: clean
	go run 0_gen.go > zz_gen.go.txt
	mv zz_gen.go.txt zz_gen.go
	go fmt

	go test

clean:
	-rm zz_*
