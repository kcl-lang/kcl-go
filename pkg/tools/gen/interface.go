package gen

import (
	"io"

	pb "kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
)

type Generator interface {
	GenFromSource(w io.Writer, filename string, src any) error
	GenFromTypes(w io.Writer, types ...*pb.KclType)
	GenSchema(w io.Writer, typ *pb.KclType)
	GetTypeName(typ *pb.KclType) string
}
