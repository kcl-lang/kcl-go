package testdata

// Person Example
type Person struct {
	Name         string            `kcl:"name=name,type=str"`           // kcl-type: str
	Age          int               `kcl:"name=age,type=int"`            // kcl-type: int
	Friends      []string          `kcl:"name=friends,type=[str]"`      // kcl-type: [str]
	Movies       map[string]*Movie `kcl:"name=movies,type={str:Movie}"` // kcl-type: {str:Movie}
	MapInterface map[string]map[string]interface{}
	Ep           *Employee
	Com          Company
	StarInt      *int
	StarMap      map[string]string
	Inter        interface{}
}

type Movie struct {
	Desc     string      `kcl:"nam=desc,type=str"`                                   // kcl-type: str
	Size     int         `kcl:"name=size,type=units.NumberMultiplier"`               // kcl-type: units.NumberMultiplier
	Kind     string      `kcl:"name=kind,type=str(Superhero)|str(War)|str(Unknown)"` // kcl-type: "Superhero"|"War"|"Unknown"
	Unknown1 interface{} `kcl:"name=unknown1,type=int|str"`                          // kcl-type: int|str
	Unknown2 interface{} `kcl:"name=unknown2,type=any"`                              // kcl-type: any
}

type Employee struct {
	Name        string            `kcl:"name=name,type=str"`           // kcl-type: str
	Age         int               `kcl:"name=age,type=int"`            // kcl-type: int
	Friends     []string          `kcl:"name=friends,type=[str]"`      // kcl-type: [str]
	Movies      map[string]*Movie `kcl:"name=movies,type={str:Movie}"` // kcl-type: {str:Movie}
	BankCard    int               `kcl:"name=bankCard,type=int"`       // kcl-type: int
	Nationality string            `kcl:"name=nationality,type=str"`    // kcl-type: str
}

type Company struct {
	Name      string      `kcl:"name=name,type=str"`             // kcl-type: str
	Employees []*Employee `kcl:"name=employees,type=[Employee]"` // kcl-type: [Employee]
	Persons   *Person     `kcl:"name=persons,type=Person"`       // kcl-type: Person
}
