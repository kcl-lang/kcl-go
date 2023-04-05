package kclgen

import (
	"fmt"
	"testing"
)

type user struct {
	Name string `kcl:"name=name,type=str"`
	Age  int    `kcl:"name=age,type=int"`
}
type Person struct {
	FirstName string          `kcl:"name=firstName,type=str"`
	AMap      map[string]user `kcl:"name=aMap,type={str:user}"`
	FullName  string          `kcl:"name=fullName,type=str"`
	LastName  string          `kcl:"name=lastName,type=str"`
	Age       int             `kcl:"name=age,type=int"`
}
type employee struct {
	BankCard    int             `kcl:"name=bankCard,type=int"`
	Nationality string          `kcl:"name=nationality,type=str"`
	AMap        map[string]user `kcl:"name=aMap,type={str:user}"`
	Age         int             `kcl:"name=age,type=int"`
	FullName    string          `kcl:"name=fullName,type=str"`
	LastName    string          `kcl:"name=lastName,type=str"`
	FirstName   string          `kcl:"name=firstName,type=str"`
}
type Company struct {
	Name      string      `kcl:"name=name,type=str"`
	Persons   *Person     `kcl:"name=persons,type=schema"`
	Employees []*employee `kcl:"name=employees,type=[employee]"`
}

func TestExample(t *testing.T) {
	structList := make([]interface{}, 0)
	structList = append(structList, &user{})
	structList = append(structList, &Person{})
	structList = append(structList, &employee{})
	structList = append(structList, &Company{})
	for _, sl := range structList {
		s := GenKclSchemaCode(sl)
		fmt.Println(s)
	}
}

func TestGenKclSchemaCode(t *testing.T) {
	type args struct {
		s interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"testUser",
			args{s: &user{}},
			"schema user:\n    name: str\n    age: int\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenKclSchemaCode(tt.args.s); got != tt.want {
				t.Errorf("GenKclSchemaCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
