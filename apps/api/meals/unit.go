package meals

type Unit int

//go:generate go run github.com/campoy/jsonenums -type=Unit
const (
	Number Unit = iota
	Tsp
	Tbsp
	Cup
	Oz
	Lb
	Gram
	Kg
)

func UnitFromString(s string) (Unit, bool) {
	v, ok := _UnitNameToValue[s]
	return v, ok
}
