package data

// Attributes of an account as defined in https://api-docs.form3.tech/api.html#organisation-accounts-create.
type Attributes struct {
	Country                 CountryCode   `json:"country"`
	BaseCurrency            Currency      `json:"base_currency"`
	AccountNumber           string        `json:"account_number"`
	BankID                  string        `json:"bank_id"`
	BankIDCode              string        `json:"bank_id_code"`
	BIC                     string        `json:"bic"`
	IBAN                    string        `json:"iban"`
	Name                    []string      `json:"name"`
	AlternativeNames        []string      `json:"alternative_names"`
	AccountClassification   AccountClass  `json:"account_classification"`
	JointAccount            bool          `json:"joint_account"`
	AccountMatchingOptOut   bool          `json:"account_matching_opt_out"`
	SecondaryIdentification string        `json:"secondary_identification"`
	Switched                bool          `json:"switched"`
	Status                  AccountStatus `json:"status"`
}
