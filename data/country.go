package data

import (
	"accountapi/lib"
	"encoding/json"

	"github.com/biter777/countries"
)

// CountryCode is a wrapper for countries.CountryCode type that implements UnmarshalJSON.
type CountryCode struct {
	countryCode countries.CountryCode
}

// NewCountryCode from a CountryCode from library countries.
func NewCountryCode(c countries.CountryCode) CountryCode {
	return CountryCode{
		countryCode: c,
	}
}

// String ...
func (c *CountryCode) String() string {
	return c.countryCode.Alpha2()
}

// MarshalJSON converts values to strings.
func (c *CountryCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// UnmarshalJSON ...
func (c *CountryCode) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	s, ok := v.(string)
	if !ok {
		return lib.NewErrorInvalidEnum()
	}

	cc := countries.ByName(s)
	if cc == countries.Unknown {
		return lib.NewErrorInvalidEnum()
	}
	c.countryCode = cc
	return nil
}
