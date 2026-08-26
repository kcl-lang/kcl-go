package basic

// Person is a human being.
type Person struct {
	Name string `kcl:"name=name,type=str"`
	Age  int    `kcl:"name=age,type=int"`
}
