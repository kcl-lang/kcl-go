// Copyright 2023 The KCL Authors. All rights reserved.

package gen_test

import (
	"fmt"
	"log"
	"testing"

	"kusionstack.io/kclvm-go/pkg/tools/gen"
)

var _ = gen.GenProto

func TestGenProto(t *testing.T) {
	const code = `
import units

#kclvm/genpb: option go_package = kcl_gen/_/hello
#kclvm/genpb: option pb_package = kcl_gen._.hello

type NumberMultiplier = units.NumberMultiplier

schema Person:
	"""Person Example"""
	name: str = "kcl"
	age: int = 2
	friends?: [str] = None
	movies?: {str: Movie} = None

schema Movie:
	desc: str = ""
	size: NumberMultiplier = 2M
	kind?: "Superhero" | "War" | "Unknown"
	unknown1?: int | str = None
	unknown2?: any = None
`

	pbCode, err := gen.GenProto("hello.k", code, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pbCode)

	// Output:
	// syntax = "proto3";
	//
	// package kcl_gen._.hello;
	//
	// option go_package = "kcl_gen/_/hello";
	//
	// import "google/protobuf/any.proto";
	//
	// // Person Example
	// message Person {
	//     string name = 1;               // kcl-type: str
	//     int64 age = 2;                 // kcl-type: int
	//     repeated string friends = 3;   // kcl-type: [str]
	//     map<string, Movie> movies = 4; // kcl-type: {str:Movie}
	// }
	//
	// message Movie {
	//     string desc = 1;                  // kcl-type: str
	//     int64 size = 2;                   // kcl-type: units.NumberMultiplier
	//     string kind = 3;                  // kcl-type: "Superhero"|"War"|"Unknown"
	//     google.protobuf.Any unknown1 = 4; // kcl-type: int|str
	//     google.protobuf.Any unknown2 = 5; // kcl-type: any
	// }
}
