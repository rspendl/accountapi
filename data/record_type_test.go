package data_test

import (
	"accountapi/data"
	"encoding/json"
	"strings"
	"testing"
)

// TestRecordType for testing unmarshalling JSON values.
type TestRecordType struct {
	TestType data.RecordType `json:"testType"`
}

// TestValidRecordType verifies proper constraints for consts ("enums"), parsing and unmarshalling.
func TestValidRecordType(t *testing.T) {
	rt := data.Accounts
	if !rt.IsValid() {
		t.Error("RecordType Accounts should be valid.")
		t.Fail()
	}
	if rt.String() != "accounts" {
		t.Error("Accounts string should be \"accounts\".")
		t.Fail()
	}

	jString := `{"testType":"accounts"}`
	jStruct := TestRecordType{}
	err := json.NewDecoder(strings.NewReader(jString)).Decode(&jStruct)
	if err != nil {
		t.Errorf("Can't unmarshal RecordType: %s", err.Error())
		t.Fail()
	}
	if jStruct.TestType != rt {
		t.Errorf("Expected RecordType value: '%s', got: '%s'", rt.String(), jStruct.TestType.String())
		t.Fail()
	}
	b, err := json.Marshal(&jStruct)
	if err != nil {
		t.Errorf("Can't marshal RecordType to string: %s\n", err.Error())
		t.Fail()
	} else if string(b) != jString {
		t.Errorf("Expected marshalled value: '%s', got: '%s'\n", jString, string(b))
		t.Fail()
	}

}

// TestInvalidRecordType verifies response of functions when called with invalid proper constraints for consts ("enums"), parsing and unmarshalling.
func TestInvalidRecordType(t *testing.T) {
	rt := data.AccountEvents // The last RecordType value, when it's increased it should become an invalid value.
	rt++
	if rt.IsValid() {
		t.Error("Invalid RecordType not detected")
		t.Fail()
	}

	jString := `{"testType":"fake_accounts"}`
	jStruct := TestRecordType{}
	err := json.NewDecoder(strings.NewReader(jString)).Decode(&jStruct)
	if err == nil {
		t.Error("Unmarshalling should fail for invalid enum values")
		t.Fail()
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Calling String on invalid RecordType value should panic.")
		}
	}()
	_ = rt.String() // This should panic.
}
