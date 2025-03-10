// Code generated by jsonenums -type=Unit; DO NOT EDIT.

package meals

import (
	"encoding/json"
	"fmt"
)

var (
	_UnitNameToValue = map[string]Unit{
		"Number": Number,
		"Tsp":    Tsp,
		"Tbsp":   Tbsp,
		"Cup":    Cup,
		"Oz":     Oz,
		"Lb":     Lb,
		"Gram":   Gram,
		"Kg":     Kg,
	}

	_UnitValueToName = map[Unit]string{
		Number: "Number",
		Tsp:    "Tsp",
		Tbsp:   "Tbsp",
		Cup:    "Cup",
		Oz:     "Oz",
		Lb:     "Lb",
		Gram:   "Gram",
		Kg:     "Kg",
	}
)

func init() {
	var v Unit
	if _, ok := interface{}(v).(fmt.Stringer); ok {
		_UnitNameToValue = map[string]Unit{
			interface{}(Number).(fmt.Stringer).String(): Number,
			interface{}(Tsp).(fmt.Stringer).String():    Tsp,
			interface{}(Tbsp).(fmt.Stringer).String():   Tbsp,
			interface{}(Cup).(fmt.Stringer).String():    Cup,
			interface{}(Oz).(fmt.Stringer).String():     Oz,
			interface{}(Lb).(fmt.Stringer).String():     Lb,
			interface{}(Gram).(fmt.Stringer).String():   Gram,
			interface{}(Kg).(fmt.Stringer).String():     Kg,
		}
	}
}

// MarshalJSON is generated so Unit satisfies json.Marshaler.
func (r Unit) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(r).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := _UnitValueToName[r]
	if !ok {
		return nil, fmt.Errorf("invalid Unit: %d", r)
	}
	return json.Marshal(s)
}

// UnmarshalJSON is generated so Unit satisfies json.Unmarshaler.
func (r *Unit) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Unit should be a string, got %s", data)
	}
	v, ok := _UnitNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid Unit %q", s)
	}
	*r = v
	return nil
}
