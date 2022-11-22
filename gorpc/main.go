package main

import (
	"io"
	"log"
	"net/rpc"
	"os"

	"github.com/chai2010/protorpc"
	capicall "kusionstack.io/kclvm-go/gorpc/pkg/c_api_call"
)

type rwCloser struct {
	io.ReadCloser
	io.WriteCloser
}

func (rw rwCloser) Close() error {
	err := rw.ReadCloser.Close()
	if err := rw.WriteCloser.Close(); err != nil {
		return err
	}
	return err
}

func main() {
	c := capicall.PROTOCAPI_NewKclvmServiceClient()
	srv := rpc.NewServer()
	if err := srv.RegisterName("KclvmService", c); err != nil {
		log.Fatal(err)
	}
	if err := srv.RegisterName("BuiltinService", c); err != nil {
		log.Fatal(err)
	}

	srv.ServeCodec(protorpc.NewServerCodec(rwCloser{os.Stdin, os.Stdout}))
}
