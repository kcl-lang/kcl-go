# https://github.com/princjef/gomarkdoc
# go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest

default:
	which kclvm
	kclvm -m kclvm --version

	go run ./cmds/kcl-go
	go run ./cmds/kcl-go run hello.k

doc:
	gomarkdoc . > doc.md

clean:
