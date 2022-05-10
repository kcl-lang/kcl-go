# https://github.com/princjef/gomarkdoc
# go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest

KCLVM_URL_MIRRORS:=http://127.0.0.1:8000/downloads

default:
	which kclvm
	kclvm -m kclvm --version

	go run ./cmds/kcl-go
	go run ./cmds/kcl-go run hello.k

setup-kclvm-all:
	go run ./cmds/kcl-go/ setup-kclvm -all -outdir=_build  -mirrors=${KCLVM_URL_MIRRORS}

clean:
	-rm -rf ./_build
