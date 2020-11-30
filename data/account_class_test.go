package data_test

import (
	"accountapi/data"
	"encoding/json"
	"strings"
	"testing"
)

// TestAccountClass for testing unmarshalling JSON values.
type TestAccountClass struct {
	TestAC data.AccountClass `json:"testAC"`
}

// TestValidAccountClass verifies proper constraints for consts ("enums"), parsing and unmarshalling.
func TestValidAccountClass(t *testing.T) {
	ac := data.Personal
	if !ac.IsValid() {
		t.Errorf("AccountClass Personal should be valid.")
		t.Fail()
	}
	if ac.String() != "Personal" {
		t.Errorf("Personal string should be \"Personal\".")
		t.Fail()
	}

	jString := `{"testAC":"Personal"}`
	jStruct := TestAccountClass{}
	err := json.NewDecoder(strings.NewReader(jString)).Decode(&jStruct)
	if err != nil {
		t.Errorf("Can't unmarshal AccountClass: %s", err.Error())
		t.Fail()
	}
	if jStruct.TestAC != ac {
		t.Errorf("Expected AccountClass value: '%s', got: '%s'", ac.String(), jStruct.TestAC.String())
		t.Fail()
	}
	b, err := json.Marshal(&jStruct)
	if err != nil {
		t.Errorf("Can't marshal AccountClass to string: %s\n", err.Error())
		t.Fail()
	} else if string(b) != jString {
		t.Errorf("Expected marshalled value: '%s', got: '%s'\n", jString, string(b))
		t.Fail()
	}
}

// TestInvalidAccountClass verifies response of functions when called with invalid proper constraints for consts ("enums"), parsing and unmarshalling.
func TestInvalidAccountClass(t *testing.T) {
	ac := data.Business // The last AccountClass value, when it's increased it should become an invalid value.
	ac++
	if ac.IsValid() {
		t.Errorf("Invalid AccountClass not detected")
		t.Fail()
	}

	jString := `{"testAC":"fake_account_class"}`
	jStruct := TestAccountClass{}
	err := json.NewDecoder(strings.NewReader(jString)).Decode(&jStruct)
	if err == nil {
		t.Errorf("Unmarshalling should fail for invalid enum values")
		t.Fail()
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Calling String on invalid AccountClass value should panic.")
		}
	}()
	_ = ac.String() // This should panic.
}
