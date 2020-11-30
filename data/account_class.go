package data

import (
	"accountapi/lib"
	"encoding/json"
)

// AccountClass can be "Personal" or "Business".
type AccountClass int

const (
	// Personal = "Personal".
	Personal AccountClass = iota
	// Business = "Business".
	Business
)

// IsValid ...
func (ac *AccountClass) IsValid() bool {
	switch *ac {
	case Personal, Business:
		return true
	}
	return false
}

func (ac *AccountClass) String() string {
	switch *ac {
	case Personal:
		return "Personal"
	case Business:
		return "Business"
	}
	panic("String not implemented for this AccountClass value")
}

// MarshalJSON converts values to strings.
func (ac *AccountClass) MarshalJSON() ([]byte, error) {
	return json.Marshal(ac.String())
}

// UnmarshalJSON converts string value names into const values.
func (ac *AccountClass) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	c, err := accountClassParse(s)
	if err != nil {
		return err
	}
	*ac = *c
	return nil
}

// parse convers string value names into const values, returns ErrorInvalidEnum if string is unknown.
func accountClassParse(v string) (*AccountClass, error) {
	ac := Personal
	switch v {
	case "Personal":
		ac = Personal
	case "Business":
		ac = Business
	default:
		return nil, lib.NewErrorInvalidEnum()
	}
	return &ac, nil
}
