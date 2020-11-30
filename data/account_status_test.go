package data_test

import (
	"accountapi/data"
	"encoding/json"
	"strings"
	"testing"
)

// TestAccountStatus for testing unmarshalling JSON values.
type TestAccountStatus struct {
	TestAS data.AccountStatus `json:"testAS"`
}

// TestValidAccountStatus verifies proper constraints for consts ("enums"), parsing and unmarshalling.
func TestValidAccountStatus(t *testing.T) {
	as := data.Confirmed
	if !as.IsValid() {
		t.Error("AccountStatus Confirmed should be valid.")
		t.Fail()
	}
	if as.String() != "confirmed" {
		t.Errorf("%s string should be \"confirmed\".", as.String())
		t.Fail()
	}

	jString := `{"testAS":"confirmed"}`
	jStruct := TestAccountStatus{}
	err := json.NewDecoder(strings.NewReader(jString)).Decode(&jStruct)
	if err != nil {
		t.Errorf("Can't unmarshal AccountStatus: %s\n", err.Error())
		t.Fail()
	}
	if jStruct.TestAS != as {
		t.Errorf("Expected AccountStatus value: '%s', got: '%s'\n", as.String(), jStruct.TestAS.String())
		t.Fail()
	}
	b, err := json.Marshal(&jStruct)
	if err != nil {
		t.Errorf("Can't marshal AccountStatus to string: %s\n", err.Error())
		t.Fail()
	} else if string(b) != jString {
		t.Errorf("Expected marshalled value: '%s', got: '%s'\n", jString, string(b))
		t.Fail()
	}
}

// TestInvalidAccountStatus verifies response of functions when called with invalid proper constraints for consts ("enums"), parsing and unmarshalling.
func TestInvalidAccountStatus(t *testing.T) {
	as := data.Business // The last AccountStatus value, when it's increased it should become an invalid value.
	as++
	if as.IsValid() {
		t.Error("Invalid AccountStatus not detected")
		t.Fail()
	}

	jString := `{"testAS":"fake_account_status"}`
	jStruct := TestAccountStatus{}
	err := json.NewDecoder(strings.NewReader(jString)).Decode(&jStruct)
	if err == nil {
		t.Error("Unmarshalling should fail for invalid enum values")
		t.Fail()
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Calling String on invalid AccountStatus value should panic.")
		}
	}()
	_ = as.String() // This should panic.
}
