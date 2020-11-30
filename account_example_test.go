package account_test

import (
	account "accountapi"
	"accountapi/data"
	"accountapi/lib"
	"fmt"

	"github.com/biter777/countries"
	"github.com/google/uuid"
	"golang.org/x/text/currency"
)

func Example() {
	client, err := account.New(account.Config{
		Server:             "http://127.0.0.1:8080",
		MaxConnections:     64,
		MaxIdleConnections: 16,
		Timeout:            3,
	})
	id, _ := uuid.NewRandom()
	org, _ := uuid.NewRandom()
	acc := data.Account{
		ID:             id,
		OrganisationID: org,
		Attributes: data.Attributes{
			Country:      data.NewCountryCode(countries.UnitedKingdom),
			BaseCurrency: data.NewCurrency(currency.GBP),
		},
	}
	createdAccount, err := client.Create(&acc)
	if err != nil {
		fmt.Printf("Error creating account: %v", err)
	}
	fmt.Printf("Account %s created", createdAccount.ID.String())
	fetchedAccount, err := client.Fetch(acc.ID)
	if err != nil {
		fmt.Printf("Error fetching account: %v", err)
	}
	fmt.Printf("Account %s fetched", fetchedAccount.ID.String())
	listedAccounts, err := client.List(lib.PageNumber(0), lib.PageSize(10))
	if err != nil {
		fmt.Printf("Error listing accounts: %v", err)
	}
	fmt.Printf("%d accounts listed", len(*listedAccounts))
	err = client.Delete(acc.ID, 0) // Delete version 0.
	if err != nil {
		fmt.Printf("Error deleting account: %v", err)
	}
}
