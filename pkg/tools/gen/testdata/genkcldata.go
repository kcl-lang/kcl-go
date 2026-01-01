package testdata

// Person Example
type Person struct {
	Name         string            `kcl:"name=name,type=str"`           // kcl-type: str
	Age          int               `kcl:"name=age,type=int"`            // kcl-type: int
	Friends      []string          `kcl:"name=friends,type=[str]"`      // kcl-type: [str]
	Movies       map[string]*Movie `kcl:"name=movies,type={str:Movie}"` // kcl-type: {str:Movie}
	MapInterface map[string]map[string]any
	Ep           *Employee
	Com          Company
	StarInt      *int
	StarMap      map[string]string
	Inter        any
}

type Movie struct {
	Desc     string      `kcl:"nam=desc,type=str"`                                   // kcl-type: str
	Size     int         `kcl:"name=size,type=units.NumberMultiplier"`               // kcl-type: units.NumberMultiplier
	Kind     string      `kcl:"name=kind,type=str(Superhero)|str(War)|str(Unknown)"` // kcl-type: "Superhero"|"War"|"Unknown"
	Unknown1 any `kcl:"name=unknown1,type=int|str"`                          // kcl-type: int|str
	Unknown2 any `kcl:"name=unknown2,type=any"`                              // kcl-type: any
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

// Test inline tag behavior
type TypeMeta struct {
	APIVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
}

type ObjectMeta struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type AppInline struct {
	TypeMeta   `json:",inline,omitempty"`
	ObjectMeta `json:"metadata,omitempty"`
}

type AppNoInline struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata"`
}
