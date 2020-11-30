package data

import (
	"accountapi/lib"
	"encoding/json"
)

// RecordType is type of resource as defined in https://api-docs.form3.tech/api.html#audits-entries-record-types.
// This library only supports a subset of types, related to account management: "accounts" and "account_events".
type RecordType int

const (
	// RTNone is used for empty/omitted value, but still a valid RecordType.
	RTNone RecordType = iota
	// Accounts = "accounts".
	Accounts
	// AccountEvents = "account_events". TODO: this type is missing from the list in the documentation.
	AccountEvents
)

// IsValid ...
func (rt *RecordType) IsValid() bool {
	switch *rt {
	case RTNone, Accounts, AccountEvents:
		return true
	}
	return false
}

// String returns string name, panics if the RecordType value is not valid.
func (rt *RecordType) String() string {
	switch *rt {
	case RTNone:
		return ""
	case Accounts:
		return "accounts"
	case AccountEvents:
		return "account_events"
	}
	panic("String not implemented for this RecordType value")
}

// MarshalJSON converts values to strings.
func (rt *RecordType) MarshalJSON() ([]byte, error) {
	return json.Marshal(rt.String())
}

// UnmarshalJSON converts string value names into const values.
func (rt *RecordType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	t, err := recordTypeParse(s)
	if err != nil {
		return err
	}
	*rt = *t
	return nil
}

// parse convers string value names into const values, returns ErrorInvalidEnum if string is unknown.
func recordTypeParse(v string) (*RecordType, error) {
	rt := Accounts
	switch v {
	case "accounts":
		rt = Accounts
	case "account_events":
		rt = AccountEvents
	default:
		return nil, lib.NewErrorInvalidEnum()
	}
	return &rt, nil
}
