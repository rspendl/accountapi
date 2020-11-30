package data

import (
	"accountapi/lib"
	"encoding/json"
)

// AccountStatus can be "Personal" or "Business".
type AccountStatus int

const (
	// Confirmed = "Confirmed".
	Confirmed AccountStatus = iota
	// Pending = "Pending".
	Pending
	// Failed = "Failed".
	Failed
)

// IsValid ...
func (as AccountStatus) IsValid() bool {
	switch as {
	case Confirmed, Pending, Failed:
		return true
	}
	return false
}

// String ...
func (as AccountStatus) String() string {
	switch as {
	case Confirmed:
		return "confirmed"
	case Pending:
		return "pending"
	case Failed:
		return "failed"
	}
	panic("String not implemented for this AccountStatus value")
}

// MarshalJSON converts values to strings.
func (as *AccountStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(as.String())
}

// UnmarshalJSON converts string value names into const values.
func (as *AccountStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	st, err := accountStatusParse(s)
	if err != nil {
		return err
	}
	*as = *st
	return nil
}

// parse convers string value names into const values, returns ErrorInvalidEnum if string is unknown.
func accountStatusParse(v string) (*AccountStatus, error) {
	as := Confirmed
	switch v {
	case "confirmed":
		as = Confirmed
	case "pending":
		as = Pending
	case "failed":
		as = Failed
	default:
		return nil, lib.NewErrorInvalidEnum()
	}
	return &as, nil
}
