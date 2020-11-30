package data_test

import (
	"accountapi/data"
	"encoding/json"
	"strings"
	"testing"

	"github.com/biter777/countries"
)

// TestCCode for testing unmarshalling JSON values.
type TestCCode struct {
	TestCC data.CountryCode `json:"testCC"`
}

// TestCountryCode verifies proper country codes parsing and unmarshalling.
func TestCountryCode(t *testing.T) {
	cc := data.NewCountryCode(countries.UnitedKingdom)
	if cc.String() != "GB" {
		t.Error("CountryCode string should be \"GB\".")
		t.Fail()
	}

	jString := `{"testCC":"GB"}`
	jStruct := TestCCode{}
	err := json.NewDecoder(strings.NewReader(jString)).Decode(&jStruct)
	if err != nil {
		t.Errorf("Can't unmarshal CountryCode: %s\n", err.Error())
		t.Fail()
	}
	if jStruct.TestCC != cc {
		t.Errorf("Expected CountryCode value: '%s', got: '%s'\n", cc.String(), jStruct.TestCC.String())
		t.Fail()
	}
	b, err := json.Marshal(&jStruct)
	if err != nil {
		t.Errorf("Can't marshal CountryCode to string: %s\n", err.Error())
		t.Fail()
	} else if string(b) != jString {
		t.Errorf("Expected marshalled value: '%s', got: '%s'\n", jString, string(b))
		t.Fail()
	}

	jString = `{"testCC":"fake_country"}`
	err = json.NewDecoder(strings.NewReader(jString)).Decode(&jStruct)
	if err == nil {
		t.Error("Invalid country code unmarshalling should fail.")
		t.Fail()
	}
}
