package data_test

import (
	"accountapi/data"
	"encoding/json"
	"strings"
	"testing"

	"golang.org/x/text/currency"
)

// TestCurrencyStruct for testing unmarshalling JSON values.
type TestCurrencyStruct struct {
	TestCurrency data.Currency `json:"testCurrency"`
}

// TestCurrency verifies proper currency codes parsing and unmarshalling.
func TestCurrency(t *testing.T) {
	cc := data.NewCurrency(currency.GBP)
	if cc.String() != "GBP" {
		t.Error("Currency string should be \"GBP\".")
		t.Fail()
	}

	jString := `{"testCurrency":"GBP"}`
	jStruct := TestCurrencyStruct{}
	err := json.NewDecoder(strings.NewReader(jString)).Decode(&jStruct)
	if err != nil {
		t.Errorf("Can't unmarshal Currency: %s\n", err.Error())
		t.Fail()
	}
	if jStruct.TestCurrency != cc {
		t.Errorf("Expected Currency value: '%s', got: '%s'\n", cc.String(), jStruct.TestCurrency.String())
		t.Fail()
	}
	b, err := json.Marshal(&jStruct)
	if err != nil {
		t.Errorf("Can't marshal Currency to string: %s\n", err.Error())
		t.Fail()
	} else if string(b) != jString {
		t.Errorf("Expected marshalled value: '%s', got: '%s'\n", jString, string(b))
		t.Fail()
	}

	jString = `{"testCurrency":"fake_currency"}`
	err = json.NewDecoder(strings.NewReader(jString)).Decode(&jStruct)
	if err == nil {
		t.Error("Invalid currency unmarshalling should fail.")
		t.Fail()
	}
}
