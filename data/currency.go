package data

import (
	"accountapi/lib"
	"encoding/json"

	"golang.org/x/text/currency"
)

// Currency is a wrapper for currency Unit type that implements UnmarshalJSON.
type Currency struct {
	currency currency.Unit
}

// NewCurrency from a currency.Unit.
func NewCurrency(c currency.Unit) Currency {
	return Currency{
		currency: c,
	}
}

// String ...
func (c *Currency) String() string {
	return c.currency.String()
}

// MarshalJSON converts values to strings.
func (c *Currency) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// UnmarshalJSON ...
func (c *Currency) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	s, ok := v.(string)
	if !ok {
		return lib.NewErrorInvalidEnum()
	}

	u, err := currency.ParseISO(s)
	if err != nil {
		return err
	}
	c.currency = u
	return nil
}
