package lib_test

import (
	"accountapi/lib"
	"testing"
)

func TestErrors(t *testing.T) {
	e := lib.NewErrorInvalidEnum()
	if !lib.IsErrorInvalidEnum(e) {
		t.Error("ErrorInvalidNum not recognised.")
		t.Fail()
	}

	eAPI := lib.NewErrorAPI(429, "test_api_error")
	if !lib.IsErrorAPI(eAPI) {
		t.Error("ErrorAPI not recognised.")
		t.Fail()
	}
	if eAPI.Error() != "429:test_api_error" {
		t.Errorf("Expected ErrorAPI(test_api_error), got '%s'", eAPI.Error())
		t.Fail()
	}
}
